package main

import (
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/cli"
	"github.com/alxeg/flibooks/internal/config"
)

func main() {
	fxApp := fx.New(
		config.Module,
		cli.Module,

		// fx.NopLogger,
	)
	fxApp.Run()
}
