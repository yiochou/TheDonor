package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	AppleCharityUrl       string `mapstructure:"APPLE_CHARITY_URL"`
	MongodbUri            string `mapstructure:"MONGODB_URI"`
	MongodbUsername       string `mapstructure:"MONGODB_USERNAME"`
	MongodbPassword       string `mapstructure:"MONGODB_PASSWORD"`
	Database              string `mapstructure:"DATABASE"`
	TwitterConsumerKey    string `mapstructure:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret string `mapstructure:"TWITTER_CONSUMER"`
	TwitterToken          string `mapstructure:"TWITTER_TOKEN"`
	TwitterTokenSecret    string `mapstructure:"TWITTER_TOKEN_SECRET"`
}

func LoadConfig(path string) (Config, error) {
	config := &Config{
		AppleCharityUrl: "https://tw.feature.appledaily.com/charity/projlist",
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return *config, err
	}

	err = viper.Unmarshal(config)

	log.Info("loaded config")

	return *config, err
}
