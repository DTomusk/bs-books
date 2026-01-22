package authors

import (
	"github.com/google/uuid"
)

type Author struct {
	ID   string
	Name string
}

func NewAuthor(name string) *Author {
	return &Author{
		ID:   uuid.New().String(),
		Name: name,
	}
}
