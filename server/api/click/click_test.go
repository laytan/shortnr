package click

import (
	"testing"
	"time"
)

func TestClickCanBeInserted(t *testing.T) {
	store := MemoryStorage{
		Clicks: make([]Click, 0),
	}

	created, err := Create("2", &store)
	if err != nil {
		t.Fatal(err)
	}

	if created.LinkID != "2" {
		t.Fatal("Created link does not have assigned id")
	}

	clicks := GetAll("2", &store)
	if len(clicks) != 1 {
		t.Fatal("Expected length of 1 clicks")
	}

	if clicks[0].ID != created.ID {
		t.Fatal("ID's don't match")
	}
}

func TestGetAllClicks(t *testing.T) {
	store := MemoryStorage{
		Clicks: []Click{
			{
				ID:        1,
				LinkID:    "1",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			{
				ID:        3,
				LinkID:    "1",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			{
				ID:        2,
				LinkID:    "3",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
		},
	}

	clicks := GetAll("1", &store)
	if len(clicks) != 2 {
		t.Fatal("Expected to return 2 clicks")
	}
}
