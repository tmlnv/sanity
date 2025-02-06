package validator

import (
	"strconv"
)

func ValidateNumber(s string) error {
	if s == "" || s == "0" {
		return nil
	}
	_, err := strconv.Atoi(s)
	return err
}
