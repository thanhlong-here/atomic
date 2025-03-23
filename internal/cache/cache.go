package cache

import (
	"encoding/json"
	"fmt"

	"github.com/allegro/bigcache"
)

var instance *bigcache.BigCache

func InitWithConfig(cfg CacheConfig) {
	config := bigcache.Config{
		Shards:             cfg.Shards,
		LifeWindow:         cfg.TTL,
		CleanWindow:        cfg.CleanInterval,
		MaxEntriesInWindow: cfg.EstimatedEntries,
		MaxEntrySize:       cfg.MaxEntrySize,
		HardMaxCacheSize:   cfg.MaxRAMMB,
		Verbose:            false,
	}

	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		panic(fmt.Sprintf("Cache init failed: %v", err))
	}

	instance = cache
}

// Dùng config mặc định
func Init() {
	InitWithConfig(DefaultCacheConfig())
}

// ===== Helper =====

func key(model, id string) string {
	return fmt.Sprintf("%s:%s", model, id)
}

func Set(model, id string, data map[string]interface{}) {
	bytes, _ := json.Marshal(data)
	_ = instance.Set(key(model, id), bytes)
}

func Get(model, id string) (map[string]interface{}, bool) {
	data, err := instance.Get(key(model, id))
	if err != nil {
		return nil, false
	}
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	return result, true
}

func Delete(model, id string) {
	_ = instance.Delete(key(model, id))
}

// Raw key — dùng cho view
func SetRaw(key string, value map[string]interface{}) {
	bytes, _ := json.Marshal(value)
	_ = instance.Set(key, bytes)
}

func GetRaw(key string) (map[string]interface{}, bool) {
	data, err := instance.Get(key)
	if err != nil {
		return nil, false
	}
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	return result, true
}
