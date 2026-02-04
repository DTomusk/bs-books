package reviews

import (
	"bs-books-api/internal/db"
	"context"
)

type ReviewRepo struct{}

func NewReviewRepo() *ReviewRepo {
	return &ReviewRepo{}
}

func (r *ReviewRepo) create(review *Review, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO reviews (id, rating_id, review, created_at) VALUES ($1, $2, $3, $4)", review.ID, review.RatingID, review.Text, review.CreatedAt)
	return err
}
