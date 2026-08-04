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
		expect  uint64
	}{
		{
			name:    "delta less than counter",
			initial: 100,
			delta:   40,
			expect:  60,
		},
		{
			name:    "delta equals counter",
			initial: 100,
			delta:   100,
			expect:  0,
		},
		{
			name:    "delta greater than counter clamps to zero",
			initial: 40,
			delta:   100,
			expect:  0,
		},
		{
			name:    "zero delta",
			initial: 100,
			delta:   0,
			expect:  100,
		},
		{
			name:    "zero counter",
			initial: 0,
			delta:   100,
			expect:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := new(atomic.Uint64)
			counter.Store(tt.initial)

			SafeSubAtomicUint64(counter, tt.delta)
			assert.New(t).Equal(tt.expect, counter.Load())
		})
	}
}

func TestSafeSubAtomicUint64_Concurrent(t *testing.T) {
	const (
		goroutines = 100
		delta      = 3
	)

	counter := new(atomic.Uint64)
	counter.Store(goroutines * delta)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			SafeSubAtomicUint64(counter, delta)
		}()
	}
	wg.Wait()

	assert.New(t).Equal(uint64(0), counter.Load())
}

func TestSafeSubAtomicUint64_ConcurrentClampsToZero(t *testing.T) {
	const goroutines = 100

	counter := new(atomic.Uint64)
	counter.Store(50)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for range goroutines {
		go func() {
			defer wg.Done()
			SafeSubAtomicUint64(counter, 3)
		}()
	}
	wg.Wait()

	assert.New(t).Equal(uint64(0), counter.Load())
}
