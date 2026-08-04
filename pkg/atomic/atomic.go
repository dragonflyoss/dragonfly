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

// Package atomic supplements the standard library sync/atomic with atomic
// time.Time and time.Duration types, which sync/atomic does not provide.
package atomic

import (
	"sync/atomic"
	"time"
)

// Time is an atomic time.Time. The zero value loads the zero time.Time.
type Time struct {
	p atomic.Pointer[time.Time]
}

// NewTime returns a new Time set to the given value.
func NewTime(v time.Time) *Time {
	t := new(Time)
	t.Store(v)
	return t
}

// Load atomically loads the wrapped time.Time.
func (t *Time) Load() time.Time {
	if p := t.p.Load(); p != nil {
		return *p
	}

	return time.Time{}
}

// Store atomically stores the passed time.Time.
func (t *Time) Store(v time.Time) {
	t.p.Store(&v)
}

// Duration is an atomic time.Duration. The zero value loads a zero duration.
type Duration struct {
	d atomic.Int64
}

// NewDuration returns a new Duration set to the given value.
func NewDuration(v time.Duration) *Duration {
	d := new(Duration)
	d.Store(v)
	return d
}

// Load atomically loads the wrapped time.Duration.
func (d *Duration) Load() time.Duration {
	return time.Duration(d.d.Load())
}

// Store atomically stores the passed time.Duration.
func (d *Duration) Store(v time.Duration) {
	d.d.Store(int64(v))
}
