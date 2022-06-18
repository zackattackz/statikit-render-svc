package cache

import (
	"context"
	"time"

	"github.com/sony/gobreaker"
	"github.com/zackattackz/statikit-render-svc/internal/ports"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-redis/redis/v8"
)

// Implementation of ports.Cache for redis
type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(client *redis.Client, ttl time.Duration) ports.Cache {
	return redisCache{client, ttl}
}

func (c redisCache) Get(k ports.CacheKey) (string, error) {
	e := c.makeGetEndpoint()
	e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
	resp, err := e(c.client.Context(), k)
	return resp.(string), err
}

func (c redisCache) Set(k ports.CacheKey, v string) error {
	e := c.makeSetEndpoint()
	e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
	_, err := e(c.client.Context(), kvPair{k, v})
	return err
}

func (c redisCache) makeGetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		k := request.(ports.CacheKey)
		resp := c.client.Get(c.client.Context(), k.Hash())
		return resp.Result()
	}
}

func (c redisCache) makeSetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		kvPair := request.(kvPair)
		resp := c.client.Set(c.client.Context(), kvPair.k.Hash(), kvPair.v, c.ttl)
		return nil, resp.Err()
	}
}
