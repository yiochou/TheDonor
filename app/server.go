package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func StartServer(config Config, job func() error) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := job()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Printf("listening on port %s", config.Port)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		return err
	}

	return nil
}
