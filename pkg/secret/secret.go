package secret

import (
	"context"
	"errors"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Secret struct {
	client       *secretmanager.Client
	gcpProjectId string
}

func New() (*Secret, error) {
	ctx := context.TODO()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Secret{client: client}, nil
}

func (s *Secret) AddGCPProjectId(id string) {
	s.gcpProjectId = id
}

func (s *Secret) Get(ctx context.Context, name string) (string, error) {
	if s.gcpProjectId == "" {
		return "", errors.New("gcpProjectId is required")
	}

	req := &secretmanagerpb.GetSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", s.gcpProjectId, name),
	}
	resp, err := s.client.GetSecret(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.GetName(), nil
}
