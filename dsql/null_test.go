package dsql_test

import (
	"encoding/json"
	"testing"

	"github.com/semrekkers/go-service/dsql"

	"github.com/stretchr/testify/assert"
)

func TestNullString(t *testing.T) {
	var (
		null dsql.Null[string]
		str  = dsql.NewNull("Hello, World!")
	)

	assert.False(t, null.Valid)
	assert.Empty(t, null.Some)
	assert.True(t, str.Valid)
	assert.Equal(t, "Hello, World!", str.Some)
}

func TestNullHelpers(t *testing.T) {
	var (
		null dsql.Null[string]
		str  = dsql.NewNull("Hello, World!")
	)

	str.Set("Hello!")
	nullPtr := null.Ptr()
	strPtr := str.Ptr()

	assert.Nil(t, nullPtr)
	assert.Equal(t, &str.Some, strPtr)
}

func TestNullSQLValue(t *testing.T) {
	var (
		dest, null dsql.Null[string]
		str        = dsql.NewNull("Hello, World!")
	)

	err := dest.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, dest.Valid)
	assert.Empty(t, dest.Some)

	err = dest.Scan("Hello!")
	assert.NoError(t, err)
	assert.True(t, dest.Valid)
	assert.Equal(t, "Hello!", dest.Some)

	v, _ := null.Value()
	assert.Nil(t, v)
	v, _ = str.Value()
	assert.Equal(t, str.Some, v)
}

func TestNullJSON(t *testing.T) {
	var (
		dest, null dsql.Null[string]
		jsonValue  = []byte("null")
	)

	data, err := json.Marshal(null)
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)

	err = json.Unmarshal(jsonValue, &dest)
	assert.NoError(t, err)
	assert.False(t, dest.Valid)
	assert.Empty(t, null.Some)
}
