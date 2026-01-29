package search

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/books/extraction"
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/queries"
	"context"
	"log/slog"
	"strconv"
)

type BookSearchService struct {
	db                    db.DBTX
	reader                *queries.BookReader
	bookService           *books.BooksService
	searchRepo            *BookSearchRepo
	bookExtractionService *extraction.BookExtractionService
}

func NewBookSearchService(
	db db.DBTX,
	reader *queries.BookReader,
	bookService *books.BooksService,
	searchRepo *BookSearchRepo,
	bookExtractionService *extraction.BookExtractionService,
) *BookSearchService {
	return &BookSearchService{
		db:                    db,
		reader:                reader,
		bookService:           bookService,
		searchRepo:            searchRepo,
		bookExtractionService: bookExtractionService,
	}
}

func (s *BookSearchService) SearchBooks(ctx context.Context, query string, page int, pageSize int) (*queries.BookSearchPage, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Searching books", "query", query, "page", strconv.Itoa(page), "pageSize", strconv.Itoa(pageSize))
	offset := (page - 1) * pageSize
	normalisedQuery := normaliseSearchQuery(query)
	resultPage, err := s.reader.SearchBooksQuery(normalisedQuery, page, pageSize, offset, ctx)
	if err != nil {
		return nil, err
	}

	if page == 1 && (len(resultPage.Items) == 0 || resultPage.Items[0].Similarity < 0.5) {
		// No results on first page, try to extract from external source
		logger.Info("Weak local results, extracting from external source", "query", query)
		err := s.searchExternalBooks(ctx, normalisedQuery, logger)
		if err != nil {
			logger.Error("Error extracting external books", "error", err.Error())
		}

		// Re-run search after extraction
		resultPage, err = s.reader.SearchBooksQuery(normalisedQuery, page, pageSize, offset, ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// Trigger background extraction for future searches
		// TODO: Use a proper background job queue
		logger.Info("Triggering background extraction of external books", "query", normalisedQuery)
		go func(normalisedQuery string) {
			bgCtx := context.Background()
			bgLogger := logging.FromContext(bgCtx)
			if err := s.searchExternalBooks(bgCtx, normalisedQuery, bgLogger); err != nil {
				bgLogger.Error("Error extracting external books in background", "error", err.Error())
			}
		}(normalisedQuery)
	}

	return resultPage, nil
}

func (s *BookSearchService) searchExternalBooks(ctx context.Context, normalisedQuery string, logger *slog.Logger) error {
	// Check external search query table to see if we've searched in last 24h
	searchedToday, err := s.searchRepo.QuerySearchedToday(normalisedQuery, ctx, s.db)

	if err != nil {
		return err
	}
	if searchedToday {
		logger.Info("Query already searched today", "query", normalisedQuery)
		return nil
	}

	books, err := s.bookExtractionService.ExtractExternalBooks(normalisedQuery, ctx)
	logger.Info("Extracted external books", "count", len(books))
	if err != nil {
		return err
	}

	err = s.bookService.CreateBooksWithAuthors(books, ctx)
	if err != nil {
		return err
	}

	// If there's a problem upstream, the search won't log
	err = s.searchRepo.LogExternalSearchQuery(normalisedQuery, ctx, s.db)
	logger.Info("Logged external search query", "query", normalisedQuery)
	if err != nil {
		return err
	}
	return nil
}
