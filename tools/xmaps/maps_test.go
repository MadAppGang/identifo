package xmaps_test

import (
	"fmt"
	"testing"

	"github.com/madappgang/identifo/v2/tools/xmaps"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name    string
	Age     int
	Address address
	Word    *address
	Other   []address
}

type address struct {
	Street string
	Apt    int
}

func TestFieldsToMap(t *testing.T) {
	p := person{
		Name: "John",
		Age:  30,
		Address: address{
			Street: "321 Main St",
			Apt:    123,
		},
	}
	m := xmaps.FieldsToMap(p)
	fmt.Printf("%v\n", m)
	assert.Len(t, m, 4)
	assert.Equal(t, 123, m["Address.Apt"])
	assert.Equal(t, "321 Main St", m["Address.Street"])
}

func TestFieldsToMapWithArray(t *testing.T) {
	p := person{
		Name: "John",
		Age:  30,
		Other: []address{
			{
				Street: "321 Main St",
				Apt:    123,
			},
			{
				Street: "Other street",
			},
		},
	}
	m := xmaps.FieldsToMap(p)
	fmt.Printf("%v\n", m)
	assert.Len(t, m, 5)
	assert.Equal(t, nil, m["Address.Apt"])
	assert.Equal(t, "Other street", m["Other[1].Street"])
}

func TestFilterMap(t *testing.T) {
	m := map[string]string{"a": "value_1", "b": "value_2", "c": "value_3"}
	assert.Len(t, xmaps.FilterMap(m, []string{"a"}), 1)
	assert.Equal(t, "value_1", xmaps.FilterMap(m, []string{"a"})["a"])
}
