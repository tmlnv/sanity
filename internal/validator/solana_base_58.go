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

// ValidateSolana checks if the provided prefix, suffix, and regexp pattern could potentially
// match a valid Solana address
func ValidateSolana(prefix, suffix, regexpPattern string) error {
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

	// Validate regexp if provided
	if regexpPattern != "" {
		if err := validateRegexpPattern(regexpPattern); err != nil {
			return fmt.Errorf("invalid regexp pattern: %v", err)
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

// validateRegexpPattern checks if a regexp pattern is valid and could potentially match
// a Solana address
func validateRegexpPattern(pattern string) error {
	// First check if it's a valid regexp with Go's standard regexp package
	_, err := regexp.Compile(pattern)
	if err != nil {
		return err // Pass through the original regex compilation error
	}

	// List of valid regex metacharacters and additional characters that should be allowed
	regexMetaChars := `.*+?^$()[]{}|\\-,`

	// Check each character - it should either be a base58 character or a regex metacharacter
	for _, c := range pattern {
		// Skip if it's a valid base58 character
		if strings.ContainsRune(base58Chars, c) {
			continue
		}

		// Skip if it's a valid regex metacharacter
		if strings.ContainsRune(regexMetaChars, c) {
			continue
		}

		// If we get here, the character is neither a valid base58 char nor a regex metachar
		return fmt.Errorf("character '%c' is neither a valid base58 character nor a regex metacharacter", c)
	}

	return nil
}
