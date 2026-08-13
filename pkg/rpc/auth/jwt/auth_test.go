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
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type interoperabilityFixture struct {
	SecretBase64 string `json:"secretBase64"`
	KeyID        string `json:"keyID"`
	Issuer       string `json:"issuer"`
	Audience     string `json:"audience"`
	IssuedAt     int64  `json:"issuedAt"`
	ExpiresAt    int64  `json:"expiresAt"`
	GoToken      string `json:"goToken"`
	RustToken    string `json:"rustToken"`
}

func TestPerRPCCredentials(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	authenticator := newTestAuthenticator(t, ModeRequired, func() time.Time { return now })
	credentials := authenticator.PerRPCCredentials(AudienceScheduler)

	md, err := credentials.GetRequestMetadata(context.Background())
	require.NoError(t, err)
	credential := md[authorizationMetadataKey]
	require.True(t, strings.HasPrefix(credential, bearerScheme+" "))
	assert.NoError(t, authenticator.verify(strings.TrimPrefix(credential, bearerScheme+" "), AudienceScheduler))
	assert.False(t, credentials.RequireTransportSecurity())
}

func TestInteroperabilityFixture(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("testdata", "interop.json"))
	require.NoError(t, err)

	var fixture interoperabilityFixture
	require.NoError(t, json.Unmarshal(content, &fixture))
	secret, err := base64.StdEncoding.DecodeString(fixture.SecretBase64)
	require.NoError(t, err)

	config := DefaultConfig()
	config.Mode = ModeRequired
	config.JWT.Issuer = fixture.Issuer
	config.JWT.TokenTTL = time.Duration(fixture.ExpiresAt-fixture.IssuedAt) * time.Second
	config.JWT.ActiveKeyID = fixture.KeyID
	config.JWT.Keys = []KeyConfig{{ID: fixture.KeyID, SecretFile: writeTestKey(t, string(secret))}}
	authenticator, err := New(config, WithClock(func() time.Time { return time.Unix(fixture.IssuedAt, 0) }))
	require.NoError(t, err)

	md, err := authenticator.PerRPCCredentials(fixture.Audience).GetRequestMetadata(context.Background())
	require.NoError(t, err)
	rawToken := strings.TrimPrefix(md[authorizationMetadataKey], bearerScheme+" ")
	assert.Equal(t, fixture.GoToken, rawToken)
	assert.NoError(t, authenticator.verify(fixture.GoToken, fixture.Audience))
	assert.NoError(t, authenticator.verify(fixture.RustToken, fixture.Audience))
}

func TestOverlappingKeysSupportRotation(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	oldKey := []byte(strings.Repeat("o", MinimumKeySize))
	newKey := []byte(strings.Repeat("n", MinimumKeySize))
	config := DefaultConfig()
	config.Mode = ModeRequired
	config.JWT.ActiveKeyID = "new"
	config.JWT.Keys = []KeyConfig{
		{ID: "old", SecretFile: writeTestKey(t, string(oldKey))},
		{ID: "new", SecretFile: writeTestKey(t, string(newKey))},
	}
	authenticator, err := New(config, WithClock(func() time.Time { return now }))
	require.NoError(t, err)

	claims := Claims{
		Issuer:    DefaultIssuer,
		Audience:  AudienceScheduler,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(DefaultTokenTTL).Unix(),
	}
	oldToken := signTestToken(t, claims, jwtlib.SigningMethodHS256, TokenType, "old", oldKey)
	assert.NoError(t, authenticator.verify(oldToken, AudienceScheduler))

	metadata, err := authenticator.PerRPCCredentials(AudienceScheduler).GetRequestMetadata(context.Background())
	require.NoError(t, err)
	newToken := strings.TrimPrefix(metadata[authorizationMetadataKey], bearerScheme+" ")
	parser := jwtlib.NewParser()
	header, _, err := parser.ParseUnverified(newToken, new(Claims))
	require.NoError(t, err)
	assert.Equal(t, "new", header.Header["kid"])
}

func TestDisabledPerRPCCredentials(t *testing.T) {
	authenticator, err := New(DefaultConfig())
	require.NoError(t, err)

	md, err := authenticator.PerRPCCredentials(AudienceManager).GetRequestMetadata(context.Background())
	require.NoError(t, err)
	assert.Empty(t, md)
}

func TestNilAuthenticatorIsDisabled(t *testing.T) {
	var authenticator *Authenticator
	assert.Equal(t, ModeDisabled, authenticator.Mode())
	assert.False(t, authenticator.Enabled())

	md, err := authenticator.PerRPCCredentials(AudienceManager).GetRequestMetadata(context.Background())
	require.NoError(t, err)
	assert.Empty(t, md)
}

func TestPerRPCCredentialsRequireTransportSecurity(t *testing.T) {
	config := validTestConfig(writeTestKey(t, strings.Repeat("k", MinimumKeySize)))
	config.RequireTransportSecurity = true
	authenticator, err := New(config)
	require.NoError(t, err)

	assert.True(t, authenticator.PerRPCCredentials(AudienceManager).RequireTransportSecurity())
}

func TestTokenCacheRefresh(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	clock := func() time.Time { return now }
	authenticator := newTestAuthenticator(t, ModeRequired, clock)
	credentials := authenticator.PerRPCCredentials(AudienceManager)

	first, err := credentials.GetRequestMetadata(context.Background())
	require.NoError(t, err)
	second, err := credentials.GetRequestMetadata(context.Background())
	require.NoError(t, err)
	assert.Equal(t, first, second)

	now = now.Add(DefaultTokenTTL - DefaultRefreshBefore)
	refreshed, err := credentials.GetRequestMetadata(context.Background())
	require.NoError(t, err)
	assert.NotEqual(t, first, refreshed)
}

func TestTokenProviderConcurrentAccess(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	authenticator := newTestAuthenticator(t, ModeRequired, func() time.Time { return now })
	credentials := authenticator.PerRPCCredentials(AudienceDfdaemon)

	const goroutines = 32
	results := make(chan string, goroutines)
	errors := make(chan error, goroutines)
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			md, err := credentials.GetRequestMetadata(context.Background())
			if err != nil {
				errors <- err
				return
			}
			results <- md[authorizationMetadataKey]
		}()
	}

	wg.Wait()
	close(results)
	close(errors)
	require.Empty(t, errors)

	var expected string
	for result := range results {
		if expected == "" {
			expected = result
		}
		assert.Equal(t, expected, result)
	}
}

func TestVerifyRejectsInvalidClaimsAndHeaders(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	authenticator := newTestAuthenticator(t, ModeRequired, func() time.Time { return now })
	validClaims := Claims{
		Issuer:    DefaultIssuer,
		Audience:  AudienceScheduler,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(DefaultTokenTTL).Unix(),
	}

	tests := []struct {
		name      string
		claims    Claims
		algorithm jwtlib.SigningMethod
		tokenType string
		keyID     string
		key       []byte
		reason    string
	}{
		{
			name:      "unsupported algorithm",
			claims:    validClaims,
			algorithm: jwtlib.SigningMethodHS384,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonUnsupportedAlg,
		},
		{
			name:      "invalid type",
			claims:    validClaims,
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: "JWT",
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonInvalidType,
		},
		{
			name:      "unknown key id",
			claims:    validClaims,
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "unknown",
			key:       authenticator.keys["test-key"],
			reason:    reasonUnknownKeyID,
		},
		{
			name:      "invalid signature",
			claims:    validClaims,
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       []byte(strings.Repeat("x", MinimumKeySize)),
			reason:    reasonInvalidSignature,
		},
		{
			name:      "invalid issuer",
			claims:    withClaims(validClaims, func(claims *Claims) { claims.Issuer = "other" }),
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonInvalidIssuer,
		},
		{
			name:      "invalid audience",
			claims:    withClaims(validClaims, func(claims *Claims) { claims.Audience = AudienceManager }),
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonInvalidAudience,
		},
		{
			name:      "future issued at",
			claims:    withClaims(validClaims, func(claims *Claims) { claims.IssuedAt = now.Add(DefaultClockSkew + time.Second).Unix() }),
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonInvalidIssuedAt,
		},
		{
			name:      "expired",
			claims:    withClaims(validClaims, func(claims *Claims) { claims.ExpiresAt = now.Add(-DefaultClockSkew).Unix() }),
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonExpired,
		},
		{
			name:      "ttl exceeded",
			claims:    withClaims(validClaims, func(claims *Claims) { claims.ExpiresAt = claims.IssuedAt + int64(DefaultMaxTokenTTL/time.Second) + 1 }),
			algorithm: jwtlib.SigningMethodHS256,
			tokenType: TokenType,
			keyID:     "test-key",
			key:       authenticator.keys["test-key"],
			reason:    reasonTTLExceeded,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rawToken := signTestToken(t, tc.claims, tc.algorithm, tc.tokenType, tc.keyID, tc.key)
			err := authenticator.verify(rawToken, AudienceScheduler)
			var authenticationError *authError
			require.ErrorAs(t, err, &authenticationError)
			assert.Equal(t, tc.reason, authenticationError.reason)
		})
	}
}

func TestUnaryServerInterceptorModes(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	fullMethod := "/dragonfly.scheduler.v2.Scheduler/AnnounceHost"

	tests := []struct {
		name        string
		mode        Mode
		metadata    metadata.MD
		wantCode    codes.Code
		wantHandled bool
	}{
		{name: "disabled", mode: ModeDisabled, wantCode: codes.OK, wantHandled: true},
		{name: "permissive missing", mode: ModePermissive, wantCode: codes.OK, wantHandled: true},
		{name: "permissive invalid", mode: ModePermissive, metadata: metadata.Pairs(authorizationMetadataKey, "invalid"), wantCode: codes.Unauthenticated},
		{name: "required missing", mode: ModeRequired, wantCode: codes.Unauthenticated},
		{name: "required malformed", mode: ModeRequired, metadata: metadata.Pairs(authorizationMetadataKey, "invalid"), wantCode: codes.Unauthenticated},
		{name: "required duplicate", mode: ModeRequired, metadata: metadata.MD{authorizationMetadataKey: []string{"Bearer first", "Bearer second"}}, wantCode: codes.Unauthenticated},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			authenticator := newTestAuthenticator(t, tc.mode, func() time.Time { return now })
			ctx := metadata.NewIncomingContext(context.Background(), tc.metadata)
			handled := false
			_, err := authenticator.UnaryServerInterceptor(AudienceScheduler)(ctx, nil, &grpc.UnaryServerInfo{FullMethod: fullMethod}, func(context.Context, any) (any, error) {
				handled = true
				return nil, nil
			})
			assert.Equal(t, tc.wantCode, status.Code(err))
			assert.Equal(t, tc.wantHandled, handled)
		})
	}
}

func TestUnaryServerInterceptorAcceptsValidToken(t *testing.T) {
	now := time.Unix(1_786_435_200, 0)
	authenticator := newTestAuthenticator(t, ModeRequired, func() time.Time { return now })
	clientMetadata, err := authenticator.PerRPCCredentials(AudienceScheduler).GetRequestMetadata(context.Background())
	require.NoError(t, err)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(clientMetadata))
	handled := false
	_, err = authenticator.UnaryServerInterceptor(AudienceScheduler)(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/dragonfly.scheduler.v2.Scheduler/AnnounceHost"}, func(context.Context, any) (any, error) {
		handled = true
		return nil, nil
	})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestHealthBypassesAuthentication(t *testing.T) {
	authenticator := newTestAuthenticator(t, ModeRequired, time.Now)
	handled := false
	_, err := authenticator.UnaryServerInterceptor(AudienceScheduler)(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}, func(context.Context, any) (any, error) {
		handled = true
		return nil, nil
	})
	require.NoError(t, err)
	assert.True(t, handled)
}

func TestReflectionDoesNotBypassAuthentication(t *testing.T) {
	authenticator := newTestAuthenticator(t, ModeRequired, time.Now)
	stream := &testServerStream{ctx: context.Background()}
	err := authenticator.StreamServerInterceptor(AudienceScheduler)(nil, stream, &grpc.StreamServerInfo{FullMethod: "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo"}, func(any, grpc.ServerStream) error {
		return nil
	})
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

type testServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *testServerStream) Context() context.Context {
	return s.ctx
}

func newTestAuthenticator(t *testing.T, mode Mode, clock func() time.Time) *Authenticator {
	t.Helper()

	if mode == ModeDisabled {
		authenticator, err := New(DefaultConfig(), WithClock(clock))
		require.NoError(t, err)
		return authenticator
	}

	config := validTestConfig(writeTestKey(t, strings.Repeat("k", MinimumKeySize)))
	config.Mode = mode
	authenticator, err := New(config, WithClock(clock))
	require.NoError(t, err)
	return authenticator
}

func signTestToken(t *testing.T, claims Claims, method jwtlib.SigningMethod, tokenType, keyID string, key []byte) string {
	t.Helper()

	token := jwtlib.NewWithClaims(method, claims)
	token.Header["typ"] = tokenType
	token.Header["kid"] = keyID
	signed, err := token.SignedString(key)
	require.NoError(t, err)
	return signed
}

func withClaims(claims Claims, mutate func(*Claims)) Claims {
	mutate(&claims)
	return claims
}
