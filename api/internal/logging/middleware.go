package logging

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLoggerMiddleware(c *gin.Context) {
	requestID := uuid.NewString()

	logger := slog.Default().With(
		"request_id", requestID,
		"method", c.Request.Method,
		"path", c.Request.URL.Path,
	)

	ctx := WithLogger(c.Request.Context(), logger)
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
