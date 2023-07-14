package controller

import (
	"context"
	"net/url"
	"time"

	"github.com/madappgang/identifo/v2/l"
	"github.com/madappgang/identifo/v2/model"
)

// CreateInvitation creates invitation for the user.
// if t - inviter token is empty, it means we are creating from admin panel or management api
func (c *UserStorageController) CreateInvitation(ctx context.Context, t *model.JWToken, tenant, group, role, email string, app *model.AppData) (model.Invite, *url.URL, error) {
	zi := model.Invite{} // empty value

	if email != "" && !model.EmailRegexp.MatchString(email) {
		return zi, nil, l.ErrorAPIRequestBodyEmailInvalid
	}
	uu := model.User{Email: email}

	invitedTo := map[string]model.TenantMembership{
		tenant: {
			TenantID: tenant,
			Groups:   map[string][]string{group: {role}},
		},
	}

	tenantName := ""

	// we have inviter
	if t != nil {
		inviterUD, err := c.u.UserData(ctx, t.Subject(), model.UserDataFieldTenantMembership)
		if err != nil {
			return zi, nil, l.NewError(l.ErrorInvalidInviteTokenBadInvitee, err)
		}

		invitedTo := c.filterInviteeCouldInvite(inviterUD.TenantMembership, invitedTo)
		if len(invitedTo) == 0 {
			return zi, nil, l.ErrorInvalidInviteTokenBadInvitee
		}

		inviter, err := c.u.UserByID(ctx, inviterUD.UserID)
		if err != nil {
			return zi, nil, err
		}
		uu.ID = inviter.ID
		uu.GivenName = inviter.GivenName
		tenantName = inviterUD.TenantMembership[tenant].TenantName
	} else {
		uu.ID = model.RootUserID.String()
		tenant, err := c.u.TenantByID(ctx, tenant)
		if err != nil {
			return zi, nil, err
		}
		tenantName = tenant.Name
	}
	// invitation token as subject has an inviter, not an invited person
	// claims should have a list of tenants, groups and roles as usual token has
	// "role:tenant1:group1" : ["admin", "user"]
	// "tenant:tenant1" : "Tenant's Name (with ID=1)"
	fields := model.UserFieldsetMap[model.UserFieldsetInviteToken]
	tenantData := TenantData(invitedTo, []string{model.TenantScopeAll})
	invToken, err := c.ts.NewToken(model.TokenTypeInvite, uu, t.FullClaims().Audience, fields, tenantData)
	if err != nil {
		return zi, nil, err
	}

	invitation := model.Invite{
		AppID:       t.FullClaims().Audience[0],
		InviterID:   uu.ID,
		InviterName: uu.GivenName,
		Token:       invToken.Raw,
		Email:       email,
		Role:        role,
		Tenant:      tenant,
		TenantName:  tenantName,
		Group:       group,
		CreatedBy:   uu.ID,
		CreatedAt:   time.Now(),
		ExpiresAt:   invToken.ExpiresAt(),
	}

	err = c.is.Save(ctx, invitation)
	if err != nil {
		return zi, nil, err
	}

	link, err := c.invitationLink(invitation, app)
	if err != nil {
		return zi, nil, err
	}
	return invitation, link, nil
}
