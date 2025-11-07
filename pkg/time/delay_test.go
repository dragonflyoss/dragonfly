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
			expectedMax: 100 * time.Millisecond,
		},
		{
			name:        "attempt five sleeps 50ms",
			attempt:     5,
			increment:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 10 * time.Millisecond,
			expectedMax: 155 * time.Millisecond,
		},
		{
			name:        "capped at maxDelay",
			attempt:     15,
			increment:   10 * time.Millisecond,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 55 * time.Millisecond,
			expectedMax: 155 * time.Millisecond,
		},
		{
			name:        "zero increment sleeps zero",
			attempt:     10,
			increment:   0,
			maxDelay:    100 * time.Millisecond,
			expectedMin: 0,
			expectedMax: 100 * time.Millisecond,
		},
		{
			name:        "zero maxDelay caps at zero",
			attempt:     10,
			increment:   10 * time.Millisecond,
			maxDelay:    0,
			expectedMin: 0,
			expectedMax: 100 * time.Millisecond,
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
			expectedMin: 15 * time.Millisecond,
			expectedMax: 250 * time.Millisecond,
		},
		{
			name:        "attempt one with jitter",
			attempt:     1,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 45 * time.Millisecond,
			expectedMax: 350 * time.Millisecond,
		},
		{
			name:        "attempt two with jitter",
			attempt:     2,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 100 * time.Millisecond,
			expectedMax: 900 * time.Millisecond,
		},
		{
			name:        "attempt three with jitter",
			attempt:     3,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 280 * time.Millisecond,
			expectedMax: 1300 * time.Millisecond,
		},
		{
			name:        "capped at maxDelay with jitter",
			attempt:     10,
			baseDelay:   100 * time.Millisecond,
			maxDelay:    1 * time.Second,
			expectedMin: 380 * time.Millisecond,
			expectedMax: 1500 * time.Millisecond,
		},
		{
			name:        "large attempt capped at maxDelay",
			attempt:     20,
			baseDelay:   50 * time.Millisecond,
			maxDelay:    2 * time.Second,
			expectedMin: 650 * time.Millisecond,
			expectedMax: 3000 * time.Millisecond,
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
			expectedMax: 100 * time.Millisecond,
		},
		{
			name:        "small baseDelay with exponential growth",
			attempt:     4,
			baseDelay:   10 * time.Millisecond,
			maxDelay:    5 * time.Second,
			expectedMin: 20 * time.Millisecond,
			expectedMax: 280 * time.Millisecond,
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

func TestRandomDelayWithJitter(t *testing.T) {
	tests := []struct {
		name        string
		baseDelay   time.Duration
		expectedMin time.Duration
		expectedMax time.Duration
	}{
		{
			name:        "1 second base delay",
			baseDelay:   1 * time.Second,
			expectedMin: 500 * time.Millisecond,
			expectedMax: 1500 * time.Millisecond,
		},
		{
			name:        "2 seconds base delay",
			baseDelay:   2 * time.Second,
			expectedMin: 1000 * time.Millisecond,
			expectedMax: 3000 * time.Millisecond,
		},
		{
			name:        "zero base delay",
			baseDelay:   0,
			expectedMin: 0,
			expectedMax: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run multiple times to account for jitter randomness
			const iterations = 3
			successCount := 0

			for i := range iterations {
				start := time.Now()
				RandomDelayWithJitter(tt.baseDelay)
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
