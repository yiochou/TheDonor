package main

import (
	"bytes"
	"sort"
	"text/template"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var twitterClient *twitter.Client

func init() {
	config := oauth1.NewConfig(viper.GetString("TWITTER_CONSUMER_KEY"), viper.GetString("TWITTER_CONSUMER_SECRET"))
	token := oauth1.NewToken(viper.GetString("TWITTER_TOKEN"), viper.GetString("TWITTER_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	twitterClient = twitter.NewClient(httpClient)
}

func TweetCases(cases []*Case) {
	sort.Slice(cases[:], func(i, j int) bool {
		return cases[i].PublishedAt.Before(cases[j].PublishedAt)
	})
	for _, c := range cases {
		tweet, err := caseToTweet(*c)

		if err != nil {
			log.Error(err)
			continue
		}
		_, _, err = twitterClient.Statuses.Update(tweet, nil)
		if err != nil {
			log.Error(err)
		}

		log.Info("tweeted: ", c.Title)
	}
}

func caseToTweet(c Case) (string, error) {
	tweetTemplate := `
{{ .Title}}

{{ .Link}}
	`
	tweet, err := template.New("tweet").Parse(tweetTemplate)
	if err != nil {
		return "", nil
	}

	var tweetBuffer bytes.Buffer
	if err := tweet.Execute(&tweetBuffer, c); err != nil {
		return "", err
	}

	return tweetBuffer.String(), nil
}
