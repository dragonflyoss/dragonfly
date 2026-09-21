/*
 *     Copyright 2022 The Dragonfly Authors
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

package rpc

import (
	"context"
	"strings"

	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	"d7y.io/dragonfly/v2/internal/dferrors"
)

// healthCheckMethodPrefix is the full method prefix of the standard gRPC health
// checking service (grpc.health.v1.Health). Health requests bypass rate limiting,
// otherwise kubelet probes and client-side health checks fail whenever the server
// sheds business traffic, turning a transient overload into pod restarts.
const healthCheckMethodPrefix = "/grpc.health.v1.Health/"

// RateLimiterInterceptor is the interface for ratelimit interceptor.
type RateLimiterInterceptor struct {
	// limiter is token bucket of ratelimit.
	limiter *rate.Limiter
}

// NewRateLimiterInterceptor returns a RateLimiterInterceptor instance.
func NewRateLimiterInterceptor(qps float64, burst int64) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		limiter: rate.NewLimiter(rate.Limit(qps), int(burst)),
	}
}

// Limit is the predicate which limits the requests.
func (r *RateLimiterInterceptor) Limit() bool {
	return !r.limiter.Allow()
}

func isHealthCheckMethod(fullMethod string) bool {
	return strings.HasPrefix(fullMethod, healthCheckMethodPrefix)
}

// RateLimitUnaryServerInterceptor wraps the upstream grpc_ratelimit unary server interceptor
// and lets gRPC health checking requests bypass it.
func RateLimitUnaryServerInterceptor(limiter grpc_ratelimit.Limiter) grpc.UnaryServerInterceptor {
	rateLimit := grpc_ratelimit.UnaryServerInterceptor(limiter)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if isHealthCheckMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		return rateLimit(ctx, req, info, handler)
	}
}

// RateLimitStreamServerInterceptor wraps the upstream grpc_ratelimit stream server interceptor
// and lets gRPC health checking requests bypass it.
func RateLimitStreamServerInterceptor(limiter grpc_ratelimit.Limiter) grpc.StreamServerInterceptor {
	rateLimit := grpc_ratelimit.StreamServerInterceptor(limiter)
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if isHealthCheckMethod(info.FullMethod) {
			return handler(srv, ss)
		}

		return rateLimit(srv, ss, info, handler)
	}
}

// ConvertErrorUnaryServerInterceptor returns a new unary server interceptor that convert error when trigger custom error.
func ConvertErrorUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	h, err := handler(ctx, req)
	if err != nil {
		return h, dferrors.ConvertDfErrorToGRPCError(err)
	}

	return h, nil
}

// ConvertErrorStreamServerInterceptor returns a new stream server interceptor that convert error when trigger custom error.
func ConvertErrorStreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := handler(srv, ss); err != nil {
		return dferrors.ConvertDfErrorToGRPCError(err)
	}

	return nil
}

// ConvertErrorUnaryClientInterceptor returns a new unary client interceptor that convert error when trigger custom error.
func ConvertErrorUnaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
		return dferrors.ConvertGRPCErrorToDfError(err)
	}

	return nil
}

// ConvertErrorStreamClientInterceptor returns a new stream client interceptor that convert error when trigger custom error.
func ConvertErrorStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, dferrors.ConvertGRPCErrorToDfError(err)
	}

	return s, nil
}
