package reviews

import "fmt"

var (
	ErrEmptyReviewText   = fmt.Errorf("review text cannot be empty")
	ErrReviewTextTooLong = fmt.Errorf("review text cannot exceed 500 characters")
)
