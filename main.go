package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	db, err := ConnectMongoDB(config)
	if err != nil {
		log.Fatal("cannot connect to MongoDB:", err)
	}

	store := NewStore(db, *log.New())

	TweetNewDonorCases(store, config)
}

func TweetNewDonorCases(store Store, config Config) error {
	parser := NewParser(*log.New())
	cases, err := parser.Parse(config.AppleCharityUrl)
	if err != nil {
		return err
	}

	newCases, err := store.InsertCasesIfNotExists(cases)
	if err != nil {
		return err
	}
	if newCases == nil {
		log.Info("no new cases")
		return nil
	}

	twitter := NewTwitter(config)
	twitter.TweetCases(cases)

	return nil
}
