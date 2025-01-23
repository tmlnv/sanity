package generator

import (
	"context"
	"os"
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

func Start(ctx context.Context, cfg config.Config, updateChan chan<- StatsUpdate) {
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
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
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
					generateAndCheck(matcher, cfg, updateChan, &totalFound, &totalAttempts)
				}
			}
		}()
	}

	wg.Wait()
	logger.Info("Generation completed",
		"attempts", atomic.LoadUint64(&totalAttempts),
		"found", totalFound,
	)
}

func generateAndCheck(m *matcher.Matcher, cfg config.Config, updateChan chan<- StatsUpdate, totalFound *uint64, totalAttempts *uint64) {
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
			updateChan <- StatsUpdate{
				Stats: Stats{
					Attempts: attempts,
					Found:    count,
				},
				LastResult: address,
			}
		}

		if cfg.NumAddresses > 0 && count >= uint64(cfg.NumAddresses) {
			logger.Info("Desired count reached - stopping generation")
			if updateChan != nil {
				close(updateChan) // Close the channel to signal completion
			}
			os.Exit(0)
		}
	}
}
