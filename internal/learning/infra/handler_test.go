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
	"time"

	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
	"github.com/themethaithian/self-learning/internal/learning/domain"
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

	recordEntry learningapp.AttemptRecord
	recordErr   error

	recordCalled         bool
	recordTopic          string
	recordConcept        string
	recordQuestion       string
	recordKind           string
	recordConfidence     string
	recordOutcome        string
	recordSelectedOption *string
}

func (f *fakeService) ListProgress(context.Context) ([]learningapp.ProgressEntry, error) {
	return f.entries, f.listErr
}

func (f *fakeService) SetProgress(_ context.Context, topic, concept, state string) (learningapp.ProgressEntry, error) {
	f.setCalled = true
	f.setTopic, f.setConcept, f.setState = topic, concept, state
	return f.setEntry, f.setErr
}

func (f *fakeService) RecordAttempt(_ context.Context, topic, concept, question, kind, confidence, outcome string, selectedOption *string) (learningapp.AttemptRecord, error) {
	f.recordCalled = true
	f.recordTopic, f.recordConcept, f.recordQuestion = topic, concept, question
	f.recordKind, f.recordConfidence, f.recordOutcome = kind, confidence, outcome
	f.recordSelectedOption = selectedOption
	return f.recordEntry, f.recordErr
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

// TestHandlerPutProgress_InvalidSlug pins R1(a): the handler must map
// ErrInvalidSlug to 400, not fall through to the generic 500 branch a
// slug-shape error used to hit.
func TestHandlerPutProgress_InvalidSlug(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{setErr: fmt.Errorf("learning: %w: %q", learningapp.ErrInvalidSlug, "B-Trees")}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/B-Trees", `{"state":"in_progress"}`)

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

// TestHandlerPutProgress_BodyTooLarge distinguishes an oversized body (413)
// from ordinary malformed JSON (400) so a client can tell "you sent garbage"
// from "you sent too much" instead of the same generic 400 for both.
func TestHandlerPutProgress_BodyTooLarge(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	oversized := `{"state":"` + strings.Repeat("x", maxPutBodyBytes) + `"}`
	rec := doRequest(h, http.MethodPut, "/api/v1/progress/ddia/b-trees", oversized)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
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

func TestHandlerPostAttempt_Success(t *testing.T) {
	logger, _ := newTestLogger()
	selected := "Option B"
	svc := &fakeService{recordEntry: learningapp.AttemptRecord{
		CheckKey: "abc123", Question: "What is a B-tree?", Kind: mustCheckKind(t, "mcq"), Confidence: mustConfidence(t, "confident"),
		Outcome: mustAttemptOutcome(t, "incorrect"), SelectedOption: &selected, GradedBy: domain.GradedBySelf,
		CreatedAt: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC),
	}}
	h := NewHandler(svc, logger)

	body := `{"question":"What is a B-tree?","kind":"mcq","confidence":"confident","outcome":"incorrect","selected_option":"Option B"}`
	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", body)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if !svc.recordCalled {
		t.Fatal("expected RecordAttempt to be called")
	}
	if svc.recordTopic != "ddia" || svc.recordConcept != "b-trees" || svc.recordQuestion != "What is a B-tree?" {
		t.Fatalf("RecordAttempt called with topic=%q concept=%q question=%q, want ddia/b-trees/\"What is a B-tree?\"", svc.recordTopic, svc.recordConcept, svc.recordQuestion)
	}
	if svc.recordKind != "mcq" || svc.recordConfidence != "confident" || svc.recordOutcome != "incorrect" {
		t.Fatalf("RecordAttempt called with kind=%q confidence=%q outcome=%q, want mcq/confident/incorrect", svc.recordKind, svc.recordConfidence, svc.recordOutcome)
	}
	if svc.recordSelectedOption == nil || *svc.recordSelectedOption != "Option B" {
		t.Fatalf("RecordAttempt called with selected_option=%v, want \"Option B\"", svc.recordSelectedOption)
	}
	want := `{"check_key":"abc123","question":"What is a B-tree?","kind":"mcq","confidence":"confident","outcome":"incorrect","selected_option":"Option B","graded_by":"self","created_at":"2026-01-01T12:00:00Z"}` + "\n"
	if got := rec.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerPostAttempt_LessonNotFound(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{recordErr: fmt.Errorf("learning: record attempt: %w", learningapp.ErrLessonNotFound)}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/no-lesson/attempts", `{"question":"q","kind":"short_answer","confidence":"unsure","outcome":"correct"}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

func TestHandlerPostAttempt_CheckNotInLesson(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{recordErr: fmt.Errorf("learning: record attempt: %w", learningapp.ErrCheckNotInLesson)}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", `{"question":"not a real question","kind":"short_answer","confidence":"unsure","outcome":"correct"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandlerPostAttempt_InvalidSlug(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{recordErr: fmt.Errorf("learning: %w: %q", learningapp.ErrInvalidSlug, "B-Trees")}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/B-Trees/attempts", `{"question":"q","kind":"short_answer","confidence":"unsure","outcome":"correct"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandlerPostAttempt_InvalidAttempt(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{recordErr: fmt.Errorf("learning: %w: %v", learningapp.ErrInvalidAttempt, "bad kind")}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", `{"question":"q","kind":"fill_in_blank","confidence":"unsure","outcome":"correct"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestHandlerPostAttempt_InvalidJSON(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", `not json`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandlerPostAttempt_UnknownField(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", `{"question":"q","kind":"short_answer","confidence":"unsure","outcome":"correct","extra":"field"}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if svc.recordCalled {
		t.Fatal("RecordAttempt must not be called for a body with an unknown field")
	}
}

func TestHandlerPostAttempt_TrailingJSON(t *testing.T) {
	logger, _ := newTestLogger()
	svc := &fakeService{}
	h := NewHandler(svc, logger)

	body := `{"question":"q","kind":"short_answer","confidence":"unsure","outcome":"correct"}{"question":"q2","kind":"short_answer","confidence":"unsure","outcome":"correct"}`
	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", body)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
	if svc.recordCalled {
		t.Fatal("RecordAttempt must not be called for a body with trailing JSON")
	}
}

func TestHandlerPostAttempt_BodyTooLarge(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(&fakeService{}, logger)

	oversized := `{"question":"` + strings.Repeat("x", maxAttemptBodyBytes) + `","kind":"short_answer","confidence":"unsure","outcome":"correct"}`
	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", oversized)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestHandlerPostAttempt_ServiceErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	svc := &fakeService{recordErr: errors.New("mysql: connection refused on 10.0.0.5")}
	h := NewHandler(svc, logger)

	rec := doRequest(h, http.MethodPost, "/api/v1/progress/ddia/b-trees/attempts", `{"question":"q","kind":"short_answer","confidence":"unsure","outcome":"correct"}`)

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
