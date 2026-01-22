package authors

import (
	"bs-books-api/internal/db"
	"context"
)

type AuthorsRepo struct{}

func NewAuthorsRepo() *AuthorsRepo {
	return &AuthorsRepo{}
}

func (r *AuthorsRepo) GetIDByName(name string, ctx context.Context, db db.DBTX) (string, error) {
	return "", nil
}

func (r *AuthorsRepo) GetIDByNormalisedName(normalisedName string, ctx context.Context, db db.DBTX) (string, error) {
	return "", nil
}

func (r *AuthorsRepo) GetIDByAlias(name string, ctx context.Context, db db.DBTX) (string, error) {
	return "", nil
}

func (r *AuthorsRepo) CreateAuthor(author *Author, ctx context.Context, db db.DBTX) error {
	return nil
}

func (r *AuthorsRepo) CreateAuthorAlias(authorID, aliasName string, ctx context.Context, db db.DBTX) error {
	return nil
}
