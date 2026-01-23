package authors

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
)

type authorsRepo struct{}

func NewAuthorsRepo() *authorsRepo {
	return &authorsRepo{}
}

type authorRow struct {
	ID             string
	Name           string
	NormalisedName string
	DuplicateID    sql.NullString
}

func (r *authorsRepo) getIDByName(name string, ctx context.Context, db db.DBTX) (string, error) {
	var id string
	row := db.QueryRowContext(ctx, `SELECT id FROM authors WHERE name = $1`, name)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (r *authorsRepo) getIDByNormalisedName(normalisedName string, ctx context.Context, db db.DBTX) (string, error) {
	var id string
	row := db.QueryRowContext(ctx, `SELECT id FROM authors WHERE normalised_name = $1`, normalisedName)
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return id, nil
}

func (r *authorsRepo) getIDByAlias(name string, ctx context.Context, db db.DBTX) (string, error) {
	var authorID string
	row := db.QueryRowContext(ctx, `SELECT author_id FROM author_alias WHERE alias = $1`, name)
	err := row.Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return authorID, nil
}

func (r *authorsRepo) createAuthor(author *Author, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO authors (id, name, normalised_name) VALUES ($1, $2, $3)", author.ID, author.Name, author.NormalisedName)
	return err
}

func (r *authorsRepo) createAuthorWithDuplicate(author *Author, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO authors (id, name, normalised_name, duplicate_id) VALUES ($1, $2, $3, $4)", author.ID, author.Name, author.NormalisedName, author.DuplicateID)
	return err
}

func (r *authorsRepo) createAuthorAlias(authorID, aliasName string, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO author_alias (author_id, alias) VALUES ($1, $2)", authorID, aliasName)
	return err
}

func (r *authorsRepo) searchByNormalisedName(normalisedName string, ctx context.Context, db db.DBTX) (*Author, error) {
	query := `
		SELECT 
			id,
			name,
			normalised_name,
			similarity(normalised_name, $1) AS similarity
		FROM authors
		WHERE similarity(normalised_name, $1) > $2
		ORDER BY similarity DESC
		LIMIT 1
	`
	var author authorRow
	row := db.QueryRowContext(ctx, query, normalisedName, 0.7)
	// Discard score (TODO: should we store it?)
	err := row.Scan(&author.ID, &author.Name, &author.NormalisedName, new(interface{}))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapAuthorRowToEntity(author), nil
}

func mapAuthorRowToEntity(row authorRow) *Author {
	var duplicateID *string
	if row.DuplicateID.Valid {
		duplicateID = &row.DuplicateID.String
	} else {
		duplicateID = nil
	}
	return &Author{
		ID:             row.ID,
		Name:           row.Name,
		NormalisedName: row.NormalisedName,
		DuplicateID:    duplicateID,
	}
}
