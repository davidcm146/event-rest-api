package utils

import (
	"fmt"
	"time"
)

// ParseAndFormatDate parses a date string like "02/01/2006" and returns it in "2006-01-02" format
func ParseAndFormatDate(input string) (string, error) {
	parsedDate, err := time.Parse("02/01/2006", input)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}
	return parsedDate.Format("2006-01-02"), nil
}
