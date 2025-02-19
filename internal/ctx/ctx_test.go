package ctx

import (
	"context"
	"testing"
	"time"

	"github.com/tmlnv/sanity/internal/config"
)

func TestCreateContext(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name            string
		args            args
		wantCtxDeadline bool // true if we expect context.WithTimeout, false for context.WithCancel
	}{
		{
			name: "with timeout",
			args: args{
				cfg: config.Config{
					Timeout: 5 * time.Second,
				},
			},
			wantCtxDeadline: true,
		},
		{
			name: "without timeout",
			args: args{
				cfg: config.Config{
					Timeout: 0,
				},
			},
			wantCtxDeadline: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := CreateContext(tt.args.cfg)
			defer cancel()

			// Verify the context type
			_, hasDeadline := ctx.Deadline()
			if hasDeadline != tt.wantCtxDeadline {
				t.Errorf("CreateContext() context deadline = %v, want %v", hasDeadline, tt.wantCtxDeadline)
			}

			// Verify the context can be cancelled
			cancel()
			if err := ctx.Err(); err != context.Canceled {
				t.Errorf("Context was not properly cancelled, got error: %v", err)
			}
		})
	}
}
