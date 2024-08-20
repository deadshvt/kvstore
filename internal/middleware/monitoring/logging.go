package monitoring

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/key"
	"github.com/deadshvt/kvstore/pkg/logger"
)

type Lgr struct {
	Logger zerolog.Logger
}

func NewLgr(logger zerolog.Logger) *Lgr {
	return &Lgr{
		Logger: logger,
	}
}

func (lr *Lgr) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		requestID, ok := r.Context().Value(key.RequestID).(string)
		if !ok {
			requestID = "unknown"
		}

		logger.LogWithParams(lr.Logger, "Completed request", struct {
			Method    string
			URI       string
			Took      time.Duration
			RequestID string
		}{Method: r.Method, URI: r.RequestURI, Took: time.Since(start), RequestID: requestID})
	})
}
