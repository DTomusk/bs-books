package content_moderation

import (
	"time"

	"github.com/google/uuid"
)

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

func NewContentModerationReport(userID, contentID, contentType, contentSnapshot, reason string, status string) *ContentModerationReport {
	return &ContentModerationReport{
		ID:              uuid.NewString(),
		UserID:          userID,
		ContentID:       contentID,
		ContentType:     contentType,
		Reason:          reason,
		ContentSnapshot: contentSnapshot,
		CreatedAt:       time.Now().UTC(),
		Status:          status,
	}
}
