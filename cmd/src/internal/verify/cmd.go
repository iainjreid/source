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

package verify

import (
	"context"

	"github.com/iainjreid/source/cmd/src/internal/cli"
	"github.com/iainjreid/source/internal/logger"
)

var usage = `Usage:
    src verify [--db <uri>] [-q | --quiet] [-v | --verbose] [-j | --json]
               [-h | --help] [--debug]

Before running a server, you can verify that the target database has been
initalised and is ready to accept connections from src without running into
any unexpected errors.

It's a good idea to run this command after first setting up your server, but
it can also be used as a readiness check before starting the server itself.

Options:
    --db <uri>                  Specify the database URI
    -q, --quiet                 Suppress all non-error logs
    -v, --verbose               Print all logs
    -j, --json                  Display JSON output
    -h, --help                  Display this message

Additional Options:
    --debug                     Enable debugging

Example:
    $ src verify --db postgresql://postgres@localhost
`

func Cmd(_ context.Context, args []string) {
	cmd := cli.New("src verify")

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

	cli.Fatal(1, "Not yet implemented...")
}
