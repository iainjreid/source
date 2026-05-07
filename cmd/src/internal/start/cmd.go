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

package start

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/iainjreid/source/cmd/src/internal/cli"
	"github.com/iainjreid/source/db/postgres"
	"github.com/iainjreid/source/db/postgres/storer"
	"github.com/iainjreid/source/internal/logger"
	"github.com/iainjreid/source/ssh"
	"github.com/iainjreid/source/web"
	"golang.org/x/sync/errgroup"
)

var usage = `Usage: 
    src start [--db <uri>] [-q | --quiet] [-v | --verbose] [-j | --json]
	          [-h | --help] [-i | --ssh-id-path <path>] [-I | --ssh-id <string>]
			  [--ssh-port <number>] [--http-port <number>] [--debug]

Start a Source instance against the specifed database.

Options:
    --db <uri>                  Specify the database URI
    -q, --quiet                 Suppress all non-error logs
    -v, --verbose               Print all logs
    -j, --json                  Display JSON output
    -h, --help                  Display this message

Additional Options:
    -i, --ssh-id-path <path>    The path to an SSH key
    -I, --ssh-id <string>       A plaintext SSH key
    --ssh-port <number>         Specify the SSH port
    --http-port <number>        Specify the HTTP port
    --debug                     Enable debugging

Example:
    $ src start --db postgresql://postgres@localhost
`

func Cmd(args []string) {
	flag := flag.NewFlagSet("start", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Run 'src start --help' for usage.\n")
	}

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		db        = cli.String(flag, "", "db", "")
		quiet     = cli.Bool(flag, false, "quiet", "q")
		verbose   = cli.Bool(flag, false, "verbose", "v")
		json      = cli.Bool(flag, false, "json", "j")
		help      = cli.Bool(flag, false, "help", "h")
		sshIdPath = cli.String(flag, "", "ssh-id-path", "i")
		sshId     = cli.String(flag, "", "ssh-id", "I")
		sshPort   = cli.Int(flag, 2222, "ssh-port", "")
		httpPort  = cli.Int(flag, 8080, "http-port", "")
		debug     = cli.Bool(flag, false, "debug", "")
	)

	flag.Parse(args)

	if *help {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if *db == "" {
		fmt.Fprint(os.Stderr, "--db is required")
		os.Exit(1)
	}

	if *quiet && *verbose {
		fmt.Fprint(os.Stderr, "-q/--quiet can't be used with -v/--verbose")
		os.Exit(1)
	}

	if *sshIdPath != "" && *sshId != "" {
		fmt.Fprint(os.Stderr, "-i/--ssh-id-path can't be used with -I/--ssh-id")
		os.Exit(1)
	}

	if *quiet {
		logger.Init(slog.LevelError, *json, *debug, nil)
	}

	if *sshIdPath != "" {
		data, err := os.ReadFile(*sshIdPath)
		if err != nil {
			fmt.Fprint(os.Stderr, "Unable to read SSH private key")
			os.Exit(1)
		}
		*sshId = string(data)
	}

	if *sshId != "" {
		if err := ssh.Init(*sshId); err != nil {
			fmt.Fprint(os.Stderr, "Unable to parse PEM encoded SSH private key")
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pg, err := postgres.Connect(ctx, *db)
	if err != nil {
		log.Fatalf("Error whilst connecting to DB: %s", err)
	}

	if err := pg.HardReset(ctx); err != nil {
		log.Fatalf("Failed to reset DB: %s", err)
	}

	if err := pg.EnsureReady(ctx); err != nil {
		log.Fatalf("Failed to setup DB: %s", err)
	}

	wg := new(errgroup.Group)

	storage := storer.NewStorage(pg.Pool)
	wg.Go(func() error {
		log.Println("Starting Web server")
		return web.NewServer(storage, *httpPort)
	})

	if *sshId != "" {
		wg.Go(func() error {
			log.Println("Starting SSH server")
			return ssh.NewServer(storage, *sshPort)
		})
	}

	err = wg.Wait()
	if err != nil {
		log.Println(err)
	}
}
