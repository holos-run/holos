package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPascalCase(t *testing.T) {
	assert.Empty(t, PascalCase(""))

	assert.Equal(t, "X", PascalCase("_"))

	assert.Equal(t, "XMyFieldName_2", PascalCase("_my_field_name_2"))

	assert.Equal(t, "TestSnake", PascalCase("test_snake"))
}
