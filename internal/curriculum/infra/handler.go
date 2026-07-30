package infra

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	curriculumapp "github.com/themethaithian/self-learning/internal/curriculum/app"
	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

// treeService and lessonService are the minimal capabilities Handler needs;
// app.Service satisfies both, and tests can fake either without wiring a
// repository.
type treeService interface {
	Tree(ctx context.Context) ([]domain.Topic, map[curriculumapp.ConceptPath]int, error)
}

type lessonService interface {
	Lesson(ctx context.Context, topicSlug, conceptSlug string) (domain.Lesson, error)
}

type curriculumService interface {
	treeService
	lessonService
}

// Handler is the HTTP adapter for the curriculum bounded context.
type Handler struct {
	service curriculumService
	logger  *slog.Logger
}

func NewHandler(service curriculumService, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/curriculum", h.getTree)
	mux.HandleFunc("GET /api/v1/lessons/{topicSlug}/{conceptSlug}", h.getLesson)
}

type conceptDTO struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Position int    `json:"position"`

	HasLesson  bool `json:"has_lesson"`
	EstMinutes *int `json:"est_minutes"`
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
	topics, availability, err := h.service.Tree(r.Context())
	if err != nil {
		h.logger.Error("curriculum: get tree failed", "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	h.writeJSON(w, http.StatusOK, treeResponse{Tracks: groupByTrack(topics, availability)})
}

// groupByTrack presents the flat, already (position, slug)-sorted topic list
// grouped by track, in the domain's canonical track order — topics.position
// is not scoped per track, so a flat cross-track order is meaningless to a
// client. Filtering preserves the incoming order, so no re-sort is needed.
func groupByTrack(topics []domain.Topic, availability map[curriculumapp.ConceptPath]int) []trackDTO {
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
		out = append(out, trackDTO{Track: key, Topics: toTopicDTOs(grouped, availability)})
	}
	return out
}

func toTopicDTOs(topics []domain.Topic, availability map[curriculumapp.ConceptPath]int) []topicDTO {
	out := make([]topicDTO, 0, len(topics))
	for _, t := range topics {
		out = append(out, topicDTO{
			Slug:     t.Slug().String(),
			Title:    t.Title(),
			Position: t.Position().Int(),
			Chapters: toChapterDTOs(t.Slug().String(), t.Chapters(), availability),
		})
	}
	return out
}

func toChapterDTOs(topicSlug string, chapters []domain.Chapter, availability map[curriculumapp.ConceptPath]int) []chapterDTO {
	out := make([]chapterDTO, 0, len(chapters))
	for _, c := range chapters {
		out = append(out, chapterDTO{
			Slug:     c.Slug().String(),
			Title:    c.Title(),
			Position: c.Position().Int(),
			Concepts: toConceptDTOs(topicSlug, c.Concepts(), availability),
		})
	}
	return out
}

func toConceptDTOs(topicSlug string, concepts []domain.Concept, availability map[curriculumapp.ConceptPath]int) []conceptDTO {
	out := make([]conceptDTO, 0, len(concepts))
	for _, c := range concepts {
		dto := conceptDTO{
			Slug:     c.Slug().String(),
			Title:    c.Title(),
			Position: c.Position().Int(),
		}
		path := curriculumapp.ConceptPath{TopicSlug: topicSlug, ConceptSlug: c.Slug().String()}
		if minutes, ok := availability[path]; ok {
			dto.HasLesson = true
			dto.EstMinutes = &minutes
		}
		out = append(out, dto)
	}
	return out
}

type referenceDTO struct {
	Title  string `json:"title"`
	Source string `json:"source"`
	Why    string `json:"why"`
}

type recallCheckDTO struct {
	Position       int      `json:"position"`
	Type           string   `json:"type"`
	Question       string   `json:"question"`
	ExpectedAnswer string   `json:"expected_answer"`
	Options        []string `json:"options,omitempty"`
	Explanation    string   `json:"explanation,omitempty"`
}

// lessonDTO carries each recall check's expected_answer, options, and
// explanation: grading is self-graded client-side (the user reveals the
// answer and rates themselves), and single-user bearer auth means there is
// no cheating concern in exposing any of them over the API.
type lessonDTO struct {
	Topic        string           `json:"topic"`
	Concept      string           `json:"concept"`
	TitleEn      string           `json:"title_en"`
	EstMinutes   int              `json:"est_minutes"`
	BodyMd       string           `json:"body_md"`
	References   []referenceDTO   `json:"references"`
	RecallChecks []recallCheckDTO `json:"recall_checks"`
}

func (h *Handler) getLesson(w http.ResponseWriter, r *http.Request) {
	topicSlug := r.PathValue("topicSlug")
	conceptSlug := r.PathValue("conceptSlug")

	lesson, err := h.service.Lesson(r.Context(), topicSlug, conceptSlug)
	if err != nil {
		if errors.Is(err, curriculumapp.ErrLessonNotFound) {
			h.writeJSON(w, http.StatusNotFound, map[string]string{"error": "lesson not found"})
			return
		}
		h.logger.Error("curriculum: get lesson failed", "topic", topicSlug, "concept", conceptSlug, "error", err)
		h.writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	h.writeJSON(w, http.StatusOK, toLessonDTO(topicSlug, conceptSlug, lesson))
}

func toLessonDTO(topicSlug, conceptSlug string, l domain.Lesson) lessonDTO {
	return lessonDTO{
		Topic:        topicSlug,
		Concept:      conceptSlug,
		TitleEn:      l.TitleEn(),
		EstMinutes:   l.EstMinutes().Int(),
		BodyMd:       l.BodyMd(),
		References:   toReferenceDTOs(l.References()),
		RecallChecks: toRecallCheckDTOs(l.RecallChecks()),
	}
}

func toReferenceDTOs(refs []domain.Reference) []referenceDTO {
	out := make([]referenceDTO, 0, len(refs))
	for _, ref := range refs {
		out = append(out, referenceDTO{Title: ref.Title(), Source: ref.Source(), Why: ref.Why()})
	}
	return out
}

func toRecallCheckDTOs(checks []domain.RecallCheck) []recallCheckDTO {
	out := make([]recallCheckDTO, 0, len(checks))
	for _, c := range checks {
		out = append(out, recallCheckDTO{
			Position:       c.Position().Int(),
			Type:           c.Kind().String(),
			Question:       c.Question(),
			ExpectedAnswer: c.ExpectedAnswer(),
			Options:        c.Options(),
			Explanation:    c.Explanation(),
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
