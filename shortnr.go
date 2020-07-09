package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

// Send hello world
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

type PostShortenerBody struct {
	URL string
}

// Create new Id to Url mapping
func ShortenerHandler(URLS map[string]string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body PostShortenerBody

		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid request, make sure to include a url"}`))
			return
		}

		if !IsUrl(body.URL) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "Invalid url, make sure to include https:// or http://"}`))
			return
		}

		// Generate a unique ID
		id := xid.New().String()

		// Store in our urls map
		URLS[id] = body.URL

		w.Write([]byte(fmt.Sprintf(`{"message": "URL created", "id": %q, "originalURL": %q}`, id, body.URL)))
	})
}

// Redirects request to original url
func RedirectHandler(URLS map[string]string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		// Check for the id in our urls map and redirect if present
		if redirectURL, ok := URLS[id]; ok {
			http.Redirect(w, r, redirectURL, 308)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func main() {
	r := mux.NewRouter()

	// Store all redirections in here
	URLS := make(map[string]string)
	r.HandleFunc("/{id}", RedirectHandler(URLS)).Methods(http.MethodGet)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(JsonMiddleware)
	api.HandleFunc("", HelloHandler).Methods(http.MethodGet)
	api.HandleFunc("/shortn", ShortenerHandler(URLS)).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Check if URL is ok
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
