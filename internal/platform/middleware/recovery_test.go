package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
		wantPanic  bool
	}{
		{
			name: "panic with string",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic("boom")
			},
			wantStatus: http.StatusInternalServerError,
			wantPanic:  true,
		},
		{
			name: "panic with error",
			handler: func(w http.ResponseWriter, r *http.Request) {
				panic(errors.New("boom"))
			},
			wantStatus: http.StatusInternalServerError,
			wantPanic:  true,
		},
		{
			name: "no panic passes through untouched",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusTeapot)
			},
			wantStatus: http.StatusTeapot,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			handler := Recovery(logger)(tt.handler)
			req := httptest.NewRequest(http.MethodGet, "/panic", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if !tt.wantPanic {
				return
			}

			var body map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body["error"] == "" {
				t.Error("expected non-empty error field in body")
			}
			if !bytes.Contains(buf.Bytes(), []byte("panic recovered")) {
				t.Errorf("expected log to contain \"panic recovered\", got %q", buf.String())
			}
			if !bytes.Contains(buf.Bytes(), []byte(`"stack"`)) {
				t.Errorf("expected log to contain a stack field, got %q", buf.String())
			}
		})
	}
}
