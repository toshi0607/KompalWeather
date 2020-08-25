package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/analyzer"
	"github.com/toshi0607/kompal-weather/pkg/kompal"
	"github.com/toshi0607/kompal-weather/pkg/notifier"
	"github.com/toshi0607/kompal-weather/pkg/storage"
)

//https://cloud.google.com/run/docs/triggering/using-scheduler Cloud Schedulerから呼び出すエンドポイント
//Cloud Scheduler で使用しているサービスをデプロイするときは、未承認の呼び出しを許可しないでください。

type Server struct {
	kompal    kompal.Fetcher
	storage   storage.Storage
	notifiers []notifier.Notifier
	analyzer  analyzer.Analyzer

	mux    *http.ServeMux
	server *http.Server
}

func New(f kompal.Fetcher, s storage.Storage, ns []notifier.Notifier, a analyzer.Analyzer) *Server {
	server := &Server{
		kompal:    f,
		storage:   s,
		notifiers: ns,
		analyzer:  a,

		mux: http.NewServeMux(),
	}

	server.registerHandlers()

	return server
}

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

func (s *Server) GracefulStop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) registerHandlers() {
	s.mux.Handle("/watch", s.watchHandler())
}
