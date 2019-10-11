package config

import (
	"github.com/spf13/viper"
	"os"
)

const (
	configDirName = "DCONFIG_DIR"
)

type Loader interface {
	Unmarshal(obj interface{}) error
	Load() error
}

func NewViperLoader(filename, filetype string) Loader {
	viper.SetConfigName(filename)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(os.Getenv(configDirName))

	return &viperLoader{}
}
