package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateValidity(t *testing.T) {
	var input = false

	input = IsValidDate("2024-06-01")
	assert.Equal(t, true, input)

	input = IsValidDate("invalid date")
	assert.Equal(t, false, input)

	input = IsValidDate("")
	assert.Equal(t, false, input)
}
