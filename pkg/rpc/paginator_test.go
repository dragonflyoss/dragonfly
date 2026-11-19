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

package rpc

import (
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/pkg/digest"
)

var (
	mockPageItems  = []string{"d", "b", "e", "a", "c"}
	mockPageGroups = [][]string{{"d", "b"}, {"e", "a", "c"}}
)

func TestNewPaginator(t *testing.T) {
	tests := []struct {
		name      string
		pageToken string
		expect    func(t *testing.T, paginator *Paginator[string], err error)
	}{
		{
			name: "page token is not specified",
			expect: func(t *testing.T, paginator *Paginator[string], err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Zero(paginator.group)
				assert.Empty(paginator.start)
			},
		},
		{
			name:      "page token is valid",
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":1,"start":"b"}`)),
			expect: func(t *testing.T, paginator *Paginator[string], err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(1, paginator.group)
				assert.Equal("b", paginator.start)
			},
		},
		{
			name:      "page token is invalid",
			pageToken: "foo!",
			expect: func(t *testing.T, paginator *Paginator[string], err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrInvalidPageToken)
				assert.Nil(paginator)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(2, tc.pageToken, func(item string) string {
				return item
			})
			tc.expect(t, paginator, err)
		})
	}
}

func TestPaginator_Page(t *testing.T) {
	tests := []struct {
		name      string
		groups    [][]string
		pageSize  int32
		pageToken string
		expect    func(t *testing.T, page []string, nextPageToken *string, err error)
	}{
		{
			name:   "page size is not specified",
			groups: [][]string{mockPageItems},
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"a", "b", "c", "d", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:     "page size is less than items",
			groups:   [][]string{mockPageItems},
			pageSize: 2,
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"a", "b"}, page)
				if !assert.NotNil(nextPageToken) {
					return
				}

				assert.Equal(base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"b"}`)), *nextPageToken)
			},
		},
		{
			name:      "page starts after the key in the page token",
			groups:    [][]string{mockPageItems},
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"b"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"c", "d"}, page)
				if !assert.NotNil(nextPageToken) {
					return
				}

				assert.Equal(base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"d"}`)), *nextPageToken)
			},
		},
		{
			name:      "page size is equal to remaining items",
			groups:    [][]string{mockPageItems},
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"c"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"d", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:      "page size is greater than remaining items",
			groups:    [][]string{mockPageItems},
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"d"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:      "page size is not specified with page token",
			groups:    [][]string{mockPageItems},
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"c"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"d", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:      "key in the page token has been deleted",
			groups:    [][]string{mockPageItems},
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"bb"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"c", "d"}, page)
				assert.NotNil(nextPageToken)
			},
		},
		{
			name:      "no items after the key in the page token",
			groups:    [][]string{mockPageItems},
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"e"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:   "items of the first group are before items of the second group",
			groups: mockPageGroups,
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"b", "d", "a", "c", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:     "page ends in the first group",
			groups:   mockPageGroups,
			pageSize: 1,
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"b"}, page)
				if !assert.NotNil(nextPageToken) {
					return
				}

				assert.Equal(base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"b"}`)), *nextPageToken)
			},
		},
		{
			name:     "page crosses the groups",
			groups:   mockPageGroups,
			pageSize: 3,
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"b", "d", "a"}, page)
				if !assert.NotNil(nextPageToken) {
					return
				}

				assert.Equal(base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":1,"start":"a"}`)), *nextPageToken)
			},
		},
		{
			name:      "page starts after the last key in the first group",
			groups:    mockPageGroups,
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0,"start":"d"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"a", "c"}, page)
				if !assert.NotNil(nextPageToken) {
					return
				}

				assert.Equal(base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":1,"start":"c"}`)), *nextPageToken)
			},
		},
		{
			name:      "page starts after the key in the second group",
			groups:    mockPageGroups,
			pageSize:  2,
			pageToken: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":1,"start":"a"}`)),
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"c", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(tc.pageSize, tc.pageToken, func(item string) string {
				return item
			})
			if err != nil {
				t.Fatal(err)
			}

			page, nextPageToken, err := paginator.Page(tc.groups...)
			tc.expect(t, page, nextPageToken, err)
		})
	}
}

func TestPaginator_decodePageToken(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		expect func(t *testing.T, group int, start string, err error)
	}{
		{
			name:  "page token is valid",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":1,"start":"foo"}`)),
			expect: func(t *testing.T, group int, start string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(1, group)
				assert.Equal("foo", start)
			},
		},
		{
			name:  "page token is not base64",
			token: "foo!",
			expect: func(t *testing.T, group int, start string, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrInvalidPageToken)
			},
		},
		{
			name:  "page token is not json",
			token: base64.RawURLEncoding.EncodeToString([]byte("foo")),
			expect: func(t *testing.T, group int, start string, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrInvalidPageToken)
			},
		},
		{
			name:  "page token version is unsupported",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v2","group":0,"start":"foo"}`)),
			expect: func(t *testing.T, group int, start string, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrInvalidPageToken)
			},
		},
		{
			name:  "page token group is negative",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":-1,"start":"foo"}`)),
			expect: func(t *testing.T, group int, start string, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrInvalidPageToken)
			},
		},
		{
			name:  "page token has no start key",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","group":0}`)),
			expect: func(t *testing.T, group int, start string, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, ErrInvalidPageToken)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			paginator, err := NewPaginator(0, "", func(item string) string {
				return item
			})
			if err != nil {
				t.Fatal(err)
			}

			group, start, err := paginator.decodePageToken(tc.token)
			tc.expect(t, group, start, err)
		})
	}
}

func BenchmarkPaginator_Page(b *testing.B) {
	tests := []struct {
		name     string
		count    int
		pageSize int32
	}{
		{
			name:     "page size is 10 in 10000 items",
			count:    10000,
			pageSize: 10,
		},
		{
			name:     "page size is 100 in 100000 items",
			count:    100000,
			pageSize: 100,
		},
		{
			name:     "page size is 100 in 1000000 items",
			count:    1000000,
			pageSize: 100,
		},
		{
			name:  "page size is not specified in 100000 items",
			count: 100000,
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			items := make([]string, 0, tc.count)
			for i := range tc.count {
				items = append(items, digest.SHA256FromStrings(strconv.Itoa(i)))
			}

			paginator, err := NewPaginator(tc.pageSize, "", func(item string) string {
				return item
			})
			if err != nil {
				b.Fatal(err)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for n := 0; n < b.N; n++ {
				if _, _, err := paginator.Page(items); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
