package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
)

// Recovery catches panics from downstream handlers, logs the value and stack
// trace, and responds 500 so a single bad request cannot take the server
// down.
func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// http.ErrAbortHandler is the stdlib's own sentinel for
					// "abort this handler without logging or a response" —
					// re-panic so net/http's own recovery handles it silently.
					if rec == http.ErrAbortHandler {
						panic(rec)
					}
					logger.Error("panic recovered",
						"error", rec,
						"stack", string(debug.Stack()),
						"method", r.Method,
						"path", r.URL.Path,
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{"error": "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
