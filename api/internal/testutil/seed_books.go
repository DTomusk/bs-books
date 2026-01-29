package testutil

import "database/sql"

func SeedBooks(tx *sql.Tx) {
	tx.Exec(`INSERT INTO books (id, title, normalised_title) VALUES
		('23681e21-08d4-43e1-b0b6-8d6f75a9b8b3', 'Big Fists', 'big fists'),
		('33681e21-08d4-43e1-b0b6-8d6f75a9b8b4', 'Wow, a Trampoline!', 'wow a trampoline')
	`)

	tx.Exec(`INSERT INTO book_author (book_id, author_id) VALUES
		('23681e21-08d4-43e1-b0b6-8d6f75a9b8b3', '43681e21-08d4-43e1-b0b6-8d6f75a9b8b1'),
		('33681e21-08d4-43e1-b0b6-8d6f75a9b8b4', '53681e21-08d4-43e1-b0b6-8d6f75a9b8b2')
	`)
}
