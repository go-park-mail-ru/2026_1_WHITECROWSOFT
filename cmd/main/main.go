package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ping", pingHandler).Methods("GET")
	
	srv := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:8000",
	}

	log.Println("Сервер запущен на http://127.0.0.1:8000")
	log.Fatal(srv.ListenAndServe())
}
