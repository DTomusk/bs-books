package reviews

import "fmt"

var (
	ErrEmptyReviewText = fmt.Errorf("review text cannot be empty")
)
