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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestParseImage(t *testing.T) {
	tests := []struct {
		name   string
		image  string
		expect func(t *testing.T, ref *Reference, err error)
	}{
		{
			name:  "short image reference",
			image: "nginx:latest",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("https", ref.Scheme)
				assert.Equal("registry-1.docker.io", ref.Registry)
				assert.Equal("library/nginx", ref.Repository)
				assert.Equal("latest", ref.Reference)
			},
		},
		{
			name:  "image reference without tag",
			image: "docker.io/library/nginx",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("registry-1.docker.io", ref.Registry)
				assert.Equal("library/nginx", ref.Repository)
				assert.Equal("latest", ref.Reference)
			},
		},
		{
			name:  "image reference with digest",
			image: "registry.example.com/library/nginx@sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("registry.example.com", ref.Registry)
				assert.Equal("library/nginx", ref.Repository)
				assert.Equal("sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e", ref.Reference)
			},
		},
		{
			name:  "image reference with registry port",
			image: "localhost:5000/myrepo:v1.0.0",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("localhost:5000", ref.Registry)
				assert.Equal("myrepo", ref.Repository)
				assert.Equal("v1.0.0", ref.Reference)
			},
		},
		{
			name:  "invalid image reference",
			image: "invalid image reference!!",
			expect: func(t *testing.T, ref *Reference, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.ErrorContains(err, "invalid image reference")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := ParseImage(tc.image)
			tc.expect(t, ref, err)
		})
	}
}

func TestResolveImage(t *testing.T) {
	assert := assert.New(t)

	// Mock registry: bearer challenge on /v2/, token endpoint and a schema2
	// manifest with a config and one layer.
	manifest := `{
		"schemaVersion": 2,
		"mediaType": "application/vnd.docker.distribution.manifest.v2+json",
		"config": {
			"mediaType": "application/vnd.docker.container.image.v1+json",
			"size": 7023,
			"digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7"
		},
		"layers": [
			{
				"mediaType": "application/vnd.docker.image.rootfs.diff.tar.gzip",
				"size": 32654,
				"digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f"
			}
		]
	}`

	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Bearer realm=%q,service=\"registry\"", server.URL+"/token"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"token":"test-token"}`)
	})
	mux.HandleFunc("/v2/library/nginx/manifests/latest", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Bearer realm=%q,service=\"registry\"", server.URL+"/token"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		fmt.Fprint(w, manifest)
	})
	server = httptest.NewServer(mux)
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	assert.NoError(err)

	// ParseImage always uses the https scheme, so resolve against the mock
	// registry by building the reference directly.
	ref := &Reference{
		Scheme:     "http",
		Registry:   serverURL.Host,
		Repository: "library/nginx",
		Reference:  "latest",
	}

	httpClient := &http.Client{Transport: http.DefaultTransport}
	client, err := NewAuthClient(ref, httpClient, "", "")
	assert.NoError(err)

	manifests, err := client.ResolveManifests(context.Background(), ref, make(http.Header), specs.Platform{OS: "linux", Architecture: "amd64"})
	assert.NoError(err)
	assert.Len(manifests, 1)

	var blobURLs []string
	for _, m := range manifests {
		for _, desc := range m.References() {
			blobURLs = append(blobURLs, ref.BlobURL(desc.Digest.String()))
		}
	}

	assert.Equal([]string{
		fmt.Sprintf("http://%s/v2/library/nginx/blobs/sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7", serverURL.Host),
		fmt.Sprintf("http://%s/v2/library/nginx/blobs/sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f", serverURL.Host),
	}, blobURLs)
	assert.Equal("Bearer test-token", client.AuthToken())
}

func TestResolveImageInvalidReference(t *testing.T) {
	assert := assert.New(t)

	_, _, err := ResolveImage(context.Background(), "invalid image reference!!")
	assert.Error(err)
	assert.ErrorContains(err, "invalid image reference")
}

func TestResolveImageInvalidPlatform(t *testing.T) {
	assert := assert.New(t)

	_, _, err := ResolveImage(context.Background(), "127.0.0.1:1/library/nginx:latest", WithPlatform("linux-amd64"))
	assert.Error(err)
	assert.ErrorContains(err, "invalid platform format")
}
