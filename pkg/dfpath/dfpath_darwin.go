//go:build darwin

/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dfpath

import (
	"os"
	"path/filepath"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

var (
	DefaultWorkHome  = filepath.Join(homeDir(), ".dragonfly")
	DefaultConfigDir = filepath.Join(DefaultWorkHome, "config")
	DefaultLogDir    = filepath.Join(DefaultWorkHome, "logs")
	DefaultPluginDir = filepath.Join(DefaultWorkHome, "plugins")
)

// homeDir returns the current user's home directory.
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Warnf("Failed to get user home directory: %s. Use / as HomeDir", err.Error())
		return "/"
	}

	return home
}
