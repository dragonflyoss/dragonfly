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
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"

	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, conn *grpc.ClientConn)
	}{
		{
			name: "health check bypasses exhausted rate limiter",
			expect: func(t *testing.T, conn *grpc.ClientConn) {
				assert := assert.New(t)
				resp, err := healthpb.NewHealthClient(conn).Check(context.Background(), &healthpb.HealthCheckRequest{})
				assert.NoError(err)
				assert.Equal(healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
			},
		},
		{
			name: "unary request is rejected by exhausted rate limiter",
			expect: func(t *testing.T, conn *grpc.ClientConn) {
				assert := assert.New(t)
				_, err := schedulerv2.NewSchedulerClient(conn).AnnounceHost(context.Background(), &schedulerv2.AnnounceHostRequest{})
				assert.Equal(codes.ResourceExhausted, status.Code(err))
			},
		},
		{
			name: "stream request is rejected by exhausted rate limiter",
			expect: func(t *testing.T, conn *grpc.ClientConn) {
				assert := assert.New(t)
				stream, err := schedulerv2.NewSchedulerClient(conn).AnnouncePeer(context.Background())
				assert.NoError(err)

				_, err = stream.Recv()
				assert.Equal(codes.ResourceExhausted, status.Code(err))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lis, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}

			svr := New(&schedulerv1.UnimplementedSchedulerServer{}, &schedulerv2.UnimplementedSchedulerServer{}, 1)
			go func() {
				if err := svr.Serve(lis); err != nil {
					logger.Errorf("failed to serve the scheduler: %v", err)
				}
			}()
			t.Cleanup(svr.Stop)

			conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { conn.Close() })

			if _, err := schedulerv2.NewSchedulerClient(conn).AnnounceHost(context.Background(), &schedulerv2.AnnounceHostRequest{}); status.Code(err) == codes.ResourceExhausted {
				t.Fatal(err)
			}

			tc.expect(t, conn)
		})
	}
}
