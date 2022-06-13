package parse

import (
	"github.com/jinzhu/copier"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/config"
	"github.com/alxeg/flibooks/internal/db"
	"github.com/alxeg/flibooks/internal/db/orm"
	"github.com/alxeg/flibooks/pkg/inpx"
	"github.com/alxeg/flibooks/pkg/inpx/models"
)

type Cmd struct {
	Command *cobra.Command

	AppConfig *config.App
	DB        db.DataStorer
}

// Module is a fx module
var Module = fx.Options(
	db.Module,
	fx.Provide(NewCommand),
)

// Params represents the module input params
type Params struct {
	fx.In

	AppConfig *config.App
	DB        db.DataStorer
}

func NewCommand(p Params) (*Cmd, error) {

	c := &Cmd{
		Command: &cobra.Command{
			Use:   "parse <inpx file>",
			Short: "Parse the inpx file",
			Args:  cobra.ExactArgs(1),
		},

		AppConfig: p.AppConfig,
		DB:        p.DB,
	}

	c.Command.Run = func(cmd *cobra.Command, args []string) {
		c.DoParse(args[0])
	}
	return c, nil
}

func (cmd *Cmd) DoParse(inpxFile string) {
	inpx.ReadInpxFile(inpxFile, cmd)
}

func (cmd *Cmd) ProcessBook(book *models.Book) error {
	b := &orm.Book{}
	copier.CopyWithOption(b, book, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	exists, _ := cmd.DB.IsBookExist(b)
	if !exists {
		cmd.DB.PutBook(b)
	}
	return nil
}

func (cmd *Cmd) FinishProcessing() {

}
