/*
 *     Copyright 2025 The Dragonfly Authors
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

package redis

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxy_StopConcurrent(t *testing.T) {
	tests := []struct {
		name    string
		callers int
	}{
		{
			name:    "stop is idempotent when called once",
			callers: 1,
		},
		{
			name:    "stop is idempotent under concurrent callers",
			callers: 16,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewProxy("127.0.0.1:0", "127.0.0.1:0")

			var wg sync.WaitGroup
			wg.Add(tc.callers)
			for i := 0; i < tc.callers; i++ {
				go func() {
					defer wg.Done()
					p.Stop()
				}()
			}

			assert.NotPanics(t, wg.Wait)
		})
	}
}
