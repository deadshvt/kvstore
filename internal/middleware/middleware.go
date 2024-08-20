package middleware

import "net/http"

func ChainMiddleware(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := 0; i < len(middlewares); i++ {
			final = middlewares[i](final)
		}

		return final
	}
}
