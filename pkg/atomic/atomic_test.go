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

package atomic

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	assert := assert.New(t)

	var zero Time
	assert.Equal(time.Time{}, zero.Load())

	now := time.Now()
	at := NewTime(now)
	assert.Equal(now, at.Load())

	next := now.Add(time.Hour)
	at.Store(next)
	assert.Equal(next, at.Load())
}

func TestDuration(t *testing.T) {
	assert := assert.New(t)

	var zero Duration
	assert.Equal(time.Duration(0), zero.Load())

	d := NewDuration(time.Second)
	assert.Equal(time.Second, d.Load())

	d.Store(time.Minute)
	assert.Equal(time.Minute, d.Load())
}
