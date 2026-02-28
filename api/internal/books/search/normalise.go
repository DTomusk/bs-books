package search

import (
	"strings"
)

// Predictively normalise search query to limit external searches
func normaliseSearchQuery(query string) string {
	query = strings.ToLower(query)
	query = strings.TrimSpace(query)
	replacer := strings.NewReplacer(
		"-", " ",
		":", " ",
		",", " ",
		".", " ",
		"!", " ",
	)
	query = replacer.Replace(query)
	return strings.Join(strings.Fields(query), " ")
}
