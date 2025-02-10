package generator

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/ctx"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/matcher"
)

func TestStart(t *testing.T) {
	logger.Init(os.DevNull)
	defer logger.Close()

	tests := []struct {
		name                 string
		cfg                  config.Config
		isTui                bool
		expectedAtLeastFound uint64
	}{
		{
			name: "Basic generation with timeout",
			cfg: config.Config{
				Timeout:     time.Second,
				Concurrency: 2,
				Prefix:      "3",
			},
			isTui:                false,
			expectedAtLeastFound: 1,
		},
		{
			name: "Generation with simple prefix (will always be found)",
			cfg: config.Config{
				Timeout:     time.Second,
				Concurrency: 2,
				Prefix:      "12",
			},
			isTui:                false,
			expectedAtLeastFound: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := ctx.CreateContext(tt.cfg)
			defer cancel()

			updateChan := make(chan StatsUpdate, 10)
			done := make(chan struct{})

			// Start generator in a goroutine
			go func() {
				Start(ctx, tt.cfg, updateChan, tt.isTui)
				close(done)
			}()

			// Monitor updates and check for completion
			var lastStats Stats
			for {
				select {
				case update, ok := <-updateChan:
					if !ok {
						t.Error("Update channel closed unexpectedly")
						return
					}
					lastStats = update.Stats
					if update.IsFinished {
						if lastStats.Found < tt.expectedAtLeastFound {
							t.Errorf("Expected at least %d found addresses, got %d", tt.expectedAtLeastFound, lastStats.Found)
						}
						return
					}
				case <-done:
					return
				case <-time.After(2 * time.Second):
					t.Error("Test timed out")
					return
				}
			}
		})
	}
}

func Test_worker(t *testing.T) {
	var (
		totalAttempts uint64
		totalFound    uint64
		wg            sync.WaitGroup
	)

	cfg := config.Config{
		Timeout:      time.Second,
		Concurrency:  1,
		Prefix:       "3",
		NumAddresses: 1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	updateChan := make(chan StatsUpdate, 10)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	m := matcher.NewMatcher(cfg)

	wg.Add(1)
	go worker(ctx, &wg, cfg, ticker, updateChan, &totalAttempts, &totalFound, m, false)

	// Wait for either completion or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Worker completed successfully
	case <-time.After(2 * time.Second):
		t.Error("Worker test timed out")
	}
}

func Test_generateAndCheck(t *testing.T) {
	var (
		totalAttempts uint64
		totalFound    uint64
	)

	cfg := config.Config{
		Prefix: "test", // Add a prefix to test matching
	}

	m := matcher.NewMatcher(cfg)
	updateChan := make(chan StatsUpdate, 1)

	// Test single generation
	generateAndCheck(m, updateChan, &totalFound, &totalAttempts, false)

	// Verify that attempts were incremented
	if totalAttempts == 0 {
		t.Error("Expected attempts to be greater than 0")
	}
}
