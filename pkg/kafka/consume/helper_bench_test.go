package consume

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// -------------------- bench helpers --------------------

type benchSource struct {
	closed atomic.Bool
	ch     chan *kgo.Record
}

func newBenchSource(buf int) *benchSource {
	return &benchSource{ch: make(chan *kgo.Record, buf)}
}

func (s *benchSource) Poll(ctx context.Context) ([]*kgo.Record, error) {
	if s.closed.Load() {
		return nil, nil
	}
	select {
	case <-ctx.Done():
		return nil, nil
	case r := <-s.ch:
		if r == nil {
			// sentinel
			return nil, nil
		}
		return []*kgo.Record{r}, nil
	}
}

func (s *benchSource) Close()         { s.closed.Store(true) }
func (s *benchSource) IsClosed() bool { return s.closed.Load() }

type noopCommitter struct{}

func (c *noopCommitter) CommitOffsets(ctx context.Context, offsets map[string]map[int32]kgo.EpochOffset) error {
	return nil
}
func (c *noopCommitter) Close() {}

// type procFunc func(ctx context.Context, r *kgo.Record) error
//
// func (f procFunc) Process(ctx context.Context, r *kgo.Record) error { return f(ctx, r) }

// -------------------- benchmarks --------------------

func BenchmarkSafeProcess_Success(b *testing.B) {
	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",

		Mode: ModeSync,
		Commit: CommitPolicy{
			Mode: CommitManual,
		},
		ProcessTimeout: 5 * time.Second,
		Retry: RetryPolicy{
			MaxAttempts:   1,
			BaseBackoff:   1 * time.Millisecond,
			MaxBackoff:    1 * time.Millisecond,
			Jitter:        0,
			PoisonBackoff: 1 * time.Millisecond,
		},
		JoinFailFast: false,
	}

	h, err := NewWith(cfg, procFunc(func(ctx context.Context, r *kgo.Record) error { return nil }), nil, &fakeSource{}, &fakeCommitter{})
	if err != nil {
		b.Fatalf("NewWith: %v", err)
	}

	r := &kgo.Record{Topic: "t", Partition: 0, Offset: 1}

	ctx := context.Background()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = h.safeProcess(ctx, r)
	}
}

func BenchmarkSafeProcess_FailThenSuccess(b *testing.B) {
	var calls uint64

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",

		Mode: ModeSync,
		Commit: CommitPolicy{
			Mode: CommitManual,
		},
		ProcessTimeout: 5 * time.Second,
		Retry: RetryPolicy{
			MaxAttempts:   3,
			BaseBackoff:   10 * time.Microsecond,
			MaxBackoff:    10 * time.Microsecond,
			Jitter:        0,
			PoisonBackoff: 1 * time.Millisecond,
		},
		JoinFailFast: false,
	}

	p := procFunc(func(ctx context.Context, r *kgo.Record) error {
		n := atomic.AddUint64(&calls, 1)
		if n%2 == 1 {
			return errors.New("transient")
		}
		return nil
	})

	h, err := NewWith(cfg, p, nil, &fakeSource{}, &fakeCommitter{})
	if err != nil {
		b.Fatalf("NewWith: %v", err)
	}

	r := &kgo.Record{Topic: "t", Partition: 0, Offset: 1}
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = h.safeProcess(ctx, r)
	}
}

func BenchmarkPartitionCoordinator_Dispatch(b *testing.B) {
	logger := func(context.Context, LogLevel, string, ...any) {}
	pc := newPartitionCoordinator(logger, 1*time.Millisecond, 50000)

	// Use buffered channels to avoid blocking effects dominating the bench.
	jobs := make(chan *kgo.Record, 4096)
	done := make(chan workDone, 4096)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		pc.run(ctx, jobs, done)
	}()

	// prefill N records per partition
	partitions := int32(8)
	perPart := 2000
	for p := int32(0); p < partitions; p++ {
		for i := 0; i < perPart; i++ {
			_ = pc.enqueue(&kgo.Record{Topic: "t", Partition: p, Offset: int64(i)})
		}
	}

	// Drain jobs and immediately mark done success, emulating fast workers.
	// We measure coordinator scheduling + bookkeeping overhead.
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		select {
		case r := <-jobs:
			done <- workDone{r: r, err: nil}
		default:
			// if empty, nudge by letting run loop continue
			time.Sleep(0)
		}
	}

	cancel()
	pc.stop()
	wg.Wait()
}

func BenchmarkAsyncPipeline_TwoPartitions(b *testing.B) {
	// End-to-end in-process pipeline (no Kafka):
	// benchSource feeds records; helper runs Async; handler is no-op; manual commit enabled.
	src := newBenchSource(1 << 16)
	cmt := &noopCommitter{}

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",

		Mode:      ModeAsync,
		Workers:   8,
		QueueSize: 1 << 16,

		Commit: CommitPolicy{
			Mode:        CommitManual,
			ManualEvery: 200 * time.Millisecond,
			ManualBatch: 2000,
		},
		CommitTimeout: 5 * time.Second,

		ProcessTimeout: 5 * time.Second,
		Retry: RetryPolicy{
			MaxAttempts:   1,
			BaseBackoff:   1 * time.Millisecond,
			MaxBackoff:    1 * time.Millisecond,
			Jitter:        0,
			PoisonBackoff: 1 * time.Millisecond,
		},
		JoinFailFast:        false,
		MaxPartitionQueue:   50000,
		PollErrorBackoff:    1 * time.Millisecond,
		PollErrorMaxBackoff: 1 * time.Millisecond,
	}

	p := procFunc(func(ctx context.Context, r *kgo.Record) error { return nil })

	h, err := NewWith(cfg, p, nil, src, cmt)
	if err != nil {
		b.Fatalf("NewWith: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run helper in background for the duration of the benchmark.
	doneRun := make(chan struct{})
	go func() {
		_ = h.Run(ctx)
		close(doneRun)
	}()

	// Feed records on two partitions to mimic real stream.
	// We measure throughput of whole pipeline (poll->coord->worker->commitAck bookkeeping).
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		part := int32(i & 1)
		src.ch <- &kgo.Record{Topic: "t", Partition: part, Offset: int64(i)}
	}

	// Stop and wait.
	cancel()
	<-doneRun
}
