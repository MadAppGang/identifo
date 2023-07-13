package model

import "context"

// UserAuthStorage store all user auth sessions, devices, challenges
// UserAuthStorage is a storage which keep all auth information for user.
// All login strategies must implement this interface.
// All 2FA strategies must implement this interface.
type UserAuthStorage interface {
	Storage
	ImportableStorage

	AddChallenge(ctx context.Context, challenge UserAuthChallenge) (UserAuthChallenge, error)
	GetLatestChallenge(ctx context.Context, strategy AuthStrategy, userID string) (UserAuthChallenge, error)
	MarkChallengeAsSolved(ctx context.Context, ch UserAuthChallenge) error

	// AddAuthEnrolment
	// RemoveAuthEnrolment
	// Add2FAEnrolment
	// Remove2FAEnrolment
	// Solve2FAChallenge
	// Solve2Challenge
}
