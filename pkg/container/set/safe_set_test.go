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

package set

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const N = 1000

func TestSafeSetAdd(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, oks []bool, s SafeSet[string])
	}{
		{
			name:   "add value",
			values: []string{"foo"},
			expect: func(t *testing.T, oks []bool, s SafeSet[string]) {
				assert := assert.New(t)
				assert.Equal([]bool{true}, oks)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
		{
			name:   "add duplicate value",
			values: []string{"foo", "foo"},
			expect: func(t *testing.T, oks []bool, s SafeSet[string]) {
				assert := assert.New(t)
				assert.Equal([]bool{true, false}, oks)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSafeSet[string]()
			oks := make([]bool, 0, len(tc.values))
			for _, v := range tc.values {
				oks = append(oks, s.Add(v))
			}

			tc.expect(t, oks, s)
		})
	}
}

func TestSafeSetAdd_Concurrent(t *testing.T) {
	assert := assert.New(t)
	runtime.GOMAXPROCS(2)

	s := NewSafeSet[int]()
	nums := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(nums))
	for i := range len(nums) {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	assert.True(s.Contains(nums...))
}

func TestSafeSetDelete(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		expect func(t *testing.T, s SafeSet[string])
	}{
		{
			name:  "delete value",
			value: "foo",
			expect: func(t *testing.T, s SafeSet[string]) {
				assert := assert.New(t)
				assert.Equal(uint(0), s.Len())
			},
		},
		{
			name:  "delete value does not exist",
			value: "bar",
			expect: func(t *testing.T, s SafeSet[string]) {
				assert := assert.New(t)
				assert.Equal(uint(1), s.Len())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSafeSet[string]()
			s.Add("foo")
			s.Delete(tc.value)
			tc.expect(t, s)
		})
	}
}

func TestSafeSetDelete_Concurrent(t *testing.T) {
	assert := assert.New(t)
	runtime.GOMAXPROCS(2)

	s := NewSafeSet[int]()
	nums := rand.Perm(N)
	for _, v := range nums {
		s.Add(v)
	}

	var wg sync.WaitGroup
	wg.Add(len(nums))
	for _, v := range nums {
		go func(i int) {
			s.Delete(i)
			wg.Done()
		}(v)
	}

	wg.Wait()

	assert.Equal(uint(0), s.Len())
}

func TestSafeSetContains(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, ok bool)
	}{
		{
			name:   "contains value",
			values: []string{"foo"},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:   "contains value does not exist",
			values: []string{"baz"},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name:   "contains all of multiple values",
			values: []string{"foo", "bar"},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:   "contains fails when one of multiple values is missing",
			values: []string{"foo", "baz"},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name:   "contains no values",
			values: nil,
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSafeSet[string]()
			s.Add("foo")
			s.Add("bar")
			tc.expect(t, s.Contains(tc.values...))
		})
	}
}

func TestSafeSetContains_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSafeSet[int]()
	nums := rand.Perm(N)
	for _, v := range nums {
		s.Add(v)
	}

	var wg sync.WaitGroup
	for range nums {
		wg.Add(1)
		go func() {
			assert := assert.New(t)
			defer wg.Done()
			assert.True(s.Contains(nums...))
		}()
	}

	wg.Wait()
}

func TestSetSafeLen(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, n uint)
	}{
		{
			name:   "get length",
			values: []string{"foo"},
			expect: func(t *testing.T, n uint) {
				assert := assert.New(t)
				assert.Equal(uint(1), n)
			},
		},
		{
			name:   "get empty set length",
			values: nil,
			expect: func(t *testing.T, n uint) {
				assert := assert.New(t)
				assert.Equal(uint(0), n)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSafeSet[string]()
			for _, v := range tc.values {
				s.Add(v)
			}

			tc.expect(t, s.Len())
		})
	}
}

func TestSafeSetLen_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSafeSet[int]()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		assert := assert.New(t)
		defer wg.Done()
		elems := s.Len()
		for range N {
			newElems := s.Len()
			assert.GreaterOrEqual(newElems, elems)
			elems = newElems
		}
	}()

	for range N {
		s.Add(rand.Int())
	}

	wg.Wait()
}

func TestSafeSetValues(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, values []string)
	}{
		{
			name:   "get values",
			values: []string{"foo"},
			expect: func(t *testing.T, values []string) {
				assert := assert.New(t)
				assert.Equal([]string{"foo"}, values)
			},
		},
		{
			name:   "get empty values",
			values: nil,
			expect: func(t *testing.T, values []string) {
				assert := assert.New(t)
				assert.Equal([]string(nil), values)
			},
		},
		{
			name:   "get multi values",
			values: []string{"foo", "bar"},
			expect: func(t *testing.T, values []string) {
				assert := assert.New(t)
				assert.ElementsMatch([]string{"foo", "bar"}, values)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSafeSet[string]()
			for _, v := range tc.values {
				s.Add(v)
			}

			tc.expect(t, s.Values())
		})
	}
}

func TestSafeSetValues_Concurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSafeSet[int]()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		assert := assert.New(t)
		defer wg.Done()
		elems := s.Values()
		for range N {
			newElems := s.Values()
			assert.GreaterOrEqual(len(newElems), len(elems))
			elems = newElems
		}
	}()

	for i := range N {
		s.Add(i)
	}

	wg.Wait()
}

func TestSafeSetClear(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, cleared []string, readded bool, s SafeSet[string])
	}{
		{
			name:   "clear empty set",
			values: nil,
			expect: func(t *testing.T, cleared []string, readded bool, s SafeSet[string]) {
				assert := assert.New(t)
				assert.Equal([]string(nil), cleared)
				assert.True(readded)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
		{
			name:   "clear set",
			values: []string{"foo"},
			expect: func(t *testing.T, cleared []string, readded bool, s SafeSet[string]) {
				assert := assert.New(t)
				assert.Equal([]string(nil), cleared)
				assert.True(readded)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSafeSet[string]()
			for _, v := range tc.values {
				s.Add(v)
			}

			s.Clear()
			cleared := s.Values()
			readded := s.Add("foo")
			tc.expect(t, cleared, readded, s)
		})
	}
}

func TestSafeSetClear_Concurrent(t *testing.T) {
	assert := assert.New(t)
	runtime.GOMAXPROCS(2)

	s := NewSafeSet[int]()
	nums := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(nums))
	for i := range len(nums) {
		go func(i int) {
			s.Add(i)
			s.Clear()
			wg.Done()
		}(i)
	}

	wg.Wait()
	assert.Equal(uint(0), s.Len())
}
