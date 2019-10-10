package config

import (
	"github.com/spf13/viper"
	"os"
)

const (
	configDirName = "DCONFIG_DIR"
)

type Loader interface {
	GetStringSlice(key string)
	GetInt(key string)
	GetString(key string)
	Load()
}

func NewViperLoader(filename, filetype string) Loader {
	viper.SetConfigName(filename)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(os.Getenv(configDirName))

	return &viperLoader{}
}
