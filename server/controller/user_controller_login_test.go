package controller_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/server/controller"
	"github.com/madappgang/identifo/v2/storage/mock"
	"github.com/stretchr/testify/assert"
)

func TestAddTenantData(t *testing.T) {
	ud := model.UserData{
		TenantMembership: []model.TenantMembership{
			{
				TenantID:   "tenant1",
				TenantName: "I am a tenant 1",
				Groups:     map[string]string{"default": "admin", "group1": "user"},
			},
			{
				TenantID:   "tenant2",
				TenantName: "Apple corporation",
				Groups:     map[string]string{"default": "guest", "group33": "admin"},
			},
		},
	}

	flattenData := controller.TenantData(ud.TenantMembership, nil)
	assert.Len(t, flattenData, 0)

	flattenData = controller.TenantData(ud.TenantMembership, []string{"tenant:all"})
	assert.Contains(t, flattenData, "tenant1:default")
	assert.Contains(t, flattenData, "tenant1:group1")
	assert.Contains(t, flattenData, "tenant2:default")
	assert.Contains(t, flattenData, "tenant2:group33")
	assert.Contains(t, flattenData, "tenant:tenant1")
	assert.Contains(t, flattenData, "tenant:tenant2")
	assert.Equal(t, flattenData["tenant1:default"], "admin")
	assert.Equal(t, flattenData["tenant1:group1"], "user")
	assert.Equal(t, flattenData["tenant2:default"], "guest")
	assert.Equal(t, flattenData["tenant2:group33"], "admin")
	assert.Equal(t, flattenData["tenant:tenant1"], "I am a tenant 1")
	assert.Equal(t, flattenData["tenant:tenant2"], "Apple corporation")
	assert.Len(t, flattenData, 6)

	scopes := []string{"tenant:tenant1", "tenant:tenant3"}
	flattenData = controller.TenantData(ud.TenantMembership, scopes)
	assert.Contains(t, flattenData, "tenant1:default")
	assert.Contains(t, flattenData, "tenant1:group1")
	assert.NotContains(t, flattenData, "tenant2:default")
	assert.NotContains(t, flattenData, "tenant2:group33")
	assert.Contains(t, flattenData, "tenant:tenant1")
	assert.NotContains(t, flattenData, "tenant:tenant2")
	assert.Equal(t, flattenData["tenant1:default"], "admin")
	assert.Equal(t, flattenData["tenant1:group1"], "user")
	assert.Equal(t, flattenData["tenant:tenant1"], "I am a tenant 1")
	assert.Len(t, flattenData, 3)

	scopes = []string{"tenant:tenant1", "tenant:tenant3", "tenant:all"}
	flattenData = controller.TenantData(ud.TenantMembership, scopes)
	assert.Len(t, flattenData, 6)
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

func TestRequestJWT(t *testing.T) {
	tm := []model.TenantMembership{
		{
			TenantID:   "tenant1",
			TenantName: "I am a tenant 1",
			Groups:     map[string]string{"default": "admin", "group1": "user"},
		},
		{
			TenantID:   "tenant2",
			TenantName: "Apple corporation",
			Groups:     map[string]string{"default": "guest", "group33": "admin"},
		},
	}
	u := mock.UserStorage{
		UData: map[string]model.UserData{
			"user1": {
				UserID:           "user1",
				TenantMembership: tm,
			},
		},
	}

	c := controller.NewUserStorageController(
		&u,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		model.ServerSettings{},
	)
}
