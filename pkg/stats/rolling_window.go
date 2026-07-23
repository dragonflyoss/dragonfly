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
)

// RollingWindow is a thread-safe sliding window over the most recent samples.
// It maintains rolling sums so snapshots are O(1) and memory stays bounded by
// the capacity regardless of how many samples are added.
type RollingWindow struct {
	// mu protects the fields below.
	mu sync.Mutex

	// capacity is the maximum number of samples retained.
	capacity int

	// samples is a ring buffer of samples.
	samples []float64

	// head is the overwrite position once the ring buffer is full,
	// which is also the position of the oldest sample.
	head int

	// last is the most recently added sample.
	last float64

	// sum is the rolling sum of the samples in the window.
	sum float64

	// sumSquares is the rolling sum of the squared samples in the window.
	sumSquares float64

	// evictions is the number of ring evictions since the rolling sums were
	// last recomputed from the retained samples.
	evictions int
}

// NewRollingWindow returns a new RollingWindow retaining at most capacity samples.
func NewRollingWindow(capacity int) *RollingWindow {
	if capacity <= 0 {
		panic("stats: rolling window capacity must be positive")
	}

	return &RollingWindow{capacity: capacity}
}

// Add adds a sample to the window, evicting the oldest sample once full.
func (w *RollingWindow) Add(sample float64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.samples) < w.capacity {
		w.samples = append(w.samples, sample)
	} else {
		evicted := w.samples[w.head]
		w.sum -= evicted
		w.sumSquares -= evicted * evicted
		w.samples[w.head] = sample
		w.head = (w.head + 1) % w.capacity
		w.evictions++
	}

	w.last = sample
	w.sum += sample
	w.sumSquares += sample * sample

	// Incrementally adding and subtracting floating point sums drifts over
	// time. Recompute the sums from the retained samples once per window of
	// evictions, bounding the drift while keeping Add amortized O(1).
	if w.evictions >= w.capacity {
		w.evictions = 0
		w.sum = 0
		w.sumSquares = 0
		for _, sample := range w.samples {
			w.sum += sample
			w.sumSquares += sample * sample
		}
	}
}

// Values returns a copy of the samples in the window, ordered from oldest to newest.
func (w *RollingWindow) Values() []float64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	values := make([]float64, 0, len(w.samples))
	for i := range w.samples {
		values = append(values, w.samples[(w.head+i)%len(w.samples)])
	}

	return values
}

// Snapshot returns a point-in-time view of the statistics of the window in O(1).
func (w *RollingWindow) Snapshot() Snapshot {
	w.mu.Lock()
	defer w.mu.Unlock()

	snapshot := Snapshot{
		Count:      len(w.samples),
		Last:       w.last,
		sum:        w.sum,
		sumSquares: w.sumSquares,
	}
	if snapshot.Count == 0 {
		return snapshot
	}

	// Population variance via E[x^2] - E[x]^2, clamped to zero against
	// floating point error.
	n := float64(snapshot.Count)
	snapshot.Mean = w.sum / n
	if variance := w.sumSquares/n - snapshot.Mean*snapshot.Mean; variance > 0 {
		snapshot.StdDev = math.Sqrt(variance)
	}

	return snapshot
}

// Snapshot is a point-in-time view of the statistics of a RollingWindow.
type Snapshot struct {
	// Count is the number of samples in the window.
	Count int

	// Last is the most recently added sample.
	Last float64

	// Mean is the mean of all samples in the window.
	Mean float64

	// StdDev is the population standard deviation of all samples in the window.
	StdDev float64

	// sum is the rolling sum of the samples in the window.
	sum float64

	// sumSquares is the rolling sum of the squared samples in the window.
	sumSquares float64
}

// MeanExcludingLast returns the mean of the window excluding the latest
// sample, supporting comparison of the latest sample against the history
// before it. It returns 0 if the window holds fewer than two samples.
func (s Snapshot) MeanExcludingLast() float64 {
	if s.Count < 2 {
		return 0
	}

	return (s.sum - s.Last) / float64(s.Count-1)
}

// StdDevExcludingLast returns the population standard deviation of the window
// excluding the latest sample. It returns 0 if the window holds fewer than
// two samples.
func (s Snapshot) StdDevExcludingLast() float64 {
	if s.Count < 2 {
		return 0
	}

	// Population variance via E[x^2] - E[x]^2, clamped to zero against
	// floating point error.
	n := float64(s.Count - 1)
	mean := (s.sum - s.Last) / n
	if variance := (s.sumSquares-s.Last*s.Last)/n - mean*mean; variance > 0 {
		return math.Sqrt(variance)
	}

	return 0
}
