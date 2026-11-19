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

package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// fakeJobRateLimiter is a job rate limiter that allows or rejects every request
// and records the scheduler cluster ids that it is asked about.
type fakeJobRateLimiter struct {
	allowed    bool
	clusterIDs []uint
}

func (f *fakeJobRateLimiter) AllowByClusterID(_ context.Context, clusterID uint) bool {
	f.clusterIDs = append(f.clusterIDs, clusterID)
	return f.allowed
}

func (f *fakeJobRateLimiter) AllowByClusterIDs(_ context.Context, clusterIDs []uint) bool {
	f.clusterIDs = append(f.clusterIDs, clusterIDs...)
	return f.allowed
}

func (f *fakeJobRateLimiter) Serve() {}

func (f *fakeJobRateLimiter) Stop() {}

func TestCreateJobRateLimiter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		allowed bool
		body    string
		expect  func(t *testing.T, w *httptest.ResponseRecorder, limiter *fakeJobRateLimiter, handled bool)
	}{
		{
			name:    "allowed request reaches the handler",
			allowed: true,
			body:    `{"type":"preheat","scheduler_cluster_ids":[1,2]}`,
			expect: func(t *testing.T, w *httptest.ResponseRecorder, limiter *fakeJobRateLimiter, handled bool) {
				assert := assert.New(t)
				assert.True(handled)
				assert.Equal(http.StatusOK, w.Code)
				assert.Equal([]uint{1, 2}, limiter.clusterIDs)
			},
		},
		{
			name:    "rate limited request is rejected before the handler",
			allowed: false,
			body:    `{"type":"preheat","scheduler_cluster_ids":[1]}`,
			expect: func(t *testing.T, w *httptest.ResponseRecorder, limiter *fakeJobRateLimiter, handled bool) {
				assert := assert.New(t)
				assert.False(handled)
				assert.Equal(http.StatusTooManyRequests, w.Code)
				assert.Equal("rate limit exceeded", w.Body.String())
			},
		},
		{
			name:    "invalid request is rejected once and does not reach the handler",
			allowed: true,
			body:    `{"scheduler_cluster_ids":[1]}`,
			expect: func(t *testing.T, w *httptest.ResponseRecorder, limiter *fakeJobRateLimiter, handled bool) {
				assert := assert.New(t)
				assert.False(handled)
				assert.Equal(http.StatusUnprocessableEntity, w.Code)
				assert.True(json.Valid(w.Body.Bytes()), "body must be a single JSON document, got %q", w.Body.String())
				assert.Empty(limiter.clusterIDs)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			limiter := &fakeJobRateLimiter{allowed: tc.allowed}

			var handled bool
			r := gin.New()
			r.POST("/jobs", CreateJobRateLimiter(limiter), func(c *gin.Context) {
				handled = true
				c.JSON(http.StatusOK, gin.H{"handled": true})
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/jobs", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			tc.expect(t, w, limiter, handled)
		})
	}
}
