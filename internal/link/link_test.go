package link

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/link/handlers"
	"github.com/laytan/shortnr/internal/storage"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"
)

func TestCreateInserts(t *testing.T) {
	// URL storage handlers will use
	storage := storage.MemoryStorage{InternalMap: make(map[string]string)}

	// Set up router
	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)
	r.HandleFunc("/api/v1/shortn", handlers.Create(storage))

	// Set up server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Perform request to shorten url
	url := ts.URL + "/api/v1/shortn"
	_, err := http.Post(url, "application/json", strings.NewReader(`{"url": "https://google.com"}`))
	if err != nil {
		t.Fatal(err)
	}

	// Value should be https://google.com
	if !storage.Contains("https://google.com") {
		t.Errorf("URL did not get inserted")
	}
}
