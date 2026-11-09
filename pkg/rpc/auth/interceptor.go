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
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// claimsKey is the private context key used to carry verified claims.
type claimsKey struct{}

// SubjectFromContext returns the token subject attached by the server
// interceptors. The boolean is false when the request was not authenticated.
func SubjectFromContext(ctx context.Context) (string, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	if !ok || claims == nil {
		return "", false
	}

	return claims.Subject, true
}

// ClaimsFromContext returns the verified token claims attached by the server
// interceptors.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok && claims != nil
}

// tokenFromMetadata extracts the raw token from the incoming "authorization"
// metadata value, which must use the "Bearer" scheme.
func tokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ErrMissingToken
	}

	values := md.Get(MetadataKey)
	if len(values) == 0 {
		return "", ErrMissingToken
	}

	// Reject ambiguous credentials silently carrying more than one token.
	if len(values) > 1 {
		return "", ErrMalformedToken
	}

	value := strings.TrimSpace(values[0])
	if len(value) <= len(BearerPrefix) || !strings.EqualFold(value[:len(BearerPrefix)], BearerPrefix) {
		return "", ErrMalformedToken
	}

	token := strings.TrimSpace(value[len(BearerPrefix):])
	if token == "" {
		return "", ErrMalformedToken
	}

	return token, nil
}

// authenticate verifies the token carried by ctx and returns the claims.
func authenticate(ctx context.Context, tm *TokenManager) (*Claims, error) {
	token, err := tokenFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return tm.ParseToken(token)
}

// isSkipped reports whether the method is exempt from authentication.
func isSkipped(method string, skippedMethods []string) bool {
	for _, skipped := range skippedMethods {
		if skipped == method || (strings.HasSuffix(skipped, "/") && strings.HasPrefix(method, skipped)) {
			return true
		}
	}

	return false
}

// UnaryServerInterceptor returns a gRPC unary server interceptor that verifies
// the bearer token in the authorization metadata. Requests without a valid
// token are rejected with codes.Unauthenticated.
//
// Methods listed in skippedMethods are not authenticated, which is useful for
// the gRPC health checking service (liveness/readiness probes). A skipped
// method may be a full method name ("/grpc.health.v1.Health/Check") or a
// service prefix ending with a slash ("/grpc.health.v1.Health/").
func UnaryServerInterceptor(tm *TokenManager, skippedMethods ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if isSkipped(info.FullMethod, skippedMethods) {
			return handler(ctx, req)
		}

		claims, err := authenticate(ctx, tm)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "failed to authenticate request: %v", err)
		}

		return handler(context.WithValue(ctx, claimsKey{}, claims), req)
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor that
// verifies the bearer token in the authorization metadata. Streams without a
// valid token are rejected with codes.Unauthenticated.
func StreamServerInterceptor(tm *TokenManager, skippedMethods ...string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if isSkipped(info.FullMethod, skippedMethods) {
			return handler(srv, ss)
		}

		claims, err := authenticate(ss.Context(), tm)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "failed to authenticate request: %v", err)
		}

		return handler(srv, &authenticatedServerStream{
			ServerStream: ss,
			ctx:          context.WithValue(ss.Context(), claimsKey{}, claims),
		})
	}
}

// authenticatedServerStream overrides the context of the wrapped stream so
// that handlers can read the verified claims.
type authenticatedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the context carrying the verified claims.
func (s *authenticatedServerStream) Context() context.Context {
	return s.ctx
}

// AppendTokenToOutgoingContext returns a context carrying the token as bearer
// authorization metadata for an outgoing call.
func AppendTokenToOutgoingContext(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, MetadataKey, BearerPrefix+token)
}

// NewUnaryClientInterceptor returns a unary client interceptor that attaches
// the provider's token to every outgoing call. The token is issued once and
// cached until it approaches expiry; see TokenProvider.
func NewUnaryClientInterceptor(provider *TokenProvider) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := provider.Token()
		if err != nil {
			return status.Errorf(codes.Internal, "failed to issue token: %v", err)
		}

		return invoker(AppendTokenToOutgoingContext(ctx, token), method, req, reply, cc, opts...)
	}
}

// NewStreamClientInterceptor returns a stream client interceptor that attaches
// the provider's token to every outgoing stream. The token is issued once and
// cached until it approaches expiry; see TokenProvider.
func NewStreamClientInterceptor(provider *TokenProvider) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		token, err := provider.Token()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to issue token: %v", err)
		}

		return streamer(AppendTokenToOutgoingContext(ctx, token), desc, cc, method, opts...)
	}
}
