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
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authenticationErrorMessage = "invalid authentication credentials"
	healthMethodPrefix         = "/grpc.health.v1.Health/"
)

const (
	reasonNone             = "none"
	reasonMissing          = "missing"
	reasonMalformed        = "malformed"
	reasonUnsupportedAlg   = "unsupported_alg"
	reasonInvalidType      = "invalid_type"
	reasonUnknownKeyID     = "unknown_kid"
	reasonInvalidSignature = "invalid_signature"
	reasonInvalidIssuer    = "invalid_issuer"
	reasonInvalidAudience  = "invalid_audience"
	reasonExpired          = "expired"
	reasonInvalidIssuedAt  = "invalid_iat"
	reasonTTLExceeded      = "ttl_exceeded"
)

type authError struct {
	reason string
	cause  error
}

var errPermissiveMissing = &authError{
	reason: reasonMissing,
	cause:  errors.New("authorization metadata is missing in permissive mode"),
}

func (e *authError) Error() string {
	return e.cause.Error()
}

func (e *authError) Unwrap() error {
	return e.cause
}

// UnaryServerInterceptor authenticates unary gRPC requests for audience.
func (a *Authenticator) UnaryServerInterceptor(audience string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if a.skipServerAuthentication(info.FullMethod) {
			return handler(ctx, req)
		}

		if err := a.authenticate(ctx, audience); errors.Is(err, errPermissiveMissing) {
			a.recordPermissiveMissing(audience)
			return handler(ctx, req)
		} else if err != nil {
			a.recordServerResult(audience, err)
			return nil, status.Error(codes.Unauthenticated, authenticationErrorMessage)
		}

		a.recordServerResult(audience, nil)
		return handler(ctx, req)
	}
}

// StreamServerInterceptor authenticates streaming gRPC requests for audience
// when a stream is established.
func (a *Authenticator) StreamServerInterceptor(audience string) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if a.skipServerAuthentication(info.FullMethod) {
			return handler(srv, stream)
		}

		if err := a.authenticate(stream.Context(), audience); errors.Is(err, errPermissiveMissing) {
			a.recordPermissiveMissing(audience)
			return handler(srv, stream)
		} else if err != nil {
			a.recordServerResult(audience, err)
			return status.Error(codes.Unauthenticated, authenticationErrorMessage)
		}

		a.recordServerResult(audience, nil)
		return handler(srv, stream)
	}
}

func (a *Authenticator) skipServerAuthentication(fullMethod string) bool {
	return !a.Enabled() || strings.HasPrefix(fullMethod, healthMethodPrefix)
}

func (a *Authenticator) authenticate(ctx context.Context, audience string) error {
	if err := validateAudience(audience); err != nil {
		return &authError{reason: reasonInvalidAudience, cause: err}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return a.handleMissingCredentials()
	}

	values := md.Get(authorizationMetadataKey)
	if len(values) == 0 {
		return a.handleMissingCredentials()
	}

	if len(values) != 1 {
		return &authError{reason: reasonMalformed, cause: errors.New("multiple authorization metadata values")}
	}

	credential := values[0]
	if len(credential) > maxCredentialLength {
		return &authError{reason: reasonMalformed, cause: errors.New("authorization metadata is too large")}
	}

	parts := strings.Fields(credential)
	if len(parts) != 2 || !strings.EqualFold(parts[0], bearerScheme) || parts[1] == "" {
		return &authError{reason: reasonMalformed, cause: errors.New("malformed bearer credential")}
	}

	return a.verify(parts[1], audience)
}

func (a *Authenticator) handleMissingCredentials() error {
	if a.Mode() == ModePermissive {
		return errPermissiveMissing
	}

	return &authError{reason: reasonMissing, cause: errors.New("authorization metadata is missing")}
}

func (a *Authenticator) recordPermissiveMissing(audience string) {
	serverRequests.WithLabelValues(audience, string(a.Mode()), "allowed", reasonMissing).Inc()
}

func (a *Authenticator) verify(rawToken, audience string) error {
	parser := jwtlib.NewParser(jwtlib.WithValidMethods([]string{jwtlib.SigningMethodHS256.Alg()}), jwtlib.WithoutClaimsValidation())
	unverifiedClaims := new(Claims)
	unverifiedToken, _, err := parser.ParseUnverified(rawToken, unverifiedClaims)
	if err != nil {
		return &authError{reason: reasonMalformed, cause: fmt.Errorf("parse jwt header: %w", err)}
	}

	if unverifiedToken.Method != jwtlib.SigningMethodHS256 {
		return &authError{reason: reasonUnsupportedAlg, cause: errors.New("jwt algorithm is not HS256")}
	}

	tokenType, ok := unverifiedToken.Header["typ"].(string)
	if !ok || tokenType != TokenType {
		return &authError{reason: reasonInvalidType, cause: errors.New("jwt type is invalid")}
	}

	keyID, ok := unverifiedToken.Header["kid"].(string)
	if !ok || keyID == "" {
		return &authError{reason: reasonUnknownKeyID, cause: errors.New("jwt key id is missing")}
	}

	key, ok := a.keys[keyID]
	if !ok {
		return &authError{reason: reasonUnknownKeyID, cause: errors.New("jwt key id is not trusted")}
	}

	claims := new(Claims)
	token, err := parser.ParseWithClaims(rawToken, claims, func(token *jwtlib.Token) (any, error) {
		if token.Method != jwtlib.SigningMethodHS256 {
			return nil, errors.New("jwt algorithm is not HS256")
		}

		return key, nil
	})
	if err != nil || !token.Valid {
		return &authError{reason: reasonInvalidSignature, cause: errors.New("jwt signature is invalid")}
	}

	if claims.Issuer != a.config.JWT.Issuer {
		return &authError{reason: reasonInvalidIssuer, cause: errors.New("jwt issuer is invalid")}
	}

	if claims.Audience != audience {
		return &authError{reason: reasonInvalidAudience, cause: errors.New("jwt audience is invalid")}
	}

	if claims.IssuedAt <= 0 {
		return &authError{reason: reasonInvalidIssuedAt, cause: errors.New("jwt issued-at time is missing")}
	}

	if claims.ExpiresAt <= 0 {
		return &authError{reason: reasonExpired, cause: errors.New("jwt expiration time is missing")}
	}

	now := a.now().Unix()
	skew := int64(a.config.JWT.ClockSkew / time.Second)
	if claims.IssuedAt > now+skew {
		return &authError{reason: reasonInvalidIssuedAt, cause: errors.New("jwt issued-at time is in the future")}
	}

	if claims.ExpiresAt <= now-skew {
		return &authError{reason: reasonExpired, cause: errors.New("jwt is expired")}
	}

	if claims.ExpiresAt <= claims.IssuedAt {
		return &authError{reason: reasonTTLExceeded, cause: errors.New("jwt lifetime is not positive")}
	}

	maxLifetimeSeconds := int64(a.config.JWT.MaxTokenTTL / time.Second)
	if claims.ExpiresAt-claims.IssuedAt > maxLifetimeSeconds {
		return &authError{reason: reasonTTLExceeded, cause: errors.New("jwt lifetime exceeds the configured maximum")}
	}

	return nil
}

func (a *Authenticator) recordServerResult(audience string, err error) {
	if err == nil {
		serverRequests.WithLabelValues(audience, string(a.Mode()), "success", reasonNone).Inc()
		return
	}

	reason := reasonMalformed
	var authenticationError *authError
	if errors.As(err, &authenticationError) {
		reason = authenticationError.reason
	}

	serverRequests.WithLabelValues(audience, string(a.Mode()), "failure", reason).Inc()
}
