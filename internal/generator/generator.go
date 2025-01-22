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

var (
	totalGenerated uint64
	totalFound     uint64
)

func Start(ctx context.Context, cfg config.Config) {
	var wg sync.WaitGroup
	matcher := matcher.NewMatcher(cfg)

	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go worker(ctx, &wg, matcher, cfg)
	}

	wg.Wait()
	logger.Info("Generation completed",
		"total", atomic.LoadUint64(&totalGenerated),
		"found", atomic.LoadUint64(&totalFound),
	)
}

func worker(ctx context.Context, wg *sync.WaitGroup, m *matcher.Matcher, cfg config.Config) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			generateAndCheck(m, cfg)
		}
	}
}

func generateAndCheck(m *matcher.Matcher, cfg config.Config) {
	wallet := solana.NewWallet()
	address := wallet.PublicKey().String()

	atomic.AddUint64(&totalGenerated, 1)

	if m.Match(address) {
		handleMatch(wallet, cfg)
	}
}

func handleMatch(wallet *solana.Wallet, cfg config.Config) {
	count := atomic.AddUint64(&totalFound, 1)

	logger.Info("Vanity address found",
		"address", wallet.PublicKey(),
		"private", wallet.PrivateKey,
		"attempts", atomic.LoadUint64(&totalGenerated),
	)

	if cfg.NumAddresses > 0 && count >= uint64(cfg.NumAddresses) {
		logger.Info("Desired count reached - stopping generation")
		os.Exit(0)
	}
}
