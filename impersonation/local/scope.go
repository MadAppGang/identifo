package local

import (
	"context"
	"slices"

	"github.com/madappgang/identifo/v2/model"
)

type ScopeImpersonator struct {
	allowedScopes []string
}

func NewScopeImpersonator(allowedScopes []string) *ScopeImpersonator {
	return &ScopeImpersonator{allowedScopes: allowedScopes}
}

func (si *ScopeImpersonator) CanImpersonate(ctx context.Context, appID string, adminUser model.User, user model.User) (bool, error) {
	if !adminUser.Active || adminUser.Anonymous {
		return false, nil
	}

	for _, scope := range si.allowedScopes {
		if slices.Contains(adminUser.Scopes, scope) {
			return true, nil
		}
	}

	return false, nil
}
