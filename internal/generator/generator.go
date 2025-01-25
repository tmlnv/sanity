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
	Stats         Stats
	LastGenerated string
	LastMatch     string
	IsFinished    bool
}

func Start(ctx context.Context, cfg config.Config, updateChan chan<- StatsUpdate, isTui bool) {
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
		go worker(ctx, &wg, cfg, ticker, updateChan, &totalAttempts, &totalFound, matcher, isTui)
	}

	wg.Wait()
	if !isTui {
		logger.Info("Generation completed",
			"attempts", atomic.LoadUint64(&totalAttempts),
			"found", totalFound,
		)
	}

	if updateChan != nil {
		updateChan <- StatsUpdate{
			Stats: Stats{
				Attempts: atomic.LoadUint64(&totalAttempts),
				Found:    atomic.LoadUint64(&totalFound),
			},
			IsFinished: true,
		}
		close(updateChan)
	}
}

// In the worker function, remove sending IsFinished on cancellation
func worker(ctx context.Context, wg *sync.WaitGroup, cfg config.Config, ticker *time.Ticker, updateChan chan<- StatsUpdate, totalAttempts *uint64, totalFound *uint64, matcher *matcher.Matcher, isTui bool) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			if !isTui {
				logger.Debug("Generation stopped due to cancellation", "Context", ctx)
			}
			return // Stop when context is canceled
		case <-ticker.C:
			// Send periodic update
			if updateChan != nil {
				updateChan <- StatsUpdate{
					Stats: Stats{
						Attempts: atomic.LoadUint64(totalAttempts),
						Found:    atomic.LoadUint64(totalFound),
					},
				}
			}
		default:
			atomic.AddUint64(totalAttempts, 1)
			generateAndCheck(matcher, updateChan, totalFound, totalAttempts, isTui)
			if cfg.NumAddresses > 0 && atomic.LoadUint64(totalFound) >= uint64(cfg.NumAddresses) {
				if !isTui {
					logger.Debug("Desired count reached - stopping generation")
				}
				return
			}
		}
	}
}

func generateAndCheck(m *matcher.Matcher, updateChan chan<- StatsUpdate, totalFound *uint64, totalAttempts *uint64, isTui bool) {
	var stats StatsUpdate

	wallet := solana.NewWallet()
	address := wallet.PublicKey().String()

	count := atomic.LoadUint64(totalFound)
	attempts := atomic.LoadUint64(totalAttempts)

	stats.Stats = Stats{
		Attempts: attempts,
		Found:    count,
	}
	stats.LastGenerated = address

	if m.Match(address) {

		count = atomic.AddUint64(totalFound, 1)
		stats.Stats.Found = count
		stats.LastMatch = address

		if !isTui {
			logger.Info("Vanity address found",
				"address", address,
				"private", wallet.PrivateKey.String(),
				"attempts", attempts,
			)
		}
	}

	if updateChan != nil {
		// Send update before checking for exit condition
		updateChan <- stats
	}
}
