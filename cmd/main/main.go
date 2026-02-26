package main

import (
	"log"
	"log/slog"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func main() {
	host := "127.0.0.1:8000"
	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    host,
	}

	slog.Info("Server started", "host", host)
	log.Fatal(srv.ListenAndServe())
}
