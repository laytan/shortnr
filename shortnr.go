package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Hello, world!"}`))
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("", HelloHandler).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
