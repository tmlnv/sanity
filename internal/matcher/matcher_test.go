package matcher

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/tmlnv/sanity/internal/config"
)

func TestNewMatcher(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name string
		args args
		want *Matcher
	}{
		{
			name: "empty config",
			args: args{cfg: config.Config{}},
			want: &Matcher{prefix: "", suffix: "", regexp: nil},
		},
		{
			name: "config with prefix and suffix",
			args: args{cfg: config.Config{Prefix: "ABC", Suffix: "XYZ"}},
			want: &Matcher{prefix: "ABC", suffix: "XYZ", regexp: nil},
		},
		{
			name: "config with regexp",
			args: args{cfg: config.Config{Regexp: "^[A-Za-z0-9]{44}$"}},
			want: &Matcher{prefix: "", suffix: "", regexp: regexp.MustCompile("^[A-Za-z0-9]{44}$")},
		},
		{
			name: "config with all fields",
			args: args{cfg: config.Config{Prefix: "ABC", Suffix: "XYZ", Regexp: "^ABC.*XYZ$"}},
			want: &Matcher{prefix: "ABC", suffix: "XYZ", regexp: regexp.MustCompile("^ABC.*XYZ$")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMatcher(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatcher_Match(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		m    *Matcher
		args args
		want bool
	}{
		{
			name: "empty matcher matches any address",
			m:    &Matcher{},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: true,
		},
		{
			name: "prefix match",
			m:    &Matcher{prefix: "5vG"},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: true,
		},
		{
			name: "prefix no match",
			m:    &Matcher{prefix: "ABC"},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: false,
		},
		{
			name: "suffix match",
			m:    &Matcher{suffix: "AmR"},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: true,
		},
		{
			name: "suffix no match",
			m:    &Matcher{suffix: "XYZ"},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: false,
		},
		{
			name: "regexp match",
			m:    &Matcher{regexp: regexp.MustCompile("^[1-9A-HJ-NP-Za-km-z]{44}$")},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: true,
		},
		{
			name: "regexp no match",
			m:    &Matcher{regexp: regexp.MustCompile("^[1-9A-HJ-NP-Za-km-z]{32}$")},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: false,
		},
		{
			name: "all conditions match",
			m:    &Matcher{prefix: "5vG", suffix: "AmR", regexp: regexp.MustCompile("^[1-9A-HJ-NP-Za-km-z]{44}$")},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: true,
		},
		{
			name: "all conditions no match - wrong prefix",
			m:    &Matcher{prefix: "ABC", suffix: "AmR", regexp: regexp.MustCompile("^[1-9A-HJ-NP-Za-km-z]{44}$")},
			args: args{address: "5vGRmNxXzZApNE8YwTuTxPzYdnmZ3GXio89hTF8VTAmR"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Match(tt.args.address); got != tt.want {
				t.Errorf("Matcher.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
