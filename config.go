package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	for k, v := range getDefaultConfig() {
		viper.SetDefault(k, v)
	}
}

func getDefaultConfig() map[string]string {
	return map[string]string{
		"APPLE_CHARITY_URL": "https://tw.feature.appledaily.com/charity/projlist",
	}
}
