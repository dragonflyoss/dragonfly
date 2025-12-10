package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UnaryServerJWTInterceptor returns a unary server interceptor that validates JWT in metadata.
func UnaryServerJWTInterceptor(key string, expectedAudience string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := validateFromMetadata(ctx, key, expectedAudience); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamServerJWTInterceptor returns a stream server interceptor that validates JWT in metadata.
func StreamServerJWTInterceptor(key string, expectedAudience string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := validateFromMetadata(ss.Context(), key, expectedAudience); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func validateFromMetadata(ctx context.Context, key string, expectedAudience string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization")
	}
	parts := strings.Fields(vals[0])
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return status.Error(codes.Unauthenticated, "invalid authorization header")
	}
	token := parts[1]
	if _, err := ValidateHS256(key, token, expectedAudience); err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	return nil
}
