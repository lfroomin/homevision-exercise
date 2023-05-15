package main

import (
	"github.com/spf13/viper"
	"log"
)

type config struct {
	HouseServiceUrl string `mapstructure:"HOUSE_SERVICE_URL"`
	NumPages        int    `mapstructure:"NUM_PAGES"`
	NumPerPage      int    `mapstructure:"NUM_PER_PAGE"`
}

// loadConfig reads configuration from file or environment variables.
func loadConfig(path string) config {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error on reading configuration file: %s", err)
	}

	var cfg config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("error on parsing configuration file: %s", err)
	}

	return cfg
}
