package main

import (
	"fmt"
	"log/slog"
	"nats/internal/broker"
	"nats/internal/db"
	"os"
	"os/signal"
	"syscall"
)

const init_db = false

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("Start program")

	db.DBConnection, _ = db.NewDBConnection()

	defer db.DBConnection.Close()

	db.TestItem(db.DBConnection)

	if init_db {
		db.ExecMigration(db.DBConnection, "../internal/db/migrations/db_init.sql")
	}

	broker.Publish()
	broker.Subscribe_to_channel()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sygType := <-stop

	logger.Info(fmt.Sprintf("End signal: %s", sygType))
}
