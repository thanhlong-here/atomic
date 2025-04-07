package cache

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/allegro/bigcache/v3"
)

var memoryCache *bigcache.BigCache

// Init khởi tạo BigCache với TTL (seconds)
func Init(ttlSeconds int) error {
	config := bigcache.DefaultConfig(time.Duration(ttlSeconds) * time.Second)
	config.CleanWindow = 30 * time.Second // dọn dẹp mỗi 30s
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return err
	}
	memoryCache = cache
	return nil
}

// Set: ghi dữ liệu vào cache (dạng JSON)
func Set(key string, value interface{}) error {
	if memoryCache == nil {
		return errors.New("cache not initialized")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return memoryCache.Set(key, data)
}

// Get: lấy dữ liệu từ cache (decode về interface)
func Get(key string, dest interface{}) error {
	if memoryCache == nil {
		return errors.New("cache not initialized")
	}
	data, err := memoryCache.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete: xóa một key khỏi cache
func Delete(key string) error {
	if memoryCache == nil {
		return errors.New("cache not initialized")
	}
	return memoryCache.Delete(key)
}

// Has: kiểm tra key có tồn tại không
func Has(key string) bool {
	if memoryCache == nil {
		return false
	}
	_, err := memoryCache.Get(key)
	return err == nil
}

// AutoCache: kiểm tra cache, nếu chưa có thì gọi fetch() và lưu lại
func AutoCache(key string, ttlSeconds int, fetch func() (interface{}, error), result interface{}) error {
	// Nếu đã có cache thì trả luôn
	if err := Get(key, result); err == nil {
		return nil
	}

	// Nếu chưa có, gọi fetch
	data, err := fetch()
	if err != nil {
		return err
	}

	// Ghi vào cache
	if err := Set(key, data); err != nil {
		return err
	}

	// Gán kết quả cho result
	b, _ := json.Marshal(data)
	return json.Unmarshal(b, result)
}
