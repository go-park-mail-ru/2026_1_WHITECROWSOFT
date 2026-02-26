package main

import (
	"log/slog"
	"net/http"
	"os"

	"wcs/internal/auth"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong")) // прописывать 200 статус не нужно, т.к. w.Write() устанавливает его сам если он ещё не установлен
}

func main() {
	host := "127.0.0.1:8000"

	authHandler := &authHandlers.AuthHandler{
		JWTSecret: "haha-secret-key-open", 
		UserSet:   authHandlers.NewUserSet(),
	}

	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)
	r.HandleFunc("POST /signup", authHandler.SignupUser)

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
