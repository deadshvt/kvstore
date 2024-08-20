package protection

import (
	"fmt"
	"net/http"

	"github.com/deadshvt/kvstore/internal/errs"

	"github.com/rs/zerolog"
)

type Recoverer struct {
	Logger zerolog.Logger
}

func NewRecoverer(logger zerolog.Logger) *Recoverer {
	return &Recoverer{
		Logger: logger,
	}
}

func (rr *Recoverer) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				var desc error
				switch e := err.(type) {
				case error:
					desc = e
				default:
					desc = fmt.Errorf("%v", e)
				}

				msg := errs.WrapError(errs.ErrPanic, desc).Error()
				rr.Logger.Error().Msg(msg)
				http.Error(w, msg, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
