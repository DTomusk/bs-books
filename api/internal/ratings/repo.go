package ratings

import (
	"bs-books-api/internal/db"
	"context"
	"database/sql"
)

type ratingRepo struct{}

func NewRatingRepo() *ratingRepo {
	return &ratingRepo{}
}

func (r *ratingRepo) create(rating *Rating, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, "INSERT INTO ratings (id, user_id, book_id, heart_score, poo_score) VALUES ($1, $2, $3, $4, $5)", rating.ID, rating.UserID, rating.BookID, rating.HeartScore, rating.PooScore)
	return err
}

func (r *ratingRepo) getRatingByUserAndBook(userID, bookID string, ctx context.Context, db db.DBTX) (*Rating, error) {
	var rating Rating
	row := db.QueryRowContext(ctx, "SELECT id, user_id, book_id, heart_score, poo_score FROM ratings WHERE user_id = $1 AND book_id = $2", userID, bookID)
	err := row.Scan(&rating.ID, &rating.UserID, &rating.BookID, &rating.HeartScore, &rating.PooScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rating, nil
}
