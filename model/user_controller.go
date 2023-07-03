package model

import (
	"context"
)

// UserController is a business logic around user storage.
type UserController interface {
	// Get users
	UserByID(ctx context.Context, userID string) (User, error)
	GetUsers(ctx context.Context, filter string, skip, limit int) ([]User, int, error)

	// Mutations
	CreateUserWithPassword(ctx context.Context, u User, password, locale string) (User, error)
	UpdateUserPassword(ctx context.Context, userID, password, locale string) error
	ChangeBlockStatus(ctx context.Context, userID, reason, locale, whoName, whoID string, blocked bool) error
	UpdateUser(ctx context.Context, u User, fields []string, locale string) (User, error)
	DeleteUser(ctx context.Context, userID string) error
}
