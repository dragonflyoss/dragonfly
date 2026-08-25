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

package oci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseManifestURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		expect func(t *testing.T, ref *Reference, err error)
	}{
		{
			name: "manifest url with tag",
			url:  "https://registry.example.com/v2/library/nginx/manifests/latest",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https", ref.Scheme)
				assert.Equal("registry.example.com", ref.Registry)
				assert.Equal("library/nginx", ref.Repository)
				assert.Equal("latest", ref.Reference)
			},
		},
		{
			name: "manifest url with digest",
			url:  "http://localhost:5000/v2/myrepo/manifests/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("http", ref.Scheme)
				assert.Equal("localhost:5000", ref.Registry)
				assert.Equal("myrepo", ref.Repository)
				assert.Equal("sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e", ref.Reference)
			},
		},
		{
			name: "invalid manifest url",
			url:  "https://registry.example.com/v2/library/nginx/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := ParseManifestURL(tc.url)
			tc.expect(t, ref, err)
		})
	}
}

func TestReferenceURLs(t *testing.T) {
	assert := assert.New(t)
	ref := &Reference{
		Scheme:     "https",
		Registry:   "registry.example.com",
		Repository: "library/nginx",
		Reference:  "latest",
	}

	assert.Equal("https://registry.example.com/v2/library/nginx/manifests/latest", ref.ManifestURL())
	assert.Equal(
		"https://registry.example.com/v2/library/nginx/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
		ref.BlobURL("sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e"),
	)
}
