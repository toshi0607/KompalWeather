package http

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const handlerName = "watchHandler"

// This application is intended to be hosted by Cloud Run which doesn't allow unauthenticated.
// Called from Cloud Scheduler. Service account OIDC token with roles/run.invoker is required.
func (s *Server) watchHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		s.log.SetHandlerName(handlerName)
		s.log.Info("%s started", handlerName)

		f, err := s.kompal.Fetch(ctx)
		if err != nil {
			s.log.Error("failed to fetch kompal status", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.log.Info("fetched: %+v", *f)

		st, err := s.storage.Save(ctx, f)
		if err != nil {
			s.log.Error("failed to save status", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.log.Info("saved: %+v", *st)

		if err := s.monitor.CreatePoint(ctx, st); err != nil {
			s.log.Error("failed to create point", err)
			// Keep processing
		}

		result, err := s.analyzer.Analyze(ctx)
		if err != nil {
			s.log.Error("failed to analyze", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		time.Now()
		s.log.Info("result: %+v", *result)

		eg, ctx := errgroup.WithContext(ctx)
		for _, n := range s.notifiers {
			n := n
			eg.Go(func() error {
				s.log.Info("notification type: %v", n.Type())
				return n.Notify(ctx, result)
			})
		}
		if err := eg.Wait(); err != nil {
			s.log.Error("failed to notify", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			s.log.Error("failed to encode", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
