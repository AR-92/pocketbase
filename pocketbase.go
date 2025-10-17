package pocketbase

import (
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type PocketBase struct {
	*core.BaseApp
}

type Config struct {
	DefaultDataDir string
	DefaultDev     bool
}

func NewWithConfig(config Config) *PocketBase {
	app := core.NewBaseApp(core.BaseAppConfig{
		DataDir: config.DefaultDataDir,
		IsDev:   config.DefaultDev,
	})
	return &PocketBase{BaseApp: app}
}

func New() *PocketBase {
	return NewWithConfig(Config{})
}

func (app *PocketBase) Start() error {
	return apis.Serve(app.BaseApp, apis.ServeConfig{
		HttpAddr: "127.0.0.1:8090",
	})
}
