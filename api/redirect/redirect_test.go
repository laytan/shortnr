package redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/laytan/shortnr/internal/storage"
)

func TestRequestGetsRedirected(t *testing.T) {
	storage := storage.MemoryStorage{InternalMap: map[string]string{"12345": "https://www.google.com/"}}

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/{id}", Handler(storage))

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
	if resURL := res.Request.URL.String(); resURL != "https://www.google.com/" {
		t.Errorf("Not redirected or redirected to the wrong url got redirected to: %q", resURL)
	}
}
