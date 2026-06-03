package consume

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

type fakeResolver struct {
	topics []string
	err    error
}

func (f *fakeResolver) ListTopics(ctx context.Context) ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return append([]string(nil), f.topics...), nil
}
func (f *fakeResolver) Close() {}

func TestResolveTopicsByPrefix_Basic(t *testing.T) {
	r := &fakeResolver{
		topics: []string{
			"general-push-notification",
			"general-email-notification",
			"user-push-notification",
			"__consumer_offsets",
			"general-push-notification-retry",
		},
	}

	got, err := ResolveTopicsByPrefix(context.Background(), r, []string{"general-", "user-"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	want := []string{
		"general-email-notification",
		"general-push-notification",
		"general-push-notification-retry",
		"user-push-notification",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestResolveTopicsByPrefix_DedupAndSort(t *testing.T) {
	r := &fakeResolver{
		topics: []string{
			"ab-1",
			"ab-2",
			"ab-1", // duplicate
			"zz",
		},
	}

	got, err := ResolveTopicsByPrefix(context.Background(), r, []string{"ab-", "ab-"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	want := []string{"ab-1", "ab-2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestResolveTopicsByPrefix_EmptyPrefixes(t *testing.T) {
	r := &fakeResolver{topics: []string{"a"}}
	_, err := ResolveTopicsByPrefix(context.Background(), r, []string{"", "   "})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveTopicsByPrefix_NoMatches(t *testing.T) {
	r := &fakeResolver{topics: []string{"a", "b"}}
	_, err := ResolveTopicsByPrefix(context.Background(), r, []string{"x-"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveTopicsByPrefix_ResolverError(t *testing.T) {
	r := &fakeResolver{err: errors.New("boom")}
	_, err := ResolveTopicsByPrefix(context.Background(), r, []string{"a"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCheckPrefixTopicsChanged(t *testing.T) {
	// We do not use the Kafka-backed resolver here; instead we validate the comparison behavior.
	prev := []string{"a-1", "a-2"}

	// Use ResolveTopicsByPrefix directly with fake resolver.
	r := &fakeResolver{topics: []string{"a-1", "a-2", "a-3"}}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cur, err := ResolveTopicsByPrefix(ctx, r, []string{"a-"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if reflect.DeepEqual(cur, prev) {
		t.Fatalf("expected changed topics")
	}
}
