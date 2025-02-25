package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client  *redis.Client
	prefix  string
	context context.Context
}

func NewRedisStore(client *redis.Client, prefix string) *RedisStore {
	return &RedisStore{
		client:  client,
		prefix:  prefix,
		context: context.Background(),
	}
}

func (r *RedisStore) Get(key string) ([]byte, bool) {
	prefixedKey := r.prefixKey(key)
	val, err := r.client.Get(r.context, prefixedKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false
		}

		fmt.Printf("redis Get error: %v\n", err)
		return nil, false
	}
	return val, true
}

func (r *RedisStore) Set(key string, data []byte, ttl time.Duration) {
	prefixedKey := r.prefixKey(key)
	err := r.client.Set(r.context, prefixedKey, data, ttl).Err()
	if err != nil {

		fmt.Printf("redis Set error: %v\n", err)
	}
}

func (r *RedisStore) Del(key string) {
	prefixedKey := r.prefixKey(key)
	err := r.client.Del(r.context, prefixedKey).Err()
	if err != nil {

		fmt.Printf("redis Del error: %v\n", err)
	}
}

func (r *RedisStore) Update(key string, data []byte) error {
	prefixedKey := r.prefixKey(key)
	exists, err := r.client.Exists(r.context, prefixedKey).Result()
	if err != nil {
		return fmt.Errorf("redis exists error: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("key not exists")
	}
	err = r.client.Set(r.context, prefixedKey, data, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}
	return nil
}

func (r *RedisStore) prefixKey(key string) string {
	if r.prefix != "" {
		return r.prefix + ":" + key
	}
	return key
}
