package domain

// maxSlugShapeLen bounds any slug-shaped identifier this package validates.
const maxSlugShapeLen = 100

// IsValidSlugShape reports whether raw has curriculum.Slug's shape —
// lowercase ASCII letters, digits, and single internal hyphens, 1-100
// characters — the shape rule shared by LessonRef (a concept slug) and, at
// the app layer, a topic slug: this bounded context never imports
// curriculum's own Slug type, so the two callers share this instead.
func IsValidSlugShape(raw string) bool {
	n := len(raw)
	if n == 0 || n > maxSlugShapeLen {
		return false
	}
	for i := 0; i < n; i++ {
		c := raw[i]
		switch {
		case c >= 'a' && c <= 'z', c >= '0' && c <= '9':
			continue
		case c == '-':
			if i == 0 || i == n-1 || raw[i-1] == '-' {
				return false
			}
		default:
			return false
		}
	}
	return true
}
