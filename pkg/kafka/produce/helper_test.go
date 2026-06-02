package produce

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

type fakeClient struct {
	mu sync.Mutex

	failN   int
	errFail error

	calls      int
	flushCalls int
	closed     atomic.Bool
}

func (f *fakeClient) Produce(ctx context.Context, r *kgo.Record, promise func(*kgo.Record, error)) {
	f.mu.Lock()
	f.calls++
	call := f.calls
	failN := f.failN
	errFail := f.errFail
	f.mu.Unlock()

	go func() {
		select {
		case <-ctx.Done():
			promise(r, ctx.Err())
			return
		default:
		}

		if call <= failN {
			promise(r, errFail)
			return
		}

		rr := *r
		rr.Partition = 3
		rr.Offset = 777
		rr.Timestamp = time.Unix(123, 0)
		promise(&rr, nil)
	}()
}

func (f *fakeClient) Flush(ctx context.Context) error {
	f.mu.Lock()
	f.flushCalls++
	f.mu.Unlock()
	return nil
}

func (f *fakeClient) Close() { f.closed.Store(true) }

func TestProduce_ValidateRecord(t *testing.T) {
	cl := &fakeClient{}
	h := NewWithClient(Config{Brokers: []string{"x"}, Mode: ModeSync, Logger: func(string, ...any) {}}, cl)
	defer h.Close()

	if _, err := h.Produce(context.Background(), nil); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := h.Produce(context.Background(), &kgo.Record{}); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSyncProduce_Success(t *testing.T) {
	cl := &fakeClient{failN: 0}
	cfg := Config{
		Brokers:          []string{"x"},
		Mode:             ModeSync,
		ProduceTimeout:   500 * time.Millisecond,
		RetryMaxAttempts: 1,
		Logger:           func(string, ...any) {},
	}

	h := NewWithClient(cfg, cl)
	defer h.Close()

	res, err := h.Produce(context.Background(), &kgo.Record{Topic: "t", Value: []byte("v")})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Topic != "t" || res.Partition != 3 || res.Offset != 777 {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestSyncProduce_RetryThenSuccess(t *testing.T) {
	// NotLeaderForPartition must be treated retriable by our allowlist.
	cl := &fakeClient{failN: 2, errFail: kerr.NotLeaderForPartition}

	cfg := Config{
		Brokers:          []string{"x"},
		Mode:             ModeSync,
		ProduceTimeout:   500 * time.Millisecond,
		RetryMaxAttempts: 5,
		RetryBaseBackoff: 1 * time.Millisecond,
		RetryMaxBackoff:  2 * time.Millisecond,
		RetryJitter:      0,
		Logger:           func(string, ...any) {},
	}

	h := NewWithClient(cfg, cl)
	defer h.Close()

	if _, err := h.Produce(context.Background(), &kgo.Record{Topic: "t"}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	cl.mu.Lock()
	calls := cl.calls
	cl.mu.Unlock()
	if calls < 3 {
		t.Fatalf("expected retries, calls=%d", calls)
	}
}

func TestSyncProduce_StopOnNonRetriable(t *testing.T) {
	cl := &fakeClient{failN: 10, errFail: kerr.TopicAuthorizationFailed}

	cfg := Config{
		Brokers:          []string{"x"},
		Mode:             ModeSync,
		ProduceTimeout:   200 * time.Millisecond,
		RetryMaxAttempts: 3,
		RetryBaseBackoff: 1 * time.Millisecond,
		RetryMaxBackoff:  2 * time.Millisecond,
		RetryJitter:      0,
		Logger:           func(string, ...any) {},
	}

	h := NewWithClient(cfg, cl)
	defer h.Close()

	if _, err := h.Produce(context.Background(), &kgo.Record{Topic: "t"}); err == nil {
		t.Fatalf("expected error")
	}

	cl.mu.Lock()
	calls := cl.calls
	cl.mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected 1 call (no retry), calls=%d", calls)
	}
}

func TestAsyncProduce_AwaitAck(t *testing.T) {
	cl := &fakeClient{failN: 0}
	cfg := Config{
		Brokers:          []string{"x"},
		Mode:             ModeAsync,
		Workers:          2,
		QueueSize:        32,
		ProduceTimeout:   500 * time.Millisecond,
		RetryMaxAttempts: 1,
		Logger:           func(string, ...any) {},
	}

	h := NewWithClient(cfg, cl)
	defer h.Close()

	res, err := h.Produce(context.Background(), &kgo.Record{Topic: "t"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if res.Offset != 777 {
		t.Fatalf("unexpected res: %+v", res)
	}
}

func TestAsyncProduce_NoWait_Backpressure(t *testing.T) {
	cl := &fakeClient{failN: 0}
	cfg := Config{
		Brokers:          []string{"x"},
		Mode:             ModeAsync,
		Workers:          1,
		QueueSize:        1,
		ProduceTimeout:   500 * time.Millisecond,
		RetryMaxAttempts: 1,
		Logger:           func(string, ...any) {},
	}

	h := NewWithClient(cfg, cl)
	defer h.Close()

	ctx := context.Background()
	if err := h.ProduceAsyncNoWait(ctx, &kgo.Record{Topic: "t"}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := h.ProduceAsyncNoWait(ctx, &kgo.Record{Topic: "t"}); err == nil {
		t.Fatalf("expected queue full error")
	}
}

func TestClose_FlushAndCloseClient(t *testing.T) {
	cl := &fakeClient{}
	cfg := Config{
		Brokers: []string{"x"},
		Mode:    ModeAsync,
		Workers: 1,
		Logger:  func(string, ...any) {},
	}

	h := NewWithClient(cfg, cl)
	h.Close()

	cl.mu.Lock()
	flushCalls := cl.flushCalls
	cl.mu.Unlock()

	if flushCalls == 0 {
		t.Fatalf("expected Flush called")
	}
	if !cl.closed.Load() {
		t.Fatalf("expected client Close called")
	}
}

func TestIsRetriable_Kerr(t *testing.T) {
	if !isRetriableProduceError(kerr.NotLeaderForPartition) {
		t.Fatalf("expected retriable")
	}
	if isRetriableProduceError(errors.New("permanent failure")) {
		t.Fatalf("expected non-retriable")
	}
}
