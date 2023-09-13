package lib

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"manuel71sj/go-api-template/constants"
	"manuel71sj/go-api-template/errors"
	"time"
)

type Redis struct {
	cache  *cache.Cache
	client *redis.Client
	prefix string
}

// NewRedis creates a new redis client instance
func NewRedis(config Config, logger Logger) Redis {
	addr := config.Redis.Addr()

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       constants.RedisMainDB,
		Password: config.Redis.Password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Zap.Fatalf("Error to open redis[%s] connection: %v", addr, err)
	}

	logger.Zap.Info("Redis connection established")

	return Redis{
		client: client,
		prefix: config.Redis.KeyPrefix,
		cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
	}
}

func (r Redis) wrapperKey(key string) string {
	return fmt.Sprintf("%s:%s", r.prefix, key)
}

func (r Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:            context.TODO(),
		Key:            r.wrapperKey(key),
		Value:          value,
		TTL:            expiration,
		SkipLocalCache: true,
	})
}

func (r Redis) Get(key string, value interface{}) error {
	err := r.cache.Get(context.TODO(), r.wrapperKey(key), value)
	if errors.Is(err, cache.ErrCacheMiss) {
		err = errors.RedisKeyNoExist
	}

	return err
}

func (r Redis) Delete(keys ...string) (bool, error) {
	wrapperKeys := make([]string, len(keys))
	for index, key := range keys {
		wrapperKeys[index] = r.wrapperKey(key)
	}

	cmd := r.client.Del(context.TODO(), wrapperKeys...)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

func (r Redis) Check(keys ...string) (bool, error) {
	wrapperKeys := make([]string, len(keys))
	for index, key := range keys {
		wrapperKeys[index] = r.wrapperKey(key)
	}

	cmd := r.client.Exists(context.TODO(), wrapperKeys...)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val() > 0, nil
}

func (r Redis) Close() error {
	return r.client.Close()
}

func (r Redis) GetClient() *redis.Client {
	return r.client
}
