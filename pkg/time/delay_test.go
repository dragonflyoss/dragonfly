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

func TestExponentialDelayWithJitter(t *testing.T) {
	tests := []struct {
		name        string
		attempt     uint
		baseDelay   time.Duration
		maxDelay    time.Duration
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{
			name:        "attempt zero with jitter",
			attempt:     0,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 45 * time.Millisecond,
			expectedMax: 200 * time.Millisecond,
		},
		{
			name:        "attempt one with jitter",
			attempt:     1,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 95 * time.Millisecond,
			expectedMax: 300 * time.Millisecond,
		},
		{
			name:        "attempt two with jitter",
			attempt:     2,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 190 * time.Millisecond,
			expectedMax: 800 * time.Millisecond,
		},
		{
			name:        "attempt three with jitter",
			attempt:     3,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 380 * time.Millisecond,
			expectedMax: 1200 * time.Millisecond,
		},
		{
			name:        "capped at maxDelay with jitter",
			attempt:     10,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    1 * time.Second,
			expectedMin: 480 * time.Millisecond,
			expectedMax: 1400 * time.Millisecond,
		},
		{
			name:        "large attempt capped at maxDelay",
			attempt:     20,
			baseDelay:   50 * time.Millisecond,
			maxDelay:    2 * time.Second,
			expectedMin: 950 * time.Millisecond,
			expectedMax: 2800 * time.Millisecond,
		},
		{
			name:        "zero baseDelay with jitter",
			attempt:     5,
			baseDelay:   0,
			maxDelay:    1 * time.Second,
			expectedMin: 0,
			expectedMax: 100 * time.Millisecond,
		},
		{
			name:        "zero maxDelay caps at zero",
			attempt:     5,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    0,
			expectedMin: 0,
			expectedMax: 50 * time.Millisecond,
		},
		{
			name:        "small baseDelay with exponential growth",
			attempt:     4,
			baseDelay:   10 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 70 * time.Millisecond,
			expectedMax: 200 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to account for jitter randomness
			const iterations = 3
			successCount := 0

			for i := range iterations {
				start := time.Now()
				if err := ExponentialDelayWithJitter(context.TODO(), tt.attempt, tt.baseDelay, tt.maxDelay); err != nil {
					t.Fatalf("ExponentialDelayWithJitter returned error: %v", err)
				}

				duration := time.Since(start)

				// Allow some failures due to system scheduling
				if duration >= tt.expectedMin && duration <= tt.expectedMax {
					successCount++
				} else {
					t.Logf("Iteration %d out of range: got %v, want [%v, %v]", i+1, duration, tt.expectedMin, tt.expectedMax)
				}
			}

			// At least 2 out of 3 iterations should be in range
			if successCount < 2 {
				t.Errorf("Too many iterations out of range: %d/%d successful", successCount, iterations)
			}
		})
	}
}
