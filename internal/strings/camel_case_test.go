package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCamelCase(t *testing.T) {
	assert.Empty(t, CamelCase(""))

	assert.Equal(t, "x", CamelCase("_"))

	assert.Equal(t, "xMyFieldName_2", CamelCase("_my_field_name_2"))

	assert.Equal(t, "testSnake", CamelCase("test_snake"))
}
