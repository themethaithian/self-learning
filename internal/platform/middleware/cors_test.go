package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const allowedOrigin = "http://localhost:3000"

func TestCORS(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		origin                 string
		wantStatus             int
		wantAllowOrig          string
		wantAllowedCORSHeaders bool
		wantNextCalled         bool
	}{
		{
			name:                   "preflight from allowed origin",
			method:                 http.MethodOptions,
			origin:                 allowedOrigin,
			wantStatus:             http.StatusNoContent,
			wantAllowOrig:          allowedOrigin,
			wantAllowedCORSHeaders: true,
		},
		{
			name:       "preflight from disallowed origin",
			method:     http.MethodOptions,
			origin:     "http://evil.example",
			wantStatus: http.StatusNoContent,
		},
		{
			name:           "normal request from allowed origin",
			method:         http.MethodGet,
			origin:         allowedOrigin,
			wantStatus:     http.StatusOK,
			wantAllowOrig:  allowedOrigin,
			wantNextCalled: true,
		},
		{
			name:           "normal request from disallowed origin",
			method:         http.MethodGet,
			origin:         "http://evil.example",
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
		{
			name:           "no Origin header (non-browser or same-origin request)",
			method:         http.MethodGet,
			origin:         "",
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			handler := CORS(allowedOrigin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest(tt.method, "/anything", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := rec.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowOrig {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowOrig)
			}
			if got := rec.Header().Get("Vary"); got != "Origin" {
				t.Errorf(`Vary = %q, want "Origin" (a shared cache must not replay this response cross-origin)`, got)
			}

			// wantMethods is a literal, not corsAllowedMethods: asserting
			// against the same constant that produces the header would let a
			// mutation of that constant pass silently.
			wantMethods, wantHeaders, wantMaxAge := "", "", ""
			if tt.wantAllowedCORSHeaders {
				wantMethods, wantHeaders, wantMaxAge = "GET, POST, PATCH, PUT, OPTIONS", corsAllowedHeaders, corsMaxAge
			}
			if got := rec.Header().Get("Access-Control-Allow-Methods"); got != wantMethods {
				t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, wantMethods)
			}
			if got := rec.Header().Get("Access-Control-Allow-Headers"); got != wantHeaders {
				t.Errorf("Access-Control-Allow-Headers = %q, want %q", got, wantHeaders)
			}
			if got := rec.Header().Get("Access-Control-Max-Age"); got != wantMaxAge {
				t.Errorf("Access-Control-Max-Age = %q, want %q", got, wantMaxAge)
			}
			if nextCalled != tt.wantNextCalled {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNextCalled)
			}
		})
	}
}
