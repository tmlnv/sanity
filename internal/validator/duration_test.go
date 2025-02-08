package validator

import (
	"reflect"
	"testing"
	"time"
)

func TestValidateTimeout(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name:    "empty string",
			args:    args{s: ""},
			want:    0,
			wantErr: false,
		},
		{
			name:    "zero value",
			args:    args{s: "0"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "valid seconds as number",
			args:    args{s: "30"},
			want:    30 * time.Second,
			wantErr: false,
		},
		{
			name:    "valid duration string",
			args:    args{s: "5m"},
			want:    5 * time.Minute,
			wantErr: false,
		},
		{
			name:    "invalid duration format",
			args:    args{s: "invalid"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "complex duration",
			args:    args{s: "1h30m"},
			want:    90 * time.Minute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateTimeout(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDuration(t *testing.T) {
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
			name:    "valid duration",
			args:    args{s: "30s"},
			wantErr: false,
		},
		{
			name:    "invalid duration",
			args:    args{s: "invalid"},
			wantErr: true,
		},
		{
			name:    "complex duration",
			args:    args{s: "2h45m"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDuration(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDuration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
