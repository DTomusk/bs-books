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

func (r *ContentModerationRepo) CreateReport(ctx context.Context, db db.DBTX, report *ContentModerationReport) error {
	const query = `
	INSERT INTO moderation_reports (id, reporter_id, content_id, content_type, content_snapshot, reason, created_at, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := db.ExecContext(ctx, query, report.ID, report.UserID, report.ContentID, report.ContentType, report.ContentSnapshot, report.Reason, report.CreatedAt, report.Status)
	return err
}
