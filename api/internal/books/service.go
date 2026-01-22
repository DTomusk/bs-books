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
	_, err := s.provider.SearchBooks(query)
	if err != nil {
		return err
	}

	return nil
}
