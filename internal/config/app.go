package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/go-viper/encoding/javaproperties"
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

type Fb2C struct {
	Path string
}

type App struct {
	Database Database
	Server   Server
	Data     Data
	Fb2C     Fb2C

	StaticsDir   string
	StaticsRoute string
}

var Module = fx.Provide(NewAppConfig)

func getConfigData(filePath string) (string, string, string) {
	dir, file := filepath.Split(filePath)
	base := filepath.Base(file)
	ext := filepath.Ext(base)

	confPath, _ := filepath.Abs(dir)
	confName := strings.TrimSuffix(base, ext)
	confType := strings.Trim(ext, ".")

	return confPath, confName, confType
}

// NewAppConfig instantiates the main app config
func NewAppConfig() (*App, error) {
	mainConfig := &App{}

	configPaths := []string{ConfigPathDefault}

	configPath := pflag.String(ConfigPathArg, "", "")
	pflag.Parse()

	configPaths = append(configPaths, *configPath)

	codecRegistry := viper.NewCodecRegistry()
	codec := &javaproperties.Codec{}
	codecRegistry.RegisterCodec("properties", codec)
	codecRegistry.RegisterCodec("props", codec)
	codecRegistry.RegisterCodec("prop", codec)

	conf := viper.New()

	for num, path := range configPaths {
		if len(path) < 1 {
			continue
		}
		tempConf := viper.NewWithOptions(
			viper.WithCodecRegistry(codecRegistry),
		)

		confPath, confName, confType := getConfigData(path)
		tempConf.AddConfigPath(confPath)
		tempConf.SetConfigName(confName)
		tempConf.SetConfigType(confType)

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
