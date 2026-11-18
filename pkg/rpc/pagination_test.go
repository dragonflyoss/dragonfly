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
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockPageItems = []string{"d", "b", "e", "a", "c"}
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name     string
		pageSize int32
		start    string
		expect   func(t *testing.T, page []string, nextPageToken *string, err error)
	}{
		{
			name: "page size is not specified",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"a", "b", "c", "d", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:     "page size is less than items",
			pageSize: 2,
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"a", "b"}, page)

				start, err := decodePageToken(*nextPageToken)
				assert.NoError(err)
				assert.Equal("b", start)
			},
		},
		{
			name:     "page starts after the key in the page token",
			pageSize: 2,
			start:    "b",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"c", "d"}, page)

				start, err := decodePageToken(*nextPageToken)
				assert.NoError(err)
				assert.Equal("d", start)
			},
		},
		{
			name:     "page size is equal to remaining items",
			pageSize: 2,
			start:    "c",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"d", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:     "page size is greater than remaining items",
			pageSize: 2,
			start:    "d",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:  "page size is not specified with page token",
			start: "c",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"d", "e"}, page)
				assert.Nil(nextPageToken)
			},
		},
		{
			name:     "key in the page token has been deleted",
			pageSize: 2,
			start:    "bb",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]string{"c", "d"}, page)
				assert.NotNil(nextPageToken)
			},
		},
		{
			name:     "no items after the key in the page token",
			pageSize: 2,
			start:    "e",
			expect: func(t *testing.T, page []string, nextPageToken *string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(page)
				assert.Nil(nextPageToken)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var token string
			if tc.start != "" {
				var err error
				token, err = encodePageToken(tc.start)
				assert.NoError(t, err)
			}

			page, nextPageToken, err := Paginate(mockPageItems, func(item string) string {
				return item
			}, tc.pageSize, token)
			tc.expect(t, page, nextPageToken, err)
		})
	}
}

func TestDecodePageToken(t *testing.T) {
	tests := []struct {
		name   string
		token  string
		expect func(t *testing.T, start string, err error)
	}{
		{
			name:  "page token is valid",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1","start":"foo"}`)),
			expect: func(t *testing.T, start string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("foo", start)
			},
		},
		{
			name:  "page token is not base64",
			token: "foo!",
			expect: func(t *testing.T, start string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "page token is not json",
			token: base64.RawURLEncoding.EncodeToString([]byte("foo")),
			expect: func(t *testing.T, start string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "page token version is unsupported",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v2","start":"foo"}`)),
			expect: func(t *testing.T, start string, err error) {
				assert := assert.New(t)
				assert.EqualError(err, `unsupported page token version "v2"`)
			},
		},
		{
			name:  "page token has no start key",
			token: base64.RawURLEncoding.EncodeToString([]byte(`{"v":"v1"}`)),
			expect: func(t *testing.T, start string, err error) {
				assert := assert.New(t)
				assert.EqualError(err, "page token has no start key")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			start, err := decodePageToken(tc.token)
			tc.expect(t, start, err)
		})
	}
}
