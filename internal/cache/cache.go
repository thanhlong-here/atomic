package cache

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/allegro/bigcache/v3"
)

var client *bigcache.BigCache

func Init() {
	cfg := DefaultConfig()
	config := bigcache.Config{
		Shards:             cfg.Shards,
		LifeWindow:         cfg.LifeWindow,
		CleanWindow:        cfg.CleanWindow,
		MaxEntrySize:       cfg.MaxEntrySize,
		Verbose:            cfg.Verbose,
		HardMaxCacheSize:   cfg.HardMaxCacheSizeMB,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	var err error
	client, err = bigcache.NewBigCache(config)
	if err != nil {
		log.Fatal("❌ Failed to init BigCache:", err)
	}
}

// -------------------- Model-level cache --------------------

func Set(model, id string, data map[string]interface{}) {
	key := model + ":" + id
	b, _ := json.Marshal(data)
	client.Set(key, b)
}

func Get(model, id string) (map[string]interface{}, bool) {
	key := model + ":" + id
	b, err := client.Get(key)
	if err != nil {
		return nil, false
	}
	var out map[string]interface{}
	_ = json.Unmarshal(b, &out)
	return out, true
}

func Delete(model, id string) {
	key := model + ":" + id
	client.Delete(key)
}

// -------------------- Raw API (dùng cho view/tracking) --------------------

func SetRaw(key string, data map[string]interface{}) {
	b, _ := json.Marshal(data)
	client.Set(key, b)
}

func GetRaw(key string) (map[string]interface{}, bool) {
	b, err := client.Get(key)
	if err != nil {
		return nil, false
	}
	var out map[string]interface{}
	_ = json.Unmarshal(b, &out)
	return out, true
}

func SetRawString(key string, val string) {
	client.Set(key, []byte(val))
}

func GetRawString(key string) (string, bool) {
	b, err := client.Get(key)
	if err != nil {
		return "", false
	}
	return string(b), true
}

// -------------------- Tooling --------------------

// Liệt kê các key theo prefix (dành cho tracking, flush)
func KeysPrefix(prefix string) []string {
	var keys []string
	iterator := client.Iterator()

	for iterator.SetNext() {
		entry, err := iterator.Value()
		if err == nil && strings.HasPrefix(entry.Key(), prefix) {
			keys = append(keys, entry.Key())
		}
	}
	return keys
}

// Flush toàn bộ key theo prefix (vd: model posts → xóa hết posts:*)
func FlushModelCache(model string) {
	prefix := model + ":"
	iterator := client.Iterator()
	for iterator.SetNext() {
		entry, err := iterator.Value()
		if err == nil && strings.HasPrefix(entry.Key(), prefix) {
			client.Delete(entry.Key())
		}
	}
}
