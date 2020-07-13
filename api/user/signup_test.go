package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/laytan/shortnr/internal/user/storage"
)

type messageResponse struct {
	Message string
}

func TestSignup(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Perform request
	url := ts.URL + "/signup"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345679"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body messageResponse

	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Signup status code wrong. Got: %d, Want: %d", res.StatusCode, 200)
	}

	if body.Message != "Signed up!" {
		t.Fatalf("Signup response message not as expected. Got: %s, Want: %s", body.Message, "Signed up!")
	}

	_, exists := store.GetByEmail("test@test.com")
	if !exists {
		t.Fatal("Signup did not save user to store")
	}
}

func TestSignupRequiresValidEmail(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Perform request
	url := ts.URL + "/signup"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "testtest", "password": "12345679"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body messageResponse
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 400 {
		t.Fatalf("Expected status code: %d, got: %d", 400, res.StatusCode)
	}

	if body.Message != "Validation failed" {
		t.Fatalf("Expected message: %q, got: %q", "Validation failed", body.Message)
	}

	if _, exists := store.GetByEmail("testtest"); exists {
		t.Fatalf("Email invalid but still inserted!")
	}
}

func TestSignupRequiresLongPassword(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Perform request
	url := ts.URL + "/signup"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "1"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body messageResponse
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 400 {
		t.Fatalf("Expected status code: %d, got: %d", 400, res.StatusCode)
	}

	if body.Message != "Validation failed" {
		t.Fatalf("Expected message: %q, got: %q", "Validation failed", body.Message)
	}

	if _, exists := store.GetByEmail("test@test.com"); exists {
		t.Fatalf("Password invalid but still inserted!")
	}
}

func TestSignupChecksEmailAlreadyUsed(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Add existing user
	store.Users = append(store.Users, storage.User{Email: "test@test.com"})

	// Perform request
	url := ts.URL + "/signup"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345679"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body messageResponse
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 409 {
		t.Fatalf("Expected status code: %d, got: %d", 400, res.StatusCode)
	}

	if body.Message != "Email in use" {
		t.Fatalf("Expected message: %q, got: %q", "Email in use", body.Message)
	}
}

func TestSignupPasswordGetsHashed(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Perform request
	url := ts.URL + "/signup"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345679"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body messageResponse
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected status code: %d, got: %d", 200, res.StatusCode)
	}

	if body.Message != "Signed up!" {
		t.Fatalf("Expected message: %q, got: %q", "Signed up!", body.Message)
	}

	storedUser, exists := store.GetByEmail("test@test.com")
	if !exists {
		t.Fatalf("User did not get inserted into store")
	}

	compareErr := bcrypt.CompareHashAndPassword([]byte(storedUser.Hash), []byte("12345679"))
	if compareErr != nil {
		t.Fatal("Hashed password does not compare to actual password")
	}
}
