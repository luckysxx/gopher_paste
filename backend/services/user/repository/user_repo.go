package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"project/services/user/db"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, parms *db.CreateUserParams) (*db.User, error)
	GetByUsername(ctx context.Context, username string) (*db.User, error)
}

type userRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q: q}
}

func (r *userRepository) Create(ctx context.Context, parms *db.CreateUserParams) (*db.User, error) {
	User, err := r.q.CreateUser(ctx, *parms)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &User, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*db.User, error) {
	user, err := r.q.GetUserByUsername(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	return &user, err
}
