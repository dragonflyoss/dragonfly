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

package version

import (
	"fmt"
)

var (
	Major      = "2"
	Minor      = "3"
	GitVersion = "v2.3.0"
	GitCommit  = "unknown"
	Platform   = osArch
	BuildTime  = "unknown"
	GoVersion  = "unknown"
	Gotags     = "unknown"
	Gogcflags  = "unknown"
)

func Version() string {
	return fmt.Sprintf("Major: %s, Minor: %s, GitVersion: %s, GitCommit: %s, Platform: %s, BuildTime: %s, GoVersion: %s, Gotags: %s, Gogcflags: %s", Major,
		Minor, GitVersion, GitCommit, Platform, BuildTime, GoVersion, Gotags, Gogcflags)
}
