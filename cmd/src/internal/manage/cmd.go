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

package manage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/iainjreid/source/cmd/src/internal/cli"
	"github.com/iainjreid/source/db/postgres"
	"github.com/iainjreid/source/db/postgres/storer"
	"github.com/iainjreid/source/db/sql/shared"
	"github.com/iainjreid/source/git"
	"github.com/iainjreid/source/internal/logger"
)

var usage = `Usage:
    src manage [--db <uri>] [-q | --quiet] [-v | --verbose] [-j | --json]
               [-h | --help] [--setup] [--clone <uri>] [--debug]

Perform actions that can initialise and or modify your Source installation. All
changes are made directly through the storage layer, and do not need an active
instance running to be executed.

Options:
    --db <uri>                  Specify the database URI
    -q, --quiet                 Suppress all non-error logs
    -v, --verbose               Print all logs
    -j, --json                  Display JSON output
    -h, --help                  Display this message

Actions:
    --setup                     Ensures that the database contains the required
                                tables to operate. This flag can be used in
                                conjunction with all other actions.

    --clone <uri>               Specifies a remote repository to be cloned.

Additional Options:
    --debug                     Enable debugging

Example:
    $ src manage --db postgresql://postgres@localhost --setup
`

func Cmd(ctx context.Context, args []string) {
	cmd := cli.New("src manage")

	if len(args) == 0 {
		cmd.Usage()
	}

	var (
		db      = cmd.String("", "db", "")
		quiet   = cmd.Bool(false, "quiet", "q")
		verbose = cmd.Bool(false, "verbose", "v")
		json    = cmd.Bool(false, "json", "j")
		help    = cmd.Bool(false, "help", "h")
		debug   = cmd.Bool(false, "debug", "")
		setup   = cmd.Bool(false, "setup", "")
		clone   = cmd.String("", "clone", "c")
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

	pg, err := postgres.Connect(ctx, *db)

	if err != nil {
		cli.Fatal(cli.Failure, err.Error())
	}

	if *setup {
		pg.EnsureReady(ctx)
	}

	if *clone != "" {
		name, err := git.GetRepoName(*clone)
		if err != nil {
			panic(err)
		}

		slog.InfoContext(ctx, "creating postgres tables")
		if _, err := pg.Pool.Exec(ctx, fmt.Sprintf(shared.CreateRepoQuery, name, name)); err != nil {
			cli.Error("error whilst ensuring tables exist " + err.Error())
		}

		storage := storer.NewStorage(pg.Pool, name)

		repo := git.CloneRepo(storage, *clone)
		if err := repo.Error(); err != nil {
			cli.Fatal(cli.Failure, err.Error())
		}
	}
}
