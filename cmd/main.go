package main

import (
	"log/slog"
	"nats/internal/broker"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("Start program")

	broker.Publish()
	broker.Subscribe_to_channel()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sygType := <-stop

	logger.Info("End program", sygType)
}
