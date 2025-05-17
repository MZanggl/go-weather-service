package utils

import "fmt"

func FormatFloat(value float64, format string) string {
	return fmt.Sprintf("%.2f%s", value, format)
}
