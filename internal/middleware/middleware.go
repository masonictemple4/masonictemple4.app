package middleware

import "net/http"

// TODO: Auth, Trace, Observable, Metrics.

// TODO: Probably want to secure this later.
func CorsMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	}
}
