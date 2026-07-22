package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMuxHealthzStaysPublic(t *testing.T) {
	mux := NewMux(fakePinger{}, "secret-token")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (no Authorization header sent)", rec.Code)
	}
}

// A route need not even exist yet for this to hold: anything other than
// /healthz falls under the auth-wrapped "/" mount in NewMux, so a future
// route is protected the moment it is registered on protectedMux.
func TestNewMuxOtherRoutesRequireAuth(t *testing.T) {
	mux := NewMux(fakePinger{}, "secret-token")

	req := httptest.NewRequest(http.MethodGet, "/anything", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401 (non-healthz routes must require auth)", rec.Code)
	}
}

func TestNewMuxOtherRoutesPassAuth(t *testing.T) {
	mux := NewMux(fakePinger{}, "secret-token")

	req := httptest.NewRequest(http.MethodGet, "/anything", nil)
	req.Header.Set("Authorization", "Bearer secret-token")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code == http.StatusUnauthorized {
		t.Fatalf("status = %d, correct token must not be rejected", rec.Code)
	}
}
