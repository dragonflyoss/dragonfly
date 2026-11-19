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
	"net/http"
	"net/url"

	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	typesregistry "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
)

// Option is an option for AuthClient.
type Option func(c *AuthClient)

// WithIssuedToken sets the issued token for AuthClient.
func WithIssuedToken(token string) Option {
	return func(c *AuthClient) {
		c.issuedToken = token
	}
}

// AuthClient is an HTTP client that negotiates registry v2 authentication.
type AuthClient struct {
	// issuedToken is the issued token specified in header from user request,
	// there is no need to go through v2 authentication to get a new token
	// if the token is not empty, just use this token to access v2 API directly.
	issuedToken string

	// httpClient is the http client.
	httpClient *http.Client

	// authConfig is the auth config.
	authConfig *typesregistry.AuthConfig

	// authorizer sets the authorization of the registry on the request.
	authorizer transport.RequestModifier

	// ref is the reference of the image to authorize.
	ref *Reference
}

// NewAuthClient creates a new AuthClient for the registry of the given
// reference. If no username and password are provided, anonymous access is
// used.
func NewAuthClient(ref *Reference, httpClient *http.Client, username, password string, opts ...Option) (*AuthClient, error) {
	c := &AuthClient{
		httpClient: httpClient,
		authConfig: &typesregistry.AuthConfig{Username: username, Password: password},
		ref:        ref,
	}

	for _, opt := range opts {
		opt(c)
	}

	if len(c.issuedToken) > 0 {
		return c, nil
	}

	// New a challenge manager for the supported authentication types.
	challengeManager, err := registry.PingV2Registry(&url.URL{Scheme: ref.Scheme, Host: ref.Registry}, httpClient.Transport)
	if err != nil {
		return nil, err
	}

	// New a credential store which always returns the same credential values.
	creds := registry.NewStaticCredentialStore(c.authConfig)

	// Authorizer with the bearer token and the basic credentials, the token is
	// cached and renewed when it expires.
	c.authorizer = auth.NewAuthorizer(
		challengeManager,
		auth.NewTokenHandlerWithOptions(auth.TokenHandlerOptions{
			Transport:   httpClient.Transport,
			Credentials: creds,
			Scopes: []auth.Scope{auth.RepositoryScope{
				Repository: ref.Repository,
				Actions:    []string{"pull"},
			}},
			ClientID: registry.AuthClientID,
		}),
		auth.NewBasicHandler(creds),
	)

	// Copy the http client, so the transport of the given client is kept.
	client := *httpClient
	client.Transport = transport.NewTransport(httpClient.Transport, c.authorizer)
	c.httpClient = &client

	return c, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (c *AuthClient) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// Authorization returns the authorization of the registry, the bearer token or the
// basic credentials. It is empty if the registry needs no authorization. The request
// is authorized like the authorizer of containerd does, refer to
// https://github.com/containerd/containerd/blob/main/core/remotes/docker/authorizer.go.
func (c *AuthClient) Authorization() (string, error) {
	if len(c.issuedToken) > 0 {
		return c.issuedToken, nil
	}

	req, err := http.NewRequest(http.MethodGet, c.ref.ManifestURL(), nil)
	if err != nil {
		return "", err
	}

	if err := c.authorizer.ModifyRequest(req); err != nil {
		return "", err
	}

	return req.Header.Get("Authorization"), nil
}
