package infra

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/themethaithian/self-learning/internal/curriculum/domain"
)

type conceptFileDTO struct {
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Outline  string `json:"outline"`
	Position int    `json:"position"`
}

type chapterFileDTO struct {
	Slug     string           `json:"slug"`
	Title    string           `json:"title"`
	Position int              `json:"position"`
	Concepts []conceptFileDTO `json:"concepts"`
}

type topicFileDTO struct {
	Track    string           `json:"track"`
	Slug     string           `json:"slug"`
	Title    string           `json:"title"`
	Position int              `json:"position"`
	Chapters []chapterFileDTO `json:"chapters"`
}

// FileLoader adapts LoadTopic to app.Loader.
type FileLoader struct{}

func (FileLoader) LoadTopic(path string) (domain.Topic, error) { return LoadTopic(path) }

// LoadTopic reads one curriculum content file (the shape documented in
// content/curriculum/ddd.json) and builds the aggregate through the domain
// constructors, so a bad file fails the same way a bad DB row would.
func LoadTopic(path string) (domain.Topic, error) {
	f, err := os.Open(path)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("infra: load curriculum file %s: %w", path, err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	var dto topicFileDTO
	if err := dec.Decode(&dto); err != nil {
		return domain.Topic{}, fmt.Errorf("infra: load curriculum file %s: decode: %w", path, err)
	}
	if dec.More() {
		return domain.Topic{}, fmt.Errorf("infra: load curriculum file %s: unexpected trailing content after JSON value", path)
	}

	topic, err := dto.toDomain()
	if err != nil {
		return domain.Topic{}, fmt.Errorf("infra: load curriculum file %s: %w", path, err)
	}
	return topic, nil
}

func (d topicFileDTO) toDomain() (domain.Topic, error) {
	track, err := domain.NewTrack(d.Track)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("track: %w", err)
	}
	slug, err := domain.NewSlug(d.Slug)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("topic: %w", err)
	}
	position, err := domain.NewPosition(d.Position)
	if err != nil {
		return domain.Topic{}, fmt.Errorf("topic %q: position: %w", d.Slug, err)
	}

	chapters := make([]domain.Chapter, 0, len(d.Chapters))
	for _, chDTO := range d.Chapters {
		chapter, err := chDTO.toDomain()
		if err != nil {
			return domain.Topic{}, fmt.Errorf("topic %q: %w", d.Slug, err)
		}
		chapters = append(chapters, chapter)
	}

	return domain.NewTopic(track, slug, d.Title, position, chapters)
}

func (d chapterFileDTO) toDomain() (domain.Chapter, error) {
	slug, err := domain.NewSlug(d.Slug)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("chapter: %w", err)
	}
	position, err := domain.NewPosition(d.Position)
	if err != nil {
		return domain.Chapter{}, fmt.Errorf("chapter %q: position: %w", d.Slug, err)
	}

	concepts := make([]domain.Concept, 0, len(d.Concepts))
	for _, coDTO := range d.Concepts {
		concept, err := coDTO.toDomain()
		if err != nil {
			return domain.Chapter{}, fmt.Errorf("chapter %q: %w", d.Slug, err)
		}
		concepts = append(concepts, concept)
	}

	return domain.NewChapter(slug, d.Title, position, concepts)
}

func (d conceptFileDTO) toDomain() (domain.Concept, error) {
	slug, err := domain.NewSlug(d.Slug)
	if err != nil {
		return domain.Concept{}, fmt.Errorf("concept: %w", err)
	}
	position, err := domain.NewPosition(d.Position)
	if err != nil {
		return domain.Concept{}, fmt.Errorf("concept %q: position: %w", d.Slug, err)
	}
	return domain.NewConcept(slug, d.Title, d.Outline, position)
}
