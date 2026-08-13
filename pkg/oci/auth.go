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

// AuthClient is an HTTP client that negotiates registry v2 authentication and
// intercepts the issued bearer token.
type AuthClient struct {
	// issuedToken is the issued token specified in header from user request,
	// there is no need to go through v2 authentication to get a new token
	// if the token is not empty, just use this token to access v2 API directly.
	issuedToken string

	// httpClient is the http client.
	httpClient *http.Client

	// authConfig is the auth config.
	authConfig *typesregistry.AuthConfig

	// tokenHandler is the token interceptor.
	tokenHandler *tokenHandler
}

// NewAuthClient creates a new AuthClient for the registry of the given
// reference. If no username and password are provided, anonymous access is
// used.
func NewAuthClient(ref *Reference, httpClient *http.Client, username, password string, opts ...Option) (*AuthClient, error) {
	c := &AuthClient{
		httpClient:   httpClient,
		authConfig:   &typesregistry.AuthConfig{Username: username, Password: password},
		tokenHandler: newTokenHandler(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if len(c.issuedToken) > 0 {
		return c, nil
	}

	// New a challenge manager for the supported authentication types.
	challengeManager, err := registry.PingV2Registry(&url.URL{Scheme: ref.Scheme, Host: ref.Registry}, c.httpClient.Transport)
	if err != nil {
		return nil, err
	}

	// New a credential store which always returns the same credential values.
	creds := registry.NewStaticCredentialStore(c.authConfig)

	// Transport with authentication.
	c.httpClient.Transport = transport.NewTransport(
		c.httpClient.Transport,
		auth.NewAuthorizer(
			challengeManager,
			auth.NewTokenHandlerWithOptions(auth.TokenHandlerOptions{
				Transport:   c.httpClient.Transport,
				Credentials: creds,
				Scopes: []auth.Scope{auth.RepositoryScope{
					Repository: ref.Repository,
					Actions:    []string{"pull"},
				}},
				ClientID: registry.AuthClientID,
			}),
			c.tokenHandler,
			auth.NewBasicHandler(creds),
		),
	)

	return c, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (c *AuthClient) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// AuthToken returns the bearer token.
func (c *AuthClient) AuthToken() string {
	if len(c.issuedToken) > 0 {
		return c.issuedToken
	}

	return c.tokenHandler.token
}

// tokenHandler is a token interceptor intercept bearer token from auth handler.
type tokenHandler struct {
	auth.AuthenticationHandler
	token string
}

// newTokenHandler returns a new tokenHandler.
func newTokenHandler() *tokenHandler {
	return &tokenHandler{}
}

// Scheme returns the authentication scheme.
func (h *tokenHandler) Scheme() string {
	return "bearer"
}

// AuthorizeRequest saves the Authorization header from the request.
func (h *tokenHandler) AuthorizeRequest(req *http.Request, params map[string]string) error {
	h.token = req.Header.Get("Authorization")
	return nil
}
