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
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		expect  func(t *testing.T, d Dfpath, err error)
	}{
		{
			name:    "new dfpath failed",
			options: []Option{WithLogDir("")},
			expect: func(t *testing.T, d Dfpath, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "new dfpath",
			expect: func(t *testing.T, d Dfpath, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(DefaultLogDir, d.LogDir())
				assert.Equal(DefaultPluginDir, d.PluginDir())
			},
		},
		{
			name:    "new dfpath by logDir",
			options: []Option{WithLogDir("foo")},
			expect: func(t *testing.T, d Dfpath, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", d.LogDir())
				assert.Equal(DefaultPluginDir, d.PluginDir())
			},
		},
		{
			name:    "new dfpath by pluginDir",
			options: []Option{WithPluginDir("foo")},
			expect: func(t *testing.T, d Dfpath, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(DefaultLogDir, d.LogDir())
				assert.Equal("foo", d.PluginDir())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache.Once = sync.Once{}
			cache.err = nil
			d, err := New(tc.options...)
			tc.expect(t, d, err)
		})
	}
}
