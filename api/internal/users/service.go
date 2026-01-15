package users

import (
	"bs-books-api/internal/db"
	"context"
)

type UserService struct {
	db       db.DBTX
	userRepo *userRepo
}

func NewUserService(db db.DBTX, userRepo *userRepo) *UserService {
	return &UserService{
		db:       db,
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByEmail(email string, ctx context.Context) (*User, error) {
	return s.userRepo.GetByEmail(email, ctx, s.db)
}

func (s *UserService) CreateUser(email, passwordHash string, ctx context.Context) error {
	new_user := NewUser(email, passwordHash)
	return s.userRepo.Create(new_user, ctx, s.db)
}
