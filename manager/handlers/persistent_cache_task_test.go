/*
 *     Copyright 2024 The Dragonfly Authors
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

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/manager/service/mocks"
	"d7y.io/dragonfly/v2/manager/types"
)

var (
	mockPersistentCacheTaskID                 = "task-id-123"
	mockSchedulerClusterID               uint = 1
	mockUnprocessableSchedulerQuery           = ""
	mockDestroyPersistentCacheQuery           = fmt.Sprintf("scheduler_cluster_id=%d", mockSchedulerClusterID)
	mockGetPersistentCacheQuery               = fmt.Sprintf("scheduler_cluster_id=%d", mockSchedulerClusterID)
	mockGetPersistentCachesQuery              = fmt.Sprintf("scheduler_cluster_id=%d", mockSchedulerClusterID)
	mockGetPersistentCachesQueryWithPage      = fmt.Sprintf("scheduler_cluster_id=%d&page=1&per_page=10", mockSchedulerClusterID)
	mockPersistentCacheTask                   = types.PersistentCacheTask{
		ID:                     mockPersistentCacheTaskID,
		PersistentReplicaCount: 3,
		Tag:                    "v1.0.0",
		Application:            "app1",
		PieceLength:            1024,
		ContentLength:          4096,
		TotalPieceCount:        4,
		State:                  "Success",
		TTL:                    time.Hour * 24,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	mockPersistentCacheTasks = []types.PersistentCacheTask{mockPersistentCacheTask}
)

func mockPersistentCacheTaskRouter(h *Handlers) *gin.Engine {
	r := gin.Default()
	apiv1 := r.Group("/api/v1")
	task := apiv1.Group("/persistent-cache-tasks")
	task.DELETE(":id", h.DestroyPersistentCacheTask)
	task.GET(":id", h.GetPersistentCacheTask)
	task.GET("", h.GetPersistentCacheTasks)
	return r
}

func TestHandlers_DestroyPersistentCacheTask(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity uri",
			req:  httptest.NewRequest(http.MethodDelete, "/api/v1/persistent-cache-tasks/", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusNotFound, w.Code)
			},
		},
		{
			name: "unprocessable entity query",
			req:  httptest.NewRequest(http.MethodDelete, "/api/v1/persistent-cache-tasks/"+mockPersistentCacheTaskID+"?"+mockUnprocessableSchedulerQuery, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodDelete, "/api/v1/persistent-cache-tasks/"+mockPersistentCacheTaskID+"?"+mockDestroyPersistentCacheQuery, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.DestroyPersistentCacheTask(gomock.Any(), gomock.Eq(mockSchedulerClusterID), gomock.Eq(mockPersistentCacheTaskID)).Return(nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockPersistentCacheTaskRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_GetPersistentCacheTask(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity query",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/persistent-cache-tasks/"+mockPersistentCacheTaskID+"?"+mockUnprocessableSchedulerQuery, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/persistent-cache-tasks/"+mockPersistentCacheTaskID+"?"+mockGetPersistentCacheQuery, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetPersistentCacheTask(gomock.Any(), gomock.Eq(mockSchedulerClusterID), gomock.Eq(mockPersistentCacheTaskID)).Return(mockPersistentCacheTask, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				var task types.PersistentCacheTask
				err := json.Unmarshal(w.Body.Bytes(), &task)
				assert.NoError(err)
				assert.Equal(mockPersistentCacheTaskID, task.ID)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockPersistentCacheTaskRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_GetPersistentCacheTasks(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity query",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/persistent-cache-tasks?"+mockUnprocessableSchedulerQuery, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "get without pagination",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/persistent-cache-tasks?"+mockGetPersistentCachesQuery, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetPersistentCacheTasks(gomock.Any(), gomock.Any()).Return(mockPersistentCacheTasks, int64(len(mockPersistentCacheTasks)), nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				var tasks []types.PersistentCacheTask
				err := json.Unmarshal(w.Body.Bytes(), &tasks)
				assert.NoError(err)
				assert.Len(tasks, 1)
				assert.Equal(mockPersistentCacheTaskID, tasks[0].ID)
			},
		},
		{
			name: "get with pagination",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/persistent-cache-tasks?"+mockGetPersistentCachesQueryWithPage, nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetPersistentCacheTasks(gomock.Any(), gomock.Any()).Return(mockPersistentCacheTasks, int64(len(mockPersistentCacheTasks)), nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				var tasks []types.PersistentCacheTask
				err := json.Unmarshal(w.Body.Bytes(), &tasks)
				assert.NoError(err)
				assert.Len(tasks, 1)
				assert.Equal(mockPersistentCacheTaskID, tasks[0].ID)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			w := httptest.NewRecorder()
			h := New(svc)
			mockRouter := mockPersistentCacheTaskRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}
