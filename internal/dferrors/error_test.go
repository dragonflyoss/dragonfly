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

package dferrors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
)

func TestDfError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    *DfError
		expect func(t *testing.T, msg string)
	}{
		{
			name: "code and message",
			err:  New(commonv1.Code_BadRequest, "bad request"),
			expect: func(t *testing.T, msg string) {
				assert := assert.New(t)
				assert.Equal("[1400]bad request", msg)
			},
		},
		{
			name: "empty message",
			err:  New(commonv1.Code_PeerTaskNotFound, ""),
			expect: func(t *testing.T, msg string) {
				assert := assert.New(t)
				assert.Equal("[1404]", msg)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.err.Error())
		})
	}
}

func TestConvertDfErrorToGRPCError(t *testing.T) {
	plainErr := errors.New("plain error")
	statusErr := status.Error(codes.NotFound, "not found")
	wrappedDfErr := fmt.Errorf("wrapped: %w", New(commonv1.Code_BadRequest, "bad request"))

	tests := []struct {
		name   string
		err    error
		expect func(t *testing.T, err error)
	}{
		{
			name: "nil",
			err:  nil,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "dferror attaches grpc detail",
			err:  New(commonv1.Code_BadRequest, "bad request"),
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				st, ok := status.FromError(err)
				assert.True(ok)
				assert.Equal(codes.Unknown, st.Code())
				if !assert.Len(st.Details(), 1) {
					return
				}

				detail, ok := st.Details()[0].(*commonv1.GrpcDfError)
				assert.True(ok)
				assert.Equal(commonv1.Code_BadRequest, detail.Code)
				assert.Equal("bad request", detail.Message)
			},
		},
		{
			name: "dferror with empty message",
			err:  New(commonv1.Code_SchedNeedBackSource, ""),
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				st, ok := status.FromError(err)
				assert.True(ok)
				assert.Equal(codes.Unknown, st.Code())
				if !assert.Len(st.Details(), 1) {
					return
				}

				detail, ok := st.Details()[0].(*commonv1.GrpcDfError)
				assert.True(ok)
				assert.Equal(commonv1.Code_SchedNeedBackSource, detail.Code)
				assert.Equal("", detail.Message)
			},
		},
		{
			name: "plain error is unchanged",
			err:  plainErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(plainErr, err)
			},
		},
		{
			name: "status error is unchanged",
			err:  statusErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(statusErr, err)
			},
		},
		{
			name: "wrapped dferror is unchanged",
			err:  wrappedDfErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(wrappedDfErr, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, ConvertDfErrorToGRPCError(tc.err))
		})
	}
}

func TestConvertGRPCErrorToDfError(t *testing.T) {
	plainErr := errors.New("plain error")
	statusErr := status.Error(codes.NotFound, "not found")

	tests := []struct {
		name   string
		err    error
		expect func(t *testing.T, err error)
	}{
		{
			name: "nil",
			err:  nil,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "round trip recovers code and message",
			err:  ConvertDfErrorToGRPCError(New(commonv1.Code_BadRequest, "bad request")),
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				var dfErr *DfError
				if !assert.ErrorAs(err, &dfErr) {
					return
				}

				assert.Equal(commonv1.Code_BadRequest, dfErr.Code)
				assert.Equal("bad request", dfErr.Message)
			},
		},
		{
			name: "round trip with empty message",
			err:  ConvertDfErrorToGRPCError(New(commonv1.Code_ClientContextCanceled, "")),
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				var dfErr *DfError
				if !assert.ErrorAs(err, &dfErr) {
					return
				}

				assert.Equal(commonv1.Code_ClientContextCanceled, dfErr.Code)
				assert.Equal("", dfErr.Message)
			},
		},
		{
			name: "status error without details is unchanged",
			err:  statusErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(statusErr, err)
			},
		},
		{
			name: "plain error is unchanged",
			err:  plainErr,
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Same(plainErr, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, ConvertGRPCErrorToDfError(tc.err))
		})
	}
}
