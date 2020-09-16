package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/toshi0607/kompal-weather/pkg/report"
)

const trendHandlerName = "trendHandler"

// This application is intended to be hosted by Cloud Run which doesn't allow unauthenticated.
// Called from Cloud Scheduler. Service account OIDC token with roles/run.invoker is required.
func (s *CoreServer) trendHandler() http.Handler {
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
		if req.ReportKind != report.WeekAgoReport {
			err := fmt.Errorf("invalid report type: %v", req.ReportKind)
			s.log.Info(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		ctx := r.Context()
		s.log.SetHandlerName(trendHandlerName)
		s.log.Info("%s started", trendHandlerName)

		if err := s.controller.Trend(ctx, req.ReportKind); err != nil {
			s.log.Error("failed to trend", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
