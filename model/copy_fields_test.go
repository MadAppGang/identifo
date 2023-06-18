package model_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testUser struct {
	ID       string
	Name     string
	Phone    string
	Address  string
	Password string
}

func TestCopyFields(t *testing.T) {
	u := testUser{
		ID:       "1",
		Name:     "JAck",
		Phone:    "+61450396664",
		Address:  "7.9 Bona Vista ave",
		Password: "Some hash",
	}
	fields := []string{"ID", "Name", "Password"}

	result := model.CopyFields(u, fields)
	assert.Empty(t, result.Phone)
	assert.Empty(t, result.Address)
	assert.Equal(t, u.ID, result.ID)
	assert.Equal(t, u.Name, result.Name)
	assert.Equal(t, u.Password, result.Password)
}

type testShortUser struct {
	ID       string
	Name     string
	Password string
	Other    string
}

func TestCopyDstFields(t *testing.T) {
	u := testUser{
		ID:       "1",
		Name:     "JAck",
		Phone:    "+61450396664",
		Address:  "7.9 Bona Vista ave",
		Password: "Some hash",
	}

	dst := testShortUser{}
	err := model.CopyDstFields(u, &dst)
	require.NoError(t, err)

	assert.Empty(t, dst.Other)
	assert.Equal(t, u.ID, dst.ID)
	assert.Equal(t, u.Name, dst.Name)
	assert.Equal(t, u.Password, dst.Password)
}
