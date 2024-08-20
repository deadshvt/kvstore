package protection

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog"

	"github.com/deadshvt/kvstore/internal/errs"
)

type Content struct {
	Type   string
	Logger zerolog.Logger
}

func NewContent(t string, logger zerolog.Logger) *Content {
	return &Content{
		Type:   t,
		Logger: logger,
	}
}

func (c *Content) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("Content-Type")
		if !strings.Contains(t, c.Type) {
			msg := errs.WrapError(errs.ErrMediaType, nil).Error()
			c.Logger.Error().Msg(msg)
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}
