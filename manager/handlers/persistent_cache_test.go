/*
 *     Copyright 2025 The Dragonfly Authors
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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/manager/service/mocks"
	"d7y.io/dragonfly/v2/manager/types"
)

var (
	mockPersistentCacheParams = types.PersistentCacheParams{
		SchedulerClusterID: 1,
		TaskID:             "task-1",
	}
	mockGetPersistentCacheResponse = &types.GetPersistentCacheResponse{
		TaskID: "task-1",
		State:  "SUCCESS",
	}
	mockGetPersistentCachesResponse = []types.GetPersistentCacheResponse{
		{
			TaskID: "task-1",
			State:  "SUCCESS",
		},
	}
)

func mockPersistentCacheRouter(h *Handlers) *gin.Engine {
	r := gin.Default()
	r.DELETE("/persistent-caches/:scheduler_cluster_id/:task_id", h.DestroyPersistentCache)
	r.GET("/persistent-caches/:scheduler_cluster_id/:task_id", h.GetPersistentCache)
	r.GET("/persistent-caches", h.GetPersistentCaches)
	return r
}

func TestHandlers_DestroyPersistentCache(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodDelete, "/persistent-caches/invalid/task-1", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodDelete, "/persistent-caches/1/task-1", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.DestroyPersistentCache(gomock.Any(), gomock.Eq(mockPersistentCacheParams)).Return(nil).Times(1)
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
			mockRouter := mockPersistentCacheRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_GetPersistentCache(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodGet, "/persistent-caches/invalid/task-1", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodGet, "/persistent-caches/1/task-1", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetPersistentCache(gomock.Any(), gomock.Eq(mockPersistentCacheParams)).Return(mockGetPersistentCacheResponse, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Contains(w.Body.String(), `"task_id":"task-1"`)
				assert.Contains(w.Body.String(), `"state":"SUCCESS"`)
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
			mockRouter := mockPersistentCacheRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}

func TestHandlers_GetPersistentCaches(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodGet, "/persistent-caches?page=1&per_page=10", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetPersistentCaches(gomock.Any(), gomock.Any()).Return(mockGetPersistentCachesResponse, int64(1), nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Contains(w.Body.String(), `"task_id":"task-1"`)
				assert.Contains(w.Body.String(), `"state":"SUCCESS"`)
				assert.Contains(w.Header().Get("Link"), "page=1")
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
			mockRouter := mockPersistentCacheRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}
