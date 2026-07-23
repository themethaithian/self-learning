package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	"github.com/themethaithian/self-learning/internal/curriculum/domain"
	"github.com/themethaithian/self-learning/internal/platform/httpserver"
)

var _ httpserver.Registrar = (*Handler)(nil)

type fakeTreeService struct {
	topics []domain.Topic
	err    error

	lesson    domain.Lesson
	lessonErr error
}

func (f fakeTreeService) Tree(context.Context) ([]domain.Topic, error) {
	return f.topics, f.err
}

func (f fakeTreeService) Lesson(context.Context, string, string) (domain.Lesson, error) {
	return f.lesson, f.lessonErr
}

func newTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	return slog.New(slog.NewTextHandler(&buf, nil)), &buf
}

func buildTestTopic(t *testing.T, track, slug, title string, position int) domain.Topic {
	t.Helper()
	tr, err := domain.NewTrack(track)
	if err != nil {
		t.Fatalf("NewTrack(%q): %v", track, err)
	}
	pos, err := domain.NewPosition(position)
	if err != nil {
		t.Fatalf("NewPosition(%d): %v", position, err)
	}
	conceptSlug, err := domain.NewSlug(slug + "-concept")
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	concept, err := domain.NewConcept(conceptSlug, "Concept "+title, "an outline", pos)
	if err != nil {
		t.Fatalf("NewConcept: %v", err)
	}
	chapterSlug, err := domain.NewSlug(slug + "-chapter")
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	chapter, err := domain.NewChapter(chapterSlug, "Chapter "+title, pos, []domain.Concept{concept})
	if err != nil {
		t.Fatalf("NewChapter: %v", err)
	}
	topicSlug, err := domain.NewSlug(slug)
	if err != nil {
		t.Fatalf("NewSlug: %v", err)
	}
	topic, err := domain.NewTopic(tr, topicSlug, title, pos, []domain.Chapter{chapter})
	if err != nil {
		t.Fatalf("NewTopic: %v", err)
	}
	return topic
}

// TestHandlerGetTree_Success pins the wire format byte-for-byte, the same
// way TestHandlerGetTree_EmptyIsArrayNotNull pins the empty case: decoding
// into treeResponse and comparing structs (the previous approach) is
// symmetric with the encoder, so a struct tag rename round-trips clean and
// the test can never catch it. Two tracks, fed in reverse-canonical arrival
// order, also proves grouping order comes from domain.Tracks(), not from
// the order topics arrived in.
func TestHandlerGetTree_Success(t *testing.T) {
	dsaTopic := buildTestTopic(t, "dsa", "dsa-arrays", "Arrays", 1)
	dddTopic := buildTestTopic(t, "ddd", "ddd-aggregates", "Aggregates", 1)
	logger, _ := newTestLogger()
	h := NewHandler(fakeTreeService{topics: []domain.Topic{dsaTopic, dddTopic}}, logger)

	rec := doGetTree(h, http.MethodGet)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}

	const want = `{"tracks":[{"track":"ddd","topics":[{"slug":"ddd-aggregates","title":"Aggregates","position":1,"chapters":[{"slug":"ddd-aggregates-chapter","title":"Chapter Aggregates","position":1,"concepts":[{"slug":"ddd-aggregates-concept","title":"Concept Aggregates","position":1}]}]}]},{"track":"dsa","topics":[{"slug":"dsa-arrays","title":"Arrays","position":1,"chapters":[{"slug":"dsa-arrays-chapter","title":"Chapter Arrays","position":1,"concepts":[{"slug":"dsa-arrays-concept","title":"Concept Arrays","position":1}]}]}]}]}
`
	if got := rec.Body.String(); got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func TestHandlerGetTree_EmptyIsArrayNotNull(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(fakeTreeService{topics: []domain.Topic{}}, logger)

	rec := doGetTree(h, http.MethodGet)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != `{"tracks":[]}`+"\n" {
		t.Fatalf("body = %q, want %q", got, `{"tracks":[]}`+"\n")
	}
}

func TestHandlerGetTree_RepositoryErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	h := NewHandler(fakeTreeService{err: errors.New("mysql: connection refused on 10.0.0.5")}, logger)

	rec := doGetTree(h, http.MethodGet)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	body := rec.Body.String()
	if body == "" {
		t.Fatal("expected a JSON error body, got empty")
	}
	for _, leak := range []string{"mysql", "10.0.0.5", "connection refused"} {
		if strings.Contains(body, leak) {
			t.Fatalf("body = %q leaks internal error detail %q", body, leak)
		}
	}
	if !strings.Contains(logBuf.String(), "connection refused") {
		t.Fatalf("log output = %q, want it to contain the real error", logBuf.String())
	}
}

func TestHandlerRejectsNonGET(t *testing.T) {
	logger, _ := newTestLogger()
	h := NewHandler(fakeTreeService{topics: []domain.Topic{}}, logger)

	rec := doGetTree(h, http.MethodPost)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func doGetTree(h *Handler, method string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(method, "/api/v1/curriculum", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func doGetLesson(h *Handler, topicSlug, conceptSlug string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lessons/"+topicSlug+"/"+conceptSlug, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

// TestHandlerGetLesson_Success is the DTO's one behavioral pin that matters
// most for this ticket: expected_answer and mcq options must reach the
// response body, since self-grading happens client-side (see lessonDTO's
// doc comment).
func TestHandlerGetLesson_Success(t *testing.T) {
	logger, _ := newTestLogger()
	lesson := newTestLesson(t, "aggregate", 1)
	h := NewHandler(fakeTreeService{lesson: lesson}, logger)

	rec := doGetLesson(h, "domain-driven-design", "aggregate")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}

	var got lessonDTO
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Topic != "domain-driven-design" || got.Concept != "aggregate" {
		t.Errorf("topic/concept = %q/%q, want domain-driven-design/aggregate", got.Topic, got.Concept)
	}
	if got.TitleEn != lesson.TitleEn() || got.EstMinutes != lesson.EstMinutes().Int() || got.BodyMd != lesson.BodyMd() {
		t.Errorf("lesson fields = %+v, want matching lesson %q", got, lesson.TitleEn())
	}
	if len(got.References) != 2 {
		t.Fatalf("references = %d, want 2", len(got.References))
	}
	if len(got.RecallChecks) != 3 {
		t.Fatalf("recall_checks = %d, want 3", len(got.RecallChecks))
	}
	if got.RecallChecks[0].ExpectedAnswer != "answer 1" {
		t.Errorf("recall_checks[0].expected_answer = %q, want %q (grading is self-graded, the client needs it)", got.RecallChecks[0].ExpectedAnswer, "answer 1")
	}
	if got.RecallChecks[1].Type != "mcq" || len(got.RecallChecks[1].Options) != 2 {
		t.Errorf("recall_checks[1] = %+v, want mcq with 2 options", got.RecallChecks[1])
	}
	if len(got.RecallChecks[0].Options) != 0 {
		t.Errorf("recall_checks[0].options = %v, want omitted for short_answer", got.RecallChecks[0].Options)
	}
}

func TestHandlerGetLesson_NotFound(t *testing.T) {
	logger, _ := newTestLogger()
	notFoundErr := fmt.Errorf("infra: lesson not found: %w", curriculumapp.ErrLessonNotFound)
	h := NewHandler(fakeTreeService{lessonErr: notFoundErr}, logger)

	rec := doGetLesson(h, "domain-driven-design", "aggregates")

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want application/json; charset=utf-8", ct)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("expected a JSON error body, got empty")
	}
}

func TestHandlerGetLesson_RepositoryErrorLeaksNoInternals(t *testing.T) {
	logger, logBuf := newTestLogger()
	h := NewHandler(fakeTreeService{lessonErr: errors.New("mysql: connection refused on 10.0.0.5")}, logger)

	rec := doGetLesson(h, "domain-driven-design", "aggregate")

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

// failingWriter fails every Write, simulating a client that disconnects
// mid-response — the one way to make json.Encoder.Encode itself return an
// error in a test.
type failingWriter struct {
	header http.Header
}

func (w *failingWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *failingWriter) Write([]byte) (int, error) { return 0, errors.New("stub: write failed") }
func (w *failingWriter) WriteHeader(int)           {}

func TestHandlerWriteJSON_LogsEncodeError(t *testing.T) {
	logger, logBuf := newTestLogger()
	h := NewHandler(fakeTreeService{}, logger)

	h.writeJSON(&failingWriter{}, http.StatusOK, map[string]string{"ok": "true"})

	if !strings.Contains(logBuf.String(), "encode response failed") {
		t.Fatalf("log output = %q, want it to mention the encode failure", logBuf.String())
	}
}
