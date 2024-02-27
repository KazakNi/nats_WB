package main

import (
	"flag"
	"fmt"
	"log/slog"
	"nats/api"
	"nats/internal/broker"
	"nats/internal/db"
	"nats/internal/storage"
	"os"
	"os/signal"
	"syscall"
)

const init_db = false

var pathToMigrations string

func main() {

	flag.StringVar(&pathToMigrations, "path to migrations", "../internal/db/migrations/db_init.sql", "migrations folder path for db init")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	logger.Info("Start program")
	db.DBConnection, _ = db.NewDBConnection()

	storage.AppCache = storage.New()

	go api.NewRouter(storage.AppCache)

	defer db.DBConnection.Close()

	if init_db {
		db.ExecMigration(db.DBConnection, pathToMigrations)
	}

	go broker.Publish()
	broker.Subscribe_to_channel()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	sygType := <-stop

	logger.Info(fmt.Sprintf("End signal: %s", sygType))
}
