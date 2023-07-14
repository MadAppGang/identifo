package controller

import (
	"context"
	"strings"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/tools/xslices"
	"golang.org/x/exp/slices"
)

// AddUserToTenantWithInvitationToken add user to tenants and groups from invitation token
// we need to get inviters information
// check if he can invite someone to the group/tenant
// add user to new group/groups
// return the new user's membership
func (c *UserStorageController) AddUserToTenantWithInvitationToken(ctx context.Context, u model.User, t *model.JWToken) (model.UserData, error) {
	if t == nil || t.Type() != model.TokenTypeInvite {
		return model.UserData{}, l.ErrorInvalidInviteToken
	}

	inviter, err := c.u.UserByID(ctx, t.Subject())
	if err != nil {
		return model.UserData{}, l.NewError(l.ErrorInvalidInviteTokenBadInvitee, err)
	}

	inviterUD, err := c.u.UserData(ctx, inviter.ID, model.UserDataFieldTenantMembership)
	if err != nil {
		return model.UserData{}, l.NewError(l.ErrorInvalidInviteTokenBadInvitee, err)
	}

	// let's parse where the user has been invited
	invitedTo := getInvitationFromClaim(t.FullClaims().Payload)

	// get the tenants and groups, which invitee is able to invite himself

	userMembership := c.filterInviteeCouldInvite(inviterUD.TenantMembership, invitedTo)
	md := model.UserData{UserID: u.ID, TenantMembership: userMembership}
	if len(userMembership) > 0 {
		md, err = c.ums.UpdateUserData(ctx, u.ID, md, model.UserDataFieldTenantMembership)
		if err != nil {
			return model.UserData{}, err
		}
	} else {
		return model.UserData{}, l.ErrorInvalidInviteTokenBadInvitee
	}

	return md, nil
}

// returns a map of invitations
func getInvitationFromClaim(claims map[string]any) map[string]model.TenantMembership {
	// tenants map
	tenants := map[string]model.TenantMembership{}
	for k, v := range claims {
		vs, ok := v.([]any)
		if !ok {
			continue
		}

		if strings.HasPrefix(k, model.RoleScopePrefix) && len(k) > len(model.RoleScopePrefix) {
			// "role:tenant_id:group_id" : "role"
			parts := strings.Split(k, ":")
			if len(parts) != 3 {
				continue
			}

			// convert slice of any to slice of unique strings of roles
			vss := []string{}
			for _, v := range vs {
				vv, ok := v.(string)
				if ok && slices.Contains(vss, vv) {
					vss = append(vss, vv)
				}
			}

			// update our membership map
			membership, ok := tenants[parts[1]]
			// create a tenant in the map
			if !ok {
				membership = model.TenantMembership{
					TenantID: parts[1],
					Groups:   map[string][]string{parts[2]: vss},
				}
				tenants[parts[1]] = membership
			} else { // we have tenant, let's add group membership
				group, ok := membership.Groups[parts[2]]
				// new group
				if !ok {
					membership.Groups[parts[2]] = vss
				} else { // add new role to existent group
					membership.Groups[parts[2]] = xslices.Unique(append(group, vss...))
				}
			}
		}
	}
	return tenants
}

func (c *UserStorageController) filterInviteeCouldInvite(inviterMembership, invitedTo map[string]model.TenantMembership) map[string]model.TenantMembership {
	userMembership := map[string]model.TenantMembership{}

	for k, v := range invitedTo {
		it, ok := inviterMembership[k]
		// invitee is not belongs to the tenant, he could not invite anyone there
		if !ok {
			// TODO: add log information in log stream about that with ERROR
			continue
		}
		// let's check the groups
		userGroups := map[string][]string{}
		for gk, g := range v.Groups {
			invg, ok := it.Groups[gk]
			// invitee does not belongs to the group, he could not invite anyone there
			if !ok {
				// TODO: add log information in log stream about that with ERROR
				continue
			}
			rrr := xslices.Intersect(invg, c.s.TenantMembershipManagementRole)
			// invitee has one or more roles in this group which allows him to make an invitations
			if len(rrr) > 0 {
				userGroups[gk] = g
			}
		}
		// we can invite to one or more groups of the tenant
		if len(userGroups) > 0 {
			userMembership[k] = model.TenantMembership{
				TenantID: v.TenantID,
				Groups:   userGroups,
			}
		}
	}

	return userMembership
}
