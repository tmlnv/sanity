package matcher

import (
	"regexp"
	"strings"

	"github.com/tmlnv/sanity/internal/config"
)

type Matcher struct {
	prefix string
	suffix string
	regex  *regexp.Regexp
}

func NewMatcher(cfg config.Config) *Matcher {
	var re *regexp.Regexp
	if cfg.Regex != "" {
		re = regexp.MustCompile(cfg.Regex)
	}

	return &Matcher{
		prefix: cfg.Prefix,
		suffix: cfg.Suffix,
		regex:  re,
	}
}

func (m *Matcher) Match(address string) bool {
	if m.prefix != "" && !strings.HasPrefix(address, m.prefix) {
		return false
	}

	if m.suffix != "" && !strings.HasSuffix(address, m.suffix) {
		return false
	}

	if m.regex != nil && !m.regex.MatchString(address) {
		return false
	}

	return true
}
