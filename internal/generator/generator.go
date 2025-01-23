package generator

import (
	"context"
	"os"
	"sync"
	"sync/atomic"

	"github.com/gagliardetto/solana-go"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/matcher"
)

type Stats struct {
	Attempts uint64
	Found    uint64
}

type StatsUpdate struct {
	Stats      Stats
	LastResult string
}

func Start(ctx context.Context, cfg config.Config, updateChan chan<- StatsUpdate) {
	var (
		wg         sync.WaitGroup
		matcher    = matcher.NewMatcher(cfg)
		totalFound uint64
	)

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					generateAndCheck(matcher, cfg, updateChan, &totalFound)
				}
			}
		}()
	}

	wg.Wait()
	logger.Info("Generation completed", "found", totalFound)
}

func generateAndCheck(m *matcher.Matcher, cfg config.Config, updateChan chan<- StatsUpdate, totalFound *uint64) {
	wallet := solana.NewWallet()
	address := wallet.PublicKey().String()

	if m.Match(address) {
		count := atomic.AddUint64(totalFound, 1)
		logger.Info("Vanity address found",
			"address", address,
			"private", wallet.PrivateKey,
			"attempts", count,
		)

		if updateChan != nil {
			updateChan <- StatsUpdate{
				Stats: Stats{
					Attempts: count,
					Found:    count,
				},
				LastResult: address,
			}
		}

		if cfg.NumAddresses > 0 && count >= uint64(cfg.NumAddresses) {
			logger.Info("Desired count reached - stopping generation")
			os.Exit(0)
		}
	}
}
