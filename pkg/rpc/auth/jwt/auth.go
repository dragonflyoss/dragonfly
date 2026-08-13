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

package jwt

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/credentials"
)

const (
	// TokenType is the explicit JWT type for Dragonfly inter-component gRPC
	// authentication.
	TokenType = "dragonfly-grpc+jwt"

	// AudienceManager is the expected audience of Manager gRPC servers.
	AudienceManager = "urn:dragonfly:grpc:manager"

	// AudienceScheduler is the expected audience of Scheduler gRPC servers.
	AudienceScheduler = "urn:dragonfly:grpc:scheduler"

	// AudienceDfdaemon is the expected audience of dfdaemon gRPC servers.
	AudienceDfdaemon = "urn:dragonfly:grpc:dfdaemon"

	authorizationMetadataKey = "authorization"
	bearerScheme             = "Bearer"
	maxCredentialLength      = 4 * 1024
)

var knownAudiences = map[string]struct{}{
	AudienceManager:   {},
	AudienceScheduler: {},
	AudienceDfdaemon:  {},
}

// Claims contains only the claims used by the first version of Dragonfly gRPC
// authentication.
type Claims struct {
	Issuer    string `json:"iss"`
	Audience  string `json:"aud"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// Valid is intentionally empty because Authenticator applies Dragonfly's
// validation rules with an injected clock and configured clock skew.
func (Claims) Valid() error {
	return nil
}

type cachedToken struct {
	value     string
	expiresAt time.Time
}

// Option configures an Authenticator.
type Option func(*Authenticator)

// WithClock overrides the clock used to generate and validate tokens.
func WithClock(now func() time.Time) Option {
	return func(authenticator *Authenticator) {
		authenticator.now = now
	}
}

// Authenticator signs and verifies inter-component gRPC JWTs.
type Authenticator struct {
	config Config
	keys   map[string][]byte
	now    func() time.Time

	cacheMu sync.Mutex
	cache   map[string]cachedToken
}

// New returns an Authenticator using the configured shared keyring.
func New(config Config, options ...Option) (*Authenticator, error) {
	keys, err := validateAndLoadKeys(config)
	if err != nil {
		return nil, err
	}

	config.Mode = config.EffectiveMode()
	authenticator := &Authenticator{
		config: config,
		keys:   keys,
		now:    time.Now,
		cache:  make(map[string]cachedToken),
	}

	for _, option := range options {
		option(authenticator)
	}

	if authenticator.now == nil {
		return nil, errors.New("grpc auth clock must not be nil")
	}

	return authenticator, nil
}

// Mode returns the configured effective authentication mode.
func (a *Authenticator) Mode() Mode {
	if a == nil {
		return ModeDisabled
	}

	return a.config.Mode
}

// Enabled returns whether authentication is enabled.
func (a *Authenticator) Enabled() bool {
	if a == nil {
		return false
	}

	return a.config.Enabled()
}

// PerRPCCredentials returns credentials that generate JWTs for audience.
func (a *Authenticator) PerRPCCredentials(audience string) credentials.PerRPCCredentials {
	return &perRPCCredentials{authenticator: a, audience: audience}
}

type perRPCCredentials struct {
	authenticator *Authenticator
	audience      string
}

// GetRequestMetadata returns authorization metadata for a new RPC.
func (c *perRPCCredentials) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	if !c.authenticator.Enabled() {
		return nil, nil
	}

	token, err := c.authenticator.token(c.audience)
	if err != nil {
		clientTokenEvents.WithLabelValues(c.audience, tokenEventGenerationFailed).Inc()
		return nil, err
	}

	return map[string]string{authorizationMetadataKey: bearerScheme + " " + token}, nil
}

// RequireTransportSecurity reports whether gRPC must refuse plaintext
// transports before attaching the JWT.
func (c *perRPCCredentials) RequireTransportSecurity() bool {
	return c.authenticator.Enabled() && c.authenticator.config.RequireTransportSecurity
}

func (a *Authenticator) token(audience string) (string, error) {
	if err := validateAudience(audience); err != nil {
		return "", err
	}

	now := a.now()
	cacheKey := audience + "\x00" + a.config.JWT.ActiveKeyID

	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()

	if cached, ok := a.cache[cacheKey]; ok && cached.expiresAt.Sub(now) > a.config.JWT.RefreshBefore {
		clientTokenEvents.WithLabelValues(audience, tokenEventCacheHit).Inc()
		return cached.value, nil
	}

	clientTokenEvents.WithLabelValues(audience, tokenEventCacheMiss).Inc()
	key, ok := a.keys[a.config.JWT.ActiveKeyID]
	if !ok {
		return "", fmt.Errorf("grpc auth active key id %q is unavailable", a.config.JWT.ActiveKeyID)
	}

	issuedAt := now.Unix()
	expiresAt := now.Add(a.config.JWT.TokenTTL)
	claims := Claims{
		Issuer:    a.config.JWT.Issuer,
		Audience:  audience,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt.Unix(),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	token.Header["typ"] = TokenType
	token.Header["kid"] = a.config.JWT.ActiveKeyID

	signed, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("sign grpc auth jwt: %w", err)
	}

	a.cache[cacheKey] = cachedToken{value: signed, expiresAt: expiresAt}
	clientTokenEvents.WithLabelValues(audience, tokenEventGenerated).Inc()
	return signed, nil
}

func validateAudience(audience string) error {
	if _, ok := knownAudiences[audience]; !ok {
		return fmt.Errorf("grpc auth has unsupported audience %q", audience)
	}

	return nil
}
