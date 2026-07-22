package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Logging must wrap Recovery, not the reverse, so the logged status
// reflects the 500 Recovery writes after catching a panic rather than the
// unset status at the moment the panic propagates.
func TestLoggingWrapsRecovery(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	panicking := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	handler := Logging(logger)(Recovery(logger)(panicking))

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}

	var accessLog map[string]any
	for _, line := range bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n")) {
		var entry map[string]any
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Fatalf("unmarshal log line: %v", err)
		}
		if entry["msg"] == "http request" {
			accessLog = entry
		}
	}
	if accessLog == nil {
		t.Fatal("no access log line found")
	}
	if int(accessLog["status"].(float64)) != http.StatusInternalServerError {
		t.Errorf("access log status = %v, want 500", accessLog["status"])
	}
}
