package jsonmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHeaderGetsCalledWithParams(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	emptyHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	Middleware(emptyHandler).ServeHTTP(rr, r)
	if ctype := rr.Header().Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content type header not set to JSON, got: %q", ctype)
	}
}

func TestContentTypeJsonGetsSetOnRequest(t *testing.T) {
	// Set up router
	r := mux.NewRouter()
	// Middleware we test
	r.Use(Middleware)
	// Empty handler
	r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Set up server
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if ctype := res.Header.Get("Content-Type"); ctype != "application/json" {
		t.Errorf("Content type JSON header not set, got: %q", ctype)
	}
}
