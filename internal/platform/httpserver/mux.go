package httpserver

import "net/http"

// NewMux registers every route. pinger backs GET /healthz.
func NewMux(pinger Pinger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler(pinger))
	return mux
}
