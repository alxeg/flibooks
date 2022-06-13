package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/cli/parse"
	"github.com/alxeg/flibooks/internal/cli/serve"
	"github.com/alxeg/flibooks/internal/config"
)

// App represents the application
type Cli struct {
	rootCmd *cobra.Command
}

// Module is a fx module
var Module = fx.Options(
	parse.Module,
	serve.Module,
	fx.Provide(NewCli),

	fx.Invoke(StartCli),
)

// Params represents the module input params
type Params struct {
	fx.In

	AppConfig     *config.App
	ParserCommand *parse.Cmd
	ServeCommand  *serve.Cmd

	Shutdowner fx.Shutdowner
}

// NewCli instantiates the main app
func NewCli(p Params) (*Cli, error) {
	rootCmd := &cobra.Command{
		Use:   "flibooks",
		Short: "The flibooks inpx app for processing ebooks inpx archives",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if cmd.Use != "serve" {
				p.Shutdowner.Shutdown()
			}
		},
	}
	var configPath string
	rootCmd.PersistentFlags().StringVar(&configPath, config.ConfigPathArg, "", "Path to config file")
	rootCmd.AddCommand(p.ParserCommand.Command)
	rootCmd.AddCommand(p.ServeCommand.Command)

	mainApp := &Cli{
		rootCmd: rootCmd,
	}

	return mainApp, nil
}

// RegisterRoutes registers the api routes and starts the http server
func StartCli(app *Cli, lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			go func() {
				err := app.rootCmd.Execute()
				if err != nil {
					os.Exit(-1)
				}
			}()
			return nil
		},
		OnStop: func(c context.Context) error {
			return nil
		},
	})

}
