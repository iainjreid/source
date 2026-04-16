package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/iainjreid/source/db/postgres"
	"github.com/iainjreid/source/db/postgres/storer"
	"github.com/iainjreid/source/internal/logger"
	"github.com/iainjreid/source/ssh"
	"github.com/iainjreid/source/web"
	"golang.org/x/sync/errgroup"
)

// config holds the program configuration settings that are read from
// the command-line flags.
type config struct {
	dbUri string
}

func main() {
	var cfg config

	// Define the accepted command-line flags and read them into the config struct
	flag.StringVar(&cfg.dbUri, "db-uri", "", "Database URI to connect to")

	// Parse the command-line flags
	flag.Parse()

	if cfg.dbUri == "" {
		invalidConfig("Please specify a database URI")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Init(slog.LevelDebug, false, nil)

	db, err := postgres.Connect(ctx, cfg.dbUri)
	if err != nil {
		log.Fatalf("Error whilst connecting to DB: %s", err)
	}

	if err := db.HardReset(ctx); err != nil {
		log.Fatalf("Failed to reset DB: %s", err)
	}

	if err := db.EnsureReady(ctx); err != nil {
		log.Fatalf("Failed to setup DB: %s", err)
	}

	wg := new(errgroup.Group)

	storage := storer.NewStorage(db.Pool)
	wg.Go(func() error {
		log.Println("Starting Web server")
		return web.NewServer(storage)
	})

	wg.Go(func() error {
		log.Println("Starting SSH server")
		return ssh.NewServer(storage)
	})

	err = wg.Wait()
	if err != nil {
		log.Println(err)
	}
}

func invalidConfig(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	flag.Usage()
	os.Exit(1)
}
