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

package cli

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/iainjreid/source/internal/config"
)

var (
	LoggerFlagsError = errors.New("-q/--quiet can't be used with -v/--verbose")
)

// ResolveLoggerFlags takes flags associated with configuring the logger and
// determines the correct [slog.Level] for the selection.
func ResolveLoggerFlags(quiet, verbose *bool) (slog.Level, error) {
	var level slog.Level
	var err error

	switch {
	case *quiet && *verbose:
		err = LoggerFlagsError

	case *quiet:
		level = slog.LevelError

	case *verbose:
		level = slog.LevelDebug
	}

	return level, err
}

type Cmd struct {
	*flag.FlagSet
}

func New(name string) Cmd {
	flag := flag.NewFlagSet(name, flag.ExitOnError)

	flag.Usage = func() {
		Fatalf(BadUsage, "Run '%s --help' for usage.\n", name)
	}

	return Cmd{flag}
}

func (c *Cmd) ExplainUsage(msg string) {
	Error(msg)
	c.Usage()
}

func (c *Cmd) String(value string, long, short string) *string {
	out := new(string)
	c.StringVar(out, long, value, "")
	if short != "" {
		c.StringVar(out, short, value, "")
	}
	return out
}

func (c *Cmd) Int(value int, long, short string) *int {
	out := new(int)
	c.IntVar(out, long, value, "")
	if short != "" {
		c.IntVar(out, short, value, "")
	}
	return out
}

func (c *Cmd) Bool(value bool, long, short string) *bool {
	out := new(bool)
	c.BoolVar(out, long, value, "")
	if short != "" {
		c.BoolVar(out, short, value, "")
	}
	return out
}

func (c *Cmd) StringSlice(long, short string) *StringSlice {
	out := new(StringSlice)
	c.Var(out, long, "")
	if short != "" {
		c.Var(out, short, "")
	}
	return out
}

func (c *Cmd) Flag(opt *config.Option, long, short string) {
	c.Var(opt, long, "")
	if short != "" {
		c.Var(opt, short, "")
	}
}

func (c *Cmd) VisitUnvisited(fn func(*flag.Flag) error) error {
	visited := make(map[string]struct{})

	c.Visit(func(f *flag.Flag) {
		visited[f.Name] = struct{}{}
	})

	var err error

	c.VisitAll(func(f *flag.Flag) {
		if err != nil {
			return
		}

		if _, ok := visited[f.Name]; ok {
			return
		}

		err = fn(f)
	})

	return err
}

func (c *Cmd) ResolveConfig(args []string) error {
	if err := c.Parse(args); err != nil {
		return err
	}

	return c.VisitUnvisited(func(f *flag.Flag) error {
		if v, ok := f.Value.(*config.Option); ok {
			return v.LoadFromEnv()
		}

		return nil
	})
}

var (
	Success  = 0
	Failure  = 1
	BadUsage = 64
)

func Error(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

func Fatal(code int, msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(code)
}

func Fatalf(code int, format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(code)
}
