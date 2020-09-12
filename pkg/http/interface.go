package http

import (
	"context"
	"net"
)

// Server is an interface of a server
type Server interface {
	Serve(ln net.Listener) error
	GracefulStop(ctx context.Context) error
}
