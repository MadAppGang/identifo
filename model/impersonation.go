package model

import "context"

type ImpersonationProvider interface {
	CanImpersonate(ctx context.Context, appID string, adminUser, impersonatedUser User) (bool, error)
}
