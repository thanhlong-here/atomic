package cache

import "time"

type CacheConfig struct {
	TTL              time.Duration
	CleanInterval    time.Duration
	MaxRAMMB         int
	Shards           int
	MaxEntrySize     int
	EstimatedEntries int
}

// DefaultCacheConfig: dùng cho dev hoặc EC2 free tier
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		TTL:              10 * time.Minute,
		CleanInterval:    5 * time.Minute,
		MaxRAMMB:         32,
		Shards:           16,
		MaxEntrySize:     512,
		EstimatedEntries: 2000,
	}
}
