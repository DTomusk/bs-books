package books

import "fmt"

var (
	ErrNoAuthorsProvided    = fmt.Errorf("no authors provided for the book")
	ErrNotAllAuthorsPresent = fmt.Errorf("not all authors are present in the system")
)
