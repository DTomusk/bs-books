package search

import (
	"bs-books-api/internal/books"
	"bs-books-api/internal/logging"
	"bs-books-api/internal/queries"
	"context"
	"strconv"
)

type BookSearchService struct {
	reader      *queries.BookReader
	bookService *books.BooksService
}

func NewBookSearchService(reader *queries.BookReader, bookService *books.BooksService) *BookSearchService {
	return &BookSearchService{reader: reader, bookService: bookService}
}

func (s *BookSearchService) SearchBooks(ctx context.Context, query string, page int, pageSize int) (*queries.BookSearchPage, error) {
	logger := logging.FromContext(ctx)
	logger.Info("Searching books", "query", query, "page", strconv.Itoa(page), "pageSize", strconv.Itoa(pageSize))
	offset := (page - 1) * pageSize
	resultPage, err := s.reader.SearchBooksQuery(query, page, pageSize, offset, ctx)
	if err != nil {
		return nil, err
	}

	if page == 1 && len(resultPage.Items) == 0 {
		// No results on first page, try to extract from external source
		logger.Info("No local results, extracting from external source", "query", query)
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
	}

	// Background extracting from external source

	return resultPage, nil
}
