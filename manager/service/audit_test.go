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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
)

func TestService_AsyncCreateAudit(t *testing.T) {
	s := newTestService(t)
	tests := []struct {
		name   string
		ctx    func() context.Context
		expect func(t *testing.T, err error)
	}{
		{
			name: "cancelled context",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, context.Canceled)
			},
		},
		{
			name: "enqueues audit",
			ctx: func() context.Context {
				return context.Background()
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, s.AsyncCreateAudit(tc.ctx(), &types.CreateAuditRequest{
				ActorType:  models.ActorTypeUser,
				ActorName:  "foo",
				EventType:  models.EventTypeAPI,
				Operation:  http.MethodPost,
				OperatedAt: time.Now(),
				State:      models.AuditStateSuccess,
				Path:       "/api/v1/jobs",
				StatusCode: http.StatusOK,
			}))
		})
	}
}

func TestService_GetAudits(t *testing.T) {
	s := newTestService(t)
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	for i, audit := range []models.Audit{
		{ActorType: models.ActorTypeUser, ActorName: "foo", EventType: models.EventTypeAPI, Operation: http.MethodPost, State: models.AuditStateSuccess, Path: "/api/v1/jobs", StatusCode: http.StatusOK},
		{ActorType: models.ActorTypeUser, ActorName: "foo", EventType: models.EventTypeAPI, Operation: http.MethodDelete, State: models.AuditStateFailure, Path: "/api/v1/jobs/1", StatusCode: http.StatusNotFound},
		{ActorType: models.ActorTypePat, ActorName: "ci", EventType: models.EventTypeAPI, Operation: http.MethodGet, State: models.AuditStateSuccess, Path: "/api/v1/jobs", StatusCode: http.StatusOK},
	} {
		audit.CreatedAt = now.Add(time.Duration(i) * time.Minute)
		if err := s.db.Create(&audit).Error; err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		query  types.GetAuditsQuery
		expect func(t *testing.T, audits []models.Audit, count int64, err error)
	}{
		{
			name:  "newest first with pagination",
			query: types.GetAuditsQuery{Page: 1, PerPage: 2},
			expect: func(t *testing.T, audits []models.Audit, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(audits, 2)
				assert.Equal(int64(3), count)
				assert.Equal("ci", audits[0].ActorName)
				assert.Equal(http.MethodDelete, audits[1].Operation)
			},
		},
		{
			name:  "last page",
			query: types.GetAuditsQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, audits []models.Audit, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(audits, 1)
				assert.Equal(http.MethodPost, audits[0].Operation)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by actor type",
			query: types.GetAuditsQuery{ActorType: models.ActorTypePat, Page: 1, PerPage: 10},
			expect: func(t *testing.T, audits []models.Audit, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(audits, 1)
				assert.Equal("ci", audits[0].ActorName)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by state and status code",
			query: types.GetAuditsQuery{State: models.AuditStateFailure, StatusCode: http.StatusNotFound, Page: 1, PerPage: 10},
			expect: func(t *testing.T, audits []models.Audit, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(audits, 1)
				assert.Equal("/api/v1/jobs/1", audits[0].Path)
			},
		},
		{
			name:  "filter by path and actor name",
			query: types.GetAuditsQuery{Path: "/api/v1/jobs", ActorName: "foo", Page: 1, PerPage: 10},
			expect: func(t *testing.T, audits []models.Audit, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(audits, 1)
				assert.Equal(http.MethodPost, audits[0].Operation)
			},
		},
		{
			name:  "no match",
			query: types.GetAuditsQuery{Operation: http.MethodPatch, Page: 1, PerPage: 10},
			expect: func(t *testing.T, audits []models.Audit, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(audits)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			audits, count, err := s.GetAudits(context.Background(), tc.query)
			tc.expect(t, audits, count, err)
		})
	}
}
