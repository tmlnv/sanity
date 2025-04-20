package validator

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
)

// --------------------------------- constants ---------------------------------

// Base58 alphabet used by Solana addresses (no 0 O I l)
const base58Chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Length of a Solana public key in base‑58
const SolanaAddressLength = 44

// --------------------------------- public API --------------------------------

// ValidateSolana checks that prefix, suffix and (optionally) a regexp pattern
// could still match a valid 44‑char Solana public key.
func ValidateSolana(prefix, suffix, pattern string) error {
	if err := validateBase58String(prefix); err != nil {
		return fmt.Errorf("invalid prefix: %w", err)
	}
	if err := validateBase58String(suffix); err != nil {
		return fmt.Errorf("invalid suffix: %w", err)
	}
	if len(prefix)+len(suffix) > SolanaAddressLength {
		return fmt.Errorf("combined prefix + suffix length (%d) exceeds %d",
			len(prefix)+len(suffix), SolanaAddressLength)
	}
	if pattern != "" {
		if err := validateRegexpPattern(pattern); err != nil {
			return fmt.Errorf("invalid regexp pattern: %w", err)
		}
	}
	return nil
}

// --------------------------------- internals ---------------------------------

// validateBase58String returns an error if s contains anything outside base‑58.
func validateBase58String(s string) error {
	for _, r := range s {
		if !strings.ContainsRune(base58Chars, r) {
			return fmt.Errorf("character %q is not valid base‑58", r)
		}
	}
	return nil
}

// validateRegexpPattern guarantees that p
//  1. compiles under github.com/dlclark/regexp2 (so back‑refs are ok) and
//  2. contains no literal runes outside the base‑58 set (meta‑chars allowed).
//
// Digits 0‑9 and ',' are permitted *only* inside a {min,max} quantifier;
// anything goes inside [...] character classes; escapes are respected.
func validateRegexpPattern(p string) error {
	// 1 ─ compile with full .NET‑style engine
	if _, err := regexp2.Compile(p, regexp2.None); err != nil {
		return err
	}

	// 2 ─ walk the pattern rune‑by‑rune to vet literals
	const meta = `.*+?^$()[]{}|\-\\`
	inClass, escaped, inQuant := false, false, false

	for _, r := range p {
		switch {
		// ───────────── handle escapes ─────────────
		case escaped:
			escaped = false
			continue
		case r == '\\':
			escaped = true
			continue

		// ─────────── inside [...] class ───────────
		case inClass:
			if r == ']' {
				inClass = false
			}
			continue // everything is allowed here
		case r == '[':
			inClass = true
			continue

		// ────────── inside {min,max} quantifier ──────────
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

		// ─────────── literal outside special zones ───────────
		if strings.ContainsRune(base58Chars, r) {
			continue // valid literal
		}
		if strings.ContainsRune(meta, r) {
			continue // accepted meta‑char
		}
		return fmt.Errorf("character %q is not base‑58 or regex meta", r)
	}
	return nil
}
