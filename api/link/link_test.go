package link

import (
	"testing"

	"github.com/laytan/shortnr/api/link/storage"
)

func TestCreateInserts(t *testing.T) {
	// URL storage handlers will use
	storage := storage.Memory{InternalMap: make(map[string]string)}

	err := Create(Link{URL: "https://google.com"}, storage)
	if err != nil {
		t.Fatal(err)
	}

	// Value should be https://google.com
	if !storage.Contains("https://google.com") {
		t.Errorf("URL did not get inserted")
	}
}
