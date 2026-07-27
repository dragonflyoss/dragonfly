package auth

import (
	"context"

	"google.golang.org/grpc/credentials"
)

// PerRPCCreds attaches a Bearer JWT to outgoing gRPC calls.
type PerRPCCreds struct {
	token string
	// If needed later, add refresh hooks.
}

// NewPerRPCCreds constructs credentials with a given token value.
func NewPerRPCCreds(token string) credentials.PerRPCCredentials {
	return &PerRPCCreds{token: token}
}

func (c *PerRPCCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + c.token,
	}, nil
}

// RequireTransportSecurity returns false for backward compatibility with existing deployments.
// In production, configure TLS separately via server.TLS config to secure JWT transmission.
func (c *PerRPCCreds) RequireTransportSecurity() bool { return false }
