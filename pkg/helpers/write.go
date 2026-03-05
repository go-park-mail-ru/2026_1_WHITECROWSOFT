package helpers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func JSONErrorResponse(w http.ResponseWriter, status int, err error) {
	JSONResponse(w, http.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}
