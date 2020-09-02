package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/monitoring"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/storage"
)

// Server represents a server
type Server struct {
	kompal    kompal.Fetcher
	storage   storage.Storage
	notifiers []notifier.Notifier
	analyzer  analyzer.Analyzer
	monitor   *monitoring.Monitor
	log       logger.Logger

	mux    *http.ServeMux
	server *http.Server
}

// New build a new Server
func New(f kompal.Fetcher, s storage.Storage, ns []notifier.Notifier, a analyzer.Analyzer, m *monitoring.Monitor, l logger.Logger) *Server {
	server := &Server{
		kompal:    f,
		storage:   s,
		notifiers: ns,
		analyzer:  a,
		monitor:   m,
		log:       l,

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
