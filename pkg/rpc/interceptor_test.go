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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"

	"d7y.io/dragonfly/v2/internal/dferrors"
)

type fakeClientStream struct {
	grpc.ClientStream
}

func TestRateLimiterInterceptor_Limit(t *testing.T) {
	assert := assert.New(t)
	limiter := NewRateLimiterInterceptor(100, 2)
	assert.False(limiter.Limit())
	assert.False(limiter.Limit())
	assert.True(limiter.Limit())

	time.Sleep(50 * time.Millisecond)
	assert.False(limiter.Limit())
}

func TestConvertErrorUnaryServerInterceptor(t *testing.T) {
	plainErr := errors.New("plain error")

	tests := []struct {
		name       string
		handlerErr error
		expect     func(t *testing.T, resp any, err error)
	}{
		{
			name: "handler succeeds",
			expect: func(t *testing.T, resp any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("resp", resp)
			},
		},
		{
			name:       "plain error is unchanged",
			handlerErr: plainErr,
			expect: func(t *testing.T, resp any, err error) {
				assert := assert.New(t)
				assert.Same(plainErr, err)
				assert.Equal("resp", resp)
			},
		},
		{
			name:       "dferror is converted to grpc error with details",
			handlerErr: dferrors.New(commonv1.Code_BadRequest, "bad request"),
			expect: func(t *testing.T, resp any, err error) {
				assert := assert.New(t)
				st, ok := status.FromError(err)
				if !assert.True(ok) || !assert.Len(st.Details(), 1) {
					return
				}

				detail, ok := st.Details()[0].(*commonv1.GrpcDfError)
				assert.True(ok)
				assert.Equal(commonv1.Code_BadRequest, detail.Code)
				assert.Equal("bad request", detail.Message)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(ctx context.Context, req any) (any, error) { return "resp", tc.handlerErr }
			resp, err := ConvertErrorUnaryServerInterceptor(context.Background(), "req", &grpc.UnaryServerInfo{}, handler)
			tc.expect(t, resp, err)
		})
	}
}

func TestConvertErrorStreamServerInterceptor(t *testing.T) {
	plainErr := errors.New("plain error")

	tests := []struct {
		name       string
		handlerErr error
		expect     func(t *testing.T, err error)
	}{
		{
			name: "handler succeeds",
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:       "plain error is unchanged",
			handlerErr: plainErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(plainErr, err)
			},
		},
		{
			name:       "dferror is converted to grpc error with details",
			handlerErr: dferrors.New(commonv1.Code_SchedNeedBackSource, "need back source"),
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				st, ok := status.FromError(err)
				if !assert.True(ok) || !assert.Len(st.Details(), 1) {
					return
				}

				detail, ok := st.Details()[0].(*commonv1.GrpcDfError)
				assert.True(ok)
				assert.Equal(commonv1.Code_SchedNeedBackSource, detail.Code)
				assert.Equal("need back source", detail.Message)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(srv any, stream grpc.ServerStream) error { return tc.handlerErr }
			tc.expect(t, ConvertErrorStreamServerInterceptor("srv", nil, &grpc.StreamServerInfo{}, handler))
		})
	}
}

func TestConvertErrorUnaryClientInterceptor(t *testing.T) {
	plainErr := errors.New("plain error")
	statusErr := status.Error(codes.NotFound, "not found")

	tests := []struct {
		name       string
		invokerErr error
		expect     func(t *testing.T, err error)
	}{
		{
			name: "invoker succeeds",
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:       "plain error is unchanged",
			invokerErr: plainErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(plainErr, err)
			},
		},
		{
			name:       "status error without details is unchanged",
			invokerErr: statusErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(statusErr, err)
			},
		},
		{
			name:       "grpc error with details is converted to dferror",
			invokerErr: dferrors.ConvertDfErrorToGRPCError(dferrors.New(commonv1.Code_BadRequest, "bad request")),
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				var dfErr *dferrors.DfError
				if !assert.ErrorAs(err, &dfErr) {
					return
				}

				assert.Equal(commonv1.Code_BadRequest, dfErr.Code)
				assert.Equal("bad request", dfErr.Message)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			invoker := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return tc.invokerErr
			}
			tc.expect(t, ConvertErrorUnaryClientInterceptor(context.Background(), "/foo/Bar", "req", nil, nil, invoker))
		})
	}
}

func TestConvertErrorStreamClientInterceptor(t *testing.T) {
	plainErr := errors.New("plain error")
	stream := &fakeClientStream{}

	tests := []struct {
		name        string
		streamerErr error
		expect      func(t *testing.T, s grpc.ClientStream, err error)
	}{
		{
			name: "streamer succeeds",
			expect: func(t *testing.T, s grpc.ClientStream, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Same(stream, s)
			},
		},
		{
			name:        "plain error is unchanged",
			streamerErr: plainErr,
			expect: func(t *testing.T, s grpc.ClientStream, err error) {
				assert := assert.New(t)
				assert.Same(plainErr, err)
				assert.Nil(s)
			},
		},
		{
			name:        "grpc error with details is converted to dferror",
			streamerErr: dferrors.ConvertDfErrorToGRPCError(dferrors.New(commonv1.Code_ClientContextCanceled, "canceled")),
			expect: func(t *testing.T, s grpc.ClientStream, err error) {
				assert := assert.New(t)
				assert.Nil(s)
				var dfErr *dferrors.DfError
				if !assert.ErrorAs(err, &dfErr) {
					return
				}

				assert.Equal(commonv1.Code_ClientContextCanceled, dfErr.Code)
				assert.Equal("canceled", dfErr.Message)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				if tc.streamerErr != nil {
					return nil, tc.streamerErr
				}

				return stream, nil
			}
			s, err := ConvertErrorStreamClientInterceptor(context.Background(), &grpc.StreamDesc{}, nil, "/foo/Bar", streamer)
			tc.expect(t, s, err)
		})
	}
}
