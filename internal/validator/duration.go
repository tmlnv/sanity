package validator

import (
	"fmt"
	"strconv"
	"time"
)

// ValidateTimeout checks if the provided timeout string is valid
// It accepts either a pure number (interpreted as seconds) or a duration string format
func ValidateTimeout(s string) (time.Duration, error) {
	if s == "" || s == "0" {
		return 0, nil
	}

	// Try parsing as pure number (seconds)
	if seconds, err := strconv.Atoi(s); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	// Try parsing as duration string
	duration, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format (e.g., 30s, 5m): %v", err)
	}

	return duration, nil
}

func ValidateDuration(s string) error {
	_, err := ValidateTimeout(s)
	return err
}
