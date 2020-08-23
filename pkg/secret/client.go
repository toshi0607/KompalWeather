package secret

import "context"

type secret struct {
	// secretManager client
}

func (s secret) Get(ctx context.Context, secretName string) (string, error) {
	panic("implement me")
}
