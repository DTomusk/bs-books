package content_moderation

const (
	EventReviewReported = "review.reported"
	EventUserReported   = "user.reported"
)

type ReviewReportedEventPayload struct {
	ReviewID string `json:"review_id"`
}

type UserReportedEventPayload struct {
	UserID string `json:"user_id"`
}
