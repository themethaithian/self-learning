package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"
)

// BearerAuth requires an "Authorization: Bearer <token>" header equal to
// token. It never takes a logger, so the token can never end up in a log
// line, and it responds with the same generic body whether the header was
// missing, malformed, or simply wrong, so a caller cannot distinguish those
// cases from the response. It panics if token is empty: NewMux runs at
// process startup, before ListenAndServe, so an empty secret fails the
// boot instead of silently authenticating every request in production.
func BearerAuth(token string) func(http.Handler) http.Handler {
	if token == "" {
		panic("middleware: BearerAuth requires a non-empty token")
	}
	want := sha256.Sum256([]byte(token))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scheme, rest, found := strings.Cut(r.Header.Get("Authorization"), " ")
			// EqualFold is deliberately not constant-time: the scheme name
			// is public (RFC 6750), only the token comparison below is a
			// secret comparison.
			if !found || !strings.EqualFold(scheme, "bearer") {
				unauthorized(w)
				return
			}

			// Hashing both sides first makes the compare operands
			// fixed-length, so ConstantTimeCompare never short-circuits
			// on length and only ever compares digest bytes.
			candidate := strings.TrimLeft(rest, " ")
			got := sha256.Sum256([]byte(candidate))
			if subtle.ConstantTimeCompare(got[:], want[:]) != 1 {
				unauthorized(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("WWW-Authenticate", "Bearer")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}
