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

// Package auth provides token (JWT) based authentication for the gRPC
// interactions between dragonfly components. It is intentionally standalone so
// that it can be wired into the manager and scheduler servers without changing
// their service implementations.
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// MetadataKey is the gRPC metadata key that carries the bearer token.
	MetadataKey = "authorization"

	// BearerPrefix is normalized to lowercase; extraction uses strings.EqualFold
	// so the standard "Bearer " spelling and other casings are accepted.
	BearerPrefix = "bearer "

	// SigningAlgorithm is the only algorithm accepted by this package.
	SigningAlgorithm = "HS256"

	// DefaultTokenTTL is used by IssueToken when ttl is not positive.
	DefaultTokenTTL = 10 * time.Minute
)

var (
	// ErrEmptySigningKey is returned when the HMAC signing key is not configured.
	ErrEmptySigningKey = errors.New("signing key is empty")

	// ErrEmptySubject is returned when issuing a token without a subject.
	ErrEmptySubject = errors.New("token subject is empty")

	// ErrMissingToken is returned when the authorization metadata is absent.
	ErrMissingToken = errors.New("authorization metadata is missing")

	// ErrMalformedToken is returned when the authorization metadata is not a
	// well-formed bearer token.
	ErrMalformedToken = errors.New("authorization metadata is not a valid bearer token")

	// ErrInvalidToken is returned when the token cannot be verified.
	ErrInvalidToken = errors.New("token is invalid")
)

// Claims is the payload of a dragonfly gRPC token. The standard registered
// claims carry the issuer, subject, audience and validity window.
type Claims struct {
	jwt.RegisteredClaims
}

// Config configures a TokenManager. Key must come from the component
// configuration (or a secret mounted from the environment); it must never be
// hardcoded in source.
type Config struct {
	// Key is the HMAC SHA-256 signing key. Required.
	Key string

	// Issuer is the expected "iss" claim. Required for verification and echoed
	// into issued tokens.
	Issuer string

	// Audience is the optional expected "aud" claim. When empty the audience
	// is not checked.
	Audience string
}

// TokenManager issues and verifies HS256 tokens.
type TokenManager struct {
	key      []byte
	issuer   string
	audience string
}

// NewTokenManager returns a TokenManager for the given configuration.
func NewTokenManager(cfg Config) (*TokenManager, error) {
	if cfg.Key == "" {
		return nil, ErrEmptySigningKey
	}

	if cfg.Issuer == "" {
		return nil, errors.New("token issuer is empty")
	}

	return &TokenManager{
		key:      []byte(cfg.Key),
		issuer:   cfg.Issuer,
		audience: cfg.Audience,
	}, nil
}

// IssueToken returns a signed HS256 token for the given subject. A non-positive
// ttl falls back to DefaultTokenTTL.
func (m *TokenManager) IssueToken(subject string, ttl time.Duration) (string, error) {
	if subject == "" {
		return "", ErrEmptySubject
	}

	if ttl <= 0 {
		ttl = DefaultTokenTTL
	}

	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	if m.audience != "" {
		claims.Audience = jwt.ClaimStrings{m.audience}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.key)
}

// ParseToken verifies the signature, algorithm, issuer, audience and validity
// window of a raw token and returns its claims.
func (m *TokenManager) ParseToken(raw string) (*Claims, error) {
	if raw == "" {
		return nil, ErrMissingToken
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: unexpected signing method %q", ErrInvalidToken, t.Header["alg"])
		}

		return m.key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// golang-jwt v4 does not expose parser options for these checks, so they
	// are enforced explicitly.
	if claims.Issuer != m.issuer {
		return nil, fmt.Errorf("%w: unexpected issuer %q", ErrInvalidToken, claims.Issuer)
	}

	// Expiration is mandatory so that tokens cannot outlive their validity
	// window by omitting the exp claim.
	if claims.ExpiresAt == nil {
		return nil, fmt.Errorf("%w: expiration is missing", ErrInvalidToken)
	}

	if m.audience != "" && !claims.VerifyAudience(m.audience, true) {
		return nil, fmt.Errorf("%w: unexpected audience %v", ErrInvalidToken, claims.Audience)
	}

	if claims.Subject == "" {
		return nil, fmt.Errorf("%w: subject is empty", ErrInvalidToken)
	}

	return claims, nil
}

// Issuer returns the configured issuer.
func (m *TokenManager) Issuer() string {
	return m.issuer
}
