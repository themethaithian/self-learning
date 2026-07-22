package httpserver

import (
	"net/http"

	"github.com/themethaithian/self-learning/internal/platform/middleware"
)

// NewMux structurally guarantees a route can only skip auth by being
// registered on mux directly instead of protectedMux — this file does
// that for nothing but /healthz, so future routes added to protectedMux
// are authenticated by construction, not by an allowlist someone has to
// remember to update.
func NewMux(pinger Pinger, bearerToken string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler(pinger))
	mux.Handle("/", middleware.BearerAuth(bearerToken)(protectedMux()))
	return mux
}

func protectedMux() *http.ServeMux {
	return http.NewServeMux()
}
