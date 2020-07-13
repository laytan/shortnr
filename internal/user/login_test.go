package user

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/laytan/shortnr/internal/user/storage"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Message string
	Token   string
}

func TestLoginReturnsToken(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Insert a user to login to
	hash, hashErr := bcrypt.GenerateFromPassword([]byte("12345678"), 1)
	if hashErr != nil {
		t.Fatal(hashErr)
	}
	userToLogin := storage.User{ID: 1, Email: "test@test.com", Hash: string(hash)}
	store.Users = append(store.Users, userToLogin)

	// Perform request
	url := ts.URL + "/login"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345678"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body LoginResponse

	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if body.Message != "Logged in!" {
		t.Fatalf("Expected message: %q, got message: %q", "Logged in!", body.Message)
	}

	if len(body.Token) < 2 {
		t.Fatalf("Expected token string to be longer then 2 characters, got: %q", body.Token)
	}
}
