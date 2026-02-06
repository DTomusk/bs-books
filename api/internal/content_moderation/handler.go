package content_moderation

import (
	"bs-books-api/internal/delivery/response"

	"github.com/gin-gonic/gin"
)

type ContentModerationHandler struct {
	service *ContentModerationService
}

func NewContentModerationHandler(service *ContentModerationService) *ContentModerationHandler {
	return &ContentModerationHandler{
		service: service,
	}
}

// ReportContent godoc
// @Summary Report inappropriate content
// @Description Report a review or user for inappropriate content or behavior
// @Tags Content Moderation
// @Accept json
// @Produce json
// @Param report body ReportContentRequest true "Report content request"
// @Success 200
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /moderation/report [post]
// @Security BearerAuth
func (h *ContentModerationHandler) ReportContent(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, response.ErrUnauthorized)
		return
	}

	var req ReportContentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.ErrInvalidRequest)
		return
	}

	err := h.service.ReportContent(ctx, req.ContentID, req.ContentType, req.Reason, userID.(string))
	if err != nil {
		c.JSON(500, response.NewInternalServerError("Failed to report content"))
		return
	}

	c.JSON(200, response.Ok())
}
