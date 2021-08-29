package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port                  string `mapstructure:"PORT"`
	AppleCharityUrl       string `mapstructure:"APPLE_CHARITY_URL"`
	MongodbUri            string `mapstructure:"MONGODB_URI"`
	MongodbUsername       string `mapstructure:"MONGODB_USERNAME"`
	MongodbPassword       string `mapstructure:"MONGODB_PASSWORD"`
	Database              string `mapstructure:"DATABASE"`
	TwitterConsumerKey    string `mapstructure:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret string `mapstructure:"TWITTER_CONSUMER_SECRET"`
	TwitterToken          string `mapstructure:"TWITTER_TOKEN"`
	TwitterTokenSecret    string `mapstructure:"TWITTER_TOKEN_SECRET"`
}

func bindEnv() {
	viper.BindEnv("PORT")
	viper.BindEnv("APPLE_CHARITY_URL")
	viper.BindEnv("MONGODB_URI")
	viper.BindEnv("MONGODB_USERNAME")
	viper.BindEnv("MONGODB_PASSWORD")
	viper.BindEnv("DATABASE")
	viper.BindEnv("TWITTER_CONSUMER_KEY")
	viper.BindEnv("TWITTER_CONSUMER_SECRET")
	viper.BindEnv("TWITTER_TOKEN")
	viper.BindEnv("TWITTER_TOKEN_SECRET")
}

func LoadConfig(path string) (Config, error) {
	config := &Config{
		Port:            "8080",
		AppleCharityUrl: "https://tw.feature.appledaily.com/charity/projlist",
		Database:        "the_donor",
	}

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return *config, err
		}
	}

	/**
	 * viper can't automatically unmarshal env variables to config
	 * issue: https://github.com/spf13/viper/issues/761
	 */
	bindEnv()

	err := viper.Unmarshal(config)

	log.Info("loaded config")

	return *config, err
}
