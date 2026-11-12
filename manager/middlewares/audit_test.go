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

package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/service/mocks"
	"d7y.io/dragonfly/v2/manager/types"
)

func TestAudit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		set    func(c *gin.Context)
		status int
		mock   func(m *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest)
	}{
		{
			name:   "unauthenticated request is recorded as unknown actor",
			set:    func(c *gin.Context) {},
			status: http.StatusOK,
			mock:   func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal(models.ActorTypeUnknown, audit.ActorType)
				assert.Equal("unknown", audit.ActorName)
				assert.Equal(models.EventTypeAPI, audit.EventType)
				assert.Equal(http.MethodPost, audit.Operation)
				assert.Equal("/api/v1/jobs", audit.Path)
				assert.Equal(models.AuditStateSuccess, audit.State)
				assert.Equal(http.StatusOK, audit.StatusCode)
				assert.False(audit.OperatedAt.IsZero())
			},
		},
		{
			name: "personal access token is recorded as pat actor",
			set: func(c *gin.Context) {
				c.Set("pat", &models.PersonalAccessToken{Name: "foo"})
			},
			status: http.StatusOK,
			mock:   func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(models.ActorTypePat, audit.ActorType)
				assert.Equal("foo", audit.ActorName)
				assert.Equal(models.AuditStateSuccess, audit.State)
			},
		},
		{
			name: "pat of unexpected type is ignored",
			set: func(c *gin.Context) {
				c.Set("pat", "foo")
			},
			status: http.StatusOK,
			mock:   func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(models.ActorTypeUnknown, audit.ActorType)
				assert.Equal("unknown", audit.ActorName)
			},
		},
		{
			name: "jwt identity is recorded as user actor",
			set: func(c *gin.Context) {
				c.Set("id", float64(1))
			},
			status: http.StatusOK,
			mock: func(m *mocks.MockServiceMockRecorder) {
				m.GetUser(gomock.Any(), uint(1)).Return(&models.User{Name: "bar"}, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(models.ActorTypeUser, audit.ActorType)
				assert.Equal("bar", audit.ActorName)
			},
		},
		{
			name: "jwt identity takes precedence over pat",
			set: func(c *gin.Context) {
				c.Set("pat", &models.PersonalAccessToken{Name: "foo"})
				c.Set("id", float64(1))
			},
			status: http.StatusOK,
			mock: func(m *mocks.MockServiceMockRecorder) {
				m.GetUser(gomock.Any(), uint(1)).Return(&models.User{Name: "bar"}, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(models.ActorTypeUser, audit.ActorType)
				assert.Equal("bar", audit.ActorName)
			},
		},
		{
			name: "failed user lookup falls back to unknown actor",
			set: func(c *gin.Context) {
				c.Set("id", float64(1))
			},
			status: http.StatusOK,
			mock: func(m *mocks.MockServiceMockRecorder) {
				m.GetUser(gomock.Any(), uint(1)).Return(nil, errors.New("record not found")).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(models.ActorTypeUnknown, audit.ActorType)
				assert.Equal("unknown", audit.ActorName)
			},
		},
		{
			name: "id of unexpected type is ignored",
			set: func(c *gin.Context) {
				c.Set("id", "1")
			},
			status: http.StatusOK,
			mock:   func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(models.ActorTypeUnknown, audit.ActorType)
				assert.Equal("unknown", audit.ActorName)
			},
		},
		{
			name:   "non 2xx response is recorded as failure",
			set:    func(c *gin.Context) {},
			status: http.StatusInternalServerError,
			mock:   func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder, audit *types.CreateAuditRequest) {
				assert := assert.New(t)
				assert.Equal(http.StatusInternalServerError, w.Code)
				assert.Equal(models.AuditStateFailure, audit.State)
				assert.Equal(http.StatusInternalServerError, audit.StatusCode)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)

			var audit *types.CreateAuditRequest
			svc.EXPECT().AsyncCreateAudit(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, req *types.CreateAuditRequest) error {
				audit = req
				return nil
			}).Times(1)
			tc.mock(svc.EXPECT())

			r := gin.New()
			r.Use(func(c *gin.Context) {
				tc.set(c)
				c.Next()
			})
			r.Use(Audit(svc))
			r.POST("/api/v1/jobs", func(c *gin.Context) {
				c.Status(tc.status)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", nil)
			r.ServeHTTP(w, req)

			tc.expect(t, w, audit)
		})
	}
}

func TestHTTPStatusCodeToState(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expect     func(t *testing.T, state string)
	}{
		{
			name:       "199 is below the success range",
			statusCode: 199,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateFailure, state)
			},
		},
		{
			name:       "200 starts the success range",
			statusCode: http.StatusOK,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateSuccess, state)
			},
		},
		{
			name:       "204 is within the success range",
			statusCode: http.StatusNoContent,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateSuccess, state)
			},
		},
		{
			name:       "299 ends the success range",
			statusCode: 299,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateSuccess, state)
			},
		},
		{
			name:       "300 is above the success range",
			statusCode: http.StatusMultipleChoices,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateFailure, state)
			},
		},
		{
			name:       "404 is a failure",
			statusCode: http.StatusNotFound,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateFailure, state)
			},
		},
		{
			name:       "500 is a failure",
			statusCode: http.StatusInternalServerError,
			expect: func(t *testing.T, state string) {
				assert := assert.New(t)
				assert.Equal(models.AuditStateFailure, state)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, httpStatusCodeToState(tc.statusCode))
		})
	}
}
