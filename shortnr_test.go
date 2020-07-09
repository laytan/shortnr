package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
