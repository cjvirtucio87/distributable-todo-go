package config

import (
	"github.com/spf13/viper"
	"log"
)

type viperLoader struct {
}

func (v *viperLoader) GetStringSlice(key string) {
}
func (v *viperLoader) GetInt(key string) {
}
func (v *viperLoader) GetString(key string) {
}
func (v *viperLoader) Load() {
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}
}
