package config

import (
	"github.com/spf13/viper"
	"log"
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

	configDir := os.Getenv(configDirName)

	if len(configDir) != 0 {
		viper.AddConfigPath(configDir)
	} else {
		if workingDir, err := os.Getwd(); err != nil {
			log.Fatalf(
				"Error attempting to load configuration from working directory\n%s",
				err.Error(),
			)
		} else {
			viper.AddConfigPath(workingDir)
		}
	}

	return &viperLoader{}
}
