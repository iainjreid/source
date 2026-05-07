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
	"flag"
	"fmt"
	"os"

	"github.com/iainjreid/source/cmd/src/internal/cli"
)

var usage = `Usage:
    src verify [--db <uri>] [-q | --quiet] [-v | --verbose] [-j | --json]
	           [-h | --help]

Before running a server, you can verify that the target database has been
initalised and is ready to accept connections from src without running into
any unexpected errors.

It's a good idea to run this command after first setting up your server, but
it can also be used as a readiness check before starting the server itself.

Options:
    --db <uri>                  Specify the database URI
    -q, --quiet                 Suppress all non-error logs
    -v, --verbose				Print all logs
    -j, --json                  Display JSON output
    -h, --help                  Display this message

Example:
    $ src verify --db postgresql://postgres@localhost
`

func Cmd(args []string) {
	flag := flag.NewFlagSet("start", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Run 'src verify --help' for usage.\n")
	}

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		_       = cli.String(flag, "", "db", "")
		quiet   = cli.Bool(flag, false, "quiet", "q")
		verbose = cli.Bool(flag, false, "verbose", "v")
		_       = cli.Bool(flag, false, "json", "j")
		help    = cli.Bool(flag, false, "help", "h")
	)

	flag.Parse(args)

	if *help {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if *quiet && *verbose {
		fmt.Fprint(os.Stderr, "-q/--quiet can't be used with -v/--verbose")
		os.Exit(1)
	}

	fmt.Fprint(os.Stderr, "Not yet implemented...\n")
	os.Exit(1)
}
