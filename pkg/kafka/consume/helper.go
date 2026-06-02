package consume

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"vn.vato.zora.be.api/pkg/logs"
)

// //////////////////////////////////////////////////////////////////////////////
// Public API
// //////////////////////////////////////////////////////////////////////////////

// Mode controls processing concurrency.
//
//   - ModeSync: sequential in polling goroutine.
//   - ModeAsync: concurrent workers while preserving per-partition ordering.
type Mode int

const (
	ModeSync Mode = iota
	ModeAsync
)

var helperTracer = otel.Tracer("vn.vato.zora.be.api/pkg/kafka/consume")

// LogLevel is a logger level.
type LogLevel int8

const (
	// LevelDebug is logger debug level.
	LevelDebug LogLevel = iota
	// LevelInfo is logger info level.
	LevelInfo
	// LevelWarn is logger warn level.
	LevelWarn
	// LevelError is logger error level.
	LevelError
)

// CommitMode controls offset commit strategy.
//
//   - CommitManual: commit only after successful processing (or DLQ success). Safe default.
//   - CommitAuto: franz-go auto commit (unsafe with async unless explicitly allowed).
type CommitMode int

const (
	CommitAuto CommitMode = iota
	CommitManual
)

// FailurePolicy controls what to do when processing still fails
// after RetryPolicy attempts.
type FailurePolicy int

const (
	FailurePolicyStop FailurePolicy = iota
	FailurePolicyRetryForever
	FailurePolicyCommitOnError
	FailurePolicyDLQThenCommit
)

type CommitPolicy struct {
	Mode CommitMode

	// ManualEvery flush interval for manual commits.
	ManualEvery time.Duration

	// ManualBatch flush after N successes across all partitions.
	// Set to 0 to disable batch flush.
	ManualBatch int

	// OnFailure controls the terminal behavior when processing still fails
	// after RetryPolicy attempts.
	OnFailure FailurePolicy
}

// RetryPolicy controls processing retries.
type RetryPolicy struct {
	MaxAttempts   int
	BaseBackoff   time.Duration
	MaxBackoff    time.Duration
	Jitter        time.Duration
	PoisonBackoff time.Duration // when poison is blocking a partition
}

// Processor is the user-provided handler.
type Processor interface {
	Process(ctx context.Context, r *kgo.Record) error
}

type ProcessFunc func(ctx context.Context, r *kgo.Record) error

func (f ProcessFunc) Process(ctx context.Context, r *kgo.Record) error {
	return f(ctx, r)
}

// DLQHandler handles poison after retries.
type DLQHandler interface {
	HandleDLQ(ctx context.Context, r *kgo.Record, cause error) error
}

// Config is the main configuration. It is designed to be "simple by default":
// only Brokers, Topics, Group are required; the rest has safe defaults.
type Config struct {
	// Required:
	Brokers []string
	Topics  []string
	Group   string

	// Optional:
	ClientID string

	Mode   Mode
	Commit CommitPolicy

	// Guardrail: Async + AutoCommit can lose messages on crash.
	AllowRiskyAsyncAutoCommit bool

	// Async:
	Workers   int
	QueueSize int

	// Processing:
	ProcessTimeout time.Duration

	// Poll / Fetch:
	MaxPollRecords int
	FetchMaxBytes  int
	FetchMaxWait   time.Duration

	// Connectivity / stability:
	// Backoff applied when poll returns dial/DNS/join-group errors, to prevent log spam / hot loop.
	PollErrorBackoff    time.Duration
	PollErrorMaxBackoff time.Duration

	// JoinFailFast closes the consumer if we cannot obtain any assignment within JoinFailFastTimeout.
	// This makes local/dev failures (DNS/VPN/private routing/auth) obvious rather than silently idle.
	JoinFailFast        bool
	JoinFailFastTimeout time.Duration

	// CommitTimeout bounds time a commit attempt may block.
	CommitTimeout time.Duration

	// Memory control (Async partition coordinator):
	// Limits buffered records per partition. If exceeded, coordinator will backpressure via poll loop.
	MaxPartitionQueue int

	// Security:
	EnableTLS     bool
	TLSConfig     *tls.Config
	SASLMechanism string // "SCRAM-SHA-256" | "SCRAM-SHA-512" | ""
	SASLUser      string
	SASLPass      string

	Retry RetryPolicy

	// LogHelper is preferred logger for level-aware + trace-aware logging.
	LogHelper *logs.Helper

	// Logger is a legacy printf-style fallback logger.
	Logger func(ctx context.Context, format string, args ...any)
}

// SimpleConfig is the minimal config surface most services should use.
type SimpleConfig struct {
	Brokers []string
	Topics  []string
	Group   string
}

// NewSimple is the recommended constructor for most services.
// It applies safe defaults: Async + Manual commit + join fail-fast + backoff + retry.
func NewSimple(sc SimpleConfig, p Processor, dlq DLQHandler) (*Helper, error) {
	cfg := Config{
		Brokers:      sc.Brokers,
		Topics:       sc.Topics,
		Group:        sc.Group,
		JoinFailFast: true,
	}
	// Defaults will be applied in validate().
	return New(cfg, p, dlq)
}

func (c *Config) setDefaults() {
	if c.LogHelper == nil && c.Logger == nil {
		c.Logger = printf
	}
	if c.ClientID == "" {
		c.ClientID = "consumer"
	}

	// Defaults: safe-by-default
	if c.Mode != ModeSync && c.Mode != ModeAsync {
		c.Mode = ModeAsync
	}
	if c.Commit.Mode != CommitAuto && c.Commit.Mode != CommitManual {
		c.Commit.Mode = CommitManual
	}
	if c.Commit.OnFailure < FailurePolicyStop || c.Commit.OnFailure > FailurePolicyDLQThenCommit {
		c.Commit.OnFailure = FailurePolicyStop
	}
	if c.Workers <= 0 {
		c.Workers = 16
	}
	if c.QueueSize <= 0 {
		c.QueueSize = 8192
	}
	if c.ProcessTimeout <= 0 {
		c.ProcessTimeout = 30 * time.Second
	}

	if c.MaxPollRecords <= 0 {
		c.MaxPollRecords = 2000
	}
	if c.FetchMaxBytes <= 0 {
		c.FetchMaxBytes = 32 << 20
	}
	if c.FetchMaxWait <= 0 {
		c.FetchMaxWait = 250 * time.Millisecond
	}

	// Manual commit defaults
	if c.Commit.ManualEvery <= 0 {
		c.Commit.ManualEvery = 2 * time.Second
	}
	if c.Commit.ManualBatch <= 0 {
		c.Commit.ManualBatch = 2000
	}
	if c.CommitTimeout <= 0 {
		c.CommitTimeout = 5 * time.Second
	}

	// Connectivity defaults
	if c.PollErrorBackoff <= 0 {
		c.PollErrorBackoff = 2 * time.Second
	}
	if c.PollErrorMaxBackoff <= 0 {
		c.PollErrorMaxBackoff = 10 * time.Second
	}
	if c.JoinFailFastTimeout <= 0 {
		c.JoinFailFastTimeout = 20 * time.Second
	}

	// Coordinator defaults
	if c.MaxPartitionQueue <= 0 {
		c.MaxPartitionQueue = 5000
	}

	if c.TLSConfig == nil {
		c.TLSConfig = new(tls.Config)
	}

	// Retry defaults
	if c.Retry.MaxAttempts <= 0 {
		c.Retry.MaxAttempts = 5
	}
	if c.Retry.BaseBackoff <= 0 {
		c.Retry.BaseBackoff = 200 * time.Millisecond
	}
	if c.Retry.MaxBackoff <= 0 {
		c.Retry.MaxBackoff = 5 * time.Second
	}
	if c.Retry.Jitter <= 0 {
		c.Retry.Jitter = 50 * time.Millisecond
	}
	if c.Retry.PoisonBackoff <= 0 {
		c.Retry.PoisonBackoff = 2 * time.Second
	}

	// JoinFailFast is configured explicitly by caller.
}

func (c *Config) validate() error {
	c.setDefaults()

	if len(c.Brokers) == 0 {
		return errors.New("Brokers required")
	}
	if len(c.Topics) == 0 {
		return errors.New("Topics required")
	}
	if c.Group == "" {
		return errors.New("Group required")
	}
	if c.Mode == ModeAsync && c.Commit.Mode == CommitAuto && !c.AllowRiskyAsyncAutoCommit {
		return errors.New("async + auto commit is risky (can lose messages). Set AllowRiskyAsyncAutoCommit=true to proceed")
	}
	return nil
}

func (c *Config) log(ctx context.Context, level LogLevel, format string, args ...any) {
	if c.LogHelper != nil {
		switch level {
		case LevelInfo:
			c.LogHelper.Infof(ctx, format, args...)
		case LevelWarn:
			c.LogHelper.Warnf(ctx, format, args...)
		case LevelError:
			c.LogHelper.Errorf(ctx, format, args...)
		default:
			c.LogHelper.Debugf(ctx, format, args...)
		}
		return
	}
	if c.Logger == nil {
		return
	}
	c.Logger(ctx, format, args...)
}

// //////////////////////////////////////////////////////////////////////////////
// Testability abstractions
// //////////////////////////////////////////////////////////////////////////////

// Source abstracts polling records. Production uses kgo.Client.
type Source interface {
	Poll(ctx context.Context) (records []*kgo.Record, fatal error)
	Close()
	IsClosed() bool
}

// Committer abstracts committing offsets. Production uses kgo.Client.
type Committer interface {
	CommitOffsets(ctx context.Context, offsets map[string]map[int32]kgo.EpochOffset) error
	Close()
}

// //////////////////////////////////////////////////////////////////////////////
// Helper
// //////////////////////////////////////////////////////////////////////////////

type Helper struct {
	cfg Config
	src Source
	cmt Committer
	p   Processor
	dlq DLQHandler

	// Async pipeline:
	jobCh     chan *kgo.Record
	doneCh    chan workDone
	commitAck chan *kgo.Record

	shutdownCh chan struct{}
	closeOnce  sync.Once

	processed uint64
	failed    uint64
	panics    uint64

	// assignment diagnostics
	assignedOnce    *atomic.Bool
	joinFailFastErr atomic.Value // stores error
}

type workDone struct {
	r   *kgo.Record
	err error
}

func New(cfg Config, p Processor, dlq DLQHandler) (*Helper, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("processor required")
	}
	if cfg.Commit.OnFailure == FailurePolicyDLQThenCommit && dlq == nil {
		return nil, errors.New("commit.onFailure=FailurePolicyDLQThenCommit requires DLQ handler")
	}

	assignedOnce := &atomic.Bool{}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.Group),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.ClientID(cfg.ClientID),

		// Cooperative sticky reduces rebalance churn for many workloads.
		kgo.Balancers(kgo.CooperativeStickyBalancer()),

		// Assignment logs: invaluable for ops.
		kgo.OnPartitionsAssigned(func(cctx context.Context, _ *kgo.Client, m map[string][]int32) {
			assignedOnce.Store(true)
			cfg.log(cctx, LevelInfo, "kafka assigned group=%s topics=%v", cfg.Group, m)
		}),
		kgo.OnPartitionsRevoked(func(cctx context.Context, _ *kgo.Client, m map[string][]int32) {
			cfg.log(cctx, LevelWarn, "kafka revoked group=%s topics=%v", cfg.Group, m)
		}),

		kgo.FetchMaxBytes(int32(cfg.FetchMaxBytes)),
		kgo.FetchMaxWait(cfg.FetchMaxWait),

		// Default reset is end; services that need replay should override by changing this opt.
		kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()),
	}

	// Manual commit: disable auto commit.
	if cfg.Commit.Mode == CommitManual {
		opts = append(opts, kgo.DisableAutoCommit())
	}

	// TLS
	if cfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(cfg.TLSConfig))
	}

	// SASL SCRAM
	if cfg.SASLMechanism != "" && cfg.SASLUser != "" {
		var mech sasl.Mechanism
		switch strings.ToUpper(cfg.SASLMechanism) {
		case "SCRAM-SHA-256":
			mech = scram.Auth{User: cfg.SASLUser, Pass: cfg.SASLPass}.AsSha256Mechanism()
		case "SCRAM-SHA-512":
			mech = scram.Auth{User: cfg.SASLUser, Pass: cfg.SASLPass}.AsSha512Mechanism()
		default:
			return nil, fmt.Errorf("unsupported SASL mechanism: %s", cfg.SASLMechanism)
		}
		opts = append(opts, kgo.SASL(mech))
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	src := &kgoSource{
		cl:          cl,
		logger:      cfg.log,
		maxPoll:     cfg.MaxPollRecords,
		backoffBase: cfg.PollErrorBackoff,
		backoffMax:  cfg.PollErrorMaxBackoff,
	}
	cmt := &kgoCommitter{cl: cl}

	h, err := NewWith(cfg, p, dlq, src, cmt)
	if err != nil {
		return nil, err
	}
	h.assignedOnce = assignedOnce
	return h, nil
}

func NewWith(cfg Config, p Processor, dlq DLQHandler, src Source, cmt Committer) (*Helper, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("processor required")
	}
	if cfg.Commit.OnFailure == FailurePolicyDLQThenCommit && dlq == nil {
		return nil, errors.New("commit.onFailure=FailurePolicyDLQThenCommit requires DLQ handler")
	}
	if src == nil || cmt == nil {
		return nil, errors.New("source and Committer required")
	}

	return &Helper{
		cfg:          cfg,
		src:          src,
		cmt:          cmt,
		p:            p,
		dlq:          dlq,
		jobCh:        make(chan *kgo.Record, cfg.QueueSize),
		doneCh:       make(chan workDone, cfg.QueueSize),
		commitAck:    make(chan *kgo.Record, cfg.QueueSize),
		shutdownCh:   make(chan struct{}),
		assignedOnce: &atomic.Bool{},
	}, nil
}

func (h *Helper) Stats() (processed, failed, panics uint64) {
	return atomic.LoadUint64(&h.processed), atomic.LoadUint64(&h.failed), atomic.LoadUint64(&h.panics)
}

func (h *Helper) Close() {
	h.closeOnce.Do(func() {
		h.src.Close()
		h.cmt.Close()
	})
}

func (h *Helper) Run(ctx context.Context) error {
	defer h.Close()

	if h.assignedOnce == nil {
		h.assignedOnce = &atomic.Bool{}
	}

	// Join fail-fast: close source if no assignment within timeout.
	if h.cfg.JoinFailFast {
		timer := time.NewTimer(h.cfg.JoinFailFastTimeout)
		go func() {
			defer timer.Stop()
			select {
			case <-ctx.Done():
				return
			case <-h.shutdownCh:
				return
			case <-timer.C:
				if !h.assignedOnce.Load() {
					err := fmt.Errorf("kafka join fail-fast: no partitions assigned within %s (check DNS/VPN/private routing/auth)", h.cfg.JoinFailFastTimeout)
					h.joinFailFastErr.Store(err)
					h.logPrinter(ctx, LevelError, "%v", err)
					h.src.Close()
				}
			}
		}()
	}

	// Manual commit loop
	var commitErrCh chan error
	if h.cfg.Commit.Mode == CommitManual {
		commitErrCh = make(chan error, 1)
		go func() { commitErrCh <- h.commitLoop(ctx) }()
	}

	switch h.cfg.Mode {
	case ModeSync:
		err := h.runSync(ctx)
		close(h.shutdownCh)
		if h.cfg.Commit.Mode == CommitManual {
			close(h.commitAck)
			_ = <-commitErrCh
		}
		if v := h.joinFailFastErr.Load(); v != nil {
			if e, ok := v.(error); ok {
				return e
			}
		}
		return err

	case ModeAsync:
		err := h.runAsync(ctx, commitErrCh)
		if v := h.joinFailFastErr.Load(); v != nil {
			if e, ok := v.(error); ok {
				return e
			}
		}
		return err

	default:
		close(h.shutdownCh)
		if h.cfg.Commit.Mode == CommitManual {
			close(h.commitAck)
			_ = <-commitErrCh
		}
		if v := h.joinFailFastErr.Load(); v != nil {
			if e, ok := v.(error); ok {
				return e
			}
		}
		return fmt.Errorf("unknown mode")
	}
}

// //////////////////////////////////////////////////////////////////////////////
// Sync mode
// //////////////////////////////////////////////////////////////////////////////

func (h *Helper) runSync(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-h.shutdownCh:
			return nil
		default:
		}

		recs, fatal := h.src.Poll(ctx)
		if fatal != nil {
			if errors.Is(fatal, context.Canceled) || errors.Is(fatal, context.DeadlineExceeded) {
				return nil
			}
			return fatal
		}
		if h.src.IsClosed() {
			return nil
		}

		for _, r := range recs {
			if r == nil {
				h.logPrinter(ctx, LevelWarn, "skip nil record from source (sync)")
				continue
			}

			shouldCommit, err := h.processSyncRecord(ctx, r)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				return err
			}

			if h.cfg.Commit.Mode == CommitManual && shouldCommit {
				select {
				case h.commitAck <- r:
				case <-ctx.Done():
					return nil
				case <-h.shutdownCh:
					return nil
				}
			}
		}
	}
}

func (h *Helper) processSyncRecord(ctx context.Context, r *kgo.Record) (bool, error) {
	for {
		err := h.safeProcess(ctx, r)
		if err == nil {
			return true, nil
		}

		shouldCommit, retry, terminalErr := h.handleProcessFailure(ctx, r, err)
		if retry {
			select {
			case <-time.After(h.cfg.Retry.PoisonBackoff):
				continue
			case <-ctx.Done():
				return false, ctx.Err()
			case <-h.shutdownCh:
				return false, nil
			}
		}
		if terminalErr != nil {
			h.logPrinter(
				ctx,
				LevelError,
				"sync processing failed topic=%s partition=%d offset=%d: %v",
				r.Topic, r.Partition, r.Offset, terminalErr,
			)
			return false, terminalErr
		}
		return shouldCommit, nil
	}
}

// //////////////////////////////////////////////////////////////////////////////
// Async mode
// //////////////////////////////////////////////////////////////////////////////

func (h *Helper) runAsync(ctx context.Context, commitErrCh chan error) error {
	var wg sync.WaitGroup
	pc := newPartitionCoordinator(h.cfg.log, h.cfg.Retry.PoisonBackoff, h.cfg.MaxPartitionQueue)
	workerErrCh := make(chan error, 1)

	for i := 0; i < h.cfg.Workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			h.worker(ctx, id, workerErrCh)
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		pc.run(ctx, h.jobCh, h.doneCh)
	}()

	pollErrCh := make(chan error, 1)
	go func() { pollErrCh <- h.pollToCoordinator(ctx, pc) }()

	var err error
	select {
	case <-ctx.Done():
	case err = <-pollErrCh:
	case err = <-workerErrCh:
	case err = <-commitErrCh:
	}

	close(h.shutdownCh)
	pc.stop()
	h.src.Close()
	wg.Wait()

	if h.cfg.Commit.Mode == CommitManual {
		close(h.commitAck)
		_ = <-commitErrCh
	}
	return err
}

func (h *Helper) pollToCoordinator(ctx context.Context, pc *partitionCoordinator) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-h.shutdownCh:
			return nil
		default:
		}

		recs, fatal := h.src.Poll(ctx)
		if fatal != nil {
			if errors.Is(fatal, context.Canceled) || errors.Is(fatal, context.DeadlineExceeded) {
				return nil
			}
			return fatal
		}
		if h.src.IsClosed() {
			return nil
		}

		for _, r := range recs {
			if r == nil {
				h.logPrinter(ctx, LevelWarn, "skip nil record from source (async poll)")
				continue
			}

			// Safety-first: never drop fetched records. Keep retrying enqueue on backpressure.
			for {
				if err := pc.enqueue(r); err != nil {
					if !errors.Is(err, ErrPartitionQueueFull) {
						return fmt.Errorf("enqueue record topic=%s partition=%d offset=%d: %w", r.Topic, r.Partition, r.Offset, err)
					}
					h.logPrinter(ctx, LevelWarn, "partition queue full, backing off: %v", err)
					select {
					case <-time.After(200 * time.Millisecond):
					case <-ctx.Done():
						return nil
					case <-h.shutdownCh:
						return nil
					}
					continue
				}
				break
			}
		}
	}
}

func (h *Helper) worker(ctx context.Context, id int, workerErrCh chan<- error) {
	_ = id
	for {
		select {
		case <-ctx.Done():
			return
		case <-h.shutdownCh:
			return
		case r, ok := <-h.jobCh:
			if !ok {
				return
			}
			if r == nil {
				h.logPrinter(ctx, LevelWarn, "skip nil record from job queue")
				continue
			}

			processErr := h.safeProcess(ctx, r)
			shouldCommit := processErr == nil
			doneErr := processErr

			if processErr != nil {
				commitOnFailure, retry, terminalErr := h.handleProcessFailure(ctx, r, processErr)
				if retry {
					doneErr = processErr
				} else if terminalErr != nil {
					select {
					case workerErrCh <- fmt.Errorf("async processing failed topic=%s partition=%d offset=%d: %w", r.Topic, r.Partition, r.Offset, terminalErr):
					default:
					}
					return
				} else {
					shouldCommit = commitOnFailure
					doneErr = nil
				}
			}

			if h.cfg.Commit.Mode == CommitManual && shouldCommit {
				select {
				case h.commitAck <- r:
				case <-ctx.Done():
					return
				case <-h.shutdownCh:
					return
				}
			}

			select {
			case h.doneCh <- workDone{r: r, err: doneErr}:
			case <-ctx.Done():
				return
			case <-h.shutdownCh:
				return
			}
		}
	}
}

func (h *Helper) handleProcessFailure(ctx context.Context, r *kgo.Record, procErr error) (commit bool, retry bool, terminalErr error) {
	switch h.cfg.Commit.OnFailure {
	case FailurePolicyRetryForever:
		h.logPrinter(
			ctx,
			LevelError,
			"processing failed topic=%s partition=%d offset=%d: %v (retry forever)",
			r.Topic, r.Partition, r.Offset, procErr,
		)
		return false, true, nil

	case FailurePolicyCommitOnError:
		h.logPrinter(
			ctx,
			LevelWarn,
			"processing failed topic=%s partition=%d offset=%d: %v (commit on error)",
			r.Topic, r.Partition, r.Offset, procErr,
		)
		return true, false, nil

	case FailurePolicyDLQThenCommit:
		if h.dlq == nil {
			return false, false, fmt.Errorf("onFailure=FailurePolicyDLQThenCommit requires DLQ handler: %w", procErr)
		}
		dlqErr := h.dlq.HandleDLQ(ctx, r, procErr)
		if dlqErr != nil {
			return false, false, fmt.Errorf("process=%v; dlq=%v", procErr, dlqErr)
		}
		h.logPrinter(
			ctx,
			LevelWarn,
			"processing failed topic=%s partition=%d offset=%d: %v (dlq success, commit)",
			r.Topic, r.Partition, r.Offset, procErr,
		)
		return true, false, nil

	case FailurePolicyStop:
		fallthrough
	default:
		return false, false, procErr
	}
}

// //////////////////////////////////////////////////////////////////////////////
// Processing
// //////////////////////////////////////////////////////////////////////////////

func (h *Helper) safeProcess(ctx context.Context, r *kgo.Record) (err error) {
	if r == nil {
		return errors.New("nil record")
	}

	ctx, span := helperTracer.Start(
		ctx,
		"kafka.consume.process",
		trace.WithAttributes(
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation", "process"),
			attribute.String("messaging.destination.name", r.Topic),
			attribute.Int64("messaging.kafka.partition", int64(r.Partition)),
			attribute.Int64("messaging.kafka.offset", r.Offset),
			attribute.String("messaging.consumer.group.name", h.cfg.Group),
		),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	defer func() {
		if rec := recover(); rec != nil {
			atomic.AddUint64(&h.panics, 1)
			h.logPrinter(
				ctx,
				LevelError,
				"PANIC recovered topic=%s partition=%d offset=%d err=%v\n%s",
				r.Topic, r.Partition, r.Offset, rec, string(debug.Stack()),
			)
			err = fmt.Errorf("panic: %v", rec)
		}
	}()

	for attempt := 1; attempt <= h.cfg.Retry.MaxAttempts; attempt++ {
		pctx, cancel := context.WithTimeout(ctx, h.cfg.ProcessTimeout)
		e := h.p.Process(pctx, r)
		cancel()

		if e == nil {
			atomic.AddUint64(&h.processed, 1)
			return nil
		}

		if attempt == h.cfg.Retry.MaxAttempts {
			atomic.AddUint64(&h.failed, 1)
			return e
		}

		sleep := expBackoff(h.cfg.Retry.BaseBackoff, h.cfg.Retry.MaxBackoff, h.cfg.Retry.Jitter, attempt)
		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return errors.New("unreachable")
}

func expBackoff(base, max, jitter time.Duration, attempt int) time.Duration {
	m := float64(base) * math.Pow(2, float64(attempt-1))
	if m > float64(max) {
		m = float64(max)
	}
	if jitter > 0 {
		j := time.Duration(time.Now().UnixNano()%int64(jitter+1)) - (jitter / 2)
		m += float64(j)
		if m < 0 {
			m = 0
		}
	}
	return time.Duration(m)
}

// //////////////////////////////////////////////////////////////////////////////
// Manual commit loop
// //////////////////////////////////////////////////////////////////////////////

func (h *Helper) commitLoop(ctx context.Context) error {
	ticker := time.NewTicker(h.cfg.Commit.ManualEvery)
	defer ticker.Stop()

	type tp struct {
		t string
		p int32
	}

	maxNext := make(map[tp]kgo.EpochOffset)
	acked := 0

	flush := func(flushCtx context.Context) error {
		if len(maxNext) == 0 {
			return nil
		}
		offsets := make(map[string]map[int32]kgo.EpochOffset)
		for k, eo := range maxNext {
			if offsets[k.t] == nil {
				offsets[k.t] = make(map[int32]kgo.EpochOffset)
			}
			offsets[k.t][k.p] = eo
		}

		cctx, cancel := context.WithTimeout(flushCtx, h.cfg.CommitTimeout)
		defer cancel()

		if err := h.cmt.CommitOffsets(cctx, offsets); err != nil {
			h.logPrinter(ctx, LevelError, "commit error: %v", err)
			// keep state; retry next flush
			return err
		}

		for k := range maxNext {
			delete(maxNext, k)
		}
		acked = 0
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			_ = flush(context.Background())
			return nil

		case <-ticker.C:
			_ = flush(ctx)

		case r, ok := <-h.commitAck:
			if !ok {
				_ = flush(context.Background())
				return nil
			}

			key := tp{t: r.Topic, p: r.Partition}
			next := r.Offset + 1

			cur, ok := maxNext[key]
			if !ok || next > cur.Offset {
				maxNext[key] = kgo.EpochOffset{Offset: next, Epoch: r.LeaderEpoch}
			}

			acked++
			if h.cfg.Commit.ManualBatch > 0 && acked >= h.cfg.Commit.ManualBatch {
				_ = flush(ctx)
			}
		}
	}
}

// //////////////////////////////////////////////////////////////////////////////
// franz-go integration
// //////////////////////////////////////////////////////////////////////////////

type kgoSource struct {
	cl     *kgo.Client
	closed atomic.Bool
	logger func(context.Context, LogLevel, string, ...any)

	maxPoll int

	backoffBase    time.Duration
	backoffMax     time.Duration
	currentBackoff time.Duration
}

func (s *kgoSource) Poll(ctx context.Context) ([]*kgo.Record, error) {
	fetches := s.cl.PollRecords(ctx, s.maxPoll)
	if fetches.IsClientClosed() {
		s.closed.Store(true)
		return nil, nil
	}

	var fatal error
	var hadErr bool
	var lastErr error

	fetches.EachError(func(t string, p int32, err error) {
		hadErr = true
		lastErr = err

		// Fatal authz/authn.
		if errors.Is(err, kerr.SaslAuthenticationFailed) || errors.Is(err, kerr.TopicAuthorizationFailed) {
			fatal = fmt.Errorf("fatal fetch error topic=%s partition=%d: %w", t, p, err)
			return
		}

		// Join/dial errors can have empty topic.
		if t == "" {
			s.logger(ctx, LevelWarn, "poll/join error: %v", err)
			return
		}
		s.logger(ctx, LevelWarn, "fetch error topic=%s partition=%d: %v", t, p, err)
	})

	if fatal != nil {
		return nil, fatal
	}

	// Backoff for connectivity/join-group errors.
	if hadErr && lastErr != nil {
		msg := lastErr.Error()
		if strings.Contains(msg, "no such host") ||
			strings.Contains(msg, "unable to dial") ||
			strings.Contains(msg, "unable to join group") {
			if s.currentBackoff <= 0 {
				s.currentBackoff = s.backoffBase
			} else {
				s.currentBackoff *= 2
				if s.currentBackoff > s.backoffMax {
					s.currentBackoff = s.backoffMax
				}
			}
			s.logger(ctx, LevelWarn, "kafka dial/join backoff=%s err=%v", s.currentBackoff, lastErr)
			select {
			case <-time.After(s.currentBackoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		} else {
			s.currentBackoff = 0
		}
	} else {
		s.currentBackoff = 0
	}

	var recs []*kgo.Record
	it := fetches.RecordIter()
	for !it.Done() {
		recs = append(recs, it.Next())
	}
	return recs, nil
}

func (s *kgoSource) Close() {
	if !s.closed.Load() {
		s.cl.Close()
		s.closed.Store(true)
	}
}

func (s *kgoSource) IsClosed() bool { return s.closed.Load() }

type kgoCommitter struct{ cl *kgo.Client }

func (c *kgoCommitter) CommitOffsets(ctx context.Context, offsets map[string]map[int32]kgo.EpochOffset) error {
	var commitErr error
	c.cl.CommitOffsetsSync(ctx, offsets, func(_ *kgo.Client, _ *kmsg.OffsetCommitRequest, resp *kmsg.OffsetCommitResponse, err error) {
		if err != nil {
			commitErr = err
			return
		}
		if resp == nil {
			return
		}
		for _, t := range resp.Topics {
			for _, p := range t.Partitions {
				if p.ErrorCode != 0 && commitErr == nil {
					commitErr = fmt.Errorf("commit error topic=%s partition=%d code=%d", t.Topic, p.Partition, p.ErrorCode)
				}
			}
		}
	})
	return commitErr
}

func (c *kgoCommitter) Close() {}

// //////////////////////////////////////////////////////////////////////////////
// Partition coordinator (Async ordering + poison blocking + memory cap)
// //////////////////////////////////////////////////////////////////////////////

var ErrPartitionQueueFull = errors.New("partition queue full")

type partitionCoordinator struct {
	logger        func(context.Context, LogLevel, string, ...any)
	poisonBackoff time.Duration
	maxQueue      int

	mu       sync.Mutex
	queues   map[string]map[int32][]*kgo.Record
	inflight map[string]map[int32]bool

	// round-robin key list for fairness
	rrKeys []tpKey
	rrPos  int

	stopCh chan struct{}

	// notifyCh wakes the run loop when new records are enqueued.
	notifyCh chan struct{}
}

type tpKey struct {
	topic string
	part  int32
}

func newPartitionCoordinator(logger func(context.Context, LogLevel, string, ...any), poisonBackoff time.Duration, maxQueue int) *partitionCoordinator {
	return &partitionCoordinator{
		logger:        logger,
		poisonBackoff: poisonBackoff,
		maxQueue:      maxQueue,
		queues:        make(map[string]map[int32][]*kgo.Record),
		inflight:      make(map[string]map[int32]bool),
		stopCh:        make(chan struct{}),
		notifyCh:      make(chan struct{}, 1),
	}
}

func (pc *partitionCoordinator) stop() {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	select {
	case <-pc.stopCh:
	default:
		close(pc.stopCh)
	}
}

func (pc *partitionCoordinator) enqueue(r *kgo.Record) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.queues[r.Topic] == nil {
		pc.queues[r.Topic] = make(map[int32][]*kgo.Record)
	}
	q := pc.queues[r.Topic][r.Partition]

	if pc.maxQueue > 0 && len(q) >= pc.maxQueue {
		return ErrPartitionQueueFull
	}
	pc.queues[r.Topic][r.Partition] = append(q, r)

	// Maintain rrKeys lazily: add key if not present.
	// O(n) scan is acceptable here because keys count = partitions count (bounded).
	k := tpKey{topic: r.Topic, part: r.Partition}
	found := false
	for _, kk := range pc.rrKeys {
		if kk == k {
			found = true
			break
		}
	}
	if !found {
		pc.rrKeys = append(pc.rrKeys, k)
	}

	// Wake dispatcher when new work arrives.
	select {
	case pc.notifyCh <- struct{}{}:
	default:
	}

	return nil
}

func (pc *partitionCoordinator) run(ctx context.Context, jobs chan<- *kgo.Record, done <-chan workDone) {
	dispatch := func() {
		pc.mu.Lock()
		defer pc.mu.Unlock()

		if len(pc.rrKeys) == 0 {
			return
		}

		// bounded attempts: at most len(rrKeys) selections per dispatch cycle
		for i := 0; i < len(pc.rrKeys); i++ {
			k := pc.rrKeys[pc.rrPos%len(pc.rrKeys)]
			pc.rrPos++

			pm := pc.queues[k.topic]
			if pm == nil {
				continue
			}
			q := pm[k.part]
			if len(q) == 0 {
				continue
			}

			if pc.inflight[k.topic] == nil {
				pc.inflight[k.topic] = make(map[int32]bool)
			}
			if pc.inflight[k.topic][k.part] {
				continue
			}

			r := q[0]
			pc.queues[k.topic][k.part] = q[1:]
			pc.inflight[k.topic][k.part] = true

			select {
			case jobs <- r:
				// scheduled
			default:
				// workers saturated: push back and abort
				pc.queues[k.topic][k.part] = append([]*kgo.Record{r}, pc.queues[k.topic][k.part]...)
				pc.inflight[k.topic][k.part] = false
				return
			}
		}
	}

	dispatch()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pc.stopCh:
			return
		case <-pc.notifyCh:
			dispatch()

		case wd, ok := <-done:
			if !ok {
				return
			}

			topic := wd.r.Topic
			part := wd.r.Partition

			if wd.err != nil {
				pc.logger(
					ctx,
					LevelError,
					"processing failed topic=%s partition=%d offset=%d: %v (blocking partition and retrying)",
					topic, part, wd.r.Offset, wd.err,
				)

				pc.mu.Lock()
				if pc.queues[topic] == nil {
					pc.queues[topic] = make(map[int32][]*kgo.Record)
				}
				// requeue front
				pc.queues[topic][part] = append([]*kgo.Record{wd.r}, pc.queues[topic][part]...)
				if pc.inflight[topic] != nil {
					pc.inflight[topic][part] = false
				}
				pc.mu.Unlock()

				select {
				case <-time.After(pc.poisonBackoff):
				case <-ctx.Done():
					return
				case <-pc.stopCh:
					return
				}
			} else {
				pc.mu.Lock()
				if pc.inflight[topic] != nil {
					pc.inflight[topic][part] = false
				}
				pc.mu.Unlock()
			}

			dispatch()
		}
	}
}

func (h *Helper) logPrinter(ctx context.Context, logLevel LogLevel, format string, args ...any) {
	if h == nil {
		printf(ctx, format, args...)
		return
	}
	h.cfg.log(ctx, logLevel, format, args...)
}

// printf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func printf(ctx context.Context, format string, v ...any) {
	log.Printf(format, v...)
}
