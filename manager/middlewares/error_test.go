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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VividCortex/mysqlerr"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"

	"d7y.io/dragonfly/v2/internal/dferrors"
)

func TestError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		err    error
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "no error keeps the handler response",
			err:  nil,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusCreated, w.Code)
				assert.Empty(w.Body.String())
			},
		},
		{
			name: "redigo nil maps to not found",
			err:  redigo.ErrNil,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusNotFound, w.Code)
				assert.JSONEq(`{"message":"Not Found"}`, w.Body.String())
			},
		},
		{
			name: "dferror with invalid resource type maps to bad request",
			err:  dferrors.New(commonv1.Code_InvalidResourceType, "foo"),
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusBadRequest, w.Code)
				assert.JSONEq(`{"message":"Bad Request"}`, w.Body.String())
			},
		},
		{
			name: "wrapped dferror with invalid resource type maps to bad request",
			err:  fmt.Errorf("create job: %w", dferrors.New(commonv1.Code_InvalidResourceType, "foo")),
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusBadRequest, w.Code)
				assert.JSONEq(`{"message":"Bad Request"}`, w.Body.String())
			},
		},
		{
			name: "dferror with other code maps to internal server error",
			err:  dferrors.New(commonv1.Code_BadRequest, "foo"),
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusInternalServerError, w.Code)
				assert.JSONEq(`{"message":"Internal Server Error"}`, w.Body.String())
			},
		},
		{
			name: "bcrypt mismatch maps to unauthorized",
			err:  bcrypt.ErrMismatchedHashAndPassword,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
			},
		},
		{
			name: "gorm record not found maps to not found with the error text",
			err:  gorm.ErrRecordNotFound,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusNotFound, w.Code)
				assert.JSONEq(`{"message":"record not found"}`, w.Body.String())
			},
		},
		{
			name: "wrapped gorm record not found keeps the wrapping text",
			err:  fmt.Errorf("find job: %w", gorm.ErrRecordNotFound),
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusNotFound, w.Code)
				assert.JSONEq(`{"message":"find job: record not found"}`, w.Body.String())
			},
		},
		{
			name: "mysql duplicate entry maps to conflict",
			err:  &mysql.MySQLError{Number: mysqlerr.ER_DUP_ENTRY, Message: "Duplicate entry"},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusConflict, w.Code)
				assert.JSONEq(`{"message":"Conflict"}`, w.Body.String())
			},
		},
		{
			name: "other mysql error maps to internal server error",
			err:  &mysql.MySQLError{Number: mysqlerr.ER_ACCESS_DENIED_ERROR, Message: "Access denied"},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusInternalServerError, w.Code)
				assert.JSONEq(`{"message":"Internal Server Error"}`, w.Body.String())
			},
		},
		{
			name: "redis nil maps to not found",
			err:  redis.Nil,
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusNotFound, w.Code)
				assert.JSONEq(`{"message":"Not Found"}`, w.Body.String())
			},
		},
		{
			name: "unknown error maps to internal server error with the error text",
			err:  errors.New("boom"),
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusInternalServerError, w.Code)
				assert.JSONEq(`{"message":"boom"}`, w.Body.String())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(Error())
			r.GET("/api/v1/jobs", func(c *gin.Context) {
				c.Status(http.StatusCreated)
				if tc.err != nil {
					c.Error(tc.err) //nolint: errcheck
				}
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
			r.ServeHTTP(w, req)

			tc.expect(t, w)
		})
	}
}
