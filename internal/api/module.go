package api

import (
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/config"
	"github.com/alxeg/flibooks/internal/db"
)

var Module = fx.Options(
	fx.Provide(NewApi),
)

type Props struct {
	fx.In

	AppConfig *config.App
	DB        db.DataStorer
}

func NewApi(p Props) (RestServer, error) {
	return NewRestService(p.AppConfig.Server.Listen, p.DB, p.AppConfig.Data.Dir), nil
}
