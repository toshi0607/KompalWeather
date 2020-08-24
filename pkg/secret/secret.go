package secret

import (
	"context"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Secret struct {
	client *secretmanager.Client
}

func New() (*Secret, error) {
	ctx := context.TODO()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Secret{client: client}, nil
}

func (s *Secret) Get(ctx context.Context, name string) (string, error) {
	req := &secretmanagerpb.GetSecretRequest{
		Name: name,
	}
	resp, err := s.client.GetSecret(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.GetName(), nil
}
