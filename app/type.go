package main

import "time"

type Case struct {
	Id          string
	Title       string
	Description string
	Link        string
	Image       string
	PublishedAt time.Time
}

type Parser interface {
	Parse(html string) ([]*Case, error)
}
