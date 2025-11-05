package config

import (
	"log"

	"github.com/spf13/viper"
)

func Read() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("config not found (using defaults)")
	}
	log.Println(viper.AllSettings())
}
