/*
 *     Copyright 2026 The Dragonfly Authors
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

package math

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeSubAtomicUint64(t *testing.T) {
	tests := []struct {
		name    string
		initial uint64
		delta   uint64
		expect  func(t *testing.T, counter *atomic.Uint64)
	}{
		{
			name:    "delta less than counter",
			initial: 100,
			delta:   40,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(60), counter.Load())
			},
		},
		{
			name:    "delta equals counter",
			initial: 100,
			delta:   100,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(0), counter.Load())
			},
		},
		{
			name:    "delta greater than counter clamps to zero",
			initial: 40,
			delta:   100,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(0), counter.Load())
			},
		},
		{
			name:    "zero delta",
			initial: 100,
			delta:   0,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(100), counter.Load())
			},
		},
		{
			name:    "zero counter",
			initial: 0,
			delta:   100,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(0), counter.Load())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			counter := new(atomic.Uint64)
			counter.Store(tc.initial)

			SafeSubAtomicUint64(counter, tc.delta)
			tc.expect(t, counter)
		})
	}
}

func TestSafeSubAtomicUint64_Concurrent(t *testing.T) {
	const goroutines = 100

	tests := []struct {
		name    string
		initial uint64
		delta   uint64
		expect  func(t *testing.T, counter *atomic.Uint64)
	}{
		{
			name:    "deltas sum to counter",
			initial: goroutines * 3,
			delta:   3,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(0), counter.Load())
			},
		},
		{
			name:    "deltas exceed counter clamps to zero",
			initial: 50,
			delta:   3,
			expect: func(t *testing.T, counter *atomic.Uint64) {
				assert := assert.New(t)
				assert.Equal(uint64(0), counter.Load())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			counter := new(atomic.Uint64)
			counter.Store(tc.initial)

			var wg sync.WaitGroup
			wg.Add(goroutines)
			for range goroutines {
				go func() {
					defer wg.Done()
					SafeSubAtomicUint64(counter, tc.delta)
				}()
			}

			wg.Wait()

			tc.expect(t, counter)
		})
	}
}
