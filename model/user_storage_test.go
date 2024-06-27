package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordGen(t *testing.T) {
	pass := randSeq(12)
	assert.Len(t, pass, 12)

	pass = randSeq(256)
	assert.Len(t, pass, 256)

	pass = randSeq(1)
	assert.Len(t, pass, 1)
}
