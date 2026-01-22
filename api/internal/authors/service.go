package authors

import (
	"bs-books-api/internal/db"
	"context"
)

type AuthorsService struct {
	db   db.DBTX
	repo *authorsRepo
}

func NewAuthorsService(db db.DBTX, repo *authorsRepo) *AuthorsService {
	return &AuthorsService{db: db, repo: repo}
}

func (s *AuthorsService) ProcessExternalAuthors(authorNames []string, ctx context.Context) map[string]string {
	namesToIDs := make(map[string]string)
	// For each author:
	for _, name := range authorNames {
		id, err := s.processExternalAuthor(name, ctx)
		if err != nil {
			// TODO: handle error case
			// we probably don't want to fail the whole batch for one error
			continue
		}
		namesToIDs[name] = id
	}
	// Check for exact match in db
	// If exists, return that ID
	// If not, check for exact match in alias table
	// If exists, return that ID
	// If not, normalise name and check against normalised name in author table
	// If exists, add name to alias table and return that author's ID
	// If it is normalised versions are highly similar, create new author entry flagged as possible duplicate
	// Keep track of all new authors to batch insert at the end
	// Return map of author names to IDs
	return nil
}

func (s *AuthorsService) processExternalAuthor(name string, ctx context.Context) (string, error) {
	// Exact name match
	existingID, err := s.repo.getIDByName(name, ctx, s.db)
	if err != nil {
		return "", err
	}
	if existingID != "" {
		return existingID, nil
	}

	// Check aliases
	existingID, err = s.repo.getIDByAlias(name, ctx, s.db)
	if err != nil {
		return "", err
	}
	if existingID != "" {
		return existingID, nil
	}

	// No alias match, check normalised name and add alias if matched
	normalisedName := normaliseAuthorName(name)
	existingID, err = s.repo.getIDByNormalisedName(normalisedName, ctx, s.db)
	if err != nil {
		return "", err
	}
	if existingID != "" {
		err = s.repo.createAuthorAlias(existingID, name, ctx, s.db)
		if err != nil {
			return existingID, err
		}
		return existingID, nil
	}

	author := NewAuthor(name)
	err = s.repo.createAuthor(author, ctx, s.db)
	if err != nil {
		return "", err
	}

	return author.ID, nil
}
