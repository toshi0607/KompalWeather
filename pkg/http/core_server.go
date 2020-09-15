package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/controller"
	"github.com/toshi0607/kompal-weather/pkg/logger"
)

// CoreServer represents a server
type CoreServer struct {
	controller controller.Controller

	log    logger.Logger
	mux    *http.ServeMux
	server *http.Server
}

// NewCore build a new CoreServer
func NewCore(c controller.Controller, l logger.Logger) Server {
	server := &CoreServer{
		controller: c,

		log: l,
		mux: http.NewServeMux(),
	}

	server.registerHandlers()

	return server
}

// Serve serves the server
func (s *CoreServer) Serve(ln net.Listener) error {
	server := &http.Server{
		Handler: s.mux,
	}
	s.server = server

	if err := server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// GracefulStop stops the server gracefully
func (s *CoreServer) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *CoreServer) registerHandlers() {
	s.mux.Handle("/watch", s.watchHandler())
	s.mux.Handle("/trend", s.trendHandler())
}
