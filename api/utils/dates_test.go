package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDateValidity(t *testing.T) {
	var input = false

	// valid date
	input = IsValidDate("2024-06-01")
	assert.Equal(t, true, input)

	// Invalid date string
	input = IsValidDate("not-a-date")
	assert.Equal(t, false, input)

	// Empty string
	input = IsValidDate("")
	assert.Equal(t, false, input)
}
func TestIsDateInFuture(t *testing.T) {
	// Date in the future
	futureDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	assert.Equal(t, true, IsDateInFuture(futureDate))

	// Date in the past
	pastDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	assert.Equal(t, false, IsDateInFuture(pastDate))

	// Today's date (should not be in the future)
	today := time.Now().Format("2006-01-02")
	assert.Equal(t, false, IsDateInFuture(today))

	// Invalid date string
	assert.Equal(t, false, IsDateInFuture("not-a-date"))

	// Empty string
	assert.Equal(t, false, IsDateInFuture(""))
}
