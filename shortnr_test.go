package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/shortenerstorage"
)

func TestOnePlusOne(t *testing.T) {
	answer := 1 + 1
	if answer != 2 {
		t.Error("One plus one is 2")
	}
}

func TestHelloWorldReturned(t *testing.T) {
	// Create request to test
	req, err := http.NewRequest("GET", "/api/v1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HelloHandler)

	// Perform request
	handler.ServeHTTP(rr, req)

	// Check status code
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"message": "Hello, world!"}`
	body := rr.Body.String()
	if body != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", body, expected)
	}
}

func TestContentTypeJsonGetsSet(t *testing.T) {
	// Set up router
	r := mux.NewRouter()
	// Middleware we test
	r.Use(JsonMiddleware)
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

func TestRequestGetsRedirected(t *testing.T) {
	URLS := shortenerstorage.MapStorage{InternalMap: map[string]string{"12345": "https://www.google.com/"}}

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/{id}", RedirectHandler(URLS))

	// Set up server
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Perform request
	url := ts.URL + "/12345"
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	// Check that we got redirected to google
	if resUrl := res.Request.URL.String(); resUrl != "https://www.google.com/" {
		t.Errorf("Not redirected or redirected to the wrong url got redirected to: %q", resUrl)
	}
}

func TestShortenedUrlGetsInserted(t *testing.T) {
	// URL storage handlers will use
	URLS := shortenerstorage.MapStorage{InternalMap: make(map[string]string)}

	// Set up router
	r := mux.NewRouter()
	r.Use(JsonMiddleware)
	r.HandleFunc("/api/v1/shortn", ShortenerHandler(URLS))

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
	if !URLS.Contains("https://google.com") {
		t.Errorf("URL did not get inserted")
	}
}
