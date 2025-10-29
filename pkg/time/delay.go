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
	"time"
)

// LinearDelay implements a linear backoff strategy for retries. It calculates delay based on
// the attempt number and sleeps for that duration, capped at maxDelay.
func LinearDelay(ctx context.Context, attempt uint, increment, maxDelay time.Duration) error {
	delay := time.Duration(attempt) * increment
	if delay > maxDelay {
		delay = maxDelay
	}

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
