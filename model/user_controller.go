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
	GetJWTTokens(ctx context.Context, app AppData, u User, scopes []string) (AuthResponse, error)
	RefreshJWTToken(ctx context.Context, refresh_token *JWToken, access string, app AppData, scopes []string) (AuthResponse, error)

	InvalidateCache()
}

// UserMutationController is a business logic around user mutation storage.
type UserMutationController interface {
	CreateUser(ctx context.Context, u User) (User, error)
	CreateUserWithPassword(ctx context.Context, u User, password string) (User, error)
	UpdateUserPassword(ctx context.Context, userID, password string) error
	ChangeBlockStatus(ctx context.Context, userID, reason, whoName, whoID string, blocked bool) error
	UpdateUser(ctx context.Context, u User, fields []string) (User, error)
	DeleteUser(ctx context.Context, userID string) error

	SendEmailConfirmation(ctx context.Context, userID string) error
	SendPhoneConfirmation(ctx context.Context, userID string) error
	SendPasswordResetEmail(ctx context.Context, userID, appID string) (ResetEmailData, error)

	AddUserToTenantWithInvitationToken(ctx context.Context, u User, t *JWToken) (UserData, error)
	CreateInvitation(ctx context.Context, invitee *JWToken, tenant, group, role, email string) (Invite, error)

	InvalidateCache()
}

type ChallengeController interface {
	RequestChallenge(ctx context.Context, challenge UserAuthChallenge, userIDValue string) (UserAuthChallenge, error)
	VerifyChallenge(ctx context.Context, challenge UserAuthChallenge, userIDValue string) (User, AppData, UserAuthChallenge, error)
}
