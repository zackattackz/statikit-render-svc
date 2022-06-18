package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zackattackz/statikit-render-svc/internal/service"
)

// Implementation of service.CacheService for redis
type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCacheService(client *redis.Client, ttl time.Duration) service.CacheService {
	return redisCache{client, ttl}
}

func (c redisCache) Get(k service.CacheKey) (string, error) {
	cmd := c.client.Get(c.client.Context(), k.Hash())
	return cmd.Result()
}

func (c redisCache) Set(k service.CacheKey, v string) error {
	cmd := c.client.Set(c.client.Context(), k.Hash(), v, c.ttl)
	return cmd.Err()
}
