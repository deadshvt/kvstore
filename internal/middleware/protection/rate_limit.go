package protection

import (
	"net/http"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"

	"github.com/deadshvt/kvstore/internal/errs"
)

type RateLimiter struct {
	Limiter *rate.Limiter
	Logger  zerolog.Logger
}

func NewRateLimiter(limiter *rate.Limiter, logger zerolog.Logger) *RateLimiter {
	return &RateLimiter{
		Limiter: limiter,
		Logger:  logger,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.Limiter.Allow() {
			msg := errs.WrapError(errs.ErrTooManyRequests, nil).Error()
			rl.Logger.Error().Msgf(msg+" from %s", r.RemoteAddr)
			http.Error(w, msg, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
