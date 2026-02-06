package reviews

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
)

type ReviewRepo struct{}

func NewReviewRepo() *ReviewRepo {
	return &ReviewRepo{}
}

func (r *ReviewRepo) create(review *Review, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO reviews (id, rating_id, review, created_at) VALUES ($1, $2, $3, $4)", review.ID, review.RatingID, review.Text, review.CreatedAt)
	return err
}

func (r *ReviewRepo) getByID(ctx context.Context, db db.DBTX, reviewID string) (*Review, error) {
	var review Review
	err := db.QueryRowContext(
		ctx,
		"SELECT id, rating_id, review, created_at FROM reviews WHERE id = $1",
		reviewID,
	).Scan(
		&review.ID,
		&review.RatingID,
		&review.Text,
		&review.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}

func (r *ReviewRepo) incrementReportCount(ctx context.Context, db db.DBTX, reviewID string, visibilityThreshold int) error {
	query := `
		UPDATE reviews
		SET report_count = report_count + 1,
		    moderation_status = CASE
		        WHEN report_count + 1 >= $2 THEN 'hidden'
		        ELSE moderation_status
		    END
		WHERE id = $1
	`
	_, err := db.ExecContext(ctx, query, reviewID, visibilityThreshold)
	if err != nil {
		return err
	}
	return nil
}
