package view

import "net/http"

type DynamicView interface {
	Name() string
	SyncModels() []string
	Rebuild(r *http.Request) (map[string]interface{}, bool, error)
}
