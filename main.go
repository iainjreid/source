// Copyright 2026 Iain J. Reid
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	debug     bool
	logLevel  logger.Level
	logFormat logger.Format
	dbUri     string
	sshKey    string
	sshPort   int
	webPort   int
}

func main() {
	var cfg config

	// Define the accepted command-line flags and read them into the config struct
	flag.StringVar(&cfg.dbUri, "db-uri", "", "Database URI to connect to (required)")
	flag.StringVar(&cfg.sshKey, "ssh-key", "", "A PEM encoded private key used to enable SSH access")
	flag.IntVar(&cfg.sshPort, "ssh-port", 2222, "The port to accept SSH connections through")

	// Logging
	flag.Var(&cfg.logLevel, "log-level", "The lowest level of logs to print")
	flag.Var(&cfg.logFormat, "log-format", "The format with which to print the logs")

	// Web
	flag.IntVar(&cfg.webPort, "web-port", 8080, "The port to serve the web interface over")

	// Debugging
	flag.BoolVar(&cfg.debug, "debug", false, "Enable debugging")

	// Parse the command-line flags
	flag.Parse()

	if cfg.dbUri == "" {
		invalidConfig("Please specify a database URI")
	}

	if cfg.sshKey != "" {
		if err := ssh.Init(cfg.sshKey); err != nil {
			fatalError("Unable to parse PEM encoded SSH private key")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Init(cfg.logLevel, cfg.logFormat, cfg.debug, nil)

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
		return web.NewServer(storage, cfg.webPort)
	})

	if cfg.sshKey != "" {
		wg.Go(func() error {
			log.Println("Starting SSH server")
			return ssh.NewServer(storage, cfg.sshPort)
		})
	}

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

func fatalError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
