package convert

import (
	"fmt"
	"log"

	"github.com/go-cmd/cmd"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/alxeg/flibooks/internal/config"
)

var Module = fx.Options(
	fx.Provide(NewConverter),
)

type Converter interface {
	Convert(src, dst, format string) error
}

type converter struct {
	AppConfig *config.App
}

type Params struct {
	fx.In

	AppConfig *config.App
}

func NewConverter(p Params) (Converter, error) {
	return &converter{
		AppConfig: p.AppConfig,
	}, nil
}

func (c *converter) Convert(src, dst, format string) error {
	cmdOptions := cmd.Options{
		Streaming: true,
		Buffered:  false,
	}

	command := cmd.NewCmdOptions(cmdOptions, c.AppConfig.Fb2C.Path, "convert", "--nodirs", "--ow", "--to", format, src, dst)
	runTimeStatus := command.Start()
	log.Printf("Executing: %s %s %s %s %s %s %s %s\n", c.AppConfig.Fb2C.Path, "convert", "--nodirs", "--ow", "--to", format, src, dst)
	go processOutput(command)

	status := <-runTimeStatus
	if status.Error != nil || !status.Complete || status.Exit != 0 {
		// Exit with error
		if status.Error != nil {
			return errors.Wrapf(status.Error, "Failed to convert %s", src)
		} else {
			return errors.New(fmt.Sprintf("Execution was failed with exit code %d", status.Exit))
		}
	}
	return nil
}

func processOutput(cmd *cmd.Cmd) {
	log.Println("Started to follow the process output")
	for cmd.Stdout != nil || cmd.Stderr != nil {
		select {
		case line, open := <-cmd.Stdout:
			if !open {
				cmd.Stdout = nil
				continue
			}
			log.Println(line)
		case line, open := <-cmd.Stderr:
			if !open {
				cmd.Stderr = nil
				continue
			}
			log.Println(line)
		}
	}
	log.Println("Finished to follow the process output")
}
