/*
 *     Copyright 2020 The Dragonfly Authors
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

package sortedmap

import (
	"d7y.io/dragonfly/v2/pkg/set"
)

type Bucket interface {
	Len() uint
	Add(uint, string) bool
	Delete(uint, string)
	Contains(uint, string) bool
	Range(func(string) bool)
	ReverseRange(fn func(string) bool)
}

type bucket struct {
	data  []set.Set
	start uint
	end   uint
}

func NewBucket(len uint) Bucket {
	return &bucket{
		data: make([]set.Set, len),
	}
}

func (b *bucket) Add(i uint, v string) bool {
	if b.Len() <= i {
		return false
	}

	if s := b.data[i]; s == nil {
		b.data[i] = set.New()

		if b.start >= i {
			b.start = i
		}

		if b.end <= i {
			b.end = i
		}
	}

	if ok := b.data[i].Add(v); !ok {
		return false
	}

	return true
}

func (b *bucket) Delete(i uint, v string) {
	if b.Len() <= i {
		return
	}

	if s := b.data[i]; s == nil {
		return
	}

	b.data[i].Delete(v)

	if b.data[i].Len() == 0 {
		b.shrink(i)
	}
}

func (b *bucket) Contains(i uint, v string) bool {
	if b.Len() <= i {
		return false
	}

	if s := b.data[i]; s == nil {
		return false
	}

	if ok := b.data[i].Contains(v); !ok {
		return false
	}

	return true
}

func (b *bucket) Len() uint {
	return uint(len(b.data))
}

func (b *bucket) Range(fn func(string) bool) {
	for i := b.start; i <= b.end; i++ {
		s := b.data[i]
		if s != nil && s.Len() > 0 {
			for _, v := range s.Values() {
				if !fn(v) {
					return
				}
			}
		}
	}
}

func (b *bucket) ReverseRange(fn func(string) bool) {
	for i := b.end; i >= b.start; i-- {
		s := b.data[i]
		if s != nil && s.Len() > 0 {
			for _, v := range s.Values() {
				if !fn(v) {
					return
				}
			}
		}

		if i == 0 {
			break
		}
	}
}

func (b *bucket) shrink(i uint) {
	if b.start < b.end && b.start == i {
		b.start++
	}

	if b.start < b.end && b.end == i {
		b.end--
	}
}
