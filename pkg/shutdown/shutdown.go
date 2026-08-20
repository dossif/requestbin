package shutdown

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// WaitForSignal spawns a goroutine that blocks until SIGINT or SIGTERM is
// received, then closes done to signal graceful shutdown.
func WaitForSignal(log *slog.Logger, done chan bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Info("os-signal", "signal", sig)
		close(done)
	}()
}
