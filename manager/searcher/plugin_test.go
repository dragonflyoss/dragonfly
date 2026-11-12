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

package searcher

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadPlugin(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		os.Remove("./testdata/d7y-manager-plugin-searcher.so")
		os.Remove("./testdata/test")
	}()

	output, err := exec.Command("go", "build", "-buildmode=plugin", "-o=./testdata/d7y-manager-plugin-searcher.so", "testdata/plugin/searcher.go").CombinedOutput()
	if err != nil {
		t.Fatal(string(output))
	}

	output, err = exec.Command("go", "build", "-o=./testdata/test", "testdata/main.go").CombinedOutput()
	if err != nil {
		t.Fatal(string(output))
	}

	wd, err := os.Getwd()
	assert.NoError(err)

	output, err = exec.Command("./testdata/test", "-plugin-dir", path.Join(wd, "testdata")).CombinedOutput()
	if err != nil {
		t.Fatal(string(output))
	}
}
