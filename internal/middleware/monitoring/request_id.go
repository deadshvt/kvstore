package monitoring

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/key"
)

type RequestIDGenerator struct {
	Logger zerolog.Logger
}

func NewRequestIDGenerator(logger zerolog.Logger) *RequestIDGenerator {
	return &RequestIDGenerator{
		Logger: logger,
	}
}

func (rg *RequestIDGenerator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), key.RequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
