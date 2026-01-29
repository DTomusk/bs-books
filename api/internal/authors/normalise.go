package authors

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

func NormaliseAuthorName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))

	t := norm.NFD.String(name)

	var b strings.Builder
	b.Grow(len(t))

	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}

		if unicode.IsPunct(r) {
			b.WriteRune(' ')
			continue
		}
		b.WriteRune(r)
	}

	return strings.Join(strings.Fields(b.String()), " ")
}
