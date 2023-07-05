package model

import (
	"context"
)

// UserController is a business logic around user storage.
type UserController interface {
	// Get users
	UserByID(ctx context.Context, userID string) (User, error)
	UserBySecondaryID(ctx context.Context, idt AuthIdentityType, id string) (User, error)
	UserByFederatedID(ctx context.Context, idt UserFederatedType, idOther, id string) (User, error)

	// Admin actions for users
	GetUsers(ctx context.Context, filter string, skip, limit int) ([]User, int, error)
	InvalidateCache()
}

// UserMutationController is a business logic around user mutation storage.
type UserMutationController interface {
	CreateUserWithPassword(ctx context.Context, u User, password string) (User, error)
	UpdateUserPassword(ctx context.Context, userID, password string) error
	ChangeBlockStatus(ctx context.Context, userID, reason, whoName, whoID string, blocked bool) error
	UpdateUser(ctx context.Context, u User, fields []string) (User, error)
	DeleteUser(ctx context.Context, userID string) error

	SendEmailConfirmation(ctx context.Context, userID string) error
	SendPhoneConfirmation(ctx context.Context, userID string) error

	InvalidateCache()
}
