package matcher

import (
	"reflect"
	"testing"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/tmlnv/sanity/internal/config"
)

func TestNewMatcher(t *testing.T) {
	makeRE := func(p string) *regexp2.Regexp {
		r := regexp2.MustCompile(p, regexp2.None)
		r.MatchTimeout = time.Second
		return r
	}

	tests := []struct {
		name string
		cfg  config.Config
		want *Matcher
	}{
		{
			name: "empty config",
			cfg:  config.Config{},
			want: &Matcher{prefix: "", suffix: "", regexp: nil},
		},
		{
			name: "config with prefix and suffix",
			cfg:  config.Config{Prefix: "ABC", Suffix: "XYZ"},
			want: &Matcher{prefix: "ABC", suffix: "XYZ", regexp: nil},
		},
		{
			name: "config with regexp",
			cfg:  config.Config{Regexp: "^[A-Za-z0-9]{44}$"},
			want: &Matcher{prefix: "", suffix: "", regexp: makeRE("^[A-Za-z0-9]{44}$")},
		},
		{
			name: "config with all fields",
			cfg:  config.Config{Prefix: "ABC", Suffix: "XYZ", Regexp: "^ABC.*XYZ$"},
			want: &Matcher{prefix: "ABC", suffix: "XYZ", regexp: makeRE("^ABC.*XYZ$")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMatcher(tt.cfg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatcher() = %#v\nwant %#v", got, tt.want)
			}
		})
	}
}

func TestMatcher_Match(t *testing.T) {
	re44 := regexp2.MustCompile("^[1-9A-HJ-NP-Za-km-z]{44}$", regexp2.None)
	re32 := regexp2.MustCompile("^[1-9A-HJ-NP-Za-km-z]{32}$", regexp2.None)

	tests := []struct {
		name    string
		matcher *Matcher
		addr    string
		want    bool
	}{
		{
			name:    "empty matcher matches any address",
			matcher: &Matcher{},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    true,
		},
		{
			name:    "prefix match",
			matcher: &Matcher{prefix: "5vG"},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    true,
		},
		{
			name:    "prefix no match",
			matcher: &Matcher{prefix: "ABC"},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    false,
		},
		{
			name:    "suffix match",
			matcher: &Matcher{suffix: "AmR"},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    true,
		},
		{
			name:    "suffix no match",
			matcher: &Matcher{suffix: "XYZ"},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    false,
		},
		{
			name:    "regexp match",
			matcher: &Matcher{regexp: re44},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    true,
		},
		{
			name:    "regexp no match",
			matcher: &Matcher{regexp: re32},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    false,
		},
		{
			name:    "all conditions match",
			matcher: &Matcher{prefix: "5vG", suffix: "AmR", regexp: re44},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    true,
		},
		{
			name:    "all conditions no match - wrong prefix",
			matcher: &Matcher{prefix: "ABC", suffix: "AmR", regexp: re44},
			addr:    "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.matcher.Match(tt.addr)
			if got != tt.want {
				t.Errorf("Matcher.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
