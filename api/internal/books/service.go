package books

import "bs-books-api/internal/db"

type BooksService struct {
	db       db.DBTX
	repo     *booksRepo
	provider BooksProvider
}

func NewBooksService(db db.DBTX, repo *booksRepo, provider BooksProvider) *BooksService {
	return &BooksService{
		db:       db,
		repo:     repo,
		provider: provider,
	}
}

func (s *BooksService) ExtractExternalBooks(query string) error {
	// Query external api and extract books into our db
	books, err := s.provider.SearchBooks(query)
	if err != nil {
		return err
	}

	// Iterate over books and process
	// Return entities that we can insert into our db and batch insert them later
	err = s.processExternalBooks(books)

	return nil
}

func (s *BooksService) processExternalBooks(books []externalBookModel) error {
	// Get unique authors from books
	// Send to author service to create any new entities and return a map of author names to IDs
	// Do the same with books, create books if needed and return their IDs
	// Then, insert book and author ids into junction table
	return nil
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
