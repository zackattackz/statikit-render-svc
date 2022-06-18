package cache

// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// 	"github.com/zackattackz/statikit-render-svc/internal/models"
// 	"github.com/zackattackz/statikit-render-svc/internal/ports"
// )

// type mockRedisStringCmd struct {
// 	redis.StringCmd
// 	res interface{}
// 	err error
// }

// func (c mockRedisStringCmd) Result() (interface{}, error) {
// 	return c.res, c.err
// }

// type mockRedisStatusCmd struct {
// 	err error
// }

// func (c mockRedisStatusCmd) Err() error {
// 	return c.err
// }

// type mockRedisClient struct {
// 	ctx        context.Context
// 	m          map[string]interface{}
// 	shouldFail bool
// }

// func newMockRedisClient(ctx context.Context) redisClient {
// 	m := make(map[string]interface{})
// 	m[initialKey.Hash()] = initialValue
// 	return mockRedisClient{ctx, m, false}
// }

// func (c mockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
// 	if c.shouldFail {
// 		return mockRedisStringCmd{nil, errFailure}
// 	}

// 	v, ok := c.m[key]
// 	if !ok {
// 		return mockRedisStringCmd{nil, errNoValue}
// 	}
// 	return mockRedisStringCmd{v, nil}
// }

// func (c mockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
// 	if c.shouldFail {
// 		return mockRedisStatusCmd{errFailure}
// 	}
// 	c.m[key] = value
// 	return mockRedisStatusCmd{nil}
// }

// func (c mockRedisClient) Context() context.Context {
// 	return c.ctx
// }

// var (
// 	errNoValue error = errors.New("no value at key")
// 	errFailure       = errors.New("shouldFail was true")

// 	initialKey   ports.CacheKey = ports.CacheKey{Schema: models.Schema{}, Contents: "initially will be in the mock cache map as a key"}
// 	initialValue string         = "initially will be in the mock cache map as a val"
// )

// func testGet(t *testing.T, cache ports.Cache, k ports.CacheKey, expectedVal string, expectedErr error) {
// 	v, err := cache.Get(k)
// 	if err != expectedErr {
// 		t.Fatalf("expected error != actual error: expected=%v , actual=%v", expectedErr, err)
// 	}
// 	if v != expectedVal {
// 		t.Fatalf("expected != actual: expected=%v , actual=%v", expectedVal, v)
// 	}
// }

// func TestRedisCache(t *testing.T) {
// 	mockClient := newMockRedisClient(context.Background())
// 	cache := NewRedisCache(mockClient, 0)

// 	// Test get works with initial key/value
// 	testGet(t, cache, initialKey, initialValue, nil)
// 	//
// }
