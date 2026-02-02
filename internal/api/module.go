package api

import (
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/config"
	"github.com/alxeg/flibooks/internal/db"
	"github.com/alxeg/flibooks/internal/services/convert"
)

var Module = fx.Options(
	convert.Module,

	fx.Provide(NewApi),
)

type Props struct {
	fx.In

	AppConfig *config.App
	Converter convert.Converter
	DB        db.DataStorer
}

func NewApi(p Props) (RestServer, error) {
	return NewRestService(
		p.AppConfig.Server.Listen,
		p.AppConfig.ApiPrefix,
		p.DB,
		p.AppConfig.Data.Dir,
		p.Converter,
		p.AppConfig.StaticsDir,
		p.AppConfig.StaticsRoute), nil
}
