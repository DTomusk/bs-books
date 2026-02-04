package testutil

import "database/sql"

func SeedUsers(tx *sql.Tx) []string {
	tx.Exec(`INSERT INTO users (id, email, password_hash) VALUES
		('73681e21-08d4-43e1-b0b6-8d6f75a9b8b3', 'test@example.com', 'hashedpassword1'),
		('83681e21-08d4-43e1-b0b6-8d6f75a9b8b4', 'peepee@poopoo.com', 'hashedpassword2')
	`)
	return []string{
		"73681e21-08d4-43e1-b0b6-8d6f75a9b8b3",
		"83681e21-08d4-43e1-b0b6-8d6f75a9b8b4",
	}
}
