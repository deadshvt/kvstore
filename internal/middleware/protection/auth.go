package protection

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/deadshvt/kvstore/internal/errs"
	"github.com/deadshvt/kvstore/internal/security"
)

const Prefix = "Bearer "

type Authenticator struct {
	JWTService *security.JWTService
	Logger     zerolog.Logger
}

func NewAuthenticator(jwts *security.JWTService, logger zerolog.Logger) *Authenticator {
	return &Authenticator{
		JWTService: jwts,
		Logger:     logger,
	}
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, Prefix) {
			msg := errs.WrapError(errs.ErrUnauthorized, nil).Error()
			a.Logger.Error().Msg(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		token := authHeader[len(Prefix):]
		err := a.JWTService.VerifyToken(token)
		if err != nil {
			msg := errs.WrapError(errs.ErrUnauthorized, err).Error()
			log.Error().Msg(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
