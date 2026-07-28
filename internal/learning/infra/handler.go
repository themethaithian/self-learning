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
}

// Handler is the HTTP adapter for the learning bounded context.
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
}

// maxPutBodyBytes caps the PUT body well above any real state value, so a
// client can't force the server to buffer an unbounded request body.
const maxPutBodyBytes = 1024

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
		default:
			h.logger.Error("learning: set progress failed", "topic", topic, "concept", concept, "error", err)
			h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
		return
	}
	h.writeJSON(w, http.StatusOK, toProgressEntryDTO(entry))
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
