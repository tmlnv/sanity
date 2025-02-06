package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Base58 alphabet used by Solana addresses
const base58Chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// SolanaAddressLength is the length of a Solana public key in base58 format
const SolanaAddressLength = 44

// ValidateSolana checks if the provided prefix, suffix, and regex pattern could potentially
// match a valid Solana address
func ValidateSolana(prefix, suffix, regexPattern string) error {
	// Check prefix
	if err := validateBase58String(prefix); err != nil {
		return fmt.Errorf("invalid prefix: %v", err)
	}

	// Check suffix
	if err := validateBase58String(suffix); err != nil {
		return fmt.Errorf("invalid suffix: %v", err)
	}

	// Check combined length
	if len(prefix)+len(suffix) > SolanaAddressLength {
		return fmt.Errorf("combined prefix and suffix length (%d) exceeds Solana address length (%d)",
			len(prefix)+len(suffix), SolanaAddressLength)
	}

	// Validate regex if provided
	if regexPattern != "" {
		if err := validateRegexPattern(regexPattern); err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	}

	return nil
}

// validateBase58String checks if a string contains only valid base58 characters
func validateBase58String(s string) error {
	if s == "" {
		return nil
	}

	for _, c := range s {
		if !strings.ContainsRune(base58Chars, c) {
			return fmt.Errorf("character '%c' is not a valid base58 character", c)
		}
	}

	return nil
}

// validateRegexPattern checks if a regex pattern is valid and could potentially match
// a Solana address
func validateRegexPattern(pattern string) error {
	// First check if it's a valid regex
	_, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid regular expression: %v", err)
	}

	// Check for obviously invalid patterns
	if strings.Contains(pattern, "[^"+base58Chars+"]") {
		return fmt.Errorf("pattern contains characters not valid in base58")
	}

	return nil
}
