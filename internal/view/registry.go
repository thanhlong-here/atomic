package view

var Views = map[string]ViewConfig{}

type ViewConfig struct {
	View       DynamicView
	SyncModels []string
}

func Register(v DynamicView) {
	Views[v.Name()] = ViewConfig{
		View:       v,
		SyncModels: v.SyncModels(),
	}
}
