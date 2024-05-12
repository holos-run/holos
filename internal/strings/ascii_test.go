package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsASCIILower(t *testing.T) {
	assert.True(t, isASCIILower('a'))
	assert.True(t, isASCIILower('z'))
	assert.False(t, isASCIILower('1'))
	assert.False(t, isASCIILower('#'))
	assert.False(t, isASCIILower('$'))
}

func TestIsASCIIDigit(t *testing.T) {
	assert.True(t, isASCIIDigit('1'))
	assert.True(t, isASCIIDigit('9'))
	assert.False(t, isASCIIDigit('a'))
	assert.False(t, isASCIIDigit('z'))
	assert.False(t, isASCIIDigit('$'))
}
