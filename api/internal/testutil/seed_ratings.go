package testutil

import (
	"database/sql"

	"github.com/google/uuid"
)

func SeedRatingsAndReviews(tx *sql.Tx, book_id, user_id string, heart_score, poo_score float64) []string {
	rating_id := uuid.NewString()
	rating_query := `INSERT INTO ratings (id, book_id, user_id, heart_score, poo_score) VALUES ($1, $2, $3, $4, $5)`
	tx.Exec(rating_query, rating_id, book_id, user_id, heart_score, poo_score)

	review_id := uuid.NewString()
	review_query := `INSERT INTO reviews (id, rating_id, review) VALUES ($1, $2, $3)`
	tx.Exec(review_query, review_id, rating_id, "Great book!")

	return []string{rating_id}
}
