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

	"github.com/iainjreid/source/cmd/src/internal/cli"
	"github.com/iainjreid/source/internal/logger"
)

var usage = `Usage:
    src manage [--db <uri>] [-q | --quiet] [-v | --verbose] [-j | --json]
	           [-h | --help] <command> [<args>]

Options:
    --db <uri>                  Specify the database URI
    -q, --quiet                 Suppress all non-error logs
    -v, --verbose               Print all logs
    -j, --json                  Display JSON output
    -h, --help                  Display this message
    --debug                     Enable debugging

Additional Options:
    --setup                     Prepare the database 

Example:
    $ src manage --db postgresql://postgres@localhost --setup
`

func Cmd(_ context.Context, args []string) {
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
