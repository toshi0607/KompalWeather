package kompal

import (
	"context"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	url := os.Getenv("KOMPAL_URL")
	k := New(url)
	s, err := k.Fetch(context.Background())
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if s.Timestamp.IsZero() {
		t.Fatal("Male is empty")
	}
}
