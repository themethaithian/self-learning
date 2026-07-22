package middleware

import "net/http"

const (
	corsAllowedMethods = "GET, POST, PATCH, OPTIONS"
	corsAllowedHeaders = "Authorization, Content-Type"
)

// CORS allows browser requests from allowedOrigin only, echoing that exact
// origin back rather than "*" since requests carry a bearer token that must
// stay invisible to other origins. It answers OPTIONS preflight itself, so
// a preflight request — which never carries Authorization — reaches this
// middleware and stops here instead of falling through to BearerAuth.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			if origin := r.Header.Get("Origin"); origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", corsAllowedMethods)
				w.Header().Set("Access-Control-Allow-Headers", corsAllowedHeaders)
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
