package search

import (
	"bs-books-api/internal/queries"
	"context"
)

type BookSearchService struct {
	reader *queries.BookReader
}

func NewBookSearchService(reader *queries.BookReader) *BookSearchService {
	return &BookSearchService{reader: reader}
}

func (s *BookSearchService) SearchBooks(ctx context.Context, query string, page int, pageSize int) (*queries.BookSearchPage, error) {
	offset := (page - 1) * pageSize
	resultPage, err := s.reader.SearchBooksQuery(query, page, pageSize, offset, ctx)
	if err != nil {
		return nil, err
	}

	if page == 1 && len(resultPage.Items) == 0 {
		// No results on first page, try to extract from external source
	}

	return resultPage, nil
}
