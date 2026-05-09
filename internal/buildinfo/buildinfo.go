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

// Package buildinfo exports build information used in the creation of the
// version string, and might be used to provide debugging information later
// down the line to help triage bugs.
//
// In such instances where the build information cannot be generated, partially
// or otherwise, it's likely that the binary was not obtained through a supported
// channel or was built without the required version control information.
package buildinfo

import (
	"fmt"
	"runtime/debug"
)

var (
	Version = "(unknown)"
)

func init() {
	info, ok := debug.ReadBuildInfo()

	// If the build information cannot be read, do nothing. See the package
	// documentation at the top of this file for a full explanation.
	if !ok {
		return
	}

	if info.Main.Version != "" {
		Version = info.Main.Version
	}

	// This snippet is here for when build information needs to be exposed as
	// part of the upcoming debugging work.
	//
	// for _, setting := range info.Settings {
	// 	switch setting.Key {
	// 	case "vcs.revision":
	// 		GitCommit = setting.Value

	// 	...
	// 	}
	// }
}

func String() string {
	return fmt.Sprintf("src version %s", Version)
}
