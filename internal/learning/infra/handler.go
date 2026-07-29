package infra

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	learningapp "github.com/themethaithian/self-learning/internal/learning/app"
)

// progressService is the minimal capability Handler needs; app.Service
// satisfies it, and tests can fake it without wiring a repository.
type progressService interface {
	ListProgress(ctx context.Context) ([]learningapp.ProgressEntry, error)
	SetProgress(ctx context.Context, topicSlug, conceptSlug, rawState string) (learningapp.ProgressEntry, error)
	RecordAttempt(ctx context.Context, topicSlug, conceptSlug, rawQuestion, rawKind, rawConfidence, rawOutcome string, rawSelectedOption *string) (learningapp.AttemptRecord, error)
}

type Handler struct {
	service progressService
	logger  *slog.Logger
}

func NewHandler(service progressService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/progress", h.getProgress)
	mux.HandleFunc("PUT /api/v1/progress/{topic}/{concept}", h.putProgress)
	mux.HandleFunc("POST /api/v1/progress/{topic}/{concept}/attempts", h.postAttempt)
}

// maxPutBodyBytes caps the PUT body well above any real state value, so a
// client can't force the server to buffer an unbounded request body.
const maxPutBodyBytes = 1024

// maxAttemptBodyBytes bounds the request body. curriculum caps a recall
// question at 1000 runes and an mcq option at 255 runes
// (internal/curriculum/domain/text.go's maxRecallQuestionRunes/
// maxRecallOptionRunes, unexported so cited here rather than imported) — at
// up to 4 bytes/rune in UTF-8, plus JSON string-escaping and field-name
// overhead, 8192 bytes comfortably covers both fields together.
const maxAttemptBodyBytes = 8192

type progressEntryDTO struct {
	Topic         string  `json:"topic"`
	Concept       string  `json:"concept"`
	State         string  `json:"state"`
	LastReadAt    *string `json:"last_read_at"`
	FirstPassedAt *string `json:"first_passed_at"`
}

type progressListResponse struct {
	Concepts []progressEntryDTO `json:"concepts"`
}

func (h *Handler) getProgress(w http.ResponseWriter, r *http.Request) {
	entries, err := h.service.ListProgress(r.Context())
	if err != nil {
		h.logger.Error("learning: list progress failed", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	h.writeJSON(w, http.StatusOK, progressListResponse{Concepts: toProgressEntryDTOs(entries)})
}

type setProgressRequest struct {
	State string `json:"state"`
}

func (h *Handler) putProgress(w http.ResponseWriter, r *http.Request) {
	topic := r.PathValue("topic")
	concept := r.PathValue("concept")

	r.Body = http.MaxBytesReader(w, r.Body, maxPutBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var body setProgressRequest
	if err := dec.Decode(&body); err != nil || dec.More() {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request body too large"})
			return
		}
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	entry, err := h.service.SetProgress(r.Context(), topic, concept, body.State)
	if err != nil {
		switch {
		case errors.Is(err, learningapp.ErrLessonNotFound):
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "lesson not found"})
		case errors.Is(err, learningapp.ErrInvalidState):
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid state"})
		case errors.Is(err, learningapp.ErrInvalidSlug):
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid topic or concept"})
		default:
			h.logger.Error("learning: set progress failed", "topic", topic, "concept", concept, "error", err)
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return
	}
	h.writeJSON(w, http.StatusOK, toProgressEntryDTO(entry))
}

type recordAttemptRequest struct {
	Question       string  `json:"question"`
	Kind           string  `json:"kind"`
	Confidence     string  `json:"confidence"`
	Outcome        string  `json:"outcome"`
	SelectedOption *string `json:"selected_option"`
}

type attemptRecordDTO struct {
	CheckKey       string  `json:"check_key"`
	Question       string  `json:"question"`
	Kind           string  `json:"kind"`
	Confidence     string  `json:"confidence"`
	Outcome        string  `json:"outcome"`
	SelectedOption *string `json:"selected_option"`
	GradedBy       string  `json:"graded_by"`
	CreatedAt      string  `json:"created_at"`
}

func (h *Handler) postAttempt(w http.ResponseWriter, r *http.Request) {
	topic := r.PathValue("topic")
	concept := r.PathValue("concept")

	r.Body = http.MaxBytesReader(w, r.Body, maxAttemptBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var body recordAttemptRequest
	if err := dec.Decode(&body); err != nil || dec.More() {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			h.writeJSON(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request body too large"})
			return
		}
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	entry, err := h.service.RecordAttempt(r.Context(), topic, concept, body.Question, body.Kind, body.Confidence, body.Outcome, body.SelectedOption)
	if err != nil {
		switch {
		case errors.Is(err, learningapp.ErrLessonNotFound):
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "lesson not found"})
		case errors.Is(err, learningapp.ErrCheckNotInLesson):
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "question not found in this lesson"})
		case errors.Is(err, learningapp.ErrInvalidSlug):
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid topic or concept"})
		case errors.Is(err, learningapp.ErrInvalidAttempt):
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid recall attempt"})
		default:
			h.logger.Error("learning: record attempt failed", "topic", topic, "concept", concept, "error", err)
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return
	}
	h.writeJSON(w, http.StatusCreated, toAttemptRecordDTO(entry))
}

func toAttemptRecordDTO(e learningapp.AttemptRecord) attemptRecordDTO {
	return attemptRecordDTO{
		CheckKey:       e.CheckKey,
		Question:       e.Question,
		Kind:           e.Kind.String(),
		Confidence:     e.Confidence.String(),
		Outcome:        e.Outcome.String(),
		SelectedOption: e.SelectedOption,
		GradedBy:       e.GradedBy.String(),
		CreatedAt:      e.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func toProgressEntryDTOs(entries []learningapp.ProgressEntry) []progressEntryDTO {
	out := make([]progressEntryDTO, 0, len(entries))
	for _, e := range entries {
		out = append(out, toProgressEntryDTO(e))
	}
	return out
}

func toProgressEntryDTO(e learningapp.ProgressEntry) progressEntryDTO {
	return progressEntryDTO{
		Topic:         e.Topic,
		Concept:       e.Concept,
		State:         e.State.String(),
		LastReadAt:    formatTime(e.LastReadAt),
		FirstPassedAt: formatTime(e.FirstPassedAt),
	}
}

func formatTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.UTC().Format(time.RFC3339)
	return &s
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error("learning: encode response failed", "error", err)
	}
}
