package model_test

import (
	"testing"

	"github.com/madappgang/identifo/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestID(t *testing.T) {
	id := model.NewID()
	require.NotEmpty(t, id)
	assert.Equal(t, 24, len(id))
	assert.True(t, primitive.IsValidObjectID(id.String()))

	// check compatibility with bson's ObjectID
	oid, err := primitive.ObjectIDFromHex(id.String())
	assert.NoError(t, err)
	assert.False(t, oid.IsZero())
	assert.Equal(t, id.String(), oid.Hex())
}
