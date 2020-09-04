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

// Server represents a server
type Server struct {
	controller controller.Controller

	log    logger.Logger
	mux    *http.ServeMux
	server *http.Server
}

// New build a new Server
func New(c controller.Controller, l logger.Logger) *Server {
	server := &Server{
		controller: c,

		log: l,
		mux: http.NewServeMux(),
	}

	server.registerHandlers()

	return server
}

// Serve serves the server
func (s *Server) Serve(ln net.Listener) error {
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
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerHandlers() {
	s.mux.Handle("/watch", s.watchHandler())
}
