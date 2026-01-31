package testutil

import "database/sql"

func SeedRatingsAndReviews(tx *sql.Tx) []string {
	tx.Exec(`INSERT INTO ratings (id, book_id, user_id, score, difficulty, comment) VALUES
		('63681e21-08d4-43e1-b0b6-8d6f75a9b8b3', '73681e21-08d4-43e1-b0b6-8d6f75a9b8b1', '83681e21-08d4-43e1-b0b6-8d6f75a9b8b1', 4.5, 3.0, 'Great book!'),
		('73681e21-08d4-43e1-b0b6-8d6f75a9b8b4', '73681e21-08d4-43e1-b0b6-8d6f75a9b8b1', '93681e21-08d4-43e1-b0b6-8d6f75a9b8b2', 3.0, 4.0, 'Good read.')`)

	tx.Exec(`INSERT INTO reviews (id, rating_id, review) VALUES
		('93681e21-08d4-43e1-b0b6-8d6f75a9b8b5', '63681e21-08d4-43e1-b0b6-8d6f75a9b8b3', 'Pee pee poo poo'),
		('a3681e21-08d4-43e1-b0b6-8d6f75a9b8b6', '73681e21-08d4-43e1-b0b6-8d6f75a9b8b4', 'I did not like it.')`)
	return []string{
		"63681e21-08d4-43e1-b0b6-8d6f75a9b8b3",
		"73681e21-08d4-43e1-b0b6-8d6f75a9b8b4",
	}
}
