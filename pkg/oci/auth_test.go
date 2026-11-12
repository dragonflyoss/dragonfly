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
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthClient(t *testing.T) {
	ref := &Reference{
		Scheme:     "http",
		Registry:   "127.0.0.1:1",
		Repository: "library/nginx",
		Reference:  "latest",
	}

	tests := []struct {
		name       string
		httpClient *http.Client
		opts       []Option
		expect     func(t *testing.T, client *AuthClient, err error)
	}{
		{
			name:       "issued token skips registry ping",
			httpClient: &http.Client{},
			opts:       []Option{WithIssuedToken("Bearer issued-token")},
			expect: func(t *testing.T, client *AuthClient, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				authorization, err := client.Authorization()
				assert.NoError(err)
				assert.Equal("Bearer issued-token", authorization)
			},
		},
		{
			name:       "unreachable registry",
			httpClient: &http.Client{Transport: http.DefaultTransport},
			expect: func(t *testing.T, client *AuthClient, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(client)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewAuthClient(ref, tc.httpClient, "", "", tc.opts...)
			tc.expect(t, client, err)
		})
	}
}

func TestAuthClientTokenChallenge(t *testing.T) {
	assert := assert.New(t)
	var server *httptest.Server
	var tokens atomic.Int32
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
		tokens.Add(1)
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

	httpClient := &http.Client{Transport: http.DefaultTransport}
	client, err := NewAuthClient(ref, httpClient, "", "")
	assert.NoError(err)
	assert.Same(http.DefaultTransport, httpClient.Transport)

	req, err := http.NewRequest(http.MethodGet, ref.ManifestURL(), nil)
	assert.NoError(err)

	resp, err := client.Do(req)
	assert.NoError(err)
	defer resp.Body.Close()

	assert.Equal(http.StatusOK, resp.StatusCode)

	authorization, err := client.Authorization()
	assert.NoError(err)
	assert.Equal("Bearer test-token", authorization)
	assert.Equal(int32(1), tokens.Load())
}
