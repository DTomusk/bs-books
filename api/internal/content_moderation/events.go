package content_moderation

const (
	EventReviewReported = "review.reported"
)

type ReviewReportedEventPayload struct {
	ReviewID string `json:"review_id"`
}
