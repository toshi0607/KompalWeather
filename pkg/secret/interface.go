package secret

import "context"

type Secret interface {
	Get(ctx context.Context, secretName string) (string, error)
}
