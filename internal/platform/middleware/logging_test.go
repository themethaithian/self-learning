package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		wantStatus int
	}{
		{
			name: "explicit status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "implicit 200 via Write",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "not found",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))

			handler := Logging(logger)(tt.handler)
			req := httptest.NewRequest(http.MethodGet, "/some/path", nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("recorded status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var logEntry map[string]any
			if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &logEntry); err != nil {
				t.Fatalf("unmarshal log line: %v, raw=%q", err, buf.String())
			}
			if logEntry["method"] != http.MethodGet {
				t.Errorf("logged method = %v, want GET", logEntry["method"])
			}
			if logEntry["path"] != "/some/path" {
				t.Errorf("logged path = %v, want /some/path", logEntry["path"])
			}
			if int(logEntry["status"].(float64)) != tt.wantStatus {
				t.Errorf("logged status = %v, want %d", logEntry["status"], tt.wantStatus)
			}
			if _, ok := logEntry["duration"]; !ok {
				t.Error("log missing duration field")
			}
		})
	}
}
