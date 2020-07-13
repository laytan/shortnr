package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/laytan/shortnr/internal/user/handlers"
	"github.com/laytan/shortnr/internal/user/storage"
)

func TestMeHandlerReturnsLoggedInUser(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	user := storage.User{
		ID:        1,
		Email:     "test@test.com",
		CreatedAt: time.Time{},
	}
	store.Users = append(store.Users, user)

	jwtToken, err := handlers.SignUserToken(user)
	if err != nil {
		t.Fatal(err)
	}

	// Perform request
	url := ts.URL + "/me"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	res, reqErr := ts.Client().Do(req)
	if reqErr != nil {
		t.Fatal(reqErr)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected status code: %d, got: %d", 200, res.StatusCode)
	}

	var returnedUser storage.User
	decodeErr := json.NewDecoder(res.Body).Decode(&returnedUser)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if returnedUser.Email != "test@test.com" {
		t.Errorf("User returned was not the expected user. Got: %v, Want: %v", returnedUser, user)
	}
}
