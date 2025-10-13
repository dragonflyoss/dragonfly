/*
 *     Copyright 2022 The Dragonfly Authors
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

import "go.uber.org/atomic"

func SafeSubAtomicUint64(counter *atomic.Uint64, delta uint64) {
	for {
		old := counter.Load()
		if old < delta {
			if counter.CompareAndSwap(old, 0) {
				return
			}

			continue
		}

		if counter.CompareAndSwap(old, old-delta) {
			return
		}
	}
}
