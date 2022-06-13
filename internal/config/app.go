package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/ogier/pflag"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	ConfigPathArg     = "config-path"
	ConfigPathDefault = "./flibooks.properties"
)

type Database struct {
	Type       string
	Connection string
	LogLevel   string
}

type Server struct {
	Listen string
}

type Data struct {
	Dir string
}

type App struct {
	Database Database
	Server   Server
	Data     Data
}

var Module = fx.Provide(NewAppConfig)

// NewAppConfig instantiates the main app config
func NewAppConfig() (*App, error) {
	mainConfig := &App{}

	configPaths := []string{ConfigPathDefault}

	configPath := pflag.String(ConfigPathArg, "", "")
	pflag.Parse()

	configPaths = append(configPaths, *configPath)

	conf := viper.New()
	for num, path := range configPaths {
		if len(path) < 1 {
			continue
		}
		tempConf := viper.New()
		tempConf.SetConfigFile(path)
		err := tempConf.ReadInConfig()
		if err != nil {
			// complain on missed non-default config
			if num > 0 {
				fmt.Printf("Can't read config from %v, Error: %v\n", path, err)
			}
		} else {
			_ = conf.MergeConfigMap(tempConf.AllSettings())
		}
	}
	err := conf.Unmarshal(mainConfig, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeDurationHookFunc(),
		)))

	return mainConfig, err
}
