package config

import (
	"github.com/spf13/viper"
)

type viperLoader struct {
}

func (v *viperLoader) Unmarshal(obj interface{}) error {
	return viper.Unmarshal(obj)
}

func (v *viperLoader) Load() error {
	return viper.ReadInConfig()
}
