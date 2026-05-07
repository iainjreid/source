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
	"flag"
	"fmt"
	"os"

	"github.com/iainjreid/source/cmd/src/internal/cli"
	"github.com/iainjreid/source/cmd/src/internal/manage"
	"github.com/iainjreid/source/cmd/src/internal/start"
	"github.com/iainjreid/source/cmd/src/internal/verify"
	"github.com/iainjreid/source/internal/debug"
)

const usage = `Usage:
    src [-h | --help] [-v | --version] <command> [<args>]

Source is an experimental Git server backed not by the filesystem, but by a
database of your choosing. It has no dependency on Git and requires only a
connection string to your database at runtime.

Options:
    -h, --help                  Display this message
    -v, --version               Display the version number

Commands:
    start                       Start a server
    manage                      Perform setup and maintenance tasks
    verify                      Ensure the database is setup and ready to go
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Run 'src --help' for usage.\n")
	}

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	var (
		help    = cli.Bool(flag.CommandLine, false, "help", "h")
		version = cli.Bool(flag.CommandLine, false, "version", "v")
	)

	flag.Parse()

	if *help {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	if *version {
		fmt.Printf("src version %s\n", debug.Version)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "start":
		start.Cmd(os.Args[2:])

	case "manage":
		manage.Cmd(os.Args[2:])

	case "verify":
		verify.Cmd(os.Args[2:])

	default:
		fmt.Printf("'%s' is not a src command.\n", os.Args[1])
		flag.Usage()
		os.Exit(1)
	}
}
