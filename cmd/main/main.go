package main

import (
	"log/slog"
	"net/http"
	"os"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong")) // прописывать 200 статус не нужно, т.к. w.Write() устанавливает его сам если он ещё не установлен
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
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
