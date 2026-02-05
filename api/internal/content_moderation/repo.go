package content_moderation

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
)

type ContentModerationRepo struct{}

func NewContentModerationRepo() *ContentModerationRepo {
	return &ContentModerationRepo{}
}

// TODO: consider adding sql row shape and mapping for entity
func (r *ContentModerationRepo) GetReportByUserByContentID(ctx context.Context, db db.DBTX, userID, contentID string) (*ContentModerationReport, error) {
	var contentReport ContentModerationReport
	const query = `
	SELECT 
		id, 
		reporter_id, 
		reason,
		content_type, 
		content_id, 
		content_snapshot,
		created_at,
		status
	FROM moderation_reports mr
	WHERE mr.reporter_id = $1
	AND content_id = $2
	`

	row := db.QueryRowContext(ctx, query, userID, contentID)
	err := row.Scan(
		&contentReport.ID,
		&contentReport.UserID,
		&contentReport.Reason,
		&contentReport.ContentType,
		&contentReport.ContentID,
		&contentReport.ContentSnapshot,
		&contentReport.CreatedAt,
		&contentReport.Status,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &contentReport, nil
}

func (r *ContentModerationRepo) CreateReport(ctx context.Context, db db.DBTX, userID, contentID, contentType, reason string) error {
	const query = `
	INSERT INTO moderation_reports (reporter_id, content_id, content_type, reason, created_at, status) VALUES ($1, $2, $3, $4, NOW(), 'pending_review')
	`

	_, err := db.ExecContext(ctx, query, userID, contentID, contentType, reason)
	return err
}
