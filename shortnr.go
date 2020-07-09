package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"message": "Hello, world!"}`))
}

// Add json response header to request
func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(JsonMiddleware)
	api.HandleFunc("", HelloHandler).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}
