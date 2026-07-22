package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

// BearerAuth requires an "Authorization: Bearer <token>" header equal to
// token. It never takes a logger, so the token can never end up in a log
// line, and it responds with the same generic body whether the header was
// missing, malformed, or simply wrong, so a caller cannot distinguish those
// cases from the response.
func BearerAuth(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got, ok := strings.CutPrefix(r.Header.Get("Authorization"), bearerPrefix)
			// subtle.ConstantTimeCompare on mismatched lengths short-circuits
			// on length, not content, so this stays timing-safe against
			// guessing the token itself.
			if !ok || subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
				unauthorized(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}
