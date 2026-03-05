package main

import (
	"net/http"
	"os"

	authHandlers "github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/auth"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/logger"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func main() {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = ":8000"
	}

	log := logger.Init()

	authHandler := authHandlers.NewAuthHandler(os.Getenv("JWT_SECRET"), authHandlers.NewUserSet())

	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)

	r.HandleFunc("POST /signup", authHandler.SignupUser)
	r.HandleFunc("POST /signin", authHandler.SigninUser)
	r.HandleFunc("POST /logout", authHandler.LogOutUser)

	r.Handle("GET /protected", authHandler.AuthMiddleware(http.HandlerFunc(authHandler.TestProtectedEndpoint)))

	handler := logger.Middleware(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    host,
	}

	log.Info("Server started", "host", host)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
