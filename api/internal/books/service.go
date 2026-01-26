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

// Fetches books by searching external provider and then saves books and authors in db
func (s *BooksService) ExtractExternalBooks(query string, ctx context.Context) ([]*Book, error) {
	externalBooks, err := s.provider.SearchBooks(query, 10, ctx)
	if err != nil {
		return nil, err
	}

	// Remove obvious duplicates
	externalBooks = deduplicateExternalBooks(externalBooks)

	authors := extractUniqueAuthors(externalBooks)

	if len(authors) == 0 {
		return nil, nil
	}

	// Get author ids, create authors as needed
	// Store regardless of book success (outside of transaction)
	authorNameIDs := s.authorService.ProcessExternalAuthors(authors, ctx)

	if len(authorNameIDs) == 0 {
		return nil, nil
	}

	createdBooks := make([]*Book, 0, len(externalBooks))

	for _, externalBook := range externalBooks {
		book, err := s.processExternalBook(externalBook, authorNameIDs, ctx)
		if err == nil {
			createdBooks = append(createdBooks, book)
		}
	}
	return createdBooks, nil
}

func (s *BooksService) processExternalBook(externalBook externalBookModel, authorNameIDs map[string]string, ctx context.Context) (*Book, error) {
	logger := logging.FromContext(ctx)

	book, err := createBookFromExternal(externalBook, authorNameIDs, logger)

	if err != nil {
		if err == ErrNotAllAuthorsPresent {
			logger.Info("Skipping book creation as not all authors are present", "title", externalBook.Title)
			return nil, err
		}
		logger.Error("Failed to create book from external data", "title", externalBook.Title, "error", err)
		return nil, err
	}

	// Skip book if a book by the same authors with a very similar title exists
	// TODO: inject threshold
	exists, err := s.repo.checkSimilarBookExists(book.Title, book.AuthorIDs, 0.7, ctx, s.db)
	if err != nil {
		logger.Error("Failed to check for similar existing book", "title", book.Title, "error", err)
		return nil, err
	}
	if exists {
		logger.Info("Skipping book creation as similar book already exists", "title", book.Title)
		return nil, nil
	}

	err = db.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		return s.CreateBookWithAuthors(book, tx, ctx)
	})

	if err != nil {
		logger.Error("Failed to create book", "title", externalBook.Title, "error", err)
		return nil, err
	}
	logger.Info("Created book", "title", externalBook.Title, "id", book.ID)
	return book, nil
}

// The service shouldn't know about how book authors are stored
// But at the same time, the repo shouldn't know about transactions
// The service coordinates the transaction
// We want the book to fail if an author association fails
func (s *BooksService) CreateBookWithAuthors(book *Book, tx *sql.Tx, ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.Info("Creating book", "title", book.Title, "id", book.ID, "authors", book.AuthorIDs)
	err := s.repo.createBook(book, ctx, tx)

	if err != nil {
		logger.Error("Failed to create book", "title", book.Title, "error", err)
		return err
	}

	err = s.repo.addAuthorsToBook(book.ID, book.AuthorIDs, ctx, tx)
	if err != nil {
		logger.Error("Failed to add authors to book", "title", book.Title, "error", err)
		return err
	}

	return nil
}
