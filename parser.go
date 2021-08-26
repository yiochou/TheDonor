package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func NewParser() Parser {
	return &AppleCharityParser{}
}

type AppleCharityParser struct{}

func (p *AppleCharityParser) Parse(url string) ([]*Case, error) {
	html, err := FetchHTML(url)
	if err != nil {
		return nil, errors.Wrap(err, "fetch failed")
	}

	const caseWorkerNum = 3
	cases := make(chan *Case, caseWorkerNum)
	caseLinks := make(chan string, caseWorkerNum)

	for i := 0; i < caseWorkerNum; i++ {
		go p.parseCaseWorker(caseLinks, cases)
	}

	links, err := p.ParseCaseLinks(html)
	if err != nil {
		return nil, errors.Wrap(err, "parse links failed")
	}
	if links == nil {
		return nil, nil
	}

	go func() {
		defer close(caseLinks)
		for _, link := range links {
			caseLinks <- link
		}
	}()

	result := []*Case{}
	for i := 0; i < len(links); i++ {
		result = append(result, <-cases)
	}
	close(cases)

	return result, nil
}

func (p *AppleCharityParser) ParseCaseLinks(html string) (links []string, err error) {
	const ActiveStatus = "未結案"

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		return nil, err
	}

	doc.Find("table tr").FilterFunction(func(i int, s *goquery.Selection) bool {
		status := s.Find(".ucb").Text()

		return status == ActiveStatus
	}).Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find(".artcatdetails").Attr("href")
		if !exists {
			return
		}

		links = append(links, link)
	})

	return links, nil
}

func (p *AppleCharityParser) parseCaseWorker(caseLinks <-chan string, cases chan<- *Case) {
	log.Info("start worker")
	defer log.Info("worker closed")

	for link := range caseLinks {
		html, err := FetchHTMLByCurl(link)
		if err != nil {
			continue
		}
		log.Info("parse case, link: ", link)

		c := p.ParseCase(html)
		cases <- c
	}
}

func (p *AppleCharityParser) ParseCase(html string) *Case {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	if err != nil {
		log.Error("parse case: ", err)
		return nil
	}

	c := &Case{}
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, exists := s.Attr("property")

		if !exists {
			return
		}
		content, exists := s.Attr("content")

		if !exists {
			return
		}

		switch property {
		case "dable:item_id":
			c.Id = "AppleCharity:" + content
		case "og:title":
			c.Title = content
		case "og:image":
			c.Image = content
		case "og:url":
			c.Link = content
		case "og:description":
			c.Description = content
		case "article:published_time":
			c.PublishedAt = ParseTimeString(content)
		default:
			return
		}
	})

	return c
}
