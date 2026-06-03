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

// -------------------- fakes --------------------

type fakeSource struct {
	mu     sync.Mutex
	recs   []*kgo.Record
	closed atomic.Bool
}

func (s *fakeSource) Poll(ctx context.Context) ([]*kgo.Record, error) {
	if s.closed.Load() {
		return nil, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.recs) == 0 {
		// behave like idle poll
		select {
		case <-time.After(10 * time.Millisecond):
		case <-ctx.Done():
		}
		return nil, nil
	}
	// return at most a small batch
	n := 10
	if len(s.recs) < n {
		n = len(s.recs)
	}
	out := s.recs[:n]
	s.recs = s.recs[n:]
	return out, nil
}

func (s *fakeSource) Close()         { s.closed.Store(true) }
func (s *fakeSource) IsClosed() bool { return s.closed.Load() }

type fakeCommitter struct {
	mu      sync.Mutex
	commits []map[string]map[int32]kgo.EpochOffset
	err     error
}

func (c *fakeCommitter) CommitOffsets(ctx context.Context, offsets map[string]map[int32]kgo.EpochOffset) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	// store a deep-ish copy to avoid mutation surprises
	cp := make(map[string]map[int32]kgo.EpochOffset)
	for t, pm := range offsets {
		cp[t] = make(map[int32]kgo.EpochOffset)
		for p, eo := range pm {
			cp[t][p] = eo
		}
	}
	c.commits = append(c.commits, cp)
	return c.err
}
func (c *fakeCommitter) Close() {}

type procFunc func(ctx context.Context, r *kgo.Record) error

func (f procFunc) Process(ctx context.Context, r *kgo.Record) error { return f(ctx, r) }

type dlqFunc func(ctx context.Context, r *kgo.Record, cause error) error

func (f dlqFunc) HandleDLQ(ctx context.Context, r *kgo.Record, cause error) error {
	return f(ctx, r, cause)
}

// -------------------- tests --------------------

func TestSyncManualCommit_CommitsOnlySuccess(t *testing.T) {
	src := &fakeSource{
		recs: []*kgo.Record{
			{Topic: "t", Partition: 0, Offset: 0},
			{Topic: "t", Partition: 0, Offset: 1},
			{Topic: "t", Partition: 0, Offset: 2},
		},
	}
	cmt := &fakeCommitter{}

	var calls int
	var seen []int64
	p := procFunc(func(ctx context.Context, r *kgo.Record) error {
		calls++
		seen = append(seen, r.Offset)
		if r.Offset == 1 {
			return errors.New("fail")
		}
		return nil
	})

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeSync,
		Commit: CommitPolicy{
			Mode:        CommitManual,
			ManualEvery: 50 * time.Millisecond,
			ManualBatch: 1,
			OnFailure:   FailurePolicyStop,
		},
		Retry: RetryPolicy{
			MaxAttempts:   1,
			BaseBackoff:   1 * time.Millisecond,
			MaxBackoff:    1 * time.Millisecond,
			Jitter:        0,
			PoisonBackoff: 1 * time.Millisecond,
		},
		JoinFailFast: false,
	}

	h, err := NewWith(cfg, p, nil, src, cmt)
	if err != nil {
		t.Fatalf("NewWith: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = h.Run(ctx)
	if err == nil {
		t.Fatalf("expected sync mode to stop on processing failure")
	}

	// Sync mode is safety-first now: stop immediately after unrecoverable failure.
	// Expected calls: offset 0 (success), offset 1 (fail), offset 2 is not processed.
	if calls != 2 || len(seen) != 2 || seen[0] != 0 || seen[1] != 1 {
		t.Fatalf("expected processed offsets [0 1], got calls=%d seen=%v", calls, seen)
	}
	// Only the success before failure should be committed.
	cmt.mu.Lock()
	defer cmt.mu.Unlock()
	if len(cmt.commits) == 0 {
		t.Fatalf("expected commits, got none")
	}
	last := cmt.commits[len(cmt.commits)-1]
	off, ok := last["t"][0]
	if !ok {
		t.Fatalf("expected commit for topic t partition 0")
	}
	if off.Offset != 1 {
		t.Fatalf("expected committed next offset=1, got %d", off.Offset)
	}
}

func TestSyncManualCommit_OnFailureCommitOnError(t *testing.T) {
	src := &fakeSource{
		recs: []*kgo.Record{
			{Topic: "t", Partition: 0, Offset: 0},
			{Topic: "t", Partition: 0, Offset: 1},
		},
	}
	cmt := &fakeCommitter{}

	var seen []int64
	p := procFunc(func(ctx context.Context, r *kgo.Record) error {
		seen = append(seen, r.Offset)
		if r.Offset == 0 {
			return errors.New("failed but should commit by policy")
		}
		return nil
	})

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeSync,
		Commit: CommitPolicy{
			Mode:        CommitManual,
			ManualEvery: 50 * time.Millisecond,
			ManualBatch: 1,
			OnFailure:   FailurePolicyCommitOnError,
		},
		Retry: RetryPolicy{
			MaxAttempts: 1,
		},
		JoinFailFast: false,
	}

	h, err := NewWith(cfg, p, nil, src, cmt)
	if err != nil {
		t.Fatalf("NewWith: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := h.Run(ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(seen) < 2 || seen[0] != 0 || seen[1] != 1 {
		t.Fatalf("expected processed offsets [0 1], got %v", seen)
	}

	cmt.mu.Lock()
	defer cmt.mu.Unlock()
	if len(cmt.commits) == 0 {
		t.Fatalf("expected commits")
	}
	last := cmt.commits[len(cmt.commits)-1]
	off := last["t"][0]
	if off.Offset != 2 {
		t.Fatalf("expected committed next offset=2, got %d", off.Offset)
	}
}

func TestSyncManualCommit_OnFailureDLQThenCommit(t *testing.T) {
	src := &fakeSource{
		recs: []*kgo.Record{
			{Topic: "t", Partition: 0, Offset: 0},
		},
	}
	cmt := &fakeCommitter{}

	var dlqCalls int
	p := procFunc(func(ctx context.Context, r *kgo.Record) error {
		return errors.New("send to dlq")
	})
	dlq := dlqFunc(func(ctx context.Context, r *kgo.Record, cause error) error {
		dlqCalls++
		return nil
	})

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeSync,
		Commit: CommitPolicy{
			Mode:        CommitManual,
			ManualEvery: 50 * time.Millisecond,
			ManualBatch: 1,
			OnFailure:   FailurePolicyDLQThenCommit,
		},
		Retry: RetryPolicy{
			MaxAttempts: 1,
		},
		JoinFailFast: false,
	}

	h, err := NewWith(cfg, p, dlq, src, cmt)
	if err != nil {
		t.Fatalf("NewWith: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := h.Run(ctx); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if dlqCalls != 1 {
		t.Fatalf("expected dlq to be called once, got %d", dlqCalls)
	}

	cmt.mu.Lock()
	defer cmt.mu.Unlock()
	if len(cmt.commits) == 0 {
		t.Fatalf("expected commits")
	}
	last := cmt.commits[len(cmt.commits)-1]
	off := last["t"][0]
	if off.Offset != 1 {
		t.Fatalf("expected committed next offset=1, got %d", off.Offset)
	}
}

func TestConfig_OnFailureDLQThenCommitRequiresDLQ(t *testing.T) {
	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeSync,
		Commit: CommitPolicy{
			Mode:      CommitManual,
			OnFailure: FailurePolicyDLQThenCommit,
		},
		JoinFailFast: false,
	}

	_, err := NewWith(cfg, procFunc(func(ctx context.Context, r *kgo.Record) error { return nil }), nil, &fakeSource{}, &fakeCommitter{})
	if err == nil {
		t.Fatalf("expected error when dlq is required but nil")
	}
}

func TestAsyncManualCommit_PerPartitionOrdering(t *testing.T) {
	// Interleaved records across two partitions; handler records processing order per partition.
	src := &fakeSource{
		recs: []*kgo.Record{
			{Topic: "t", Partition: 0, Offset: 0},
			{Topic: "t", Partition: 1, Offset: 0},
			{Topic: "t", Partition: 0, Offset: 1},
			{Topic: "t", Partition: 1, Offset: 1},
			{Topic: "t", Partition: 0, Offset: 2},
			{Topic: "t", Partition: 1, Offset: 2},
		},
	}
	cmt := &fakeCommitter{}

	var mu sync.Mutex
	seen := map[int32][]int64{}
	var cancel context.CancelFunc
	var cancelOnce sync.Once

	p := procFunc(func(ctx context.Context, r *kgo.Record) error {
		mu.Lock()
		seen[r.Partition] = append(seen[r.Partition], r.Offset)
		done := len(seen[int32(0)]) == 3 && len(seen[int32(1)]) == 3
		mu.Unlock()
		if done && cancel != nil {
			cancelOnce.Do(cancel)
		}
		time.Sleep(5 * time.Millisecond)
		return nil
	})

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeAsync,
		Workers: 4,
		Commit: CommitPolicy{
			Mode:        CommitManual,
			ManualEvery: 50 * time.Millisecond,
			ManualBatch: 1,
		},
		JoinFailFast: false,
	}

	h, err := NewWith(cfg, p, nil, src, cmt)
	if err != nil {
		t.Fatalf("NewWith: %v", err)
	}

	ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	cancel = stop
	defer stop()
	_ = h.Run(ctx)

	mu.Lock()
	defer mu.Unlock()

	if got := seen[int32(0)]; len(got) != 3 || got[0] != 0 || got[1] != 1 || got[2] != 2 {
		t.Fatalf("partition 0 order wrong: %v", got)
	}
	if got := seen[int32(1)]; len(got) != 3 || got[0] != 0 || got[1] != 1 || got[2] != 2 {
		t.Fatalf("partition 1 order wrong: %v", got)
	}
}

func TestAsyncPoisonBlocksPartition(t *testing.T) {
	// Partition 0 first record always fails; partition 1 should still process.
	src := &fakeSource{
		recs: []*kgo.Record{
			{Topic: "t", Partition: 0, Offset: 0},
			{Topic: "t", Partition: 1, Offset: 0},
			{Topic: "t", Partition: 1, Offset: 1},
			{Topic: "t", Partition: 0, Offset: 1}, // should not run because partition 0 is blocked by offset 0 poison
		},
	}
	cmt := &fakeCommitter{}

	var mu sync.Mutex
	seen := map[int32][]int64{}
	var cancel context.CancelFunc
	var cancelOnce sync.Once
	p := procFunc(func(ctx context.Context, r *kgo.Record) error {
		if r.Partition == 0 && r.Offset == 0 {
			return errors.New("poison")
		}
		mu.Lock()
		seen[r.Partition] = append(seen[r.Partition], r.Offset)
		done := len(seen[int32(1)]) == 2
		mu.Unlock()
		if done && cancel != nil {
			cancelOnce.Do(cancel)
		}
		return nil
	})

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeAsync,
		Workers: 2,
		Commit: CommitPolicy{
			Mode:        CommitManual,
			ManualEvery: 50 * time.Millisecond,
			ManualBatch: 1,
			OnFailure:   FailurePolicyRetryForever,
		},
		Retry: RetryPolicy{
			MaxAttempts:   1, // fail immediately -> poison
			PoisonBackoff: 50 * time.Millisecond,
		},
		JoinFailFast: false,
	}

	h, err := NewWith(cfg, p, nil, src, cmt)
	if err != nil {
		t.Fatalf("NewWith: %v", err)
	}

	ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	cancel = stop
	defer stop()
	_ = h.Run(ctx)

	mu.Lock()
	defer mu.Unlock()

	// partition 1 should proceed
	if got := seen[int32(1)]; len(got) != 2 || got[0] != 0 || got[1] != 1 {
		t.Fatalf("partition 1 expected [0 1], got %v", got)
	}
	// partition 0 might attempt offset 0 repeatedly; we should not see offset 1.
	if got := seen[int32(0)]; len(got) != 0 {
		t.Fatalf("partition 0 should not record successes, got %v", got)
	}
}

func TestJoinFailFast_ReturnsError(t *testing.T) {
	// Source never yields records nor assignment; fail-fast should close source and Run returns error override.
	src := &fakeSource{}
	cmt := &fakeCommitter{}

	p := procFunc(func(ctx context.Context, r *kgo.Record) error { return nil })

	cfg := Config{
		Brokers: []string{"x"},
		Topics:  []string{"t"},
		Group:   "g",
		Mode:    ModeSync,
		Commit:  CommitPolicy{Mode: CommitManual, ManualEvery: 50 * time.Millisecond, ManualBatch: 1},

		JoinFailFast:        true,
		JoinFailFastTimeout: 80 * time.Millisecond,
	}

	h, err := NewWith(cfg, p, nil, src, cmt)
	if err != nil {
		t.Fatalf("NewWith: %v", err)
	}

	// Ensure no assignment is recorded.
	h.assignedOnce = &atomic.Bool{} // false

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	err = h.Run(ctx)
	if err == nil {
		t.Fatalf("expected join fail-fast error, got nil")
	}
	if !stringsContains(err.Error(), "join fail-fast") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func stringsContains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (func() bool {
		// tiny helper to avoid importing strings in test
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})())
}
