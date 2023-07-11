package xslices_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/tools/xslices"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestSliceContains(t *testing.T) {
	assert.True(t, slices.Contains([]string{"a", "b"}, "a"))
	assert.True(t, slices.Contains([]string{"a", "b", "1"}, "1"))
	assert.False(t, slices.Contains([]string{"a", "b", "1"}, "11"))
}

func TestIntersect(t *testing.T) {
	assert.Contains(t, xslices.Intersect([]string{"a", "b"}, []string{"a"}), "a")
	assert.Len(t, xslices.Intersect([]string{"a", "b"}, []string{"a"}), 1)

	assert.Contains(t, xslices.Intersect([]string{"a", "b"}, []string{"a", "b", "c"}), "a")
	assert.Contains(t, xslices.Intersect([]string{"a", "b"}, []string{"a", "b", "c"}), "b")
	assert.Len(t, xslices.Intersect([]string{"a", "b"}, []string{"a", "b", "c"}), 2)

	assert.Len(t, xslices.Intersect([]string{"a", "b"}, []string{"c"}), 0)
}

func TestConcatUnique(t *testing.T) {
	assert.Contains(t, xslices.ConcatUnique([]string{"a", "b"}, []string{"a"}), "a")
	assert.Len(t, xslices.Intersect([]string{"a", "b"}, []string{"a"}), 1)

	r := xslices.ConcatUnique([]string{"a", "b", "c", "d"}, []string{"a"})
	assert.Len(t, r, 4)
}
