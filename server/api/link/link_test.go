package link

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/laytan/shortnr/pkg/jsonmiddleware"
)

// Loads .env, sets up memory storage, sets up a server with link handlers
func setupServer() (*httptest.Server, *MemoryStorage) {
	envErr := godotenv.Load("../../.env")
	if envErr != nil {
		panic(envErr)
	}

	store := MemoryStorage{Links: make([]Link, 0)}

	r := mux.NewRouter()
	r.Use(jsonmiddleware.Middleware)
	SetRoutes(r, &store)

	return httptest.NewServer(r), &store
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
	ts, store := setupServer()
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
