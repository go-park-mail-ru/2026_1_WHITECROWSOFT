package router

import (
	"net/http"
	"os"
	"time"

	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/auth"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/handlers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/middleware"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/mock"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/storage"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/internal/types"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/helpers"
	"github.com/go-park-mail-ru/2026_1_WHITECROWSOFT/pkg/jwt"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func TestProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(types.UserIDKey).(string)
	if !ok {
		helpers.JSONErrorResponse(w, http.StatusInternalServerError, jwt.ErrNoUserID)
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
	mockData := mock.NewMockData()
	noteHandler := handlers.NewNoteHandler(mockData)

	r := http.NewServeMux()

	r.HandleFunc("GET /ping", pingHandler)

	r.HandleFunc("POST /signup", authHandler.SignupUser)
	r.HandleFunc("POST /signin", authHandler.SigninUser)
	r.HandleFunc("POST /logout", authHandler.LogOutUser)

	r.Handle("GET /protected", middleware.Auth(http.HandlerFunc(TestProtectedEndpoint), authHandler))

	r.Handle("GET /notes", middleware.Auth(http.HandlerFunc(noteHandler.GetAllNotes), authHandler))
	r.Handle("GET /notes/{id}", middleware.Auth(http.HandlerFunc(noteHandler.GetNote), authHandler))
	r.Handle("GET /notes/{id}/blocks", middleware.Auth(http.HandlerFunc(noteHandler.GetNoteBlocks), authHandler))

	return middleware.Logger(r)
}
