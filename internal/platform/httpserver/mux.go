package httpserver

import (
	"net/http"

	"github.com/themethaithian/self-learning/internal/platform/middleware"
)

// NewMux builds the process's top-level handler. GET /healthz is public;
// every other route lives on protectedMux, which is mounted at "/" behind
// BearerAuth. A route can only skip auth by being registered on mux
// directly instead of protectedMux — this file does that for nothing but
// /healthz, so future routes added to protectedMux are authenticated by
// construction, not by an allowlist someone has to remember to update.
func NewMux(pinger Pinger, bearerToken string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler(pinger))
	mux.Handle("/", middleware.BearerAuth(bearerToken)(protectedMux()))
	return mux
}

// protectedMux registers every route that requires authentication. Later
// tickets add handlers here.
func protectedMux() *http.ServeMux {
	return http.NewServeMux()
}
