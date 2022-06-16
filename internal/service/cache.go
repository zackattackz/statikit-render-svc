package service

import (
	"crypto/sha256"
	"fmt"
)

// Stores render results
type Cache interface {
	Get(CacheKey) (string, error)
	Set(CacheKey, string) error
}

// Composite key to uniquely identify a render result
type CacheKey struct {
	Schema   Schema
	Contents string
}

func (k CacheKey) Hash() string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%v", k.Schema) + k.Contents))
	return string(sum[:])
}
