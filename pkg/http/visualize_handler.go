package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/report"
)

const visualizerHandlerName = "visualizeHandler"

type RequestBody struct {
	ReportKind report.Kind `json:"reportKind"`
}

type ResponseBody struct {
	Files []string `json:"files,omitempty"`
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
		if !req.ReportKind.IsValid() {
			err := fmt.Errorf("invalid report type: %v", req.ReportKind)
			s.log.Info(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		ctx := r.Context()
		s.log.SetHandlerName(visualizerHandlerName)
		s.log.Info("%s started", visualizerHandlerName)

		files, err := s.visualizer.Save(ctx, req.ReportKind)
		if err != nil {
			s.log.Error("failed to save", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ResponseBody{Files: files}); err != nil {
			s.log.Error("failed to encode", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
