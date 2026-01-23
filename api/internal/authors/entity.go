package authors

import (
	"github.com/google/uuid"
)

type Author struct {
	ID             string
	Name           string
	NormalisedName string
	DuplicateID    *string
}

func NewAuthor(name string) *Author {
	normalisedName := normaliseAuthorName(name)
	return &Author{
		ID:             uuid.New().String(),
		Name:           name,
		NormalisedName: normalisedName,
		DuplicateID:    nil,
	}
}

func NewAuthorWithDuplicate(name string, duplicateID string) *Author {
	normalisedName := normaliseAuthorName(name)
	return &Author{
		ID:             uuid.New().String(),
		Name:           name,
		NormalisedName: normalisedName,
		DuplicateID:    &duplicateID,
	}
}
