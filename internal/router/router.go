package router

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/auth"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/handlers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/config"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/middleware"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/mock"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func New(cfg *config.Config) http.Handler {
	authHandler := auth.NewHandler(cfg.JWT, storage.NewUserSet())
	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)

	r.HandleFunc("POST /signup", authHandler.SignupUser)
	r.HandleFunc("POST /signin", authHandler.SigninUser)
	r.HandleFunc("POST /logout", authHandler.LogOutUser)

	r.Handle("GET /notes", middleware.Auth(http.HandlerFunc(noteHandler.GetAllNotes), authHandler))
	r.Handle("GET /notes/{id}", middleware.Auth(http.HandlerFunc(noteHandler.GetNote), authHandler))
	r.Handle("GET /notes/{id}/blocks", middleware.Auth(http.HandlerFunc(noteHandler.GetNoteBlocks), authHandler))
	return middleware.Logger(r)
}
