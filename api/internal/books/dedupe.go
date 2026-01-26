package books

import "bs-books-api/internal/authors"

func deduplicateExternalBooks(books []externalBookModel) []externalBookModel {
	seen := make(map[string]externalBookModel)
	for _, book := range books {
		bookKey := generateDedupeKey(book)
		if _, exists := seen[bookKey]; !exists {
			seen[bookKey] = book
		}
	}
	deduped := make([]externalBookModel, 0, len(seen))
	for _, book := range seen {
		deduped = append(deduped, book)
	}
	return deduped
}

func generateDedupeKey(book externalBookModel) string {
	normalisedTitle := normaliseBookTitle(book.Title)

	var joinedAuthors string

	for _, author := range book.Authors {
		joinedAuthors += authors.NormaliseAuthorName(author) + ","
	}

	return normalisedTitle + "|" + joinedAuthors
}
