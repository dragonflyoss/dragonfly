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

package server

import (
	"math"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ratelimit "github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	managerv1 "d7y.io/api/v2/pkg/apis/manager/v1"
	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/rpc"
	"d7y.io/dragonfly/v2/pkg/rpc/auth"
)

const (
	// DefaultMaxConnectionIdle is default max connection idle of grpc keepalive.
	DefaultMaxConnectionIdle = 10 * time.Minute

	// DefaultMaxConnectionAge is default max connection age of grpc keepalive.
	DefaultMaxConnectionAge = 12 * time.Hour

	// DefaultMaxConnectionAgeGrace is default max connection age grace of grpc keepalive.
	DefaultMaxConnectionAgeGrace = 5 * time.Minute
)

// New returns grpc server instance and register service on grpc server.
func New(managerServerV1 managerv1.ManagerServer, managerServerV2 managerv2.ManagerServer, requestRateLimit float64, opts ...grpc.ServerOption) *grpc.Server {
	limiter := rpc.NewRateLimiterInterceptor(requestRateLimit, int64(requestRateLimit))
	jwtKey := auth.GetServerKey("manager")

	grpcServer := grpc.NewServer(append([]grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
			otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
		)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     DefaultMaxConnectionIdle,
			MaxConnectionAge:      DefaultMaxConnectionAge,
			MaxConnectionAgeGrace: DefaultMaxConnectionAgeGrace,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			auth.UnaryServerJWTInterceptor(jwtKey, "manager"),
			grpc_ratelimit.UnaryServerInterceptor(limiter),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger.GrpcLogger.Desugar()),
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			auth.StreamServerJWTInterceptor(jwtKey, "manager"),
			grpc_ratelimit.StreamServerInterceptor(limiter),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger.GrpcLogger.Desugar()),
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
	}, opts...)...)

	// Register servers on v1 version of the grpc server.
	managerv1.RegisterManagerServer(grpcServer, managerServerV1)

	// Register servers on v2 version of the grpc server.
	managerv2.RegisterManagerServer(grpcServer, managerServerV2)

	// Register health on grpc server.
	healthpb.RegisterHealthServer(grpcServer, health.NewServer())

	// Register reflection on grpc server.
	reflection.Register(grpcServer)
	return grpcServer
}
