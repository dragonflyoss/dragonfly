/*
 *     Copyright 2025 The Dragonfly Authors
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

package oauth

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	oauth2github "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		oauthName string
		expect    func(t *testing.T, o Oauth, err error)
	}{
		{
			name:      "google",
			oauthName: Google,
			expect: func(t *testing.T, o Oauth, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.IsType(&oauthGoogle{}, o)
			},
		},
		{
			name:      "github",
			oauthName: Github,
			expect: func(t *testing.T, o Oauth, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.IsType(&oauthGithub{}, o)
			},
		},
		{
			name:      "name is case sensitive",
			oauthName: "Google",
			expect: func(t *testing.T, o Oauth, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(o)
			},
		},
		{
			name:      "unknown name",
			oauthName: "gitlab",
			expect: func(t *testing.T, o Oauth, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(o)
			},
		},
		{
			name:      "empty name",
			oauthName: "",
			expect: func(t *testing.T, o Oauth, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(o)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o, err := New(tc.oauthName, "client-id", "client-secret", "https://manager.example.com/callback")
			tc.expect(t, o, err)
		})
	}
}

func TestOauth_AuthCodeURL(t *testing.T) {
	tests := []struct {
		name      string
		oauthName string
		expect    func(t *testing.T, authCodeURL string, err error)
	}{
		{
			name:      "google",
			oauthName: Google,
			expect: func(t *testing.T, authCodeURL string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				u, err := url.Parse(authCodeURL)
				assert.NoError(err)
				assert.True(strings.HasPrefix(authCodeURL, google.Endpoint.AuthURL))
				assert.Equal("code", u.Query().Get("response_type"))
				assert.Equal("client-id", u.Query().Get("client_id"))
				assert.Equal("https://manager.example.com/callback", u.Query().Get("redirect_uri"))
				assert.Equal(strings.Join(googleScopes, " "), u.Query().Get("scope"))
				assert.Len(u.Query().Get("state"), 24)
			},
		},
		{
			name:      "github",
			oauthName: Github,
			expect: func(t *testing.T, authCodeURL string, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				u, err := url.Parse(authCodeURL)
				assert.NoError(err)
				assert.True(strings.HasPrefix(authCodeURL, oauth2github.Endpoint.AuthURL))
				assert.Equal("code", u.Query().Get("response_type"))
				assert.Equal("client-id", u.Query().Get("client_id"))
				assert.Equal("https://manager.example.com/callback", u.Query().Get("redirect_uri"))
				assert.Equal(strings.Join(githubScopes, " "), u.Query().Get("scope"))
				assert.Len(u.Query().Get("state"), 24)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o, err := New(tc.oauthName, "client-id", "client-secret", "https://manager.example.com/callback")
			if err != nil {
				t.Fatal(err)
			}

			authCodeURL, err := o.AuthCodeURL()
			tc.expect(t, authCodeURL, err)
		})
	}
}

func TestOauth_AuthCodeURLStateIsRandom(t *testing.T) {
	tests := []struct {
		name      string
		oauthName string
		expect    func(t *testing.T, first, second string)
	}{
		{
			name:      "google",
			oauthName: Google,
			expect: func(t *testing.T, first, second string) {
				assert := assert.New(t)
				assert.NotEmpty(first)
				assert.NotEqual(first, second)
			},
		},
		{
			name:      "github",
			oauthName: Github,
			expect: func(t *testing.T, first, second string) {
				assert := assert.New(t)
				assert.NotEmpty(first)
				assert.NotEqual(first, second)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			o, err := New(tc.oauthName, "client-id", "client-secret", "https://manager.example.com/callback")
			if err != nil {
				t.Fatal(err)
			}

			firstURL, err := o.AuthCodeURL()
			if err != nil {
				t.Fatal(err)
			}

			secondURL, err := o.AuthCodeURL()
			if err != nil {
				t.Fatal(err)
			}

			first, err := url.Parse(firstURL)
			if err != nil {
				t.Fatal(err)
			}

			second, err := url.Parse(secondURL)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, first.Query().Get("state"), second.Query().Get("state"))
		})
	}
}
