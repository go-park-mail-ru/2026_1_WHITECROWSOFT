package main

import (
	"log/slog"
	"net/http"
	"os"
	authHandlers "wcs/internal/auth"

	"github.com/joho/godotenv"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong")) // прописывать 200 статус не нужно, т.к. w.Write() устанавливает его сам если он ещё не установлен
}

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, using system env")
	}

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "127.0.0.1:8000"
	}

	authHandler := &authHandlers.AuthHandler{
		JWTSecret: os.Getenv("JWT_SECRET"),
		UserSet:   authHandlers.NewUserSet(),
	}

	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)
	r.HandleFunc("POST /signup", authHandler.SignupUser)
	r.HandleFunc("POST /signin", authHandler.SigninUser)
	r.HandleFunc("POST /logout", authHandler.LogOutUser)

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
