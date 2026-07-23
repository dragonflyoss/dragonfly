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
	"math/bits"
	"math/rand"
	"time"
)

// LinearDelay implements a linear backoff strategy for retries. It calculates delay based on
// the attempt number and sleeps for that duration, capped at maxDelay.
func LinearDelay(ctx context.Context, attempt uint, increment, maxDelay time.Duration) error {
	delay := time.Duration(attempt) * increment
	delay = min(delay, maxDelay)

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ExponentialDelayWithJitter is an exponential backoff strategy with jitter for retries. It calculates delay based on the attempt number,
// adds jitter, and sleeps for that duration, capped at maxDelay.
func ExponentialDelayWithJitter(ctx context.Context, attempt uint, baseDelay, maxDelay time.Duration) error {
	// Compute baseDelay * 2^attempt, capped at maxDelay. The shift is only applied when it
	// cannot overflow time.Duration; otherwise the result would exceed maxDelay anyway, so
	// the cap is used directly. Without this guard a large attempt overflows int64 before the
	// cap is applied, yielding a negative or zero delay that silently skips the backoff.
	var delay time.Duration
	if baseDelay > 0 {
		delay = maxDelay
		if attempt < uint(bits.LeadingZeros64(uint64(baseDelay))) {
			delay = min(baseDelay<<attempt, maxDelay)
		}
	}

	if delay > 0 {
		jitter := time.Duration(rand.Int63n(int64(delay)))
		delay = delay/2 + jitter // delay is now between [delay/2, delay]
	}

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RandomDelayWithJitter sleeps for a duration within ±25% of baseDelay to add jitter.
// This helps prevent thundering herd when multiple clients retry simultaneously.
// Example: baseDelay=2s results in sleep time between [1.5s, 2.5s).
func RandomDelayWithJitter(baseDelay time.Duration) {
	if baseDelay <= 0 {
		return
	}

	jitter := time.Duration(rand.Int63n(int64(baseDelay) / 2))
	delay := baseDelay*3/4 + jitter
	time.Sleep(delay)
}
