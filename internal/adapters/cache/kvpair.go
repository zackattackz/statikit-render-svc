package cache

import "github.com/zackattackz/statikit-render-svc/internal/ports"

// Combines a key and value
type kvPair struct {
	k ports.CacheKey
	v string
}
