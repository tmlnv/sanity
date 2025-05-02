package generator

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
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
	// Drain the update channel to prevent blocking
	go func() {
		for range updateChan {
		}
	}()

	m := matcher.NewMatcher(cfg)

	wg.Add(1)
	go worker(ctx, &wg, cfg, m, &totalAttempts, &totalFound, updateChan, cancel, false)

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

func Test_generate(t *testing.T) {
	type args struct {
		totalAttempts *uint64
	}

	// We can't predict the exact wallet or address that will be generated
	// since they're random, so we'll test the behavior instead
	t.Run("increments attempt counter", func(t *testing.T) {
		var attempts uint64 = 0
		wallet, addr := generate(&attempts)

		// Check that we got a non-nil wallet
		if wallet == nil {
			t.Error("generate() returned nil wallet")
		}

		// Check that we got a non-empty address
		if addr == "" {
			t.Error("generate() returned empty address")
		}

		// Check that the address is a valid Solana address (base58 encoded, starts with a number or letter)
		if len(addr) < 32 || len(addr) > 44 {
			t.Errorf("generate() returned address with invalid length: %d", len(addr))
		}

		// Check that the counter was incremented
		if attempts != 1 {
			t.Errorf("generate() did not increment attempt counter, got %d, want 1", attempts)
		}

		// Generate another address and check counter again
		generate(&attempts)
		if attempts != 2 {
			t.Errorf("generate() did not increment attempt counter correctly, got %d, want 2", attempts)
		}
	})
}

func Test_onMatch(t *testing.T) {
	// Setup test logger
	logger.Init(os.DevNull)
	defer logger.Close()

	t.Run("increments found counter and sends update", func(t *testing.T) {
		// Setup
		var found uint64 = 0
		wallet := solana.NewWallet()
		addr := wallet.PublicKey().String()
		updateChan := make(chan StatsUpdate, 1)

		// Call onMatch
		onMatch(addr, wallet, true, &found, updateChan)

		// Check that found counter was incremented
		if found != 1 {
			t.Errorf("onMatch() did not increment found counter, got %d, want 1", found)
		}

		// Check that update was sent
		select {
		case update := <-updateChan:
			if update.Stats.Found != 1 {
				t.Errorf("onMatch() sent update with incorrect found count, got %d, want 1", update.Stats.Found)
			}
			if update.LastMatch != addr {
				t.Errorf("onMatch() sent update with incorrect address, got %s, want %s", update.LastMatch, addr)
			}
		default:
			t.Error("onMatch() did not send update to channel")
		}
	})

	t.Run("handles nil update channel", func(t *testing.T) {
		// Setup
		var found uint64 = 0
		wallet := solana.NewWallet()
		addr := wallet.PublicKey().String()

		// This should not panic
		onMatch(addr, wallet, true, &found, nil)

		// Check that found counter was still incremented
		if found != 1 {
			t.Errorf("onMatch() did not increment found counter with nil channel, got %d, want 1", found)
		}
	})
}

func Test_progressTicker(t *testing.T) {
	t.Run("sends stats updates", func(t *testing.T) {
		// Setup
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var attempts uint64 = 100
		var found uint64 = 5
		updateChan := make(chan StatsUpdate, 2)

		// Start progressTicker in a goroutine
		go progressTicker(ctx, updateChan, &attempts, &found)

		// Wait for an update
		select {
		case update := <-updateChan:
			if update.Stats.Attempts != 100 {
				t.Errorf("progressTicker() sent update with incorrect attempts, got %d, want 100", update.Stats.Attempts)
			}
			if update.Stats.Found != 5 {
				t.Errorf("progressTicker() sent update with incorrect found, got %d, want 5", update.Stats.Found)
			}
		case <-time.After(time.Second):
			t.Error("progressTicker() did not send update within timeout")
		}

		// Cancel context to stop the ticker
		cancel()
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		// Setup with already cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var attempts uint64 = 100
		var found uint64 = 5
		updateChan := make(chan StatsUpdate, 1)

		// This should return quickly due to cancelled context
		done := make(chan struct{})
		go func() {
			progressTicker(ctx, updateChan, &attempts, &found)
			close(done)
		}()

		// Check that it returns
		select {
		case <-done:
			// Success - function returned
		case <-time.After(time.Second):
			t.Error("progressTicker() did not respect context cancellation")
		}
	})
}

func Test_logProgress(t *testing.T) {
	// Setup test logger
	logger.Init(os.DevNull)
	defer logger.Close()

	t.Run("respects context cancellation", func(t *testing.T) {
		// Setup with already cancelled context
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		var wg sync.WaitGroup
		wg.Add(1)

		var attempts uint64 = 100
		var found uint64 = 5

		// This should return quickly due to cancelled context
		done := make(chan struct{})
		go func() {
			logProgress(ctx, &wg, 100*time.Millisecond, &attempts, &found)
			close(done)
		}()

		// Check that it returns and decrements the WaitGroup
		select {
		case <-done:
			// Success - function returned
		case <-time.After(time.Second):
			t.Error("logProgress() did not respect context cancellation")
		}

		// Create a channel to signal when WaitGroup is done
		wgDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(wgDone)
		}()

		// Check that WaitGroup was decremented
		select {
		case <-wgDone:
			// Success - WaitGroup was decremented
		case <-time.After(time.Second):
			t.Error("logProgress() did not decrement WaitGroup")
		}
	})
}
