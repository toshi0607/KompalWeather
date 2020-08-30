package kompal

import (
	"context"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external test")
	}

	url := os.Getenv("KOMPAL_URL")
	k := New(&Config{URL: url})
	s, err := k.Fetch(context.Background())
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	if s.Timestamp.IsZero() {
		t.Fatal("Timestamp should not be zero value")
	}
}
