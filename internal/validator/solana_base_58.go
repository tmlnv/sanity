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
func validateRegexpPattern(p string) error {
	_, err := regexp.Compile(p)
	if err != nil {
		return err
	}

	const meta = `.*+?^$()[]{}|\-\\` // no digits here
	inClass, escaped, inQuant := false, false, false

	for _, r := range p {
		switch {
		case escaped:
			escaped = false
			continue
		case r == '\\':
			escaped = true
			continue
		case inClass:
			if r == ']' {
				inClass = false
			}
			continue // everything allowed inside [...]
		case r == '[':
			inClass = true
			continue
		case inQuant:
			if r == '}' {
				inQuant = false
				continue
			}
			if r == ',' || ('0' <= r && r <= '9') {
				continue
			}
			return fmt.Errorf("invalid rune %q in quantifier", r)
		case r == '{':
			inQuant = true
			continue
		}

		// Outside [], {}, and escapes: only meta or base58 literals
		if strings.ContainsRune(base58Chars, r) {
			continue
		}
		if strings.ContainsRune(meta, r) {
			continue
		}

		return fmt.Errorf("character %q is not base58 or regex meta", r)
	}
	return nil
}
