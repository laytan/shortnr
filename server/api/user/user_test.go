package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"
	"golang.org/x/crypto/bcrypt"
)

// SetupUserServer returns a server and a store with the user endpoints registered, make sure to call Close on returned server
func SetupUserServer() (*httptest.Server, *MemoryStorage) {
	// Load .env
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		panic(envErr)
	}

	// Set up store
	store := MemoryStorage{Users: make([]User, 0)}

	// Set up handler
	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)
	SetRoutes(r, &store)

	// Set up server
	return httptest.NewServer(r), &store
}

type Response struct {
	Err Res
	Res Res
}

type Res struct {
	Data map[string]interface{}
	Msg  string
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

	var body Response

	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Signup status code wrong. Got: %d, Want: %d", res.StatusCode, 200)
	}

	if body.Res.Msg != "signed up" {
		t.Fatalf("Signup response message not as expected. Got: %s, Want: %s", body.Res.Msg, "signed up")
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

	var body Response
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 400 {
		t.Fatalf("Expected status code: %d, got: %d", 400, res.StatusCode)
	}

	if body.Err.Msg != "validation failed" {
		t.Fatalf("Expected message: %q, got: %q", "validation failed", body.Err.Msg)
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

	var body Response
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 400 {
		t.Fatalf("Expected status code: %d, got: %d", 400, res.StatusCode)
	}

	if body.Err.Msg != "validation failed" {
		t.Fatalf("Expected message: %q, got: %q", "validation failed", body.Err.Msg)
	}

	if _, exists := store.GetByEmail("test@test.com"); exists {
		t.Fatalf("Password invalid but still inserted!")
	}
}

func TestSignupChecksEmailAlreadyUsed(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Add existing user
	store.Users = append(store.Users, User{Email: "test@test.com"})

	// Perform request
	url := ts.URL + "/signup"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345679"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body Response
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 409 {
		t.Fatalf("Expected status code: %d, got: %d", 400, res.StatusCode)
	}

	if body.Err.Msg != "email in use" {
		t.Fatalf("Expected message: %q, got: %q", "email in use", body.Err.Msg)
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

	var body Response
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected status code: %d, got: %d", 200, res.StatusCode)
	}

	if body.Res.Msg != "signed up" {
		t.Fatalf("Expected message: %q, got: %q", "signed up", body.Res.Msg)
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

func TestJwtMiddlewareRequiresValidToken(t *testing.T) {
	ts, _ := SetupUserServer()
	defer ts.Close()

	res, err := http.Get(ts.URL + "/me")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 401 {
		t.Fatal("No 401 returned on authorized route without authorization")
	}
}

func TestMeHandlerReturnsLoggedInUser(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	user := User{
		ID:        1,
		Email:     "test@test.com",
		CreatedAt: time.Time{},
	}
	store.Users = append(store.Users, user)

	jwtToken, _, err := SignUserToken(user)
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

	var body Response
	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	// Turn map into json encoded byte string
	userBytes, marshalErr := json.Marshal(body.Res.Data)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}

	// Turn json into User
	var returnedUser User
	unmarshalErr := json.Unmarshal(userBytes, &returnedUser)
	if unmarshalErr != nil {
		t.Fatal(unmarshalErr)
	}

	if returnedUser.Email != "test@test.com" {
		t.Errorf("User returned was not the expected user. Got: %v, Want: %v", returnedUser.Email, user)
	}
}

func TestLoginReturnsToken(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Insert a user to login to
	hash, hashErr := bcrypt.GenerateFromPassword([]byte("12345678"), 1)
	if hashErr != nil {
		t.Fatal(hashErr)
	}
	userToLogin := User{ID: 1, Email: "test@test.com", Hash: string(hash)}
	store.Users = append(store.Users, userToLogin)

	// Perform request
	url := ts.URL + "/login"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345678"}`))
	if err != nil {
		t.Fatal(err)
	}

	var body Response

	decodeErr := json.NewDecoder(res.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	fmt.Println(body)

	if body.Res.Msg != "logged in" {
		t.Fatalf("Expected message: %q, got message: %q", "logged in", body.Res.Msg)
	}

	if len(body.Res.Data["token"].(string)) < 2 {
		t.Fatalf("Expected token string to be longer then 2 characters, got: %q", body.Res.Data["token"])
	}
}

func TestLoginSetsRefreshTokenCookie(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Insert a user to login to
	hash, hashErr := bcrypt.GenerateFromPassword([]byte("12345678"), 1)
	if hashErr != nil {
		t.Fatal(hashErr)
	}
	userToLogin := User{ID: 1, Email: "test@test.com", Hash: string(hash)}
	store.Users = append(store.Users, userToLogin)

	// Perform request
	url := ts.URL + "/login"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345678"}`))
	if err != nil {
		t.Fatal(err)
	}

	cookieHeader := res.Header.Get("Set-Cookie")
	if len(cookieHeader) < 2 || !strings.Contains(cookieHeader, "refreshToken=") {
		t.Fatal("refreshToken cookie not present")
	}

	if !strings.Contains(cookieHeader, "HttpOnly") {
		t.Fatal("refreshToken is not HTTP only")
	}

	if !strings.Contains(cookieHeader, "Secure") {
		t.Fatal("refreshToken cookie is not secure")
	}
}

func TestTokenCanBeRefreshed(t *testing.T) {
	ts, store := SetupUserServer()
	defer ts.Close()

	// Insert a user to login to
	hash, hashErr := bcrypt.GenerateFromPassword([]byte("12345678"), 1)
	if hashErr != nil {
		t.Fatal(hashErr)
	}
	userToLogin := User{ID: 1, Email: "test@test.com", Hash: string(hash)}
	store.Users = append(store.Users, userToLogin)

	// Perform request
	url := ts.URL + "/login"
	res, err := http.Post(url, "application/json", strings.NewReader(`{"email": "test@test.com", "password": "12345678"}`))
	if err != nil {
		t.Fatal(err)
	}

	cookieHeader := res.Header.Get("Set-Cookie")
	header := http.Header{}
	header.Add("Cookie", cookieHeader)
	request := http.Request{Header: header}
	refreshCookie, noRefreshCookieErr := request.Cookie("refreshToken")
	if noRefreshCookieErr != nil {
		t.Fatal(noRefreshCookieErr)
	}

	refreshURL := ts.URL + "/refresh"
	req, reqErr := http.NewRequest("POST", refreshURL, strings.NewReader(""))
	if reqErr != nil {
		t.Fatal(reqErr)
	}
	req.AddCookie(refreshCookie)

	refreshRes, refreshErr := ts.Client().Do(req)
	if refreshErr != nil {
		t.Fatal(refreshErr)
	}

	var body Response

	decodeErr := json.NewDecoder(refreshRes.Body).Decode(&body)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	if len(body.Res.Data["token"].(string)) < 2 {
		t.Fatal("No token came back from refresh request")
	}

	if body.Res.Msg != "token refreshed" {
		t.Fatal("Token refresh message not as expected")
	}

	cookies := refreshRes.Header.Get("Set-Cookie")
	if !strings.Contains(cookies, "refreshToken=") {
		t.Fatal("Refresh did not re set refreshToken")
	}
}

func TestLogoutRemovesRefreshToken(t *testing.T) {
	ts, _ := SetupUserServer()
	defer ts.Close()

	// Perform request
	url := ts.URL + "/logout"
	req, err := http.NewRequest("DELETE", url, strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	cookieHeader := res.Header.Get("Set-Cookie")
	if !strings.Contains(cookieHeader, "refreshToken=;") {
		t.Fatal("refreshToken not set to empty")
	}
}
