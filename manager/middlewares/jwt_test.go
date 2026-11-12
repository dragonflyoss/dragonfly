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
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/service/mocks"
	"d7y.io/dragonfly/v2/manager/types"
)

var (
	mockJWTConfig = config.JWTConfig{
		Realm:      "Dragonfly",
		Key:        "foo",
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
	}

	mockUser = &models.User{
		BaseModel: models.BaseModel{ID: 1},
		Name:      "foo",
	}
)

func TestJwt_LoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		set    func(c *gin.Context)
		body   string
		mock   func(m *mocks.MockServiceMockRecorder)
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "valid credentials return token and set http only cookie",
			set:  func(c *gin.Context) {},
			body: `{"name":"foo","password":"foobarbaz"}`,
			mock: func(m *mocks.MockServiceMockRecorder) {
				m.SignIn(gomock.Any(), types.SignInRequest{Name: "foo", Password: "foobarbaz"}).Return(mockUser, nil).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)

				var resp struct {
					Token  string `json:"token"`
					Expire string `json:"expire"`
				}
				assert.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
				assert.NotEmpty(resp.Token)
				_, err := time.Parse(time.RFC3339, resp.Expire)
				assert.NoError(err)

				cookies := w.Result().Cookies()
				assert.Len(cookies, 1)
				assert.Equal("jwt", cookies[0].Name)
				assert.Equal(resp.Token, cookies[0].Value)
				assert.True(cookies[0].HttpOnly)
			},
		},
		{
			name: "invalid request body",
			set:  func(c *gin.Context) {},
			body: `{"name":"foo"}`,
			mock: func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
				assert.Empty(w.Result().Cookies())
			},
		},
		{
			name: "failed authentication",
			set:  func(c *gin.Context) {},
			body: `{"name":"foo","password":"foobarbaz"}`,
			mock: func(m *mocks.MockServiceMockRecorder) {
				m.SignIn(gomock.Any(), types.SignInRequest{Name: "foo", Password: "foobarbaz"}).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
				assert.Empty(w.Result().Cookies())
			},
		},
		{
			name: "oauth user redirects and sets http only cookie",
			set: func(c *gin.Context) {
				c.Set("user", mockUser)
			},
			mock: func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusFound, w.Code)
				assert.Equal("/", w.Header().Get(headers.Location))

				cookies := w.Result().Cookies()
				assert.Len(cookies, 1)
				assert.Equal("jwt", cookies[0].Name)
				assert.NotEmpty(cookies[0].Value)
				assert.True(cookies[0].HttpOnly)
			},
		},
		{
			name: "oauth user of unexpected type",
			set: func(c *gin.Context) {
				c.Set("user", "foo")
			},
			mock: func(m *mocks.MockServiceMockRecorder) {},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
				assert.Empty(w.Result().Cookies())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			svc := mocks.NewMockService(ctl)
			tc.mock(svc.EXPECT())

			authMiddleware, err := Jwt(mockJWTConfig, svc)
			if err != nil {
				t.Fatal(err)
			}

			r := gin.New()
			r.Use(func(c *gin.Context) {
				tc.set(c)
				c.Next()
			})
			r.POST("/api/v1/users/signin", authMiddleware.LoginHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/signin", strings.NewReader(tc.body))
			r.ServeHTTP(w, req)

			tc.expect(t, w)
		})
	}
}

func TestJwt_MiddlewareFunc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		request func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request
		expect  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "token from cookie",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				token, _, err := authMiddleware.TokenGenerator(mockUser)
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
				req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
				return req
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.JSONEq(`{"id":1}`, w.Body.String())
			},
		},
		{
			name: "token from authorization header",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				token, _, err := authMiddleware.TokenGenerator(mockUser)
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
				req.Header.Set(headers.Authorization, "Bearer "+token)
				return req
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.JSONEq(`{"id":1}`, w.Body.String())
			},
		},
		{
			name: "token from query",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				token, _, err := authMiddleware.TokenGenerator(mockUser)
				if err != nil {
					t.Fatal(err)
				}

				return httptest.NewRequest(http.MethodGet, "/api/v1/users?token="+token, nil)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.JSONEq(`{"id":1}`, w.Body.String())
			},
		},
		{
			name: "missing token",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.Equal("JWT realm=Dragonfly", w.Header().Get(headers.WWWAuthenticate))
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
			},
		},
		{
			name: "token without identity claim",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				token, _, err := authMiddleware.TokenGenerator("foo")
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
				req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
				return req
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unavailable token: require user id"}`, w.Body.String())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			authMiddleware, err := Jwt(mockJWTConfig, mocks.NewMockService(ctl))
			if err != nil {
				t.Fatal(err)
			}

			r := gin.New()
			r.GET("/api/v1/users", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"id": c.MustGet("id")})
			})

			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.request(t, authMiddleware))

			tc.expect(t, w)
		})
	}
}

func TestJwt_LogoutHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name   string
		expect func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "clears http only cookie",
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)
				assert.Empty(w.Body.String())

				cookies := w.Result().Cookies()
				assert.Len(cookies, 1)
				assert.Equal("jwt", cookies[0].Name)
				assert.Empty(cookies[0].Value)
				assert.Equal(-1, cookies[0].MaxAge)
				assert.True(cookies[0].HttpOnly)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			authMiddleware, err := Jwt(mockJWTConfig, mocks.NewMockService(ctl))
			if err != nil {
				t.Fatal(err)
			}

			r := gin.New()
			r.POST("/api/v1/users/signout", authMiddleware.LogoutHandler)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users/signout", nil)
			r.ServeHTTP(w, req)

			tc.expect(t, w)
		})
	}
}

func TestJwt_RefreshHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name    string
		request func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request
		expect  func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name: "valid token returns new token and sets http only cookie",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				token, _, err := authMiddleware.TokenGenerator(mockUser)
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(http.MethodPost, "/api/v1/users/refresh_token", nil)
				req.AddCookie(&http.Cookie{Name: "jwt", Value: token})
				return req
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusOK, w.Code)

				var resp struct {
					Token  string `json:"token"`
					Expire string `json:"expire"`
				}
				assert.NoError(json.Unmarshal(w.Body.Bytes(), &resp))
				assert.NotEmpty(resp.Token)
				_, err := time.Parse(time.RFC3339, resp.Expire)
				assert.NoError(err)

				cookies := w.Result().Cookies()
				assert.Len(cookies, 1)
				assert.Equal("jwt", cookies[0].Name)
				assert.Equal(resp.Token, cookies[0].Value)
				assert.True(cookies[0].HttpOnly)
			},
		},
		{
			name: "missing token",
			request: func(t *testing.T, authMiddleware *jwt.GinJWTMiddleware) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/api/v1/users/refresh_token", nil)
			},
			expect: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert := assert.New(t)
				assert.Equal(http.StatusUnauthorized, w.Code)
				assert.JSONEq(`{"message":"Unauthorized"}`, w.Body.String())
				assert.Empty(w.Result().Cookies())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			authMiddleware, err := Jwt(mockJWTConfig, mocks.NewMockService(ctl))
			if err != nil {
				t.Fatal(err)
			}

			r := gin.New()
			r.POST("/api/v1/users/refresh_token", authMiddleware.RefreshHandler)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.request(t, authMiddleware))

			tc.expect(t, w)
		})
	}
}
