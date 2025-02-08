package validator

import "testing"

func TestValidateNumber(t *testing.T) {
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
			name:    "zero value",
			args:    args{s: "0"},
			wantErr: false,
		},
		{
			name:    "positive integer",
			args:    args{s: "42"},
			wantErr: false,
		},
		{
			name:    "negative integer",
			args:    args{s: "-42"},
			wantErr: false,
		},
		{
			name:    "invalid number - decimal",
			args:    args{s: "3.14"},
			wantErr: true,
		},
		{
			name:    "invalid number - alphabetic",
			args:    args{s: "abc"},
			wantErr: true,
		},
		{
			name:    "invalid number - special characters",
			args:    args{s: "123$"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateNumber(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
