package content_moderation

import "fmt"

var (
	ErrContentElementDoesntExist = fmt.Errorf("the content element doesn't exist")
	ErrAlreadyReported           = fmt.Errorf("the content has already been reported by the user")
	ErrInvalidContentType        = fmt.Errorf("invalid content type")
)
