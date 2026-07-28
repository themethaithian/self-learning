package infra

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
	"github.com/themethaithian/self-learning/internal/platform/httpserver"
)

var _ httpserver.Registrar = (*Handler)(nil)

type fakeService struct {
	entries []learningapp.ProgressEntry
	listErr error

	setEntry learningapp.ProgressEntry
	setErr   error

	setCalled  bool
	setTopic   string
	setConcept string
	setState   string
}

func (f *fakeService) ListProgress(context.Context) ([]learningapp.ProgressEntry, error) {
	return f.entries, f.listErr
}

func (f *fakeService) SetProgress(_ context.Context, topic, concept, state string) (learningapp.ProgressEntry, error) {
	f.setCalled = true
	f.setTopic, f.setConcept, f.setState = topic, concept, state
	return f.setEntry, f.setErr
}

func newTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return slog.New(slog.NewTextHandler(&buf, nil)), &buf
}

func doRequest(h *Handler, method, path, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHandlerGetProgress_Empty(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodGet, "/api/v1/progress", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got, want := rec.Body.String(), `{"concepts":[]}`+"\n"; got != want {
		t.Fatalf("body = %q, want %q (empty must marshal as [] never null)", got, want)
	}
}

func TestHandlerGetProgress_WithEntries(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{entries: []learningapp.ProgressEntry{
		{Topic: "ddia", Concept: "b-trees", State: mustChunkState(t, "in_progress")},
	}}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodGet, "/api/v1/progress", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	want := `{"concepts":[{"topic":"ddia","concept":"b-trees","state":"in_progress","last_read_at":null,"first_passed_at":null}]}` + "\n"
	if got := rec.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerGetProgress_ServiceErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	h := NewHandler(&fakeService{listErr: errors.New("mysql: connection refused on 10.0.0.5")}, logger)

	rec := doRequest(h, http.MethodGet, "/api/v1/progress", "")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body := rec.Body.String()
	for _, leak := range []string{"mysql", "10.0.0.5", "connection refused"} {
		if strings.Contains(body, leak) {
			t.Fatalf("body = %q leaks internal error detail %q", body, leak)
		}
	}
	if !strings.Contains(logBuf.String(), "connection refused") {
		t.Fatalf("log output = %q, want it to contain the real error", logBuf.String())
	}
}

func TestHandlerPutProgress_Success(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{setEntry: learningapp.ProgressEntry{Topic: "ddia", Concept: "b-trees", State: mustChunkState(t, "in_progress")}}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", `{"state":"in_progress"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if !svc.setCalled {
		t.Fatal("expected SetProgress to be called")
	}
	if svc.setTopic != "ddia" || svc.setConcept != "b-trees" || svc.setState != "in_progress" {
		t.Fatalf("SetProgress called with (%q, %q, %q), want (ddia, b-trees, in_progress)", svc.setTopic, svc.setConcept, svc.setState)
	}
	want := `{"topic":"ddia","concept":"b-trees","state":"in_progress","last_read_at":null,"first_passed_at":null}` + "\n"
	if got := rec.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerPutProgress_LessonNotFound(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{setErr: fmt.Errorf("learning: set progress: %w", learningapp.ErrLessonNotFound)}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/no-lesson", `{"state":"in_progress"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

func TestHandlerPutProgress_InvalidState(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{setErr: fmt.Errorf("learning: state %q: %w", "locked", learningapp.ErrInvalidState)}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", `{"state":"locked"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandlerPutProgress_InvalidJSON(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", `not json`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestHandlerPutProgress_UnknownField guards against a typo'd field name
// (e.g. "stat" instead of "state") decoding as a no-op State == "" and
// silently doing the wrong thing instead of failing loudly.
func TestHandlerPutProgress_UnknownField(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", `{"stat":"passed"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if svc.setCalled {
		t.Fatal("SetProgress must not be called for a body with an unknown field")
	}
}

// TestHandlerPutProgress_TrailingJSON guards dec.More(): a plain
// json.Decoder.Decode call silently stops after the first JSON value, so a
// body carrying two concatenated values would otherwise decode (and act on)
// only the first one instead of failing.
func TestHandlerPutProgress_TrailingJSON(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", `{"state":"passed"}{"state":"in_progress"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if svc.setCalled {
		t.Fatal("SetProgress must not be called for a body with trailing JSON")
	}
}

func TestHandlerPutProgress_ServiceErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	svc := &fakeService{setErr: errors.New("mysql: connection refused on 10.0.0.5")}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", `{"state":"passed"}`)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body := rec.Body.String()
	for _, leak := range []string{"mysql", "10.0.0.5", "connection refused"} {
		if strings.Contains(body, leak) {
			t.Fatalf("body = %q leaks internal error detail %q", body, leak)
		}
	}
	if !strings.Contains(logBuf.String(), "connection refused") {
		t.Fatalf("log output = %q, want it to contain the real error", logBuf.String())
	}
}

func TestHandlerRejectsUnknownMethodOnCollection(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandlerRejectsUnknownMethodOnItem(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodDelete, "/api/v1/progress/ddia/b-trees", "")

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}
