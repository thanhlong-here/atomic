package cache

import "time"

type CacheConfig struct {
	Shards             int
	LifeWindow         time.Duration
	CleanWindow        time.Duration
	MaxEntrySize       int
	HardMaxCacheSizeMB int
	Verbose            bool
}

func DefaultConfig() CacheConfig {
	return CacheConfig{
		Shards:             1024,
		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntrySize:       1024,
		HardMaxCacheSizeMB: 64,
		Verbose:            false,
	}
}
