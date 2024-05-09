package local

import (
	"context"
	"slices"

	"github.com/madappgang/identifo/v2/model"
)

type AccessRoleImpersonator struct {
	allowedAccessRoles []string
}

func NewAccessRoleImpersonator(allowedAccessRoles []string) *AccessRoleImpersonator {
	return &AccessRoleImpersonator{allowedAccessRoles: allowedAccessRoles}
}

func (si *AccessRoleImpersonator) CanImpersonate(ctx context.Context, appID string, adminUser model.User, user model.User) (bool, error) {
	if !adminUser.Active || adminUser.Anonymous {
		return false, nil
	}

	if slices.Contains(si.allowedAccessRoles, adminUser.AccessRole) {
		return true, nil
	}

	return false, nil
}
