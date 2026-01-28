package books

import "fmt"

var (
	ErrNoAuthorsProvided = fmt.Errorf("no authors provided for the book")
)
