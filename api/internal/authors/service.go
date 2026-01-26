package authors

import (
	"bs-books-api/internal/db"
	"bs-books-api/internal/logging"
	"context"
	"log/slog"
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
	logger := logging.FromContext(ctx)
	namesToIDs := make(map[string]string)
	// For each author:
	for _, name := range authorNames {
		id, err := s.processExternalAuthor(name, ctx, logger)
		if err != nil {
			// Log and skip error, see how often this happens irl
			logger.Error("Failed to process external author", "name", name, "error", err)
			continue
		}
		namesToIDs[name] = id
	}
	return namesToIDs
}

func (s *AuthorsService) processExternalAuthor(name string, ctx context.Context, logger *slog.Logger) (string, error) {
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
	normalisedName := NormaliseAuthorName(name)

	// Search for similar normalised names
	similarAuthor, err := s.repo.searchByNormalisedName(normalisedName, s.similarityThreshold, ctx, s.db)
	if err != nil {
		return "", err
	}

	if similarAuthor == nil {
		// No similar normalised name, create new author
		author := NewAuthor(name)
		logger.Info("Creating new author", "name", name, "id", author.ID)
		err = s.repo.createAuthor(author, ctx, s.db)
		if err != nil {
			return "", err
		}
		return author.ID, nil
	}

	// Similar normalised name found
	// Check for exact match
	if similarAuthor.NormalisedName == normalisedName {
		logger.Info("Adding alias for existing author", "name", name, "id", similarAuthor.ID)
		err = s.repo.createAuthorAlias(similarAuthor.ID, name, ctx, s.db)
		if err != nil {
			return "", err
		}
		return similarAuthor.ID, nil
	}

	// Similar but not exact, create author and flag as possible duplicate
	author := NewAuthorWithDuplicate(name, similarAuthor.ID)
	logger.Info("Creating author as potential duplicate", "name", name, "id", author.ID, "duplicate_of", similarAuthor.ID)
	err = s.repo.createAuthor(author, ctx, s.db)
	if err != nil {
		return "", err
	}
	return author.ID, nil
}
