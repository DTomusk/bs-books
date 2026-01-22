package books

import (
	"bs-books-api/internal/authors"
	"bs-books-api/internal/db"
	"context"
)

type BooksService struct {
	db            db.DBTX
	repo          *booksRepo
	provider      BooksProvider
	authorService *authors.AuthorsService
}

func NewBooksService(db db.DBTX, repo *booksRepo, provider BooksProvider, authorService *authors.AuthorsService) *BooksService {
	return &BooksService{
		db:            db,
		repo:          repo,
		provider:      provider,
		authorService: authorService,
	}
}

func (s *BooksService) ExtractExternalBooks(query string, ctx context.Context) error {
	externalBooks, err := s.provider.SearchBooks(query, ctx)
	if err != nil {
		return err
	}

	err = s.processExternalBooks(externalBooks, ctx)
	return err
}

func (s *BooksService) processExternalBooks(books []externalBookModel, ctx context.Context) error {
	authors := extractUniqueAuthors(books)
	s.authorService.ProcessExternalAuthors(authors, ctx)
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
