package infra

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// treeService is the minimal capability Handler needs; app.Service
// satisfies it, and tests can fake it without wiring a repository.
type treeService interface {
	Tree(ctx context.Context) ([]domain.Topic, error)
}

// Handler is the HTTP adapter for the curriculum bounded context.
type Handler struct {
	service treeService
	logger  *slog.Logger
}

func NewHandler(service treeService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/curriculum", h.getTree)
}

type conceptDTO struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Position int    `json:"position"`
}

type chapterDTO struct {
	Slug     string       `json:"slug"`
	Title    string       `json:"title"`
	Position int          `json:"position"`
	Concepts []conceptDTO `json:"concepts"`
}

type topicDTO struct {
	Slug     string       `json:"slug"`
	Title    string       `json:"title"`
	Position int          `json:"position"`
	Chapters []chapterDTO `json:"chapters"`
}

type trackDTO struct {
	Track  string     `json:"track"`
	Topics []topicDTO `json:"topics"`
}

type treeResponse struct {
	Tracks []trackDTO `json:"tracks"`
}

func (h *Handler) getTree(w http.ResponseWriter, r *http.Request) {
	topics, err := h.service.Tree(r.Context())
	if err != nil {
		h.logger.Error("curriculum: get tree failed", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	h.writeJSON(w, http.StatusOK, treeResponse{Tracks: groupByTrack(topics)})
}

// groupByTrack presents the flat, already (position, slug)-sorted topic list
// grouped by track, in the domain's canonical track order — topics.position
// is not scoped per track, so a flat cross-track order is meaningless to a
// client. Filtering preserves the incoming order, so no re-sort is needed.
func groupByTrack(topics []domain.Topic) []trackDTO {
	byTrack := make(map[string][]domain.Topic, len(topics))
	for _, t := range topics {
		key := t.Track().String()
		byTrack[key] = append(byTrack[key], t)
	}

	out := make([]trackDTO, 0, len(byTrack))
	for _, track := range domain.Tracks() {
		key := track.String()
		grouped, ok := byTrack[key]
		if !ok {
			continue
		}
		out = append(out, trackDTO{Track: key, Topics: toTopicDTOs(grouped)})
	}
	return out
}

func toTopicDTOs(topics []domain.Topic) []topicDTO {
	out := make([]topicDTO, 0, len(topics))
	for _, t := range topics {
		out = append(out, topicDTO{
			Slug:     t.Slug().String(),
			Title:    t.Title(),
			Position: t.Position().Int(),
			Chapters: toChapterDTOs(t.Chapters()),
		})
	}
	return out
}

func toChapterDTOs(chapters []domain.Chapter) []chapterDTO {
	out := make([]chapterDTO, 0, len(chapters))
	for _, c := range chapters {
		out = append(out, chapterDTO{
			Slug:     c.Slug().String(),
			Title:    c.Title(),
			Position: c.Position().Int(),
			Concepts: toConceptDTOs(c.Concepts()),
		})
	}
	return out
}

func toConceptDTOs(concepts []domain.Concept) []conceptDTO {
	out := make([]conceptDTO, 0, len(concepts))
	for _, c := range concepts {
		out = append(out, conceptDTO{
			Slug:     c.Slug().String(),
			Title:    c.Title(),
			Position: c.Position().Int(),
		})
	}
	return out
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error("curriculum: encode response failed", "error", err)
	}
}
