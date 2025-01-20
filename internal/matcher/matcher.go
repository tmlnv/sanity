package matcher

import (
	"regexp"
	"strings"
)

// Match checks if the address matches the desired pattern.
func Match(address, prefix, suffix, regex string) bool {
	if prefix != "" && !strings.HasPrefix(address, prefix) {
		return false
	}
	if suffix != "" && !strings.HasSuffix(address, suffix) {
		return false
	}
	if regex != "" {
		matched, _ := regexp.MatchString(regex, address)
		return matched
	}
	return true
}
