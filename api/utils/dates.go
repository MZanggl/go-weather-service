package utils

import "time"

func IsValidDate(dateStr string) bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, dateStr)
	return err == nil
}
