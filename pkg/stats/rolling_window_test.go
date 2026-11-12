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

package stats

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRollingWindow_NewRollingWindow(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() { NewRollingWindow(0) })
	assert.Panics(func() { NewRollingWindow(-1) })
	assert.NotNil(NewRollingWindow(1))
}

func TestRollingWindow_Snapshot(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, window *RollingWindow)
	}{
		{
			name: "window is empty",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				snapshot := window.Snapshot()
				assert.Equal(0, snapshot.Count)
				assert.Equal(8, snapshot.Capacity)
				assert.False(snapshot.IsFull())
				assert.Equal(float64(0), snapshot.Last)
				assert.Equal(float64(0), snapshot.Mean)
				assert.Equal(float64(0), snapshot.StdDev)
				assert.Equal(float64(0), snapshot.MeanExcludingLast())
				assert.Equal(float64(0), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "window has a single sample",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				window.Add(100)
				snapshot := window.Snapshot()
				assert.Equal(1, snapshot.Count)
				assert.Equal(float64(100), snapshot.Last)
				assert.Equal(float64(100), snapshot.Mean)
				assert.Equal(float64(0), snapshot.StdDev)
				assert.Equal(float64(0), snapshot.MeanExcludingLast())
				assert.Equal(float64(0), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "mean and standard deviation over the window and excluding the latest sample",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				window.Add(100)
				window.Add(200)
				window.Add(600)
				snapshot := window.Snapshot()
				assert.Equal(3, snapshot.Count)
				assert.Equal(float64(600), snapshot.Last)
				assert.Equal(float64(300), snapshot.Mean)
				assert.InDelta(math.Sqrt(140000.0/3), snapshot.StdDev, 1e-9)
				assert.Equal(float64(150), snapshot.MeanExcludingLast())
				assert.Equal(float64(50), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "window is bounded and evicts the oldest samples",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				for i := range 12 {
					window.Add(float64(i))
				}

				values := window.Values()
				assert.Equal([]float64{4, 5, 6, 7, 8, 9, 10, 11}, values)

				snapshot := window.Snapshot()
				assert.Equal(8, snapshot.Count)
				assert.Equal(float64(11), snapshot.Last)
				assert.Equal(float64(7), snapshot.MeanExcludingLast())
				assert.Equal(float64(2), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "rolling sums are periodically recomputed under sustained evictions",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				for i := range 5 * 8 {
					window.Add(1e9 + float64(i%5)*1e6)
				}

				values := window.Values()
				assert.Equal(8, len(values))

				history := values[:len(values)-1]
				var sum float64
				for _, value := range history {
					sum += value
				}

				mean := sum / float64(len(history))

				var squaredDeviations float64
				for _, value := range history {
					deviation := value - mean
					squaredDeviations += deviation * deviation
				}

				stdDev := math.Sqrt(squaredDeviations / float64(len(history)))

				snapshot := window.Snapshot()
				assert.Equal(8, snapshot.Count)
				assert.Equal(values[len(values)-1], snapshot.Last)
				assert.InEpsilon(mean, snapshot.MeanExcludingLast(), 1e-9)
				assert.InEpsilon(stdDev, snapshot.StdDevExcludingLast(), 1e-6)
			},
		},
		{
			name: "window has two samples",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				window.Add(100)
				window.Add(300)
				snapshot := window.Snapshot()
				assert.Equal(2, snapshot.Count)
				assert.Equal(float64(300), snapshot.Last)
				assert.Equal(float64(200), snapshot.Mean)
				assert.Equal(float64(100), snapshot.StdDev)
				assert.Equal(float64(100), snapshot.MeanExcludingLast())
				assert.Equal(float64(0), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "values of an empty window",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				assert.Empty(window.Values())
			},
		},
		{
			name: "values keep insertion order before the window fills and are a copy",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				window.Add(3)
				window.Add(1)
				window.Add(2)

				values := window.Values()
				assert.Equal([]float64{3, 1, 2}, values)
				assert.False(window.Snapshot().IsFull())

				values[0] = 999
				assert.Equal([]float64{3, 1, 2}, window.Values())
			},
		},
		{
			name: "window filled to exactly its capacity without eviction",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				for i := range 8 {
					window.Add(float64(i))
				}

				assert.Equal([]float64{0, 1, 2, 3, 4, 5, 6, 7}, window.Values())

				snapshot := window.Snapshot()
				assert.Equal(8, snapshot.Count)
				assert.True(snapshot.IsFull())
				assert.Equal(float64(7), snapshot.Last)
				assert.Equal(float64(3), snapshot.MeanExcludingLast())
				assert.Equal(float64(2), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "identical samples have zero standard deviation",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				for range 12 {
					window.Add(5)
				}

				snapshot := window.Snapshot()
				assert.Equal(8, snapshot.Count)
				assert.Equal(float64(5), snapshot.Last)
				assert.Equal(float64(5), snapshot.Mean)
				assert.Equal(float64(0), snapshot.StdDev)
				assert.Equal(float64(5), snapshot.MeanExcludingLast())
				assert.Equal(float64(0), snapshot.StdDevExcludingLast())
			},
		},
		{
			name: "near-constant large samples keep standard deviation negligible",
			expect: func(t *testing.T, window *RollingWindow) {
				assert := assert.New(t)
				for range 12 {
					window.Add(1e9 + 0.1)
				}

				snapshot := window.Snapshot()
				assert.InEpsilon(1e9+0.1, snapshot.Mean, 1e-9)
				assert.InDelta(0, snapshot.StdDev, 1e-6*snapshot.Mean)
				assert.InDelta(0, snapshot.StdDevExcludingLast(), 1e-6*snapshot.Mean)
			},
		},
		{
			name: "window with capacity one keeps only the latest sample",
			expect: func(t *testing.T, _ *RollingWindow) {
				assert := assert.New(t)
				window := NewRollingWindow(1)
				for i := range 10 {
					window.Add(float64(i))
				}

				assert.Equal([]float64{9}, window.Values())
				snapshot := window.Snapshot()
				assert.Equal(1, snapshot.Count)
				assert.Equal(1, snapshot.Capacity)
				assert.True(snapshot.IsFull())
				assert.Equal(float64(9), snapshot.Last)
				assert.Equal(float64(9), snapshot.Mean)
				assert.Equal(float64(0), snapshot.StdDev)
				assert.Equal(float64(0), snapshot.MeanExcludingLast())
				assert.Equal(float64(0), snapshot.StdDevExcludingLast())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, NewRollingWindow(8))
		})
	}
}

func TestRollingWindow_Concurrent(t *testing.T) {
	assert := assert.New(t)
	window := NewRollingWindow(64)

	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			for j := range 1000 {
				window.Add(float64(i*1000 + j))
			}
		})
	}

	for range 4 {
		wg.Go(func() {
			for range 1000 {
				window.Snapshot()
				window.Values()
			}
		})
	}

	wg.Wait()

	values := window.Values()
	assert.Equal(64, len(values))
	var sum float64
	for _, value := range values {
		sum += value
	}

	snapshot := window.Snapshot()
	assert.Equal(64, snapshot.Count)
	assert.InEpsilon(sum/64, snapshot.Mean, 1e-9)
}
