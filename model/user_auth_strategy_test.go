package model_test

import (
	_ "embed"
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestStrategies(t *testing.T) {
	strategies := model.Strategies()
	assert.Len(t, strategies, 3)

	// TODO: Jack implement proper unit tests
}
