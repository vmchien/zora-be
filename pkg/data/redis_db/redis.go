package redis_db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Conf        map[string]any
	Name        string
	IsConnected bool
	c           *redis.Client
}

func NewRedisClient() *RedisClient {
	return &RedisClient{
		Name:        fmt.Sprintf("redis-%d-%d", os.Getpid(), time.Now().UnixNano()),
		IsConnected: false,
	}
}

func (r *RedisClient) Connect(config map[string]any, keepAliveFunc func(err error)) error {
	opts := &redis.Options{
		Addr:       config["addr"].(string),
		ClientName: r.Name,
	}
	if v, ok := config["username"].(string); ok && len(v) > 0 {
		opts.Username = v
	}
	if v, ok := config["password"].(string); ok && len(v) > 0 {
		opts.Password = v
	}
	if v, ok := config["db"].(string); ok && len(v) > 0 {
		opts.DB, _ = strconv.Atoi(v)
	}

	// for GCP, must clear client name and disable identity
	// https://github.com/redis/go-redis/discussions/2537
	opts.DisableIdentity = true
	opts.ClientName = ""

	r.c = redis.NewClient(opts)
	r.Conf = config

	fmt.Printf("redis addr=%s db=%d\n", opts.Addr, opts.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	state, err := r.c.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Printf("%s ping result: %v \n", r.Name, state)
	r.IsConnected = true

	if keepAliveFunc != nil {
		go r.keepAliveFunc(keepAliveFunc)
	}

	return nil
}

func (r *RedisClient) Fetch(ctx context.Context, k string) ([]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}
	rs, err := r.c.Get(ctx, k).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	return rs, nil
}
func (r *RedisClient) Get(ctx context.Context, k string) ([]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}
	rs, err := r.c.Get(ctx, k).Bytes()
	if err != nil {
		return nil, err
	}
	return rs, nil
}
func (r *RedisClient) Set(ctx context.Context, k string, v []byte, expiration time.Duration) error {
	if !r.valid() {
		return fmt.Errorf("redis client not connected")
	}
	return r.c.Set(ctx, k, v, expiration).Err()
}

func (r *RedisClient) Del(ctx context.Context, k string) error {
	if !r.valid() {
		return fmt.Errorf("redis client not connected")
	}
	return r.c.Del(ctx, k).Err()
}

// // SearchKeysByPrefix returns keys matched by a prefix using Redis SCAN.
// func (r *RedisClient) SearchKeysByPrefix(ctx context.Context, prefix string, limit int64) ([]string, error) {
// 	if !r.valid() {
// 		return nil, fmt.Errorf("redis client not connected")
// 	}
//
// 	pattern := buildPrefixPattern(prefix)
// 	if limit <= 0 {
// 		limit = 100
// 	}
//
// 	keys := make([]string, 0, limit)
// 	iter := r.c.Scan(ctx, 0, pattern, limit).Iterator()
// 	for iter.Next(ctx) {
// 		keys = append(keys, iter.Val())
// 		if int64(len(keys)) >= limit {
// 			break
// 		}
// 	}
// 	if err := iter.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return keys, nil
// }
// SearchByPrefix returns raw values matched by a prefix using Redis SCAN and MGET.
// func (r *RedisClient) SearchByPrefix(ctx context.Context, prefix string, limit int64) ([][]byte, error) {
// 	keys, err := r.SearchKeysByPrefix(ctx, prefix, limit)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(keys) == 0 {
// 		return nil, nil
// 	}
//
// 	values, err := r.c.MGet(ctx, keys...).Result()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	result := make([][]byte, 0, len(values))
// 	for _, value := range values {
// 		switch typed := value.(type) {
// 		case string:
// 			result = append(result, []byte(typed))
// 		case []byte:
// 			result = append(result, typed)
// 		case nil:
// 			continue
// 		default:
// 			return nil, fmt.Errorf("unsupported redis value type: %T", value)
// 		}
// 	}
//
// 	return result, nil
// }

// SearchKeysByPrefix returns keys matched by a prefix using Redis SCAN.
func (r *RedisClient) SearchKeysByPrefix(ctx context.Context, prefix string, limit int64) ([]string, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}

	pattern := buildPrefixPattern(prefix)
	if limit <= 0 {
		limit = 100
	}

	const scanCount int64 = 5000

	keys := make([]string, 0, limit)
	var cursor uint64

	for {
		batch, nextCursor, err := r.c.Scan(ctx, cursor, pattern, scanCount).Result()
		// keys, cursor, err := r.c.ScanType(ctx, cursor, pattern, scanCount, "string").Result()
		if err != nil {
			return nil, err
		}

		if len(batch) > 0 {
			remaining := int(limit) - len(keys)
			if remaining > 0 {
				if len(batch) > remaining {
					batch = batch[:remaining]
				}
				keys = append(keys, batch...)
			}
		}

		if int64(len(keys)) >= limit {
			break
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// SearchByPrefix returns raw values matched by a prefix using Redis SCAN and MGET.
func (r *RedisClient) SearchByPrefix(ctx context.Context, prefix string, limit int64, count int64) ([][]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}

	keys, err := r.SearchKeysByPrefix(ctx, prefix, limit)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}

	values, err := r.c.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([][]byte, 0, len(values))
	for _, value := range values {
		switch typed := value.(type) {
		case string:
			result = append(result, []byte(typed))
		case []byte:
			result = append(result, typed)
		case nil:
			continue
		default:
			return nil, fmt.Errorf("unsupported redis value type: %T", value)
		}
	}

	return result, nil
}

// SearchPartnerTripKeysByPrefix scans the PartnerTrip namespace first, then filters by the target prefix in memory.
func (r *RedisClient) SearchPartnerTripKeysByPrefix(ctx context.Context, prefix string, limit int64, maxRounds int) ([]string, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}

	if limit <= 0 {
		limit = 100
	}
	if maxRounds <= 0 {
		maxRounds = 2
	}

	const (
		basePattern = "PartnerTrip:*"
		scanCount   = 5000
	)

	keys := make([]string, 0, limit)
	var cursor uint64

	for round := 0; round < maxRounds; round++ {
		batch, nextCursor, err := r.c.ScanType(ctx, cursor, basePattern, scanCount, "string").Result()
		if err != nil {
			return nil, err
		}

		if len(batch) > 0 {
			remaining := int(limit) - len(keys)
			for _, key := range batch {
				if !strings.HasPrefix(key, prefix) {
					continue
				}
				keys = append(keys, key)
				remaining--
				if remaining <= 0 {
					return keys, nil
				}
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (r *RedisClient) lock(ctx context.Context, k string, t time.Duration) {
	return
}
func (r *RedisClient) unlock(ctx context.Context, k string) error {
	return nil
}

func (r *RedisClient) valid() bool {
	return r != nil && r.IsConnected
}

func (r *RedisClient) keepAliveFunc(callback func(err error)) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if r == nil || r.c == nil {
				fmt.Printf("[%s] Redis client is nil, skipping ping. \n", r.Name)
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := r.c.Ping(ctx).Err()
			cancel()
			if err != nil {
				fmt.Printf("[%s] Redis ping failed: %v\n", r.Name, err)
				r.IsConnected = false
				callback(err)
				continue
			}
			r.IsConnected = true
			callback(nil)
		}
	}
}

// LPush pushes one or more values to the left of the list (head).
func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	if !r.valid() {
		return fmt.Errorf("redis client not connected")
	}
	return r.c.LPush(ctx, key, values...).Err()
}

func (r *RedisClient) RPush(ctx context.Context, key string, values ...interface{}) error {
	if !r.valid() {
		return fmt.Errorf("redis client not connected")
	}
	return r.c.RPush(ctx, key, values...).Err()
}

// LPop — pop (remove and return) the first element from a list
func (r *RedisClient) LPop(ctx context.Context, key string) ([]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}
	val, err := r.c.LPop(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // no element in list
		}
		return nil, err
	}
	return val, nil
}

// RPop — pop (remove and return) the last element from a list
func (r *RedisClient) RPop(ctx context.Context, key string) ([]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}
	val, err := r.c.RPop(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

// LRange — read all (or a range of) elements from a list
func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([][]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}
	values, err := r.c.LRange(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	result := make([][]byte, len(values))
	for i, v := range values {
		result[i] = []byte(v)
	}
	return result, nil
}

// This is great for workflow consumers or background workers that wait for user input or next-step data.
func (r *RedisClient) BLPop(ctx context.Context, timeout time.Duration, key string) ([]byte, error) {
	if !r.valid() {
		return nil, fmt.Errorf("redis client not connected")
	}

	// BLPop returns a slice: [key, value]
	result, err := r.c.BLPop(ctx, timeout, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // no data before timeout
		}
		return nil, err
	}

	// result[0] = key, result[1] = value
	if len(result) < 2 {
		return nil, fmt.Errorf("invalid BLPop result: %+v", result)
	}

	return []byte(result[1]), nil
}

func buildPrefixPattern(prefix string) string {
	if prefix == "" {
		return "*"
	}
	if strings.HasSuffix(prefix, "*") {
		return prefix
	}
	return prefix + "*"
}
