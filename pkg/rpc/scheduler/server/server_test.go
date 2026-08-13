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

package server

import (
	"context"
	"encoding/base64"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	grpcauth "d7y.io/dragonfly/v2/pkg/rpc/auth/jwt"
)

func TestServerAuthentication(t *testing.T) {
	authenticator := newTestAuthenticator(t)
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := NewWithAuthentication(
		&schedulerv1.UnimplementedSchedulerServer{},
		&schedulerv2.UnimplementedSchedulerServer{},
		100,
		authenticator,
	)

	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(func() {
		grpcServer.Stop()
		require.NoError(t, listener.Close())
	})

	connection := dialTestServer(t, listener)
	_, err := schedulerv2.NewSchedulerClient(connection).ListHosts(
		context.Background(),
		&schedulerv2.ListHostsRequest{},
	)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))

	response, err := healthpb.NewHealthClient(connection).Check(
		context.Background(),
		&healthpb.HealthCheckRequest{},
	)
	require.NoError(t, err)
	assert.Equal(t, healthpb.HealthCheckResponse_SERVING, response.GetStatus())

	authenticatedConnection := dialTestServer(
		t,
		listener,
		grpc.WithPerRPCCredentials(authenticator.PerRPCCredentials(grpcauth.AudienceScheduler)),
	)
	_, err = schedulerv2.NewSchedulerClient(authenticatedConnection).ListHosts(
		context.Background(),
		&schedulerv2.ListHostsRequest{},
	)
	assert.Equal(t, codes.Unimplemented, status.Code(err))
}

func dialTestServer(t *testing.T, listener *bufconn.Listener, options ...grpc.DialOption) *grpc.ClientConn {
	t.Helper()

	connection, err := grpc.NewClient(
		"passthrough:///bufnet",
		append([]grpc.DialOption{
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return listener.Dial()
			}),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}, options...)...,
	)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, connection.Close()) })
	return connection
}

func newTestAuthenticator(t *testing.T) *grpcauth.Authenticator {
	t.Helper()

	secretFile := filepath.Join(t.TempDir(), "grpc-auth-key")
	secret := base64.StdEncoding.EncodeToString([]byte(strings.Repeat("k", grpcauth.MinimumKeySize)))
	require.NoError(t, os.WriteFile(secretFile, []byte(secret), 0o600))

	config := grpcauth.DefaultConfig()
	config.Mode = grpcauth.ModeRequired
	config.JWT.ActiveKeyID = "test-key"
	config.JWT.Keys = []grpcauth.KeyConfig{{ID: "test-key", SecretFile: secretFile}}
	authenticator, err := grpcauth.New(config)
	require.NoError(t, err)
	return authenticator
}
