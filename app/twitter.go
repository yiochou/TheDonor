package main

import (
	"bytes"
	"sort"
	"text/template"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
)

type Twitter struct {
	client *twitter.Client
}

func NewTwitter(config Config) Twitter {
	OAuthConfig := oauth1.NewConfig(config.TwitterConsumerKey, config.TwitterConsumerSecret)
	OAuthToken := oauth1.NewToken(config.TwitterToken, config.TwitterTokenSecret)
	httpClient := OAuthConfig.Client(oauth1.NoContext, OAuthToken)

	t := Twitter{
		client: twitter.NewClient(httpClient),
	}

	return t
}

func (twitter *Twitter) TweetCases(cases []*Case) {
	sort.Slice(cases, func(i, j int) bool {
		return cases[i].PublishedAt.Before(cases[j].PublishedAt)
	})
	for _, c := range cases {
		tweet, err := twitter.caseToTweet(*c)

		if err != nil {
			log.Error(err)
			continue
		}
		_, _, err = twitter.client.Statuses.Update(tweet, nil)
		if err != nil {
			log.Error(err)
		}

		log.Info("tweeted: ", c.Title)
	}
}

func (twitter *Twitter) caseToTweet(c Case) (string, error) {
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
