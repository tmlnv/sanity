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
	"github.com/tmlnv/sanity/internal/saver"
)

// ─── Public types ──────────────────────────────────────────────────────────────

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

// ─── Entry point ───────────────────────────────────────────────────────────────

// Start launches the generator and blocks until the job is finished or the
// supplied ctx is cancelled.
func Start(
	parentCtx context.Context,
	cfg config.Config,
	update chan<- StatsUpdate,
	isTUI bool,
) {
	// Derived context so that the cancellation signal is owned here.
	// Cancelling the parent will still propagate.
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	var (
		totalAttempts uint64
		totalFound    uint64
		m             = matcher.NewMatcher(cfg)
		wg            sync.WaitGroup
	)

	// ── Progress / UI updater ────────────────────────────────────────────────
	if update != nil {
		go progressTicker(ctx, update, &totalAttempts, &totalFound)
	}

	// ── Periodic CLI logger  ─────────────────────────────────────────────────
	if !isTUI && cfg.LogInterval > 0 {
		wg.Add(1)
		go logProgress(ctx, &wg, cfg.LogInterval, &totalAttempts, &totalFound)
	}

	// ── Workers ──────────────────────────────────────────────────────────────
	for range cfg.Concurrency {
		wg.Add(1)
		go worker(ctx, &wg, cfg, m,
			&totalAttempts, &totalFound,
			update, cancel, isTUI)
	}

	wg.Wait()

	// ── Final stats to CLI + UI channel ──────────────────────────────────────
	if !isTUI {
		logger.Info("Generation completed",
			"attempts", atomic.LoadUint64(&totalAttempts),
			"found", atomic.LoadUint64(&totalFound))
	}
	if update != nil {
		update <- StatsUpdate{
			Stats: Stats{
				Attempts: atomic.LoadUint64(&totalAttempts),
				Found:    atomic.LoadUint64(&totalFound),
			},
			IsFinished: true,
		}
		close(update)
	}
}

// ─── Worker ───────────────────────────────────────────────────────────────────

func worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	cfg config.Config,
	m *matcher.Matcher,
	totalAttempts, totalFound *uint64,
	update chan<- StatsUpdate,
	cancel context.CancelFunc,
	isTUI bool,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			wallet, addr := generate(totalAttempts)
			if m.Match(addr) {
				onMatch(addr, wallet, isTUI, totalFound, update)
				// **Single exit decision point**
				if cfg.NumAddresses > 0 && atomic.LoadUint64(totalFound) >= uint64(cfg.NumAddresses) {
					cancel()
				}
			} else if update != nil && isTUI {
				update <- StatsUpdate{
					Stats: Stats{
						Attempts: atomic.LoadUint64(totalAttempts),
						Found:    atomic.LoadUint64(totalFound),
					},
					LastGenerated: addr,
				}
			}
		}
	}
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// generate returns the new wallet and address
func generate(totalAttempts *uint64) (*solana.Wallet, string) {
	wallet := solana.NewWallet()
	address := wallet.PublicKey().String()
	atomic.AddUint64(totalAttempts, 1)
	return wallet, address
}

// Handles everything that should happen once a match is found.
func onMatch(
	addr string,
	wallet *solana.Wallet,
	isTUI bool,
	totalFound *uint64,
	update chan<- StatsUpdate,
) {
	n := atomic.AddUint64(totalFound, 1)

	// persistant saving first
	saver.SaveKeyPair(addr, wallet.PrivateKey.String())

	if !isTUI {
		logger.Info("Vanity address found",
			"address", addr,
			"sequence", n)
	}

	if update != nil {
		update <- StatsUpdate{
			Stats: Stats{
				Attempts: 0, // will be populated by progressTicker
				Found:    n,
			},
			LastMatch: addr,
		}
	}
}

// Emits periodic updates to the UI *and* keeps Attempts / Found fresh.
func progressTicker(
	ctx context.Context,
	update chan<- StatsUpdate,
	totalAttempts, totalFound *uint64,
) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			update <- StatsUpdate{
				Stats: Stats{
					Attempts: atomic.LoadUint64(totalAttempts),
					Found:    atomic.LoadUint64(totalFound),
				},
			}
		}
	}
}

// Periodic CLI logger (for non‑TUI mode).
func logProgress(
	ctx context.Context,
	wg *sync.WaitGroup,
	interval time.Duration,
	totalAttempts, totalFound *uint64,
) {
	defer wg.Done()
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			logger.Info("Progress update",
				"attempts", atomic.LoadUint64(totalAttempts),
				"found", atomic.LoadUint64(totalFound))
		}
	}
}
