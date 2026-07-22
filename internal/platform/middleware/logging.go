// Package middleware provides hand-written http.Handler wrappers in the
// func(http.Handler) http.Handler pattern.
package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder captures the status code a handler writes. Write can be
// called without a preceding WriteHeader (net/http then defaults to 200), so
// both methods must record the status.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.wroteHeader {
		return
	}
	r.status = status
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

// Logging logs method, path, status, and duration for every request. It must
// wrap Recovery (not the reverse) so the status it reports reflects any 500
// Recovery wrote after a panic.
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()

			next.ServeHTTP(rec, r)

			logger.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration", time.Since(start),
			)
		})
	}
}
