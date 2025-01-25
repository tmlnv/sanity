package ctx

import (
	"context"

	"github.com/tmlnv/sanity/internal/config"
)

func CreateContext(cfg config.Config) (context.Context, context.CancelFunc) {
	if cfg.Timeout > 0 {
		return context.WithTimeout(context.Background(), cfg.Timeout)
	}
	return context.WithCancel(context.Background())
}
