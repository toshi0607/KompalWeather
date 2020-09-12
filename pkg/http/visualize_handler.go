package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/visualizer"
)

const visualizerHandlerName = "visualizeHandler"

type RequestBody struct {
	ReportType visualizer.ReportType `json:"ReportType"`
}

// This application is intended to be hosted by Cloud Run which doesn't allow unauthenticated.
// Called from Cloud Scheduler. Service account OIDC token with roles/run.invoker is required.
func (s *VisualizeServer) visualizeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var req RequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.log.Error("failed to decode request body", err)
			http.Error(w, "failed to read body", http.StatusBadRequest)
			return
		}
		if !req.ReportType.IsValid() {
			err := fmt.Errorf("invalid report type: %v", req.ReportType)
			s.log.Info(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		ctx := r.Context()
		s.log.SetHandlerName(visualizerHandlerName)
		s.log.Info("%s started", visualizerHandlerName)

		result, err := s.visualizer.Save(ctx, req.ReportType)
		if err != nil {
			s.log.Error("failed to save", err)
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
