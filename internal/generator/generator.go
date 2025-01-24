package generator

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/matcher"
)

type Stats struct {
	Attempts uint64 // Total attempts across all workers
	Found    uint64 // Total found addresses
}

type StatsUpdate struct {
	Stats      Stats
	LastResult string
}

func Start(ctx context.Context, calcel context.CancelFunc, cfg config.Config, updateChan chan<- StatsUpdate) {
	var (
		wg            sync.WaitGroup
		matcher       = matcher.NewMatcher(cfg)
		totalFound    uint64
		totalAttempts uint64
	)

	// Ticker for periodic updates
	ticker := time.NewTicker(500 * time.Millisecond) // Update every 500ms
	defer ticker.Stop()

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(ctx, calcel, &wg, cfg, ticker, updateChan, totalAttempts, totalFound, matcher)
	}

	wg.Wait()
	logger.Info("Generation completed",
		"attempts", atomic.LoadUint64(&totalAttempts),
		"found", totalFound,
	)

	if updateChan != nil {
		close(updateChan) // Close the channel to signal completion
	}
}

func worker(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, cfg config.Config, ticker *time.Ticker, updateChan chan<- StatsUpdate, totalAttempts uint64, totalFound uint64, matcher *matcher.Matcher) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logger.Info("Generation stopped due to cancellation", "Context", ctx)
			return // Stop when context is canceled
		case <-ticker.C:
			// Send periodic update
			if updateChan != nil {
				updateChan <- StatsUpdate{
					Stats: Stats{
						Attempts: atomic.LoadUint64(&totalAttempts),
						Found:    atomic.LoadUint64(&totalFound),
					},
				}
			}
		default:
			atomic.AddUint64(&totalAttempts, 1)
			generateAndCheck(matcher, updateChan, &totalFound, &totalAttempts)
			if cfg.NumAddresses > 0 && totalFound >= uint64(cfg.NumAddresses) {
				logger.Info("Desired count reached - stopping generation")
				cancel() // Stop all workers
				return
			}
		}
	}
}

func generateAndCheck(m *matcher.Matcher, updateChan chan<- StatsUpdate, totalFound *uint64, totalAttempts *uint64) (isFinished bool) {
	wallet := solana.NewWallet()
	address := wallet.PublicKey().String()

	if m.Match(address) {
		count := atomic.AddUint64(totalFound, 1)
		attempts := atomic.LoadUint64(totalAttempts)

		logger.Info("Vanity address found",
			"address", address,
			"private", wallet.PrivateKey.String(),
			"attempts", attempts,
		)

		if updateChan != nil {
			// Send update before checking for exit condition
			updateChan <- StatsUpdate{
				Stats: Stats{
					Attempts: attempts,
					Found:    count,
				},
				LastResult: address,
			}
		}
	}
	return
}
