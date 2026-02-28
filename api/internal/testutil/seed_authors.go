package testutil

import "database/sql"

func SeedAuthors(tx *sql.Tx) []string {
	tx.Exec(`INSERT INTO authors (id, name) VALUES
		('43681e21-08d4-43e1-b0b6-8d6f75a9b8b1', 'Morgan Bob'),
		('53681e21-08d4-43e1-b0b6-8d6f75a9b8b2', 'Alice Plop')
	`)

	return []string{
		"43681e21-08d4-43e1-b0b6-8d6f75a9b8b1",
		"53681e21-08d4-43e1-b0b6-8d6f75a9b8b2",
	}
}
