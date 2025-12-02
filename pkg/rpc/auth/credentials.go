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

// For now keep false to avoid breaking existing insecure transports; can be tightened later.
func (c *PerRPCCreds) RequireTransportSecurity() bool { return false }
