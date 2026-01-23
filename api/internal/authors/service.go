package authors

import (
	"bs-books-api/internal/db"
	"context"
)

type AuthorsService struct {
	db                  db.DBTX
	repo                *authorsRepo
	similarityThreshold float64
}

func NewAuthorsService(db db.DBTX, repo *authorsRepo, similarityThreshold float64) *AuthorsService {
	return &AuthorsService{db: db, repo: repo, similarityThreshold: similarityThreshold}
}

// TODO: batching if optimisation necessary
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
	// Keep track of all new authors to batch insert at the end
	// Return map of author names to IDs
	return namesToIDs
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

	// Check aliases, return id if exact match
	existingID, err = s.repo.getIDByAlias(name, ctx, s.db)
	if err != nil {
		return "", err
	}
	if existingID != "" {
		return existingID, nil
	}

	// No alias match, check normalised name and add alias if matched
	normalisedName := normaliseAuthorName(name)

	// Search for similar normalised names
	similarAuthor, err := s.repo.searchByNormalisedName(normalisedName, s.similarityThreshold, ctx, s.db)
	if err != nil {
		return "", err
	}

	if similarAuthor == nil {
		// No similar normalised name, create new author
		author := NewAuthor(name)
		err = s.repo.createAuthor(author, ctx, s.db)
		if err != nil {
			return "", err
		}
		return author.ID, nil
	}

	// Similar normalised name found
	// Check for exact match
	if similarAuthor.NormalisedName == normalisedName {
		err = s.repo.createAuthorAlias(similarAuthor.ID, name, ctx, s.db)
		if err != nil {
			return "", err
		}
		return similarAuthor.ID, nil
	}

	// Similar but not exact, create author and flag as possible duplicate
	author := NewAuthorWithDuplicate(name, similarAuthor.ID)
	err = s.repo.createAuthor(author, ctx, s.db)
	if err != nil {
		return "", err
	}
	return author.ID, nil
}
