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
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"
)

func TestNew(t *testing.T) {
	// Rate limit of 1 request per second with a burst of 1 makes the limiter
	// reject every business request after the first one within the same second.
	const requestRateLimit = 1

	newBufconnClient := func(t *testing.T) *grpc.ClientConn {
		listener := bufconn.Listen(1024 * 1024)
		grpcServer := New(&schedulerv1.UnimplementedSchedulerServer{}, &schedulerv2.UnimplementedSchedulerServer{}, requestRateLimit)
		go func() { _ = grpcServer.Serve(listener) }()
		t.Cleanup(grpcServer.Stop)

		conn, err := grpc.NewClient(
			"passthrough:///bufconn",
			grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			t.Fatalf("dial bufconn: %v", err)
		}
		t.Cleanup(func() { _ = conn.Close() })

		return conn
	}

	t.Run("health check is never rate limited", func(t *testing.T) {
		assert := assert.New(t)
		conn := newBufconnClient(t)
		healthClient := healthpb.NewHealthClient(conn)
		schedulerClient := schedulerv2.NewSchedulerClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Exhaust the limiter with a business request first, then verify that
		// health checks still succeed while business requests are rejected.
		_, _ = schedulerClient.AnnounceHost(ctx, &schedulerv2.AnnounceHostRequest{})

		for range 20 {
			resp, err := healthClient.Check(ctx, &healthpb.HealthCheckRequest{})
			if !assert.NoError(err) {
				return
			}
			assert.Equal(healthpb.HealthCheckResponse_SERVING, resp.GetStatus())
		}

		_, err := schedulerClient.AnnounceHost(ctx, &schedulerv2.AnnounceHostRequest{})
		assert.Equal(codes.ResourceExhausted, status.Code(err))
		assert.Contains(status.Convert(err).Message(), "/scheduler.v2.Scheduler/AnnounceHost is rejected by grpc_ratelimit middleware")
	})

	t.Run("business unary request is rate limited", func(t *testing.T) {
		assert := assert.New(t)
		conn := newBufconnClient(t)
		schedulerClient := schedulerv2.NewSchedulerClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// The first request consumes the single burst token and passes the limiter
		// (it is then rejected downstream by the validator), every following
		// request is rejected by the limiter itself.
		_, err := schedulerClient.AnnounceHost(ctx, &schedulerv2.AnnounceHostRequest{})
		assert.NotEqual(codes.ResourceExhausted, status.Code(err))

		_, err = schedulerClient.AnnounceHost(ctx, &schedulerv2.AnnounceHostRequest{})
		assert.Equal(codes.ResourceExhausted, status.Code(err))
	})

	t.Run("business stream request is rate limited", func(t *testing.T) {
		assert := assert.New(t)
		conn := newBufconnClient(t)
		schedulerClient := schedulerv2.NewSchedulerClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := schedulerClient.AnnounceHost(ctx, &schedulerv2.AnnounceHostRequest{})
		assert.NotEqual(codes.ResourceExhausted, status.Code(err))

		stream, err := schedulerClient.AnnouncePeer(ctx)
		if !assert.NoError(err) {
			return
		}

		_, err = stream.Recv()
		assert.Equal(codes.ResourceExhausted, status.Code(err))
		assert.Contains(status.Convert(err).Message(), "/scheduler.v2.Scheduler/AnnouncePeer is rejected by grpc_ratelimit middleware")
	})
}
