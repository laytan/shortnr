package link

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/laytan/shortnr/api/click"
	"github.com/laytan/shortnr/api/user"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"
)

// Loads .env, sets up memory storage, sets up a server with link handlers
func setupServer() (*httptest.Server, *MemoryStorage, *click.MemoryStorage) {
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		panic(envErr)
	}

	store := MemoryStorage{Links: make([]Link, 0)}
	clickStore := click.MemoryStorage{Clicks: make([]click.Click, 0)}

	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)
	SetRoutes(r, &store, &clickStore)

	return httptest.NewServer(r), &store, &clickStore
}

func decodeBody(body io.ReadCloser) map[string]interface{} {
	var decoded map[string]interface{}
	decodeErr := json.NewDecoder(body).Decode(&decoded)
	if decodeErr != nil {
		panic(decodeErr)
	}
	return decoded
}

func getLinkFromBody(body io.ReadCloser) Link {
	decoded := decodeBody(body)
	bodyBytes, err := json.Marshal(decoded["res"].(map[string]interface{})["data"])
	if err != nil {
		panic(err)
	}

	var link Link
	err = json.Unmarshal(bodyBytes, &link)
	if err != nil {
		panic(err)
	}

	return link
}

func TestCreateInserts(t *testing.T) {
	// URL storage handlers will use
	storage := &MemoryStorage{Links: make([]Link, 0)}

	_, err := Create(Link{URL: "https://google.com", UserID: 1}, storage)
	if err != nil {
		t.Fatal(err)
	}

	// Value should be https://google.com
	if !storage.Contains("https://google.com") {
		t.Errorf("URL did not get inserted")
	}
}

func TestGet(t *testing.T) {
	ts, store, _ := setupServer()
	defer ts.Close()

	store.Links = append(store.Links, Link{
		ID:  "12345",
		URL: "https://www.google.com/",
	})

	res, err := http.Get(ts.URL + "/12345")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Status code is not 200 but %d", res.StatusCode)
	}

	link := getLinkFromBody(res.Body)

	if link.ID != "12345" {
		t.Fatalf("Expected id '12345', got: %s", link.ID)
	}

	if link.URL != "https://www.google.com/" {
		t.Fatalf("URL returned (%s) is not URL (%s) inserted", link.URL, "https://www.google.com/")
	}
}

func TestLinkCanBeDeleted(t *testing.T) {
	ts, store, _ := setupServer()
	defer ts.Close()

	store.Create(Link{
		ID:     "0",
		UserID: 0,
	})

	authUser := user.User{
		ID: 0,
	}

	token, _, err := user.SignUserToken(authUser)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", ts.URL+"/0", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code: 200, got: %d", res.StatusCode)
	}

	if _, exists := store.Get("0"); exists {
		t.Fatal("Link still exists")
	}
}

func TestLinkCanOnlyBeDeletedByCreator(t *testing.T) {
	user := user.User{
		ID: 0,
	}

	store := MemoryStorage{[]Link{
		{
			ID:     "0",
			UserID: 1,
		},
	}}

	err := Destroy("0", user, &store)
	if err == nil {
		t.Fatal("destroy should error")
	}
}

func TestGetOneCreatesClick(t *testing.T) {
	ts, store, clickStore := setupServer()
	defer ts.Close()

	store.Links = append(store.Links, Link{
		ID:     "0",
		UserID: 1,
	})

	res, err := http.Get(ts.URL + "/0")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatal("Not OK status returned")
	}

	clicks := clickStore.GetAll("0")
	if len(clicks) != 1 {
		t.Fatal("Expected 1 click")
	}
}

func TestLinkCanBeCreatedWithCustomID(t *testing.T) {
	store := MemoryStorage{
		Links: make([]Link, 0),
	}

	link := Link{
		ID:     "custom",
		URL:    "https://www.google.com/",
		UserID: 1,
	}

	link, err := Create(link, &store)
	if err != nil {
		t.Fatal(err)
	}

	if link.ID != "custom" {
		t.Fatal("Link ID was not our custom ID")
	}
}

func TestLinkGetsRandomIDIfNoIDIsGiven(t *testing.T) {
	store := MemoryStorage{
		Links: make([]Link, 0),
	}

	link := Link{
		URL:    "https://www.google.com/",
		UserID: 1,
	}

	link, err := Create(link, &store)
	if err != nil {
		t.Fatal(err)
	}

	if len(link.ID) == 0 {
		t.Fatal("Link ID is not generated")
	}
}

func TestLinksCanNotHaveTheSameID(t *testing.T) {
	store := MemoryStorage{
		Links: []Link{
			{
				ID:     "custom",
				UserID: 2,
			},
		},
	}

	_, err := Create(Link{
		ID:     "custom",
		UserID: 1,
	}, &store)
	if err == nil {
		t.Fatal("Expected create to error because of conflict")
	}
}
