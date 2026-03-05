package router

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/auth"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/middleware"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/helpers"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func TestProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		helpers.JSONErrorResponse(w, http.StatusInternalServerError, errors.New("user_id not found in context"))
		return
	}

	response := map[string]string{
		"msg":     "This is a protected endpoint",
		"user_id": userID,
		"time":    time.Now().Format(time.RFC3339),
	}

	helpers.JSONResponse(w, http.StatusOK, response)
}

func New() http.Handler {
	authHandler := auth.NewHandler(os.Getenv("JWT_SECRET"), storage.NewUserSet())

	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)

	r.HandleFunc("POST /signup", authHandler.SignupUser)
	r.HandleFunc("POST /signin", authHandler.SigninUser)
	r.HandleFunc("POST /logout", authHandler.LogOutUser)

	r.Handle("GET /protected", middleware.Auth(http.HandlerFunc(TestProtectedEndpoint), authHandler))

	return middleware.Logger(r)
}
