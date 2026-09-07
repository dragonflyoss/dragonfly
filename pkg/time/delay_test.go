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

	"github.com/stretchr/testify/assert"
)

const iterations = 3

func TestExponentialDelayWithJitter(t *testing.T) {
	tests := []struct {
		name      string
		attempt   uint
		baseDelay time.Duration
		maxDelay  time.Duration
		expect    func(t *testing.T, durations []time.Duration)
	}{
		{
			name:      "attempt zero with jitter",
			attempt:   0,
			baseDelay: 100 * time.Millisecond,
			maxDelay:  5 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 15*time.Millisecond && duration <= 250*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "attempt one with jitter",
			attempt:   1,
			baseDelay: 100 * time.Millisecond,
			maxDelay:  5 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 45*time.Millisecond && duration <= 350*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "attempt two with jitter",
			attempt:   2,
			baseDelay: 100 * time.Millisecond,
			maxDelay:  5 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 100*time.Millisecond && duration <= 900*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "attempt three with jitter",
			attempt:   3,
			baseDelay: 100 * time.Millisecond,
			maxDelay:  5 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 280*time.Millisecond && duration <= 1300*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "capped at maxDelay with jitter",
			attempt:   10,
			baseDelay: 100 * time.Millisecond,
			maxDelay:  1 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 380*time.Millisecond && duration <= 1500*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "large attempt capped at maxDelay",
			attempt:   20,
			baseDelay: 50 * time.Millisecond,
			maxDelay:  2 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 650*time.Millisecond && duration <= 3000*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "overflow attempt stays capped at maxDelay",
			attempt:   39,
			baseDelay: 30 * time.Millisecond,
			maxDelay:  1 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 380*time.Millisecond && duration <= 1500*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "shift overflow attempt stays capped at maxDelay",
			attempt:   62,
			baseDelay: 30 * time.Millisecond,
			maxDelay:  1 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 380*time.Millisecond && duration <= 1500*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "huge attempt stays capped at maxDelay",
			attempt:   100,
			baseDelay: 30 * time.Millisecond,
			maxDelay:  1 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 380*time.Millisecond && duration <= 1500*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "zero baseDelay returns without sleeping",
			attempt:   5,
			baseDelay: 0,
			maxDelay:  1 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 0 && duration <= 100*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "zero maxDelay caps at zero",
			attempt:   5,
			baseDelay: 100 * time.Millisecond,
			maxDelay:  0,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 0 && duration <= 100*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "small baseDelay with exponential growth",
			attempt:   4,
			baseDelay: 10 * time.Millisecond,
			maxDelay:  5 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 20*time.Millisecond && duration <= 280*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			durations := make([]time.Duration, 0, iterations)
			for range iterations {
				start := time.Now()
				ExponentialDelayWithJitter(context.TODO(), tc.attempt, tc.baseDelay, tc.maxDelay)
				durations = append(durations, time.Since(start))
			}

			tc.expect(t, durations)
		})
	}
}

func TestRandomDelayWithJitter(t *testing.T) {
	tests := []struct {
		name      string
		baseDelay time.Duration
		expect    func(t *testing.T, durations []time.Duration)
	}{
		{
			name:      "1 second base delay",
			baseDelay: 1 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 500*time.Millisecond && duration <= 1500*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "2 seconds base delay",
			baseDelay: 2 * time.Second,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 1000*time.Millisecond && duration <= 3000*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "zero base delay returns without sleeping",
			baseDelay: 0,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 0 && duration <= 100*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
		{
			name:      "one nanosecond base delay has no jitter range",
			baseDelay: time.Nanosecond,
			expect: func(t *testing.T, durations []time.Duration) {
				assert := assert.New(t)
				inRange := 0
				for _, duration := range durations {
					if duration >= 0 && duration <= 100*time.Millisecond {
						inRange++
					}
				}

				assert.GreaterOrEqual(inRange, 2)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			durations := make([]time.Duration, 0, iterations)
			for range iterations {
				start := time.Now()
				RandomDelayWithJitter(context.Background(), tc.baseDelay)
				durations = append(durations, time.Since(start))
			}

			tc.expect(t, durations)
		})
	}
}

func TestDelayWithJitter_ContextCanceled(t *testing.T) {
	tests := []struct {
		name   string
		delay  func(ctx context.Context)
		expect func(t *testing.T, duration time.Duration)
	}{
		{
			name: "exponential delay",
			delay: func(ctx context.Context) {
				ExponentialDelayWithJitter(ctx, 0, time.Minute, time.Minute)
			},
			expect: func(t *testing.T, duration time.Duration) {
				assert := assert.New(t)
				assert.LessOrEqual(duration, time.Second)
			},
		},
		{
			name: "random delay",
			delay: func(ctx context.Context) {
				RandomDelayWithJitter(ctx, time.Minute)
			},
			expect: func(t *testing.T, duration time.Duration) {
				assert := assert.New(t)
				assert.LessOrEqual(duration, time.Second)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			start := time.Now()
			tc.delay(ctx)
			tc.expect(t, time.Since(start))
		})
	}
}
