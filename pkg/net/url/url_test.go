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

package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterQuery(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		filters []string
		expect  func(t *testing.T, url string, err error)
	}{
		{
			name:    "filter repeated params and keep fragment",
			rawURL:  "http://www.xx.yy/path?u=f&x=y&m=z&x=s#size",
			filters: []string{"x", "m"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("http://www.xx.yy/path?u=f#size", url)
			},
		},
		{
			name:    "no filters returns raw url",
			rawURL:  "http://www.xx.yy/path?u=f&x=y&m=z&x=s#size",
			filters: []string{},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("http://www.xx.yy/path?u=f&x=y&m=z&x=s#size", url)
			},
		},
		{
			name:    "remaining params are sorted",
			rawURL:  "https://example.com/file.txt?z=9&b=2&a=1",
			filters: []string{"z"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https://example.com/file.txt?a=1&b=2", url)
			},
		},
		{
			name:    "same key params keep order",
			rawURL:  "https://example.com/file.txt?b=2&a=1&b=1",
			filters: []string{"c"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https://example.com/file.txt?a=1&b=2&b=1", url)
			},
		},
		{
			name:    "all params filtered",
			rawURL:  "https://example.com?foo=foo",
			filters: []string{"foo"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https://example.com", url)
			},
		},
		{
			name:    "params are escaped",
			rawURL:  "https://example.com/file.txt?k=a b&m=x*y&n=c~d",
			filters: []string{"none"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https://example.com/file.txt?k=a+b&m=x%2Ay&n=c~d", url)
			},
		},
		{
			name:    "semicolon separated param is dropped",
			rawURL:  "https://example.com/file.txt?a=1;x&b=2",
			filters: []string{"none"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https://example.com/file.txt?b=2", url)
			},
		},
		{
			name:    "invalid url",
			rawURL:  ":error_url",
			filters: []string{"x", "m"},
			expect: func(t *testing.T, url string, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(url)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			url, err := FilterQueryParams(tc.rawURL, tc.filters)
			tc.expect(t, url, err)
		})
	}
}
