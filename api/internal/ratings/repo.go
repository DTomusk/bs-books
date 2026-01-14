package ratings

import (
	"bs-books-api/internal/db"
	"context"
)

type ratingRepo struct{}

func NewRatingRepo() *ratingRepo {
	return &ratingRepo{}
}

func (r *ratingRepo) create(rating *Rating, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO ratings (id, book_id, heart_score, poo_score) VALUES ($1, $2, $3, $4)", rating.ID, rating.BookID, rating.HeartScore, rating.PooScore)
	return err
}
