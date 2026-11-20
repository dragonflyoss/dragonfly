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
	"context"
	"net"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	testSigningKey  = "dragonfly-test-signing-key"
	testIssuer      = "dragonfly-test"
	testWrongKey    = "another-signing-key"
	testWrongIssuer = "another-issuer"
)

// newTestHealthClient starts a real gRPC server on bufconn with the auth
// interceptors installed and returns a client connected to it.
func newTestHealthClient(t *testing.T, tm *TokenManager, skippedMethods ...string) healthpb.HealthClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor(tm, skippedMethods...)),
		grpc.StreamInterceptor(StreamServerInterceptor(tm, skippedMethods...)),
	)

	hs := health.NewServer()
	hs.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, hs)

	go func() {
		_ = server.Serve(lis)
	}()
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	return healthpb.NewHealthClient(conn)
}

// signToken signs arbitrary claims with the given key.
func signToken(t *testing.T, key string, method jwt.SigningMethod, claims jwt.Claims) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString([]byte(key))
	require.NoError(t, err)
	return signed
}

// expiredToken returns a correctly signed token whose exp is in the past.
func expiredToken(t *testing.T, tm *TokenManager) string {
	t.Helper()

	now := time.Now()
	return signToken(t, string(tm.key), jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tm.issuer,
			Subject:   "expired-client",
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now.Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
		},
	})
}

func newTestTokenManager(t *testing.T, key, issuer string) *TokenManager {
	t.Helper()

	tm, err := NewTokenManager(Config{Key: key, Issuer: issuer})
	require.NoError(t, err)
	return tm
}

func requireUnauthenticated(t *testing.T, err error) {
	t.Helper()

	require.Error(t, err)
	require.Equal(t, codes.Unauthenticated, status.Code(err), "unexpected error: %v", err)
}

// TestUnaryServerInterceptorWithoutToken covers a request that carries no
// authorization metadata at all.
func TestUnaryServerInterceptorWithoutToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	_, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{})
	requireUnauthenticated(t, err)
}

// TestUnaryServerInterceptorWithExpiredToken covers a correctly signed token
// whose validity window has passed.
func TestUnaryServerInterceptorWithExpiredToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	ctx := AppendTokenToOutgoingContext(context.Background(), expiredToken(t, tm))
	_, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	requireUnauthenticated(t, err)
}

// TestUnaryServerInterceptorWithWrongSignature covers a token signed by a key
// the server does not trust.
func TestUnaryServerInterceptorWithWrongSignature(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	attacker := newTestTokenManager(t, testWrongKey, testIssuer)
	token, err := attacker.IssueToken("attacker", time.Minute)
	require.NoError(t, err)

	ctx := AppendTokenToOutgoingContext(context.Background(), token)
	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{})
	requireUnauthenticated(t, err)
}

// TestUnaryServerInterceptorWithWrongIssuer covers a validly signed token that
// was issued by another authority.
func TestUnaryServerInterceptorWithWrongIssuer(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	other := newTestTokenManager(t, testSigningKey, testWrongIssuer)
	token, err := other.IssueToken("other", time.Minute)
	require.NoError(t, err)

	ctx := AppendTokenToOutgoingContext(context.Background(), token)
	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{})
	requireUnauthenticated(t, err)
}

// TestUnaryServerInterceptorWithNoneAlgorithm covers an unsigned token.
func TestUnaryServerInterceptorWithNoneAlgorithm(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	token, err := jwt.NewWithClaims(jwt.SigningMethodNone, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    testIssuer,
			Subject:   "nobody",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	ctx := AppendTokenToOutgoingContext(context.Background(), token)
	_, err = client.Check(ctx, &healthpb.HealthCheckRequest{})
	requireUnauthenticated(t, err)
}

// TestUnaryServerInterceptorWithMalformedMetadata covers a non-bearer scheme.
func TestUnaryServerInterceptorWithMalformedMetadata(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	ctx := metadata.AppendToOutgoingContext(context.Background(), MetadataKey, "Basic dXNlcjpwYXNz")
	_, err := client.Check(ctx, &healthpb.HealthCheckRequest{})
	requireUnauthenticated(t, err)
}

// TestUnaryServerInterceptorWithValidToken covers the happy path: the token is
// accepted and the call reaches the handler.
func TestUnaryServerInterceptorWithValidToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	token, err := tm.IssueToken("dfdaemon", time.Minute)
	require.NoError(t, err)

	resp, err := client.Check(AppendTokenToOutgoingContext(context.Background(), token), &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

// TestUnaryServerInterceptorSkipMethod covers the health check exemption.
func TestUnaryServerInterceptorSkipMethod(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm, "/grpc.health.v1.Health/")

	resp, err := client.Check(context.Background(), &healthpb.HealthCheckRequest{})
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

// TestStreamServerInterceptorWithoutToken covers an unauthenticated stream.
func TestStreamServerInterceptorWithoutToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.Watch(ctx, &healthpb.HealthCheckRequest{})
	require.NoError(t, err)

	_, err = stream.Recv()
	requireUnauthenticated(t, err)
}

// TestStreamServerInterceptorWithValidToken covers an authenticated stream.
func TestStreamServerInterceptorWithValidToken(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	client := newTestHealthClient(t, tm)

	token, err := tm.IssueToken("scheduler", time.Minute)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.Watch(AppendTokenToOutgoingContext(ctx, token), &healthpb.HealthCheckRequest{})
	require.NoError(t, err)

	resp, err := stream.Recv()
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
}

// TestUnaryServerInterceptorClaimsInContext covers that the verified subject is
// available to handlers.
func TestUnaryServerInterceptorClaimsInContext(t *testing.T) {
	tm := newTestTokenManager(t, testSigningKey, testIssuer)
	interceptor := UnaryServerInterceptor(tm)

	token, err := tm.IssueToken("manager", time.Minute)
	require.NoError(t, err)

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(MetadataKey, BearerPrefix+token))
	subject, ok := "", false
	_, err = interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test.v1.Test/Check"}, func(ctx context.Context, _ any) (any, error) {
		subject, ok = SubjectFromContext(ctx)
		return nil, nil
	})
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "manager", subject)
}

func TestValidateSkippedMethods(t *testing.T) {
	tests := []struct {
		name           string
		skippedMethods []string
		wantErr        bool
	}{
		{
			name:           "root prefix is rejected",
			skippedMethods: []string{"/"},
			wantErr:        true,
		},
		{
			name:           "partial prefix is rejected",
			skippedMethods: []string{"/grpc."},
			wantErr:        true,
		},
		{
			name:           "missing leading slash is rejected",
			skippedMethods: []string{"grpc.health.v1.Health/Check"},
			wantErr:        true,
		},
		{
			name:           "exact method is accepted",
			skippedMethods: []string{"/grpc.health.v1.Health/Check"},
		},
		{
			name:           "service prefix is accepted",
			skippedMethods: []string{"/grpc.health.v1.Health/"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSkippedMethods(tc.skippedMethods)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
