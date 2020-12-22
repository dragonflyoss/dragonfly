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

package basic

import (
	"github.com/dragonflyoss/Dragonfly2/pkg/util/assert"
	"os"
	"os/user"
	"strings"
)

var HomeDir string
var TmpDir string

func init() {
	u, err := user.Current()

	if err != nil {
		panic(err.Error())
	}

	HomeDir = u.HomeDir
	if len(HomeDir) > 1 {
		HomeDir = strings.TrimRight(HomeDir, "/")
		if HomeDir == "" {
			HomeDir = "/"
		}
	}
	assert.PAssert(len(HomeDir) > 0, "home dir is empty")

	TmpDir = os.TempDir()
	if TmpDir == "" {
		TmpDir = "/tmp"
	}

	if len(TmpDir) > 1 {
		TmpDir = strings.TrimRight(TmpDir, "/")
	}
}
