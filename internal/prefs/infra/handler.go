package infra

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/themethaithian/self-learning/internal/prefs/domain"
)

// focusTrackService is the minimal capability Handler needs; app.Service
// satisfies it, and tests can fake it without wiring a repository.
type focusTrackService interface {
	FocusTrack(ctx context.Context) (value string, ok bool, err error)
	SetFocusTrack(ctx context.Context, raw *string) error
}

// Handler is the HTTP adapter for the prefs bounded context.
type Handler struct {
	service focusTrackService
	logger  *slog.Logger
}

func NewHandler(service focusTrackService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/prefs/focus-track", h.getFocusTrack)
	mux.HandleFunc("PUT /api/v1/prefs/focus-track", h.putFocusTrack)
}

type focusTrackDTO struct {
	Track *string `json:"track"`
}

// maxPutBodyBytes caps the PUT body well above any real slug, so a client
// can't force the server to buffer an unbounded request body.
const maxPutBodyBytes = 1024

func (h *Handler) getFocusTrack(w http.ResponseWriter, r *http.Request) {
	value, ok, err := h.service.FocusTrack(r.Context())
	if err != nil {
		h.logger.Error("prefs: get focus track failed", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	var track *string
	if ok {
		track = &value
	}
	h.writeJSON(w, http.StatusOK, focusTrackDTO{Track: track})
}

func (h *Handler) putFocusTrack(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxPutBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var body focusTrackDTO
	if err := dec.Decode(&body); err != nil || dec.More() {
		h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.SetFocusTrack(r.Context(), body.Track); err != nil {
		if errors.Is(err, domain.ErrInvalidFocusTrack) {
			h.writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid track format"})
			return
		}
		h.logger.Error("prefs: set focus track failed", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	h.writeJSON(w, http.StatusOK, focusTrackDTO{Track: body.Track})
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error("prefs: encode response failed", "error", err)
	}
}
