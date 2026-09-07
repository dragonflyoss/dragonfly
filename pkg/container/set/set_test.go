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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAdd(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, oks []bool, s Set[string])
	}{
		{
			name:   "add value",
			values: []string{"foo"},
			expect: func(t *testing.T, oks []bool, s Set[string]) {
				assert := assert.New(t)
				assert.Equal([]bool{true}, oks)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
		{
			name:   "add duplicate value",
			values: []string{"foo", "foo"},
			expect: func(t *testing.T, oks []bool, s Set[string]) {
				assert := assert.New(t)
				assert.Equal([]bool{true, false}, oks)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New[string]()
			oks := make([]bool, 0, len(tc.values))
			for _, v := range tc.values {
				oks = append(oks, s.Add(v))
			}

			tc.expect(t, oks, s)
		})
	}
}

func TestSetDelete(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		expect func(t *testing.T, s Set[string])
	}{
		{
			name:  "delete value",
			value: "foo",
			expect: func(t *testing.T, s Set[string]) {
				assert := assert.New(t)
				assert.Equal(uint(0), s.Len())
			},
		},
		{
			name:  "delete value does not exist",
			value: "bar",
			expect: func(t *testing.T, s Set[string]) {
				assert := assert.New(t)
				assert.Equal(uint(1), s.Len())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New[string]()
			s.Add("foo")
			s.Delete(tc.value)
			tc.expect(t, s)
		})
	}
}

func TestSetContains(t *testing.T) {
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
			s := New[string]()
			s.Add("foo")
			s.Add("bar")
			tc.expect(t, s.Contains(tc.values...))
		})
	}
}

func TestSetLen(t *testing.T) {
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
			s := New[string]()
			for _, v := range tc.values {
				s.Add(v)
			}

			tc.expect(t, s.Len())
		})
	}
}

func TestSetValues(t *testing.T) {
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
			s := New[string]()
			for _, v := range tc.values {
				s.Add(v)
			}

			tc.expect(t, s.Values())
		})
	}
}

func TestSetClear(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		expect func(t *testing.T, cleared []string, readded bool, s Set[string])
	}{
		{
			name:   "clear empty set",
			values: nil,
			expect: func(t *testing.T, cleared []string, readded bool, s Set[string]) {
				assert := assert.New(t)
				assert.Equal([]string(nil), cleared)
				assert.True(readded)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
		{
			name:   "clear set",
			values: []string{"foo"},
			expect: func(t *testing.T, cleared []string, readded bool, s Set[string]) {
				assert := assert.New(t)
				assert.Equal([]string(nil), cleared)
				assert.True(readded)
				assert.Equal([]string{"foo"}, s.Values())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := New[string]()
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
