package http

import (
	"net/http"
)

// TODO: для чего мидлвейр?
func FilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}
