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

package evaluator

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlugin_Load(t *testing.T) {
	assert := assert.New(t)
	defer func() {
		os.Remove("./testdata/d7y-scheduler-plugin-evaluator.so")
		os.Remove("./testdata/test")
	}()

	var (
		cmd    *exec.Cmd
		output []byte
		wd     string
		err    error
	)

	cmd = exec.Command("go", "build", "-buildmode=plugin", "-o=./testdata/d7y-scheduler-plugin-evaluator.so", "testdata/plugin/evaluator.go")
	output, err = cmd.CombinedOutput()
	assert.NoError(err)
	if err != nil {
		t.Fatal(string(output))
		return
	}

	cmd = exec.Command("go", "build", "-o=./testdata/test", "testdata/main.go")
	output, err = cmd.CombinedOutput()
	assert.NoError(err)
	if err != nil {
		t.Fatal(string(output))
		return
	}

	wd, err = os.Getwd()
	assert.NoError(err)
	wd = path.Join(wd, "testdata")

	cmd = exec.Command("./testdata/test", "-plugin-dir", wd)
	output, err = cmd.CombinedOutput()
	assert.NoError(err)
	if err != nil {
		t.Fatal(string(output))
		return
	}
}

func TestPlugin_LoadFailed(t *testing.T) {
	tests := []struct {
		name   string
		dir    string
		expect func(t *testing.T, e Evaluator, err error)
	}{
		{
			name: "plugin directory does not exist",
			dir:  "./testdata/not-exist",
			expect: func(t *testing.T, e Evaluator, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(e)
			},
		},
		{
			name: "plugin directory without plugin file",
			dir:  "./testdata/plugin",
			expect: func(t *testing.T, e Evaluator, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(e)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, err := LoadPlugin(tc.dir)
			tc.expect(t, e, err)
		})
	}
}
