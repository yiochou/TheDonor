package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"
)

const TIME_STRING_LAYOUT = "2006-01-02T15:04:05Z"

func FetchHTML(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FetchHTMLByCurl(url string) (string, error) {
	cmd := exec.Command("curl", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return string(out), nil
}

func ParseTimeString(ts string) time.Time {
	t, err := time.Parse(TIME_STRING_LAYOUT, ts)
	if err != nil {
		return time.Now()
	}

	return t
}
