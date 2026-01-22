package authors

import (
	"github.com/google/uuid"
)

type Author struct {
	ID             string
	Name           string
	NormalisedName string
}

func NewAuthor(name string) *Author {
	normalisedName := normaliseAuthorName(name)
	return &Author{
		ID:             uuid.New().String(),
		Name:           name,
		NormalisedName: normalisedName,
	}
}
