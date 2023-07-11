package controller

import (
	"context"
	"strings"

	"github.com/madappgang/identifo/v2/model"
)

// TODO! we need to add tenant related information flattered, as:
// "112233:admin:user", where 112233 - tenant ID, admin - a group, user - role in a group
func (c *UserStorageController) getJWTTokens(ctx context.Context, app model.AppData, u model.User, scopes []string) (model.AuthResponse, error) {
	// check if we are

	// TODO: implement custom payload provider for app
	resp := model.AuthResponse{}
	ap := AccessTokenScopes(scopes) // fields for access token
	apf := model.FieldsetForScopes(scopes)

	at, err := c.ts.NewToken(model.TokenTypeAccess, u, apf, nil)
	if err != nil {
		return resp, err
	}
	access, err := c.ts.SignToken(at)
	if err != nil {
		return resp, err
	}

	// id token
	var id string
	if sliceContains(scopes, model.IDTokenScope) {
		// get fields for id token
		f := model.FieldsetForScopes(scopes)
		data := map[string]any{}

		idt, err := c.ts.NewToken(model.TokenTypeID, u, f, data)
		if err != nil {
			return resp, err
		}
		id, err = c.ts.SignToken(idt)
		if err != nil {
			return resp, err
		}
	}

	// refresh token
	var refresh string
	if sliceContains(scopes, model.OfflineScope) && app.Offline {
		rt, err := c.ts.NewToken(model.TokenTypeRefresh, u, ap, nil)
		if err != nil {
			return resp, err
		}
		refresh, err = c.ts.SignToken(rt)
		if err != nil {
			return resp, err
		}
	}

	resp.AccessToken = &access
	if len(id) > 0 {
		resp.IDToken = &id
	}
	if len(refresh) > 0 {
		// TODO: attach refresh token to device
		// TODO: save refresh token to db to invalidate on logout or device deactivation
		resp.RefreshToken = &refresh
	}

	return resp, nil
}

func AccessTokenScopes(scopes []string) []string {
	result := []string{}
	for _, s := range scopes {
		if strings.HasPrefix(s, model.AccessTokenScopePrefix) && len(s) > len(model.AccessTokenScopePrefix) {
			result = append(result, s[len(model.AccessTokenScopePrefix):])
		}
	}
	return result
}

func TenantData(ud model.UserData) map[string]any {
	res := map[string]any{}
	for _, t := range ud.TenantMembership {
		tid := t.TenantID
		for k, v := range t.Groups {
			// "tenant_id:group_id" : "role"
			res[tid+":"+k] = v
		}
	}
	return res
}
