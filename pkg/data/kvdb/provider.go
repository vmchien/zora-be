package kvdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
	"vn.vato.zora.be.api/pkg/data/redis_db"
	"vn.vato.zora.be.api/pkg/encode"
)

type KvClientType int

const (
	Redis    KvClientType = 0
	MemCache KvClientType = 1
)

type KvClient interface {
	Connect(config map[string]any, keepAliveFunc func(err error)) error
}

var getOrAddSingleflight singleflight.Group

func Get[T any](cli KvClient, ctx context.Context, key string) (*T, error) {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		var rs *T
		data, err := typed.Get(ctx, key)
		if err != nil {
			return nil, err
		}
		err = encode.Decode(data, &rs)
		return rs, err
	}

	return nil, notfoundProvider()
}

func Set[T any](cli KvClient, ctx context.Context, key string, val T, expiration time.Duration) error {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		buf, err := encode.Encode(val)
		if err != nil {
			return err
		}
		return typed.Set(ctx, key, buf, expiration)
	}
	return notfoundProvider()
}

func Delete(cli KvClient, ctx context.Context, key string) error {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		return typed.Del(ctx, key)
	}
	return notfoundProvider()
}

func GetOrAdd[T any](cli KvClient, ctx context.Context, key string, expiration time.Duration, fetchFunc func() (*T, error)) (*T, error) {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		data, err := typed.Fetch(ctx, key)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			var rs *T
			err = encode.Decode(data, &rs)
			return rs, err
		}

		sharedKey := fmt.Sprintf("kvdb:get_or_add:%T:%s", *new(T), key)
		v, err, _ := getOrAddSingleflight.Do(sharedKey, func() (any, error) {
			// Re-check cache inside singleflight so only one goroutine per key actually runs fetchFunc.
			data, err := typed.Fetch(ctx, key)
			if err != nil {
				return nil, err
			}
			if len(data) > 0 {
				var cached *T
				if err = encode.Decode(data, &cached); err != nil {
					return nil, err
				}
				return cached, nil
			}

			val, err := fetchFunc()
			if err != nil {
				return nil, err
			}
			if val == nil {
				return nil, nil
			}

			data, err = encode.Encode(val)
			if err != nil {
				return nil, err
			}
			if err = typed.Set(ctx, key, data, expiration); err != nil {
				return nil, err
			}
			return val, nil
		})
		if err != nil {
			return nil, err
		}
		if v == nil {
			return nil, nil
		}

		rs, ok := v.(*T)
		if !ok {
			return nil, fmt.Errorf("singleflight result type mismatch for key %s", key)
		}
		return rs, nil
	}
	return nil, notfoundProvider()
}

func SetJSON[T any](cli KvClient, ctx context.Context, key string, val T, expiration time.Duration) error {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		buf, err := json.Marshal(val)
		if err != nil {
			return err
		}
		return typed.Set(ctx, key, buf, expiration)
	}
	return notfoundProvider()
}

func GetJSON[T any](cli KvClient, ctx context.Context, key string) (*T, error) {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		data, err := typed.Get(ctx, key)
		if err != nil {
			return nil, err
		}

		var rs T
		if err = json.Unmarshal(data, &rs); err != nil {
			return nil, err
		}
		return &rs, nil
	}

	return nil, notfoundProvider()
}

// // SearchKeysByPrefix returns keys matched by a prefix for the active KV provider.
// func SearchKeysByPrefix(cli KvClient, ctx context.Context, prefix string, limit int64) ([]string, error) {
// 	switch typed := cli.(type) {
// 	case *redis_db.RedisClient:
// 		return typed.SearchKeysByPrefix(ctx, prefix, limit)
// 	}
// 	return nil, notfoundProvider()
// }
// // SearchByPrefix returns decoded values matched by a prefix for the active KV provider.
// func SearchByPrefix[T any](cli KvClient, ctx context.Context, prefix string, limit int64) ([]*T, error) {
// 	switch typed := cli.(type) {
// 	case *redis_db.RedisClient:
// 		dataList, err := typed.SearchByPrefix(ctx, prefix, limit)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		result := make([]*T, 0, len(dataList))
// 		for _, data := range dataList {
// 			var item *T
// 			if err = encode.Decode(data, &item); err != nil {
// 				return nil, err
// 			}
// 			if item != nil {
// 				result = append(result, item)
// 			}
// 		}
// 		return result, nil
// 	}
// 	return nil, notfoundProvider()
// }

// SearchKeysByPrefix returns keys matched by a prefix for the active KV provider.
func SearchKeysByPrefix(cli KvClient, ctx context.Context, prefix string, limit int64, count int64) ([]string, error) {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		return typed.SearchKeysByPrefix(ctx, prefix, limit)
	}
	return nil, notfoundProvider()
}

// SearchByPrefix returns decoded values matched by a prefix for the active KV provider.
func SearchByPrefix[T any](cli KvClient, ctx context.Context, prefix string, limit int64, count int64) ([]*T, error) {
	switch typed := cli.(type) {
	case *redis_db.RedisClient:
		dataList, err := typed.SearchByPrefix(ctx, prefix, limit, count)
		if err != nil {
			return nil, err
		}

		result := make([]*T, 0, len(dataList))
		for _, data := range dataList {
			var item *T
			if err = encode.Decode(data, &item); err != nil {
				return nil, err
			}
			if item != nil {
				result = append(result, item)
			}
		}
		return result, nil
	}
	return nil, notfoundProvider()
}

func notfoundProvider() error {
	return fmt.Errorf("not found key-value database provider")
}
