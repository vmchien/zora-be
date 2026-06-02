package produce

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type Mode int

const (
	ModeSync Mode = iota
	ModeAsync
)

type Result struct {
	Topic     string
	Partition int32
	Offset    int64
	Timestamp time.Time
}

type Config struct {
	Brokers  []string
	ClientID string
	Mode     Mode

	ProduceTimeout time.Duration

	Workers   int
	QueueSize int

	RetryMaxAttempts int
	RetryBaseBackoff time.Duration
	RetryMaxBackoff  time.Duration
	RetryJitter      time.Duration

	Linger        time.Duration
	BatchMaxBytes int
	Compression   string // none | snappy | gzip | lz4 | zstd

	EnableTLS     bool
	TLSConfig     *tls.Config
	SASLMechanism string // "SCRAM-SHA-256" | "SCRAM-SHA-512" | ""
	SASLUser      string
	SASLPass      string

	Logger func(ctx context.Context, format string, args ...any)
}

func (c *Config) setDefaults() {
	if c.Logger == nil {
		c.Logger = printf
	}
	if c.ClientID == "" {
		c.ClientID = "producer"
	}
	if c.Mode != ModeSync && c.Mode != ModeAsync {
		c.Mode = ModeAsync
	}
	if c.ProduceTimeout <= 0 {
		c.ProduceTimeout = 10 * time.Second
	}
	if c.Workers <= 0 {
		c.Workers = 8
	}
	if c.QueueSize <= 0 {
		c.QueueSize = 8192
	}
	if c.RetryMaxAttempts <= 0 {
		c.RetryMaxAttempts = 5
	}
	if c.RetryBaseBackoff <= 0 {
		c.RetryBaseBackoff = 200 * time.Millisecond
	}
	if c.RetryMaxBackoff <= 0 {
		c.RetryMaxBackoff = 5 * time.Second
	}
	if c.RetryJitter <= 0 {
		c.RetryJitter = 50 * time.Millisecond
	}
	if c.Linger <= 0 {
		c.Linger = 20 * time.Millisecond
	}
	if c.BatchMaxBytes <= 0 {
		c.BatchMaxBytes = 1 << 20 // 1MB
	}
	if c.Compression == "" {
		c.Compression = "none"
	}
	if c.TLSConfig == nil {
		c.TLSConfig = new(tls.Config)
	}
}

func (c *Config) validate() error {
	c.setDefaults()
	if len(c.Brokers) == 0 {
		return errors.New("brokers required")
	}
	return nil
}

// Client abstracts the subset of *kgo.Client we use.
type Client interface {
	Produce(ctx context.Context, r *kgo.Record, promise func(*kgo.Record, error))
	Flush(ctx context.Context) error
	Close()
}

type Helper struct {
	cfg Config
	cl  Client

	jobCh chan job
	wg    sync.WaitGroup

	closeOnce sync.Once
	closed    atomic.Bool

	sent   uint64
	failed uint64
}

type job struct {
	ctx    context.Context
	record *kgo.Record
	ack    chan result // nil for no-wait
}

type result struct {
	res Result
	err error
}

func New(cfg Config) (*Helper, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),

		// prod-safe baseline
		kgo.RequiredAcks(kgo.AllISRAcks()),

		kgo.ProducerLinger(cfg.Linger),
		kgo.ProducerBatchMaxBytes(int32(cfg.BatchMaxBytes)),
	}

	switch strings.ToLower(strings.TrimSpace(cfg.Compression)) {
	case "none":
	case "snappy":
		opts = append(opts, kgo.ProducerBatchCompression(kgo.SnappyCompression()))
	case "gzip":
		opts = append(opts, kgo.ProducerBatchCompression(kgo.GzipCompression()))
	case "lz4":
		opts = append(opts, kgo.ProducerBatchCompression(kgo.Lz4Compression()))
	case "zstd":
		opts = append(opts, kgo.ProducerBatchCompression(kgo.ZstdCompression()))
	default:
		return nil, fmt.Errorf("unsupported compression: %s", cfg.Compression)
	}

	if cfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(cfg.TLSConfig))
	}

	if cfg.SASLMechanism != "" && cfg.SASLUser != "" {
		mech, err := scramMechanism(cfg.SASLMechanism, cfg.SASLUser, cfg.SASLPass)
		if err != nil {
			return nil, err
		}
		opts = append(opts, kgo.SASL(mech))
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return NewWithClient(cfg, cl), nil
}

func NewWithClient(cfg Config, cl Client) *Helper {
	cfg.setDefaults()
	h := &Helper{
		cfg:   cfg,
		cl:    cl,
		jobCh: make(chan job, cfg.QueueSize),
	}
	if cfg.Mode == ModeAsync {
		h.startWorkers()
	}
	return h
}

func (h *Helper) Close() {
	h.closeOnce.Do(func() {
		h.closed.Store(true)
		if h.cfg.Mode == ModeAsync {
			close(h.jobCh)
			h.wg.Wait()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		_ = h.cl.Flush(ctx)
		cancel()

		h.cl.Close()
	})
}

func (h *Helper) Produce(ctx context.Context, r *kgo.Record) (Result, error) {
	if r == nil || strings.TrimSpace(r.Topic) == "" {
		return Result{}, errors.New("record and record.Topic required")
	}
	if h.closed.Load() {
		return Result{}, errors.New("producer closed")
	}

	switch h.cfg.Mode {
	case ModeSync:
		return h.produceWithRetry(ctx, r)
	case ModeAsync:
		return h.produceAsyncAwait(ctx, r)
	default:
		return Result{}, fmt.Errorf("unknown mode")
	}
}

func (h *Helper) ProduceAsyncNoWait(ctx context.Context, r *kgo.Record) error {
	if h.cfg.Mode != ModeAsync {
		return errors.New("ProduceAsyncNoWait requires ModeAsync")
	}
	if r == nil || strings.TrimSpace(r.Topic) == "" {
		return errors.New("record and record.Topic required")
	}
	if h.closed.Load() {
		return errors.New("producer closed")
	}

	select {
	case h.jobCh <- job{ctx: ctx, record: r, ack: nil}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.New("producer queue full")
	}
}

func (h *Helper) produceWithRetry(ctx context.Context, r *kgo.Record) (Result, error) {
	var out Result

	err := h.withRetry(ctx, func(attemptCtx context.Context) error {
		done := make(chan error, 1)

		h.cl.Produce(attemptCtx, r, func(rec *kgo.Record, err error) {
			if err == nil && rec != nil {
				out = Result{
					Topic:     rec.Topic,
					Partition: rec.Partition,
					Offset:    rec.Offset,
					Timestamp: rec.Timestamp,
				}
			}
			done <- err
		})

		select {
		case err := <-done:
			return err
		case <-attemptCtx.Done():
			return attemptCtx.Err()
		}
	})

	if err != nil {
		atomic.AddUint64(&h.failed, 1)
		return Result{}, err
	}

	atomic.AddUint64(&h.sent, 1)
	return out, nil
}

func (h *Helper) withRetry(ctx context.Context, fn func(context.Context) error) error {
	var last error

	for attempt := 1; attempt <= h.cfg.RetryMaxAttempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, h.cfg.ProduceTimeout)
		err := fn(attemptCtx)
		cancel()

		if err == nil {
			return nil
		}
		last = err

		if ctx.Err() != nil {
			return ctx.Err()
		}
		if attempt == h.cfg.RetryMaxAttempts || !isRetriableProduceError(err) {
			return err
		}

		sleep := backoff(h.cfg.RetryBaseBackoff, h.cfg.RetryMaxBackoff, h.cfg.RetryJitter, attempt)
		h.cfg.Logger(ctx, "produce retry attempt=%d backoff=%s err=%v", attempt, sleep, err)

		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return last
}

func (h *Helper) startWorkers() {
	for i := 0; i < h.cfg.Workers; i++ {
		h.wg.Add(1)
		go func() {
			defer h.wg.Done()
			for j := range h.jobCh {
				res, err := h.produceWithRetry(j.ctx, j.record)

				if j.ack != nil {
					j.ack <- result{res: res, err: err}
				} else if err != nil {
					h.cfg.Logger(context.Background(), "produce async error topic=%s err=%v", j.record.Topic, err)
				}
			}
		}()
	}
}

func (h *Helper) produceAsyncAwait(ctx context.Context, r *kgo.Record) (Result, error) {
	ack := make(chan result, 1)

	select {
	case h.jobCh <- job{ctx: ctx, record: r, ack: ack}:
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}

	select {
	case res := <-ack:
		return res.res, res.err
	case <-ctx.Done():
		return Result{}, ctx.Err()
	}
}

// //////////////////////////////////////////////////////////////////////////////
// Retriable detection (robust across franz-go versions)
// //////////////////////////////////////////////////////////////////////////////

func isRetriableProduceError(err error) bool {
	// 1) Best-effort: franz-go helper.
	// Some versions may not mark all leader-move errors retriable, so we add allowlist below.
	if kerr.IsRetriable(err) {
		return true
	}

	// 2) Hard allowlist for common retry-safe broker errors (leader changes / metadata churn).
	// These are safe to retry and frequently observed in production.
	if errors.Is(err, kerr.NotLeaderForPartition) ||
		errors.Is(err, kerr.LeaderNotAvailable) ||
		errors.Is(err, kerr.UnknownTopicOrPartition) ||
		errors.Is(err, kerr.RequestTimedOut) ||
		errors.Is(err, kerr.NotCoordinator) ||
		errors.Is(err, kerr.CoordinatorNotAvailable) ||
		errors.Is(err, kerr.RebalanceInProgress) {
		return true
	}

	// 3) Fallback: network-ish transient strings.
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "unable to dial") ||
		strings.Contains(msg, "i/o timeout") ||
		strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "not leader")
}

func backoff(base, max, jitter time.Duration, attempt int) time.Duration {
	d := float64(base) * math.Pow(2, float64(attempt-1))
	if d > float64(max) {
		d = float64(max)
	}
	if jitter > 0 {
		j := time.Duration(time.Now().UnixNano()%int64(jitter+1)) - jitter/2
		d += float64(j)
		if d < 0 {
			d = 0
		}
	}
	return time.Duration(d)
}

func scramMechanism(mech, user, pass string) (sasl.Mechanism, error) {
	switch strings.ToUpper(strings.TrimSpace(mech)) {
	case "SCRAM-SHA-256":
		return scram.Auth{User: user, Pass: pass}.AsSha256Mechanism(), nil
	case "SCRAM-SHA-512":
		return scram.Auth{User: user, Pass: pass}.AsSha512Mechanism(), nil
	default:
		return nil, fmt.Errorf("unsupported sasl mechanism: %s", mech)
	}
}

// printf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func printf(ctx context.Context, format string, v ...any) {
	log.Printf(format, v)
}
