package books

import (
	"strings"
)

func NormaliseBookTitle(title string) string {
	title = strings.ToLower(title)
	title = strings.TrimSpace(title)

	replacer := strings.NewReplacer(
		"-", " ",
		":", " ",
		",", " ",
		".", " ",
		"!", " ",
	)
	title = replacer.Replace(title)

	for _, article := range []string{"the", "a", "an"} {
		if strings.HasPrefix(title, article+" ") {
			title = strings.TrimPrefix(title, article+" ")
			break
		}
	}

	return strings.Join(strings.Fields(title), " ")
}
