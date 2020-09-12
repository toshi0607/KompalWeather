package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/logger"
	"github.com/toshi0607/kompal-weather/pkg/visualizer"
)

type VisualizeServer struct {
	visualizer *visualizer.Visualizer

	log    logger.Logger
	mux    *http.ServeMux
	server *http.Server
}

// NewVisualizer build a new VisualizeServer
func NewVisualizer(v *visualizer.Visualizer, l logger.Logger) Server {
	server := &VisualizeServer{
		visualizer: v,

		log: l,
		mux: http.NewServeMux(),
	}

	server.registerHandlers()

	return server
}

// Serve serves the server
func (s *VisualizeServer) Serve(ln net.Listener) error {
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
func (s *VisualizeServer) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *VisualizeServer) registerHandlers() {
	s.mux.Handle("/visualize", s.visualizeHandler())
}
