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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthClientWithIssuedToken(t *testing.T) {
	assert := assert.New(t)
	ref := &Reference{
		Scheme:     "http",
		Registry:   "127.0.0.1:1",
		Repository: "library/nginx",
		Reference:  "latest",
	}

	client, err := NewAuthClient(ref, &http.Client{}, "", "", WithIssuedToken("Bearer issued-token"))
	assert.NoError(err)
	assert.Equal("Bearer issued-token", client.AuthToken())
}

func TestNewAuthClientUnreachableRegistry(t *testing.T) {
	assert := assert.New(t)
	ref := &Reference{
		Scheme:     "http",
		Registry:   "127.0.0.1:1",
		Repository: "library/nginx",
		Reference:  "latest",
	}

	_, err := NewAuthClient(ref, &http.Client{Transport: http.DefaultTransport}, "", "")
	assert.Error(err)
}

func TestAuthClientTokenChallenge(t *testing.T) {
	assert := assert.New(t)
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
		assert.Contains(r.URL.Query().Get("scope"), "repository:library/nginx:pull")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"token":"test-token"}`)
	})

	server = httptest.NewServer(mux)
	defer server.Close()

	serverURL, err := url.Parse(server.URL)
	assert.NoError(err)

	ref := &Reference{
		Scheme:     "http",
		Registry:   serverURL.Host,
		Repository: "library/nginx",
		Reference:  "latest",
	}

	client, err := NewAuthClient(ref, &http.Client{Transport: http.DefaultTransport}, "", "")
	assert.NoError(err)

	req, err := http.NewRequest(http.MethodGet, ref.ManifestURL(), nil)
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)
	defer resp.Body.Close()

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("Bearer test-token", client.AuthToken())
}

func TestTokenHandler(t *testing.T) {
	assert := assert.New(t)

	h := newTokenHandler()
	assert.Equal("bearer", h.Scheme())

	req, err := http.NewRequest(http.MethodGet, "http://registry.example.com/v2/", nil)
	assert.NoError(err)
	req.Header.Set("Authorization", "Bearer test-token")

	assert.NoError(h.AuthorizeRequest(req, nil))
	assert.Equal("Bearer test-token", h.token)
}
