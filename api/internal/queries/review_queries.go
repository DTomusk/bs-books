package queries

import (
	"bs-books-api/internal/db"
	"context"
)

type ReviewReader struct {
	db db.DBTX
}

func NewReviewReader(db db.DBTX) *ReviewReader {
	return &ReviewReader{db: db}
}

type ReviewPage struct {
	Items      []ReviewItem
	Total      int
	TotalPages int
	Page       int
	Size       int
}

type ReviewItem struct {
	ID         string
	HeartScore float64
	PooScore   float64
	Text       string
	// TODO: we'll want to store a username as well once we have it
	UserID    string
	CreatedAt string
}

// TODO: add created at for ordering
func (r *ReviewReader) GetReviewsByBookIDQuery(ctx context.Context, bookID string, page, pageSize, offset int) (*ReviewPage, error) {
	const reviewsQuery = `
	SELECT
		r.id,
		rating.heart_score,
		rating.poo_score,
		r.review,
		rating.user_id,
		r.created_at
	FROM reviews r
	JOIN ratings rating ON r.rating_id = rating.id
	WHERE rating.book_id = $1
	ORDER BY r.created_at DESC, r.id DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, reviewsQuery, bookID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []ReviewItem

	for rows.Next() {
		var review ReviewItem
		err := rows.Scan(&review.ID, &review.HeartScore, &review.PooScore, &review.Text, &review.UserID, &review.CreatedAt)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	const countQuery = `
	SELECT COUNT(*)
	FROM reviews r
	JOIN ratings rating ON r.rating_id = rating.id
	WHERE rating.book_id = $1
	`

	var total int
	err = r.db.QueryRowContext(ctx, countQuery, bookID).Scan(&total)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &ReviewPage{
		Items:      reviews,
		Page:       page,
		Size:       pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}
