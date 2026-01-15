package users

import (
	"bs-books-api/internal/db"
	"context"
)

type userRepo struct {
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (r *userRepo) GetByEmail(email string, ctx context.Context, db db.DBTX) (*User, error) {
	var user *User
	row, err := db.QueryContext(ctx, `SELECT id, email, password_hash FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		user = &User{}
		if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
			return nil, err
		}
		return user, nil
	}
	return nil, nil
}

func (r *userRepo) Create(user *User, ctx context.Context, db db.DBTX) error {
	_, err := db.ExecContext(ctx, `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`, user.ID, user.Email, user.PasswordHash)
	return err
}
