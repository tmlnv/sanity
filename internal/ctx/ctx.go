package ctx

import (
	"context"

	"github.com/tmlnv/sanity/internal/config"
)

// createContext creates a new context with a timeout if it is set in the config.
// If not, it creates a context with cancel function.
func CreateContext(cfg config.Config) (context.Context, context.CancelFunc) {
	if cfg.Timeout > 0 {
		return context.WithTimeout(context.Background(), cfg.Timeout)
	}
	return context.WithCancel(context.Background())
}
