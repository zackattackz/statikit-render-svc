package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/zackattackz/statikit-render-svc/internal/models"
)

// Stores render results, which can be set/retrieved via a CacheKey
type CacheService interface {
	// On success returns the render result from the cache at input key, returns error != nil if failed
	Get(CacheKey) (string, error)
	// Maps input key to input render result inside the cache, returns error != nil if failed
	Set(CacheKey, string) error
}

// Composite key used by Cache to uniquely identify a render result
type CacheKey struct {
	Schema   models.Schema
	Contents string
}

// Returns base64 endcoded sha256 sum of key.Schema + k.Contents
func (k CacheKey) Hash() string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%v", k.Schema) + k.Contents))
	return base64.StdEncoding.EncodeToString(sum[:])
}

type cacheEndpointGetRequest struct {
	k CacheKey
}

type cacheEndpointSetRequest struct {
	k CacheKey
	v string
}

func EndPointFromCacheService(s CacheService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		switch v := request.(type) {
		case cacheEndpointGetRequest:
			return s.Get(v.k)
		case cacheEndpointSetRequest:
			err := s.Set(v.k, v.v)
			return nil, err
		default:
			panic(fmt.Sprintf("CacheEndpoint: unexpected request type: %T", request))
		}
	}
}

type cacheEndpointService struct {
	ctx context.Context
	e   endpoint.Endpoint
}

func CacheServiceFromEndpoint(ctx context.Context, e endpoint.Endpoint) CacheService {
	return cacheEndpointService{ctx, e}
}

func (s cacheEndpointService) Get(k CacheKey) (string, error) {
	resp, err := s.e(s.ctx, cacheEndpointGetRequest{k})
	if err != nil {
		return "", err
	}
	return resp.(string), nil
}

func (s cacheEndpointService) Set(k CacheKey, v string) error {
	_, err := s.e(s.ctx, cacheEndpointSetRequest{k, v})
	return err
}
