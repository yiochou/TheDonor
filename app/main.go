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

	err = StartServer(config, func() error {
		log.Info("handle request")
		return TweetNewDonorCases(store, config)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func TweetNewDonorCases(store Store, config Config) error {
	parser := NewParser(*log.New())
	cases, err := parser.Parse(config.AppleCharityUrl)
	if err != nil {
		return err
	}
	log.Info("cases parsed")

	newCases, err := store.InsertCasesIfNotExists(cases)
	if err != nil {
		return err
	}
	if newCases == nil {
		log.Info("no new cases")
		return nil
	}
	log.Info("new cases inserted")

	twitter := NewTwitter(config)
	twitter.TweetCases(newCases)

	log.Info("cases tweeted")

	return nil
}
