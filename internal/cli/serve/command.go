package serve

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/api"
	"github.com/alxeg/flibooks/internal/config"
)

type Cmd struct {
	Command *cobra.Command

	AppConfig *config.App
}

// Module is a fx module
var Module = fx.Options(
	api.Module,

	fx.Provide(NewCommand),
)

// Params represents the module input params
type Params struct {
	fx.In

	AppConfig  *config.App
	RestServer api.RestServer
}

func NewCommand(p Params) (*Cmd, error) {

	c := &Cmd{
		Command: &cobra.Command{
			Use:   "serve",
			Short: "Serves the REST API",
		},

		AppConfig: p.AppConfig,
	}

	c.Command.Run = func(cmd *cobra.Command, args []string) {
		p.RestServer.StartListen()
	}
	return c, nil
}
