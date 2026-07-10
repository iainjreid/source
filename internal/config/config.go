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

// Package config is responsible for consuming environment variables to control
// the behaviour of the program, and provide sensible defaults for those that
// are left undefined.
//
// Internally, configuration values closely resemble flag values. This is
// intentional. It allows the package to integrate naturally with the standard
// library’s flag package where required, while maintaining a consistent
// approach to configuration when command-line flags are not used.
package config

import (
	"errors"
	"log/slog"
	"os"

	"github.com/iainjreid/source/internal/ssh"
)

var (
	OptDatabaseURI, DatabaseURI = NewOption("SOURCE_DB", StringValue(""))

	OptHttpPort, HttpPort = NewOption("SOURCE_HTTP_PORT", IntValue(8080))

	OptSshPort, SshPort = NewOption("SOURCE_SSH_PORT", IntValue(2222))

	// This option is mutually exclusive with the SSH ID string.
	OptSshIdPath, SshIdPath = NewOption("SOURCE_SSH_ID_PATH", StringValue(""))

	// This option is mutually exclusive with the SSH ID path.
	OptSshId, SshId = NewOption("SOURCE_SSH_ID", StringValue(""))

	OptLogLevel, LogLevel = NewOption("SOURCE_LOG_LEVEL", SlogLevelValue(slog.LevelInfo))
)

// Validate performs cross-option validation and initialisation of some
// variables under certain conditions.
func Validate() error {
	if string(*DatabaseURI) == "" {
		return errors.New("SOURCE_DB is required")
	}

	if string(*SshIdPath) != "" && string(*SshId) != "" {
		return errors.New("SOURCE_SSH_ID_PATH and SOURCE_SSH_ID are mutually exclusive")
	}

	if string(*SshIdPath) != "" {
		data, err := os.ReadFile(string(*SshIdPath))
		if err != nil {
			return errors.New("unable to read SSH private key")
		}

		SshIdPath.Set(string(data))
	}

	if string(*SshId) != "" {
		if err := ssh.Init(string(*SshId)); err != nil {
			return errors.New("unable to parse PEM encoded SSH private key")
		}
	}

	return nil
}
