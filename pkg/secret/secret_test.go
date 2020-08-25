package secret

import (
	"context"
	"os"
	"testing"
)

func TestSecret_Get(t *testing.T) {
	ctx := context.TODO()
	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	s, err := New()
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	s.AddGCPProjectId(gcpProjectID)

	if _, err := s.Get(ctx, "test"); err != nil {
		t.Fatalf("error: %s", err)
	}
}
