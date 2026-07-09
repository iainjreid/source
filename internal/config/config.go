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
