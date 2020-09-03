package secret

import (
	"context"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Secret represents a Secret
type Secret struct {
	client       *secretmanager.Client
	gcpProjectID string
}

// New builds new Secret
func New() (*Secret, error) {
	ctx := context.TODO()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Secret{client: client}, nil
}

// Close closes the client connection
func (s *Secret) Close() error {
	return s.client.Close()
}

// AddGCPProjectID adds gcpProjectID to Secret
func (s *Secret) AddGCPProjectID(id string) {
	s.gcpProjectID = id
}

// Get returns the latest secret version data by secret_id
func (s *Secret) Get(ctx context.Context, name string) (string, error) {
	if s.gcpProjectID == "" {
		return "", errors.New("gcpProjectID is required")
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.gcpProjectID, name),
	}
	resp, err := s.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}
	return string(resp.Payload.Data), nil
}
