package content_moderation

type ReportContentRequest struct {
	ContentID   string `json:"content_id" binding:"required"`
	ContentType string `json:"content_type" binding:"required,oneof=review"`
	Reason      string `json:"reason" binding:"required"`
}
