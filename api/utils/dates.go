package utils

import (
	"time"
)

func IsValidDate(dateStr string) bool {
	format := "2006-01-02"
	_, err := time.Parse(format, dateStr)
	return err == nil
}

func IsDateInFuture(dateStr string) bool {
	format := "2006-01-02"
	date, err := time.Parse(format, dateStr)
	if err != nil {
		return false
	}
	return date.After(time.Now())
}
