package redis

import (
	"context"
	"time"

	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	redisSdk "github.com/go-redis/redis/v8"
	"github.com/zackattackz/statikit-render-svc/internal/service"
)

type cache struct {
	client *redisSdk.Client
	ttl    time.Duration
}

func New(client *redisSdk.Client, ttl time.Duration) cache {
	return cache{client, ttl}
}

func (c cache) Get(k service.CacheKey) (string, error) {
	e := c.makeGetEndpoint()
	e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
	resp, err := e(c.client.Context(), k)
	return resp.(string), err
}

// Combines a key and value
type kvPair struct {
	k service.CacheKey
	v string
}

func (c cache) Set(k service.CacheKey, v string) error {
	e := c.makeSetEndpoint()
	e = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(e)
	_, err := e(c.client.Context(), kvPair{k, v})
	return err
}

func (c cache) makeGetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		k := request.(service.CacheKey)
		resp := c.client.Get(c.client.Context(), k.Hash())
		return resp.Result()
	}
}

func (c cache) makeSetEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		kvPair := request.(kvPair)
		resp := c.client.Set(c.client.Context(), kvPair.k.Hash(), kvPair.v, c.ttl)
		return nil, resp.Err()
	}
}
