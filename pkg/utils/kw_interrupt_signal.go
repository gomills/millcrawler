package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
)

// SetKeyboardInterruptSignal relays the Ctrl+C OS' signal into the program
func SetKeyboardInterruptSignal(ctx context.Context, cancel context.CancelFunc) {
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	select {
	case <-ctrlC:
		log.Print("<||> Ctrl+C pressed <||>")
		cancel()
		return
	case <-ctx.Done():
		return
	}

}
