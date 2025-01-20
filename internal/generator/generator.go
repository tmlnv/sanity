package generator

import (
	"context"
	"encoding/hex"
	"sync"

	"github.com/charmbracelet/log"
	solana "github.com/gagliardetto/solana-go"
	"github.com/tmlnv/sanity/internal/config"
	"github.com/tmlnv/sanity/internal/logger"
	"github.com/tmlnv/sanity/internal/matcher"
)

// StartGeneration starts the vanity address generation process.
func StartGeneration(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	if cfg.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
	}
	defer cancel()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.Threads)

	for {
		select {
		case <-ctx.Done():
			log.Info("Generation stopped", "reason", ctx.Err())
			return
		default:
			semaphore <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() { <-semaphore }()

				// Generate a new Solana wallet
				account := solana.NewWallet()
				publicKey := account.PublicKey().String()

				// Check if the address matches the desired pattern
				if matcher.Match(publicKey, cfg.Prefix, cfg.Suffix, cfg.Regex) {
					logger.LogResult(publicKey, hex.EncodeToString(account.PrivateKey))
				}
			}()
		}
	}
}
