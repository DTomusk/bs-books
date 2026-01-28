package extraction

import (
	"bs-books-api/internal/authors"
	"bs-books-api/internal/books"
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"context"
)

type BookExtractionService struct {
	db            db.DBTX
	provider      BooksProvider
	authorService *authors.AuthorsService
}

func NewBookExtractionService(
	db db.DBTX,
	provider BooksProvider,
	authorService *authors.AuthorsService,
) *BookExtractionService {
	return &BookExtractionService{
		db:            db,
		provider:      provider,
		authorService: authorService,
	}
}

func (s *BookExtractionService) ExtractExternalBooks(query string, ctx context.Context) ([]*books.Book, error) {
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

	extractedBooks := make([]*books.Book, 0, len(externalBooks))

	for _, externalBook := range externalBooks {
		book, err := s.processExternalBook(externalBook, authorNameIDs, ctx)
		if err == nil {
			extractedBooks = append(extractedBooks, book)
		}
	}
	return extractedBooks, nil
}

func (s *BookExtractionService) processExternalBook(externalBook externalBookModel, authorNameIDs map[string]string, ctx context.Context) (*books.Book, error) {
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

	return book, nil
}
