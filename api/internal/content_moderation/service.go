package content_moderation

import "context"

type ContentModerationService struct{}

func NewContentModerationService() *ContentModerationService {
	return &ContentModerationService{}
}

func (s *ContentModerationService) ReportContent(ctx context.Context, contentID, contentType, reason, userID string) error {
	return nil
}
