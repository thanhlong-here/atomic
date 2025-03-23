package view

import (
	"atomic/internal/cache"
	"atomic/internal/tracking"
	"fmt"
	"net/http"
)

type ViewConfig struct {
	View       DynamicView
	SyncModels []string
}

var Views = make(map[string]ViewConfig)

// Đăng ký view
func RegisterView(name string, view DynamicView) {
	Views[name] = ViewConfig{
		View:       view,
		SyncModels: view.SyncModels(),
	}
}

// Hàm gọi mỗi khi /view/{view} được request
func RebuildView(name string, r *http.Request) (map[string]interface{}, bool, error) {
	viewCfg, ok := Views[name]
	if !ok {
		return nil, false, fmt.Errorf("view '%s' not found", name)
	}

	// Nếu view có SyncModels: kiểm tra trạng thái tracking
	if len(viewCfg.SyncModels) > 0 && !tracking.ShouldRebuild(name, viewCfg.SyncModels) {
		// Nếu không cần rebuild, nhưng có cache thì dùng luôn
		data, found := cache.GetRaw("view:" + name + r.URL.RawQuery)
		if found {
			return data, true, nil
		}
		// Nếu không có cache thì rebuild để cập nhật tracking
	}

	// Nếu cần rebuild hoặc cache miss → thực hiện build
	result, shouldCache, err := viewCfg.View.Rebuild(r)
	if err != nil {
		return nil, false, err
	}

	if shouldCache {
		cache.SetRaw("view:"+name+r.URL.RawQuery, result)
	}

	// Cập nhật thời gian view đã đồng bộ
	tracking.MarkViewSynced(name, viewCfg.SyncModels)

	return result, shouldCache, nil
}
