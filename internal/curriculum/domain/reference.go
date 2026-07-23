package domain

import "fmt"

// Reference is a pointer to a primary source (book chapter, official docs,
// paper) a Lesson cites for going deeper.
type Reference struct {
	title  string
	source string
	why    string
}

func NewReference(title, source, why string) (Reference, error) {
	trimmedTitle, err := validateReferenceField(title, maxReferenceTitleRunes)
	if err != nil {
		return Reference{}, fmt.Errorf("curriculum: reference: title: %w", err)
	}
	trimmedSource, err := validateReferenceField(source, maxReferenceSourceRunes)
	if err != nil {
		return Reference{}, fmt.Errorf("curriculum: reference %q: source: %w", trimmedTitle, err)
	}
	trimmedWhy, err := validateReferenceField(why, maxReferenceWhyRunes)
	if err != nil {
		return Reference{}, fmt.Errorf("curriculum: reference %q: why: %w", trimmedTitle, err)
	}
	return Reference{title: trimmedTitle, source: trimmedSource, why: trimmedWhy}, nil
}

func (r Reference) Title() string  { return r.title }
func (r Reference) Source() string { return r.source }
func (r Reference) Why() string    { return r.why }

// IsZero reports whether r was never constructed via NewReference.
func (r Reference) IsZero() bool { return r.title == "" }
