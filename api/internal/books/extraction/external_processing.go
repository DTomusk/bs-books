package extraction

import (
	"bs-books-api/internal/books"
	"log/slog"
)

func createBookFromExternal(externalBook externalBookModel, authorNameIDs map[string]string, logger *slog.Logger) (*books.Book, error) {
	allAuthorsPresent := true
	authorIDs := make([]string, 0, len(externalBook.Authors))
	for _, authorName := range externalBook.Authors {
		authorID, exists := authorNameIDs[authorName]
		if !exists {
			allAuthorsPresent = false
			break
		}
		authorIDs = append(authorIDs, authorID)
	}
	if !allAuthorsPresent {
		return nil, ErrNotAllAuthorsPresent
	}
	book, err := books.NewBook(externalBook.Title, authorIDs)
	if err != nil {
		logger.Error("Failed to create book entity", "title", externalBook.Title, "error", err)
		return nil, err
	}
	return book, nil
}

func extractUniqueAuthors(books []externalBookModel) []string {
	seen := make(map[string]struct{})

	for _, book := range books {
		for _, author := range book.Authors {
			seen[author] = struct{}{}
		}
	}

	uniqueAuthors := make([]string, 0, len(seen))
	for author := range seen {
		uniqueAuthors = append(uniqueAuthors, author)
	}

	return uniqueAuthors
}
