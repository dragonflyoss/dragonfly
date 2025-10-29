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

package time

import (
	"context"
	"testing"
	"time"
)

func TestLinearDelay(t *testing.T) {
	tests := []struct {
		name        string
		attempt     uint
		increment   time.Duration
		maxDelay    time.Duration
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{
			name:        "attempt zero sleeps zero",
			attempt:     0,
			increment:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 0,
			expectedMax: 1 * time.Millisecond,
		},
		{
			name:        "attempt one sleeps increment",
			attempt:     1,
			increment:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 8 * time.Millisecond,
			expectedMax: 12 * time.Millisecond,
		},
		{
			name:        "attempt five sleeps 50ms",
			attempt:     5,
			increment:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 45 * time.Millisecond,
			expectedMax: 55 * time.Millisecond,
		},
		{
			name:        "capped at maxDelay",
			attempt:     15,
			increment:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 95 * time.Millisecond,
			expectedMax: 105 * time.Millisecond,
		},
		{
			name:        "zero increment sleeps zero",
			attempt:     10,
			increment:   0,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 0,
			expectedMax: 1 * time.Millisecond,
		},
		{
			name:        "zero maxDelay caps at zero",
			attempt:     10,
			increment:   10 * time.Millisecond,
			maxDelay:    0,
			expectedMin: 0,
			expectedMax: 1 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			if err := LinearDelay(context.TODO(), tt.attempt, tt.increment, tt.maxDelay); err != nil {
				t.Fatalf("LinearDelay returned error: %v", err)
			}

			duration := time.Since(start)
			if duration < tt.expectedMin {
				t.Errorf("LinearDelay slept too short: got %v, want at least %v", duration, tt.expectedMin)
			}

			if duration > tt.expectedMax {
				t.Errorf("LinearDelay slept too long: got %v, want at most %v", duration, tt.expectedMax)
			}
		})
	}
}
