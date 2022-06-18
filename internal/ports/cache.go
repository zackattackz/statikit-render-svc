package ports

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/zackattackz/statikit-render-svc/internal/models"
)

// Stores render results, which can be set/retrieved via a CacheKey
type Cache interface {
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
