package consume

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// TopicResolver abstracts topic listing to allow unit tests without Kafka.
type TopicResolver interface {
	ListTopics(ctx context.Context) ([]string, error)
	Close()
}

// ResolveTopicsByPrefix returns all topics whose name starts with any prefix.
// - Results are unique and sorted.
// - Empty prefixes are ignored.
// - If no topics match, returns an error.
func ResolveTopicsByPrefix(ctx context.Context, r TopicResolver, prefixes []string) ([]string, error) {
	if r == nil {
		return nil, fmt.Errorf("TopicResolver required")
	}

	var ps []string
	for _, p := range prefixes {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		ps = append(ps, p)
	}
	if len(ps) == 0 {
		return nil, fmt.Errorf("at least one non-empty prefix is required")
	}

	topics, err := r.ListTopics(ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{}, 64)
	var out []string
	for _, t := range topics {
		for _, p := range ps {
			if strings.HasPrefix(t, p) {
				if _, ok := seen[t]; !ok {
					seen[t] = struct{}{}
					out = append(out, t)
				}
				break
			}
		}
	}

	sort.Strings(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("no topics match prefixes=%v", ps)
	}
	return out, nil
}

// NewTopicResolverFromConfig builds a Kafka-backed TopicResolver using a subset of Config.
// This is intentionally small so it remains safe and reusable.
func NewTopicResolverFromConfig(cfg Config) (TopicResolver, error) {
	// Reuse validate defaults, but allow Topics empty here because we're discovering them.
	cfg.setDefaults()

	if len(cfg.Brokers) == 0 {
		return nil, fmt.Errorf("brokers required")
	}
	if cfg.ClientID == "" {
		cfg.ClientID = "topic-resolver"
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ClientID(cfg.ClientID),
	}

	// TLS
	if cfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(cfg.TLSConfig))
	}

	// SASL (SCRAM)
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

	return &kgoTopicResolver{
		cl:     cl,
		logger: cfg.Logger,
	}, nil
}

// ConsumeByPrefix resolves topics by prefix and then creates a Helper with those topics.
// Note: This does NOT auto-refresh for newly created topics; Kafka consumers do not support
// regex/prefix subscribe without rebuilding the client.
func ConsumeByPrefix(ctx context.Context, cfg Config, prefixes []string, p Processor, dlq DLQHandler) (*Helper, []string, error) {
	res, err := NewTopicResolverFromConfig(cfg)
	if err != nil {
		return nil, nil, err
	}
	defer res.Close()

	topics, err := ResolveTopicsByPrefix(ctx, res, prefixes)
	if err != nil {
		return nil, nil, err
	}

	cfg.Topics = topics
	h, err := New(cfg, p, dlq)
	if err != nil {
		return nil, nil, err
	}
	return h, topics, nil
}

// ConsumeByPrefixWithRefresh is an optional pattern for environments where topics are created dynamically.
// It periodically checks for new matching topics, and if the set changes, returns (nil, newTopics, ErrTopicsChanged)
// so the caller can restart the consumer (rolling restart, supervisor, etc).
//
// This helper does NOT restart automatically (to keep core behavior explicit and stable).
var ErrTopicsChanged = fmt.Errorf("topics changed; restart required")

func CheckPrefixTopicsChanged(ctx context.Context, cfg Config, prefixes []string, prev []string) ([]string, error) {
	res, err := NewTopicResolverFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	cur, err := ResolveTopicsByPrefix(ctx, res, prefixes)
	if err != nil {
		return nil, err
	}

	if !stringSlicesEqual(cur, prev) {
		return cur, ErrTopicsChanged
	}
	return cur, nil
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// -------- franz-go implementation --------

type kgoTopicResolver struct {
	cl     *kgo.Client
	logger func(context.Context, string, ...any)
}

func (r *kgoTopicResolver) Close() { r.cl.Close() }

func (r *kgoTopicResolver) ListTopics(ctx context.Context) ([]string, error) {
	req := kmsg.NewPtrMetadataRequest()
	req.AllowAutoTopicCreation = false

	respAny, err := r.cl.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	resp, ok := respAny.(*kmsg.MetadataResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", respAny)
	}

	out := make([]string, 0, len(resp.Topics))
	for _, t := range resp.Topics {
		if t.ErrorCode != 0 {
			continue // ACL / unknown topic
		}
		if t.Topic == nil || *t.Topic == "" {
			continue
		}
		out = append(out, *t.Topic)
	}

	return out, nil
}

// scramMechanism is a small internal helper to avoid duplicating logic from helper.go.
func scramMechanism(mech, user, pass string) (sasl.Mechanism, error) {
	switch strings.ToUpper(strings.TrimSpace(mech)) {
	case "SCRAM-SHA-256":
		return scram.Auth{User: user, Pass: pass}.AsSha256Mechanism(), nil
	case "SCRAM-SHA-512":
		return scram.Auth{User: user, Pass: pass}.AsSha512Mechanism(), nil
	default:
		return nil, fmt.Errorf("unsupported SASL mechanism: %s", mech)
	}
}

// Optional convenience: caller can use a short timeout when resolving topics.
func ResolveWithTimeout(cfg Config, prefixes []string, timeout time.Duration) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r, err := NewTopicResolverFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ResolveTopicsByPrefix(ctx, r, prefixes)
}
