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

type testUserPointer struct {
	ID         *string
	Name       *string
	Married    *bool
	NonPointer string
	Company    *Company
}

type Company struct {
	Name   *string
	People *int
}

func TestFilledValues(t *testing.T) {
	tu := testUserPointer{
		ID:         sp("1"),
		Name:       sp("Jack"),
		Married:    ib(true),
		NonPointer: "NonPointer",
		Company: &Company{
			Name:   sp("Apple"),
			People: ip(98),
		},
	}

	expected := []string{"ID", "Name", "Married", "Company.Name", "Company.People"}

	// pointer to struct should works
	result := model.Filled(&tu)
	assert.EqualValues(t, result, expected)

	// reference to struct should works
	result = model.Filled(tu)
	assert.EqualValues(t, result, expected)

	// let's clear the Name
	tu.Name = nil
	expected = []string{"ID", "Married", "Company.Name", "Company.People"}
	result = model.Filled(tu)
	assert.EqualValues(t, result, expected)

	// empty value should be treated as non nil
	tu.Name = sp("")
	expected = []string{"ID", "Name", "Married", "Company.Name", "Company.People"}
	result = model.Filled(tu)
	assert.EqualValues(t, result, expected)
}

func sp(s string) *string {
	return &s
}

func ip(i int) *int {
	return &i
}

func ib(b bool) *bool {
	return &b
}

func TestCopyOnlyFilledValues(t *testing.T) {
	tu := testUserPointer{
		ID:         sp("1"),
		Name:       sp("Jack"),
		Married:    ib(true),
		NonPointer: "NonPointer",
		Company: &Company{
			Name:   sp("Apple"),
			People: ip(98),
		},
	}

	dst := model.User{}
	err := model.CopyDstFields(tu, &dst)

	assert.NoError(t, err)
	assert.Equal(t, *tu.ID, dst.ID)
}

func TestContainsFields(t *testing.T) {
	tu := testUserPointer{
		ID:         sp("1"),
		Name:       sp("Jack"),
		Married:    ib(true),
		NonPointer: "NonPointer",
		Company: &Company{
			Name:   sp("Apple"),
			People: ip(98),
		},
	}

	contains := model.ContainsFields(tu, []string{"ID", "Name", "NonPointer", "Whatever", "Company.Name", "Company.People", "Company"})
	expected := []string{"ID", "Name", "NonPointer", "Company.Name", "Company.People"}
	assert.Equal(t, contains, expected)
}
