// Package httpserver builds the API's http.ServeMux and http.Server.
package httpserver

import (
	"net/http"
	"time"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	// WriteTimeout is generous: recall-check and DSA-review handlers call the
	// Anthropic API synchronously and grading can take 10-20s.
	writeTimeout = 60 * time.Second
	idleTimeout  = 120 * time.Second
)

// New constructs an http.Server with explicit timeouts so a slow or stalled
// client can never hold a connection open indefinitely.
func New(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}
