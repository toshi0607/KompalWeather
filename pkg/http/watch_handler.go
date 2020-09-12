package http

import (
	"encoding/json"
	"net/http"
)

const watchHandlerName = "watchHandler"

// This application is intended to be hosted by Cloud Run which doesn't allow unauthenticated.
// Called from Cloud Scheduler. Service account OIDC token with roles/run.invoker is required.
func (s *CoreServer) watchHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		s.log.SetHandlerName(watchHandlerName)
		s.log.Info("%s started", watchHandlerName)

		result, err := s.controller.Watch(ctx)
		if err != nil {
			s.log.Error("failed to watch", err)
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
