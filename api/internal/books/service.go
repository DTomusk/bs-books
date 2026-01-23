package books

import (
	"bs-books-api/internal/authors"
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"context"
	"database/sql"
)

type BooksService struct {
	db            *sql.DB
	repo          *booksRepo
	provider      BooksProvider
	authorService *authors.AuthorsService
}

func NewBooksService(db *sql.DB, repo *booksRepo, provider BooksProvider, authorService *authors.AuthorsService) *BooksService {
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

	authors := extractUniqueAuthors(externalBooks)
	// Get author ids, create authors as needed
	authorNameIDs := s.authorService.ProcessExternalAuthors(authors, ctx)

	for _, externalBook := range externalBooks {
		s.processExternalBook(externalBook, authorNameIDs, ctx)
	}
	return err
}

func (s *BooksService) processExternalBook(externalBook externalBookModel, authorNameIDs map[string]string, ctx context.Context) {
	logger := logging.FromContext(ctx)
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
		return
	}
	book := NewBook(externalBook.Title, authorIDs)

	err := db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		return s.CreateBookWithAuthors(book, tx, ctx)
	})

	if err != nil {
		logger.Error("Failed to create book", "title", externalBook.Title, "error", err)
		return
	}
	logger.Info("Created book", "title", externalBook.Title, "id", book.ID)
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

// The service shouldn't know about how book authors are stored
// But at the same time, the repo shouldn't know about transactions
// The service coordinates the transaction
// We want the book to fail if an author association fails
func (s *BooksService) CreateBookWithAuthors(book *Book, tx *sql.Tx, ctx context.Context) error {
	err := s.repo.createBook(book, ctx, tx)

	if err != nil {
		return err
	}

	err = s.repo.addAuthorsToBook(book.ID, book.AuthorIDs, ctx, tx)
	if err != nil {
		return err
	}

	return err
}
