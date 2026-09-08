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

package service

import (
	"context"
	"fmt"
	"testing"

	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func TestService_GetV1Preheat(t *testing.T) {
	s := mockService(t)
	job := models.Job{TaskID: "group-1", Type: internaljob.PreheatJob, State: machineryv1tasks.StateStarted, Args: models.JSONMap{}}
	if err := s.db.Create(&job).Error; err != nil {
		t.Fatal(err)
	}

	if err := s.db.First(&job, job.ID).Error; err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		id     string
		expect func(t *testing.T, resp *types.GetV1PreheatResponse, err error)
	}{
		{
			name: "id is not a number",
			id:   "abc",
			expect: func(t *testing.T, resp *types.GetV1PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Equal(codes.Unknown, status.Code(err))
				assert.Contains(status.Convert(err).Message(), "invalid syntax")
				assert.Nil(resp)
			},
		},
		{
			name: "job not found",
			id:   "99",
			expect: func(t *testing.T, resp *types.GetV1PreheatResponse, err error) {
				assert := assert.New(t)
				assert.Equal(codes.Unknown, status.Code(err))
				assert.Equal("record not found", status.Convert(err).Message())
				assert.Nil(resp)
			},
		},
		{
			name: "running job",
			id:   fmt.Sprint(job.ID),
			expect: func(t *testing.T, resp *types.GetV1PreheatResponse, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(fmt.Sprint(job.ID), resp.ID)
				assert.Equal(V1PreheatingStateRunning, resp.Status)
				assert.Equal(job.CreatedAt.String(), resp.StartTime)
				assert.Equal(job.UpdatedAt.String(), resp.FinishTime)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := s.GetV1Preheat(context.Background(), tc.id)
			tc.expect(t, resp, err)
		})
	}
}

func TestConvertState(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		expect func(t *testing.T, state string)
	}{
		{
			name:  "pending",
			state: machineryv1tasks.StatePending,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStatePending, state)
			},
		},
		{
			name:  "received",
			state: machineryv1tasks.StateReceived,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStatePending, state)
			},
		},
		{
			name:  "retry",
			state: machineryv1tasks.StateRetry,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStatePending, state)
			},
		},
		{
			name:  "started",
			state: machineryv1tasks.StateStarted,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStateRunning, state)
			},
		},
		{
			name:  "success",
			state: machineryv1tasks.StateSuccess,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStateSuccess, state)
			},
		},
		{
			name:  "failure",
			state: machineryv1tasks.StateFailure,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStateFail, state)
			},
		},
		{
			name:  "unknown",
			state: "UNKNOWN",
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(V1PreheatingStateFail, state)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, convertState(tc.state))
		})
	}
}
