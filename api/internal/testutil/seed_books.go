package testutil

import "database/sql"

func SeedBooks(tx *sql.Tx) {
	tx.Exec(`INSERT INTO books (id, author_id, title) VALUES
		('23681e21-08d4-43e1-b0b6-8d6f75a9b8b3', '43681e21-08d4-43e1-b0b6-8d6f75a9b8b1', 'Big Fists'),
		('33681e21-08d4-43e1-b0b6-8d6f75a9b8b4', '53681e21-08d4-43e1-b0b6-8d6f75a9b8b2', 'Wow, a Trampoline!')
	`)
}
