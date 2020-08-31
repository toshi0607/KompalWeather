package controller

import "context"

// Controller is an interface of a controller
type Controller interface {
	Run(ctx context.Context) error
}
