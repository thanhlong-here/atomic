package view

import (
	"atomic/internal/cache"
	"fmt"
	"net/http"
)

// Gọi rebuild thủ công
func RebuildView(name string, r *http.Request) (map[string]interface{}, bool, error) {
	viewCfg, ok := Views[name]
	if !ok {
		return nil, false, fmt.Errorf("view '%s' not found", name)
	}

	data, shouldCache, err := viewCfg.View.Rebuild(r)
	if err != nil {
		return nil, false, err
	}

	if shouldCache {
		cache.SetRaw("view:"+name+r.URL.RawQuery, data)
	}
	return data, shouldCache, nil
}
func GetCache(key string) (map[string]interface{}, bool) {
	return cache.GetRaw(key)
}

func SetCache(key string, data map[string]interface{}) {
	cache.SetRaw(key, data)
}

func emptyRequest() *http.Request {
	req, _ := http.NewRequest("GET", "/", nil)
	return req
}

func TriggerSync(model string) {
	req := emptyRequest()
	for viewName, cfg := range Views {
		for _, m := range cfg.SyncModels {
			if m == model {
				go RebuildView(viewName, req)
				break
			}
		}
	}
}
