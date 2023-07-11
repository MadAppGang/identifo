package controller_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/server/controller"
	"github.com/stretchr/testify/assert"
)

func TestAddTenantData(t *testing.T) {
	ud := model.UserData{
		TenantMembership: []model.TenantMembership{
			{
				TenantID: "tenant1",
				Groups:   map[string]string{"default": "admin", "group1": "user"},
			},
			{
				TenantID: "tenant2",
				Groups:   map[string]string{"default": "guest", "group33": "admin"},
			},
		},
	}

	flattenData := controller.TenantData(ud)
	assert.Contains(t, flattenData, "tenant1:default")
	assert.Contains(t, flattenData, "tenant1:group1")
	assert.Contains(t, flattenData, "tenant2:default")
	assert.Contains(t, flattenData, "tenant2:group33")
	assert.Equal(t, flattenData["tenant1:default"], "admin")
	assert.Equal(t, flattenData["tenant1:group1"], "user")
	assert.Equal(t, flattenData["tenant2:default"], "guest")
	assert.Equal(t, flattenData["tenant2:group33"], "admin")
	assert.Len(t, flattenData, 4)
}

func TestAccessTokenScopes(t *testing.T) {
	scopes := []string{
		"id",
		"offline",
		"access:",
		"access:",
		"access:profile",
		"access:oidc",
	}
	r := controller.AccessTokenScopes(scopes)
	assert.Len(t, r, 2)
	assert.Contains(t, r, "profile")
	assert.Contains(t, r, "oidc")
}
