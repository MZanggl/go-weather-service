package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFloat(t *testing.T) {
	input := FormatFloat(23.456789, "°C")
	assert.Equal(t, "23.46°C", input)
}
