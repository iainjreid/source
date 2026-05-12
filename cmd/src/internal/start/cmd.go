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

func Cmd(ctx context.Context, args []string) {
	cmd := cli.New("src start")

	if len(args) == 0 {
		cmd.Usage()
	}

	var (
		db        = cmd.String("", "db", "")
		quiet     = cmd.Bool(false, "quiet", "q")
		verbose   = cmd.Bool(false, "verbose", "v")
		json      = cmd.Bool(false, "json", "j")
		help      = cmd.Bool(false, "help", "h")
		sshIdPath = cmd.String("", "ssh-id-path", "i")
		sshId     = cmd.String("", "ssh-id", "I")
		sshPort   = cmd.Int(2222, "ssh-port", "")
		httpPort  = cmd.Int(8080, "http-port", "")
		debug     = cmd.Bool(false, "debug", "")
	)

	cmd.Parse(args)

	if *help {
		cli.Fatal(1, usage)
	}

	if level, err := cli.ResolveLoggerFlags(quiet, verbose); err != nil {
		cmd.ExplainUsage(err.Error())
	} else {
		logger.Init(level, *json, *debug, nil)
	}

	if *db == "" {
		cmd.ExplainUsage("--db is required")
	}

	if *sshIdPath != "" && *sshId != "" {
		cmd.ExplainUsage("-i/--ssh-id-path can't be used with -I/--ssh-id")
	}

	if *sshIdPath != "" {
		data, err := os.ReadFile(*sshIdPath)
		if err != nil {
			cli.Fatal(cli.Failure, "Unable to read SSH private key")
		}
		*sshId = string(data)
	}

	if *sshId != "" {
		if err := ssh.Init(*sshId); err != nil {
			cli.Fatal(cli.Failure, "Unable to parse PEM encoded SSH private key")
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pg, err := postgres.Connect(ctx, *db)
	if err != nil {
		cli.Fatalf(cli.Failure, "Error whilst connecting to DB: %s", err)
	}

	wg := new(errgroup.Group)

	storage := storer.NewStorage(pg.Pool, "")
	wg.Go(func() error {
		slog.Info("Starting Web server")
		return web.NewServer(storage, *httpPort)
	})

	if *sshId != "" {
		wg.Go(func() error {
			slog.Info("Starting SSH server")
			return ssh.NewServer(storage, *sshPort)
		})
	}

	err = wg.Wait()
	if err != nil {
		slog.Error(err.Error())
	}
}
