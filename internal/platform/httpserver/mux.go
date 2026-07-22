package httpserver

import (
	"net/http"

	"github.com/themethaithian/self-learning/internal/platform/middleware"
)

// Registrar lets a bounded context wire its own routes onto the
// authenticated mux, so NewMux does not grow one parameter per context and
// every registered route stays protected by construction.
type Registrar interface {
	Register(mux *http.ServeMux)
}

// NewMux structurally guarantees a route can only skip auth by being
// registered on mux directly instead of protectedMux — this file does
// that for nothing but /healthz, so every registrar's routes are
// authenticated by construction, not by an allowlist someone has to
// remember to update.
func NewMux(pinger Pinger, bearerToken string, registrars ...Registrar) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler(pinger))

	protected := protectedMux()
	for _, registrar := range registrars {
		registrar.Register(protected)
	}
	mux.Handle("/", middleware.BearerAuth(bearerToken)(protected))
	return mux
}

func protectedMux() *http.ServeMux {
	return http.NewServeMux()
}
