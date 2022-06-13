package db

import (
	"github.com/alxeg/flibooks/internal/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewDb),
)

// Params represents the module input params
type Params struct {
	fx.In

	AppConfig *config.App
}

// NewCli instantiates the main app
func NewDb(p Params) (DataStorer, error) {
	return NewDBStore(p.AppConfig.Database.Type, p.AppConfig.Database.Connection, p.AppConfig.Database.LogLevel)
}
