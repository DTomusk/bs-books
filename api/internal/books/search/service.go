package search

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/queries"
	"context"
	"strconv"
)

type BookSearchService struct {
	db          db.DBTX
	reader      *queries.BookReader
	bookService *books.BooksService
	searchRepo  *BookSearchRepo
}

func NewBookSearchService(
	db db.DBTX,
	reader *queries.BookReader,
	bookService *books.BooksService,
	searchRepo *BookSearchRepo,
) *BookSearchService {
	return &BookSearchService{
		db:          db,
		reader:      reader,
		bookService: bookService,
		searchRepo:  searchRepo,
	}
}

func (s *BookSearchService) SearchBooks(ctx context.Context, query string, page int, pageSize int) (*queries.BookSearchPage, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Searching books", "query", query, "page", strconv.Itoa(page), "pageSize", strconv.Itoa(pageSize))
	offset := (page - 1) * pageSize
	resultPage, err := s.reader.SearchBooksQuery(query, page, pageSize, offset, ctx)
	if err != nil {
		return nil, err
	}

	if page == 1 && (len(resultPage.Items) == 0 || resultPage.Items[0].Similarity < 0.5) {
		// No results on first page, try to extract from external source
		logger.Info("Weak local results, extracting from external source", "query", query)

		// Normalise search query
		normalisedQuery := normaliseSearchQuery(query)

		// Check external search query table to see if we've searched in last 24h
		searchedToday, err := s.searchRepo.QuerySearchedToday(normalisedQuery, ctx, s.db)

		if err != nil {
			logger.Error("Error checking external search query history", "error", err.Error())
			return resultPage, nil
		}
		if searchedToday {
			logger.Info("Query searched recently, skipping external extraction", "query", query)
			return resultPage, nil
		}

		books, err := s.bookService.ExtractExternalBooks(query, ctx)
		if err != nil {
			logger.Error("Error extracting external books", "error", err.Error())
			return resultPage, nil
		}
		if len(books) == 0 {
			logger.Info("No external results found", "query", query)
			return resultPage, nil
		}
		// Either use the books we have and get their authors
		// Or re-run search query to get updated results
		resultPage, err = s.reader.SearchBooksQuery(query, page, pageSize, offset, ctx)
		if err != nil {
			return nil, err
		}

		err = s.searchRepo.LogExternalSearchQuery(normalisedQuery, ctx, s.db)
		if err != nil {
			logger.Error("Error logging external search query", "error", err.Error())
		}
	}

	// Background extracting from external source

	return resultPage, nil
}
