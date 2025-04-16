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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/service/mocks"
	"d7y.io/dragonfly/v2/manager/types"
)

var (
	mockAuditModel = &models.Audit{
		BaseModel: models.BaseModel{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		ActorType:  "user",
		ActorName:  "root",
		EventType:  "API",
		Operation:  "GET",
		OperatedAt: time.Time{},
		State:      "success",
		Path:       "/api/v1/audits",
		StatusCode: 200,
	}
)

func mockAuditRouter(h *Handlers) *gin.Engine {
	r := gin.Default()
	apiv1 := r.Group("/api/v1")
	audits := apiv1.Group("/audits")
	audits.GET("", h.GetAudits)
	return r
}

func TestHandlers_GetAudits(t *testing.T) {
	tests := []struct {
		name   string
		req    *http.Request
		mock   func(ms *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "unprocessable entity",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/audits?page=-1", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
			},
		},
		{
			name: "success",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/audits?page=1&per_page=5", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetAudits(gomock.Any(), gomock.Eq(types.GetAuditsQuery{
					Page:    1,
					PerPage: 5,
				})).Return([]models.Audit{*mockAuditModel}, int64(1), nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				var audits []models.Audit
				err := json.Unmarshal(w.Body.Bytes(), &audits)
				assert.NoError(err)
				assert.Len(audits, 1)
				assert.Equal(mockAuditModel, &audits[0])

				// Check pagination header
				linkHeader := w.Header().Get("Link")
				assert.NotEmpty(linkHeader)
				links := strings.Split(linkHeader, ",")
				assert.Len(links, 4) // prev, next, first, last
				assert.Contains(linkHeader, `rel=first`)
				assert.Contains(linkHeader, `rel=last`)
				assert.Contains(linkHeader, `rel=prev`)
				assert.Contains(linkHeader, `rel=next`)
			},
		},
		{
			name: "success with default pagination",
			req:  httptest.NewRequest(http.MethodGet, "/api/v1/audits", nil),
			mock: func(ms *mocks.MockServiceMockRecorder) {
				ms.GetAudits(gomock.Any(), gomock.Eq(types.GetAuditsQuery{
					Page:    1,
					PerPage: 10,
				})).Return([]models.Audit{*mockAuditModel}, int64(1), nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				var audits []models.Audit
				err := json.Unmarshal(w.Body.Bytes(), &audits)
				assert.NoError(err)
				assert.Len(audits, 1)
				assert.Equal(mockAuditModel, &audits[0])
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
			mockRouter := mockAuditRouter(h)

			tc.mock(svc.EXPECT())
			mockRouter.ServeHTTP(w, tc.req)
			tc.expect(t, w)
		})
	}
}
