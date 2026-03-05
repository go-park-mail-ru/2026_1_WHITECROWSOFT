package main

import (
	"net/http"
	"os"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/logger"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/router"
)

func main() {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = ":8000"
	}

	log := logger.Init()

	srv := &http.Server{
		Handler: router.New(),
		Addr:    host,
	}

	log.Info("Server started", "host", host)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
