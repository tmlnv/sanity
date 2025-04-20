package matcher

import (
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/tmlnv/sanity/internal/config"
)

type Matcher struct {
	prefix string
	suffix string
	regexp *regexp2.Regexp
}

func NewMatcher(cfg config.Config) *Matcher {
	var r *regexp2.Regexp
	if cfg.Regexp != "" {
		r = regexp2.MustCompile(cfg.Regexp, regexp2.None)
		r.MatchTimeout = time.Second // DoS guard
	}
	return &Matcher{
		prefix: cfg.Prefix,
		suffix: cfg.Suffix,
		regexp: r,
	}
}

func (m *Matcher) Match(address string) bool {
	if m.prefix != "" && !strings.HasPrefix(address, m.prefix) {
		return false
	}
	if m.suffix != "" && !strings.HasSuffix(address, m.suffix) {
		return false
	}
	if m.regexp != nil {
		ok, _ := m.regexp.MatchString(address)
		return ok
	}
	return true
}
