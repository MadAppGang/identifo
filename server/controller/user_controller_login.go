package controller

import (
	"context"
	"strings"

	"github.com/madappgang/identifo/v2/model"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// tenant related information flattered, as:
// "112233:admin" : "user", where 112233 - tenant ID, admin - a group, user - role in a group
// and tenant name added as well:
// "tenant:112233": "tenant corporation"
func (c *UserStorageController) GetJWTTokens(ctx context.Context, app model.AppData, u model.User, scopes []string) (model.AuthResponse, error) {
	// check if we are
	var err error

	// TODO: implement custom payload provider for app
	resp := model.AuthResponse{}
	ud := model.UserData{}
	aud := []string{app.ID}

	ap := AccessTokenScopes(scopes) // fields for access token
	apf := model.FieldsetForScopes(ap)
	data := map[string]any{}

	// access token needs tenant data in it
	if needTenantInfo(ap) {
		ud, err = c.u.UserData(ctx, u.ID, model.UserDataFieldTenantMembership)
		if err != nil {
			return resp, err
		}
		ti := TenantData(ud.TenantMembership, ap)
		maps.Copy(data, ti)
	}
	// create access token
	at, err := c.ts.NewToken(model.TokenTypeAccess, u, aud, apf, data)
	if err != nil {
		return resp, err
	}
	access, err := c.ts.SignToken(at)
	if err != nil {
		return resp, err
	}

	// id token
	var id string
	if slices.Contains(scopes, model.IDTokenScope) {
		// get fields for id token
		f := model.FieldsetForScopes(scopes)
		data := map[string]any{}

		// if we need tenant data in id token
		if needTenantInfo(scopes) {
			// we can already have userData fetched for access token
			if len(ud.UserID) == 0 {
				ud, err = c.u.UserData(ctx, u.ID, model.UserDataFieldTenantMembership)
				if err != nil {
					return resp, err
				}
			}
			ti := TenantData(ud.TenantMembership, scopes)
			maps.Copy(data, ti)
		}

		// create id token
		idt, err := c.ts.NewToken(model.TokenTypeID, u, aud, f, data)
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
	if slices.Contains(scopes, model.OfflineScope) && app.Offline {
		rt, err := c.ts.NewToken(model.TokenTypeRefresh, u, aud, ap, nil)
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

func TenantData(ud []model.TenantMembership, scopes []string) map[string]any {
	res := map[string]any{}
	filter := []string{}
	getAll := false
	for _, s := range scopes {
		if s == model.TenantScopeAll {
			getAll = true
			break
		} else if strings.HasPrefix(s, model.TenantScopePrefix) && len(s) > len(model.TenantScopePrefix) {
			filter = append(filter, s[len(model.TenantScopePrefix):])
		}
	}
	for _, t := range ud {
		// skip the scopes we don't need to have
		if !getAll && !slices.Contains(filter, t.TenantID) {
			continue
		}
		tid := t.TenantID
		res["tenant:"+t.TenantID] = t.TenantName
		for k, v := range t.Groups {
			// "role:tenant_id:group_id" : "role"
			res["role:"+tid+":"+k] = v
		}
	}
	return res
}

func needTenantInfo(scopes []string) bool {
	for _, s := range scopes {
		if strings.HasPrefix(s, model.TenantScopePrefix) {
			return true
		}
	}
	return false
}
