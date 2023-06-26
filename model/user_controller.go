package model

import "context"

// UserController is a business logic around user storage.
type UserController interface {
	// Get users
	UserByID(ctx context.Context, userID string) (User, error)
	GetUsers(ctx context.Context, filter string, skip, limit int) ([]User, int, error)

	// Mutations
	CreateUserWithPassword(ctx context.Context, u User, password string) (User, error)
}
