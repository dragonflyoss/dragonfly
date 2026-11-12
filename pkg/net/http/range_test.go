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

package http

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_String(t *testing.T) {
	tests := []struct {
		name   string
		rg     Range
		expect func(t *testing.T, s string)
	}{
		{
			name: "bytes=0-9",
			rg:   Range{Start: 0, Length: 10},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("bytes=0-9", s)
			},
		},
		{
			name: "bytes=1-10",
			rg:   Range{Start: 1, Length: 10},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("bytes=1-10", s)
			},
		},
		{
			name: "bytes=1-0",
			rg:   Range{Start: 1, Length: 0},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("bytes=1-0", s)
			},
		},
		{
			name: "bytes=1-1",
			rg:   Range{Start: 1, Length: 1},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("bytes=1-1", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.rg.String())
		})
	}
}

func TestRange_URLMetaString(t *testing.T) {
	tests := []struct {
		name   string
		rg     Range
		expect func(t *testing.T, s string)
	}{
		{
			name: "0-9",
			rg:   Range{Start: 0, Length: 10},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("0-9", s)
			},
		},
		{
			name: "1-10",
			rg:   Range{Start: 1, Length: 10},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("1-10", s)
			},
		},
		{
			name: "1-0",
			rg:   Range{Start: 1, Length: 0},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("1-0", s)
			},
		},
		{
			name: "1-1",
			rg:   Range{Start: 1, Length: 1},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("1-1", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.rg.URLMetaString())
		})
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		size   int64
		expect func(t *testing.T, rg []Range, err error)
	}{
		{
			name: "empty header with zero size",
			s:    "",
			size: 0,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(rg)
			},
		},
		{
			name: "empty header",
			s:    "",
			size: 1000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(rg)
			},
		},
		{
			name: "missing bytes prefix",
			s:    "foo",
			size: 0,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "prefix only",
			s:    "bytes=",
			size: 0,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(rg)
			},
		},
		{
			name: "missing separator",
			s:    "bytes=7",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "missing separator with spaces",
			s:    "bytes= 7 ",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "start beyond zero size",
			s:    "bytes=1-",
			size: 0,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrNoOverlap)
				assert.Empty(rg)
			},
		},
		{
			name: "end before start",
			s:    "bytes=5-4",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "second range end before start",
			s:    "bytes=0-2,5-4",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "overlapping range end before start",
			s:    "bytes=2-5,4-3",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "negative suffix in list",
			s:    "bytes=--5,4--3",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "negative suffix",
			s:    "bytes=--5",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "negative suffix of one",
			s:    "bytes=--1",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "negative suffix with max size",
			s:    "bytes=--5",
			size: math.MaxInt64,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "separator only",
			s:    "bytes=-",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "non numeric start",
			s:    "bytes=A-",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "non numeric start with space",
			s:    "bytes=A- ",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "non numeric start and end",
			s:    "bytes=A-Z",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "non numeric suffix",
			s:    "bytes= -Z",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "non numeric end",
			s:    "bytes=5-Z",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "garbage",
			s:    "bytes=Ran-dom, garbage",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "hex numbers",
			s:    "bytes=0x01-0x02",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "whitespace only",
			s:    "bytes=         ",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(rg)
			},
		},
		{
			name: "commas and whitespace only",
			s:    "bytes= , , ,   ",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(rg)
			},
		},
		{
			name: "full range",
			s:    "bytes=0-9",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 10}}, rg)
			},
		},
		{
			name: "open ended from start",
			s:    "bytes=0-",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 10}}, rg)
			},
		},
		{
			name: "open ended from middle",
			s:    "bytes=5-",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{5, 5}}, rg)
			},
		},
		{
			name: "end clamped to size",
			s:    "bytes=0-20",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 10}}, rg)
			},
		},
		{
			name: "non overlapping range is skipped",
			s:    "bytes=15-,0-5",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 6}}, rg)
			},
		},
		{
			name: "two ranges",
			s:    "bytes=1-2,5-",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{1, 2}, {5, 5}}, rg)
			},
		},
		{
			name: "suffix and open ended",
			s:    "bytes=-2 , 7-",
			size: 11,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{9, 2}, {7, 4}}, rg)
			},
		},
		{
			name: "three ranges",
			s:    "bytes=0-0 ,2-2, 7-",
			size: 11,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 1}, {2, 1}, {7, 4}}, rg)
			},
		},
		{
			name: "suffix",
			s:    "bytes=-5",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{5, 5}}, rg)
			},
		},
		{
			name: "suffix larger than size",
			s:    "bytes=-15",
			size: 10,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 10}}, rg)
			},
		},
		{
			name: "first 500 bytes",
			s:    "bytes=0-499",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 500}}, rg)
			},
		},
		{
			name: "second 500 bytes",
			s:    "bytes=500-999",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{500, 500}}, rg)
			},
		},
		{
			name: "last 500 bytes",
			s:    "bytes=-500",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{9500, 500}}, rg)
			},
		},
		{
			name: "open ended last 500 bytes",
			s:    "bytes=9500-",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{9500, 500}}, rg)
			},
		},
		{
			name: "first and last byte",
			s:    "bytes=0-0,-1",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{0, 1}, {9999, 1}}, rg)
			},
		},
		{
			name: "adjacent ranges",
			s:    "bytes=500-600,601-999",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{500, 101}, {601, 399}}, rg)
			},
		},
		{
			name: "overlapping ranges",
			s:    "bytes=500-700,601-999",
			size: 10000,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{500, 201}, {601, 399}}, rg)
			},
		},
		{
			name: "apache laxity",
			s:    "bytes=   1 -2   ,  4- 5, 7 - 8 , ,,",
			size: 11,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{1, 2}, {4, 2}, {7, 2}}, rg)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rg, err := ParseRange(tc.s, tc.size)
			tc.expect(t, rg, err)
		})
	}
}

func TestParseURLMetaRange(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		size   int64
		expect func(t *testing.T, rg []Range, err error)
	}{
		{
			name: "full range",
			s:    "0-65575",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{Start: 0, Length: 65576}}, rg)
			},
		},
		{
			name: "single byte",
			s:    "2-2",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{Start: 2, Length: 1}}, rg)
			},
		},
		{
			name: "open ended",
			s:    "2-",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{Start: 2, Length: 65574}}, rg)
			},
		},
		{
			name: "suffix",
			s:    "-100",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{Start: 65476, Length: 100}}, rg)
			},
		},
		{
			name: "end clamped to size",
			s:    "0-66575",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]Range{{Start: 0, Length: 65576}}, rg)
			},
		},
		{
			name: "too many separators",
			s:    "0-65-575",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "non numeric end",
			s:    "0-hello",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "end before start",
			s:    "65575-0",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
		{
			name: "negative start",
			s:    "-1-8",
			size: 65576,
			expect: func(t *testing.T, rg []Range, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(rg)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rg, err := ParseURLMetaRange(tc.s, tc.size)
			tc.expect(t, rg, err)
		})
	}
}
