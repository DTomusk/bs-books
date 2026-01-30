package ratings

import "fmt"

var (
	ErrNegativeScore = fmt.Errorf("scores must be non-negative")
	ErrLargeScore    = fmt.Errorf("scores must not exceed 5.0")
	ErrBookNotFound  = fmt.Errorf("book not found")
)
