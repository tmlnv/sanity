package validator

import "testing"

func TestValidateSolana(t *testing.T) {
	type args struct {
		prefix        string
		suffix        string
		regexpPattern string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty prefix and suffix",
			args:    args{prefix: "", suffix: "", regexpPattern: ""},
			wantErr: false,
		},
		{
			name:    "valid prefix only",
			args:    args{prefix: "123456789", suffix: "", regexpPattern: ""},
			wantErr: false,
		},
		{
			name:    "valid suffix only",
			args:    args{prefix: "", suffix: "ABC123", regexpPattern: ""},
			wantErr: false,
		},
		{
			name:    "valid prefix and suffix",
			args:    args{prefix: "123", suffix: "ABC", regexpPattern: ""},
			wantErr: false,
		},
		{
			name:    "invalid prefix characters",
			args:    args{prefix: "123$456", suffix: "", regexpPattern: ""},
			wantErr: true,
		},
		{
			name:    "invalid suffix characters",
			args:    args{prefix: "", suffix: "ABC@123", regexpPattern: ""},
			wantErr: true,
		},
		{
			name:    "combined length exceeds limit",
			args:    args{prefix: "123456789012345", suffix: "987654321098765432109876543210987654321098", regexpPattern: ""},
			wantErr: true,
		},
		{
			name:    "valid regexp pattern",
			args:    args{prefix: "", suffix: "", regexpPattern: "[A-Za-z0-9]+"},
			wantErr: false,
		},
		{
			name:    "invalid regexp pattern",
			args:    args{prefix: "", suffix: "", regexpPattern: "["},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateSolana(tt.args.prefix, tt.args.suffix, tt.args.regexpPattern); (err != nil) != tt.wantErr {
				t.Errorf("ValidateSolana() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateBase58String(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{s: ""},
			wantErr: false,
		},
		{
			name:    "valid base58 string",
			args:    args{s: "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"},
			wantErr: false,
		},
		{
			name:    "invalid character - 0",
			args:    args{s: "0123456789"},
			wantErr: true,
		},
		{
			name:    "invalid character - O",
			args:    args{s: "ABCDEFGHIJKLMNO"},
			wantErr: true,
		},
		{
			name:    "invalid character - l",
			args:    args{s: "abcdefghijklmno"},
			wantErr: true,
		},
		{
			name:    "invalid character - special",
			args:    args{s: "ABC$123"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBase58String(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("validateBase58String() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateRegexpPattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty pattern",
			args:    args{pattern: ""},
			wantErr: false,
		},
		{
			name:    "valid simple pattern",
			args:    args{pattern: "[0-9]+"},
			wantErr: false,
		},
		{
			name:    "valid complex pattern",
			args:    args{pattern: "^[A-Za-z0-9]{10,20}$"},
			wantErr: false,
		},
		{
			name:    "invalid pattern - unclosed bracket",
			args:    args{pattern: "["},
			wantErr: true,
		},
		{
			name:    "invalid pattern - invalid quantifier",
			args:    args{pattern: "a{-1}(nd1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRegexpPattern(tt.args.pattern); (err != nil) != tt.wantErr {
				t.Errorf("validateRegexpPattern() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
