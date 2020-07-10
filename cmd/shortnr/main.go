package main

import (
	"log"
	"net/http"

	"github.com/laytan/shortnr/pkg/jsonmiddleware"

	"github.com/laytan/shortnr/internal/link"
	"github.com/laytan/shortnr/internal/redirect"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/laytan/shortnr/internal/storage"
)

// init runs before main
func init() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()

	// Store all links in here
	storage := storage.NewMysqlStorage()

	// General routes
	redirect.SetRoutes(r, storage)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(jsonmiddleware.Middleware)

	// API routes
	link.SetRoutes(r, storage)

	// Start up server
	log.Fatal(http.ListenAndServe(":8080", r))
}
