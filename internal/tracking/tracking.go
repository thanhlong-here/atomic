package tracking

import (
	"atomic/internal/cache"
	"time"
)

var viewLastSynced = map[string]map[string]time.Time{}

// Mark khi model được thay đổi
func MarkModelUpdated(model string) {
	key := "tracking_model:" + model
	now := time.Now().UTC().Format(time.RFC3339)
	cache.SetRawString(key, now)
}

// Lấy thời gian model được cập nhật gần nhất
func GetModelUpdatedTime(model string) time.Time {
	key := "tracking_model:" + model
	val, found := cache.GetRawString(key)
	if !found {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}
	}
	return t
}

// Lưu thời gian view sync
func MarkViewSynced(view string, models []string) {
	if viewLastSynced[view] == nil {
		viewLastSynced[view] = make(map[string]time.Time)
	}
	now := time.Now()
	for _, m := range models {
		viewLastSynced[view][m] = now
	}
}

// Kiểm tra view có cần rebuild không
func ShouldRebuild(view string, models []string) bool {
	for _, m := range models {
		lastUpdate := GetModelUpdatedTime(m)
		lastSync := viewLastSynced[view][m]
		if lastUpdate.IsZero() || lastUpdate.After(lastSync) {
			return true
		}
	}
	return false
}

// Trả về toàn bộ model đang được tracking
func GetAllModelTrackingStatus() map[string]string {
	out := map[string]string{}
	for _, k := range cache.KeysPrefix("tracking_model:") {
		val, ok := cache.GetRawString(k)
		if ok {
			model := k[len("tracking_model:"):]
			out[model] = val
		}
	}
	return out
}
