package http

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// This application is intended to be hosted by Cloud Run which doesn't allow unauthenticated.
// Called from Cloud Scheduler. Service account OIDC token with roles/run.invoker is required.
func (s *Server) watchHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		fmt.Print("request to watch handler")
		f, err := s.kompal.Fetch(ctx)
		if err != nil {
			fmt.Print(fmt.Errorf("failed to fetch kompal status: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		st, err := s.storage.Save(ctx, f)
		if err != nil {
			fmt.Print(fmt.Errorf("failed to save status: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("new status saved!: %s\n", st)

		result, err := s.analyzer.Analyze(ctx)
		if err != nil {
			fmt.Print(fmt.Errorf("failed to analyze: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("result: %s\n", result)

		eg, ctx := errgroup.WithContext(ctx)
		for _, n := range s.notifiers {
			n := n
			eg.Go(func() error {
				fmt.Printf("notify type: %s\n", n.Type())
				return n.Notify(ctx, result)
			})
		}
		if err := eg.Wait(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err != nil {
			fmt.Print(fmt.Errorf("failed to notify: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
