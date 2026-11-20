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

package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestNewTokenManager(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg:  Config{Key: testSigningKey, Issuer: testIssuer},
		},
		{
			name:    "empty key",
			cfg:     Config{Issuer: testIssuer},
			wantErr: true,
		},
		{
			name:    "empty issuer",
			cfg:     Config{Key: testSigningKey},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tm, err := NewTokenManager(tc.cfg)
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, tm)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, tm)
		})
	}
}

func TestTokenManagerIssueAndParse(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	token, err := tm.IssueToken("dfdaemon", time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := tm.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, "dfdaemon", claims.Subject)
	require.Equal(t, testIssuer, claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
	require.True(t, claims.ExpiresAt.After(time.Now()))
}

func TestTokenManagerIssueTokenDefaultTTL(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	token, err := tm.IssueToken("dfdaemon", 0)
	require.NoError(t, err)

	claims, err := tm.ParseToken(token)
	require.NoError(t, err)
	require.WithinDuration(t, time.Now().Add(DefaultTokenTTL), claims.ExpiresAt.Time, time.Minute)
}

func TestTokenManagerIssueTokenEmptySubject(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	_, err := tm.IssueToken("", time.Minute)
	require.ErrorIs(t, err, ErrEmptySubject)
}

func TestTokenManagerParseEmptyToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	_, err := tm.ParseToken("")
	require.ErrorIs(t, err, ErrMissingToken)
}

func TestTokenManagerParseExpiredToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	_, err := tm.ParseToken(expiredToken(t, tm))
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenManagerParseWrongIssuer(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	other := newTestTokenManager(t, testSigningKey, testWrongIssuer)

	token, err := other.IssueToken("dfdaemon", time.Minute)
	require.NoError(t, err)

	_, err = tm.ParseToken(token)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenManagerParseMissingExpiration(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	token := signToken(t, testSigningKey, jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  testIssuer,
			Subject: "dfdaemon",
		},
	})

	_, err := tm.ParseToken(token)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenManagerParseUnknownAlgorithm(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)

	token := signToken(t, testSigningKey, jwt.SigningMethodHS512, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    testIssuer,
			Subject:   "dfdaemon",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	})

	_, err := tm.ParseToken(token)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenManagerAudience(t *testing.T) {
	tm, err := NewTokenManager(Config{Key: testSigningKey, Issuer: testIssuer, Audience: "scheduler"})
	require.NoError(t, err)

	token, err := tm.IssueToken("dfdaemon", time.Minute)
	require.NoError(t, err)

	claims, err := tm.ParseToken(token)
	require.NoError(t, err)
	require.Equal(t, jwt.ClaimStrings{"scheduler"}, claims.Audience)

	// A manager that requires another audience must reject the token.
	other, err := NewTokenManager(Config{Key: testSigningKey, Issuer: testIssuer, Audience: "manager"})
	require.NoError(t, err)

	_, err = other.ParseToken(token)
	require.ErrorIs(t, err, ErrInvalidToken)
}
