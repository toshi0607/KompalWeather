package http

import (
	"net/http"

	"golang.org/x/sync/errgroup"
)

func (s *Server) watchHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		f, err := s.kompal.Fetch(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.storage.Save(ctx, f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := s.analyzer.Analyze(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		eg, ctx := errgroup.WithContext(ctx)
		for _, n := range s.notifiers {
			eg.Go(func() error {
				return n.Notify(ctx, result)
			})
		}
		if err := eg.Wait(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
