package matcher

import (
	"regexp"
	"strings"

	"github.com/tmlnv/sanity/internal/config"
)

type Matcher struct {
	prefix string
	suffix string
	regexp *regexp.Regexp
}

func NewMatcher(cfg config.Config) *Matcher {
	var re *regexp.Regexp
	if cfg.Regexp != "" {
		re = regexp.MustCompile(cfg.Regexp)
	}

	return &Matcher{
		prefix: cfg.Prefix,
		suffix: cfg.Suffix,
		regexp: re,
	}
}

func (m *Matcher) Match(address string) bool {
	if m.prefix != "" && !strings.HasPrefix(address, m.prefix) {
		return false
	}

	if m.suffix != "" && !strings.HasSuffix(address, m.suffix) {
		return false
	}

	if m.regexp != nil && !m.regexp.MatchString(address) {
		return false
	}

	return true
}
