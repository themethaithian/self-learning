package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// protectedMux is empty, so a request that reaches it 404s — that 404
// (rather than 401) is proof BearerAuth let the request through. A route
// need not exist yet for the auth boundary to hold: anything but GET
// /healthz falls under the auth-wrapped "/" mount in NewMux, so a future
// route is protected the moment it is registered on protectedMux.
func TestNewMux(t *testing.T) {
	mux := NewMux(fakePinger{}, "secret-token")

	tests := []struct {
		name       string
		method     string
		path       string
		authHeader string
		wantStatus int
	}{
		{name: "healthz stays public with no token", method: http.MethodGet, path: "/healthz", wantStatus: http.StatusOK},
		{name: "POST healthz has no route, method mismatch lands on auth", method: http.MethodPost, path: "/healthz", wantStatus: http.StatusUnauthorized},
		{name: "other route, no token", method: http.MethodGet, path: "/anything", wantStatus: http.StatusUnauthorized},
		{name: "other route, wrong token", method: http.MethodGet, path: "/anything", authHeader: "Bearer wrong-token", wantStatus: http.StatusUnauthorized},
		{name: "other route, correct token reaches protectedMux", method: http.MethodGet, path: "/anything", authHeader: "Bearer secret-token", wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
