package content_moderation

import "time"

// TODO: replace strings with const values
type ContentModerationReport struct {
	ID              string
	UserID          string
	ContentID       string
	Reason          string
	ContentType     string
	ContentSnapshot string
	CreatedAt       time.Time
	Status          string
}
