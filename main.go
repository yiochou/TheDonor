package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	appleCharityUrl := viper.GetString("APPLE_CHARITY_URL")

	parser := NewParser()
	cases, err := parser.Parse(appleCharityUrl)

	if err != nil {
		log.Fatal(err)
	}

	newCases, err := InsertCasesIfNotExists(cases)

	if err != nil {
		log.Fatal(err)
	}
	if newCases == nil {
		log.Info("no new cases")
		return
	}

	TweetCases(cases)
}
