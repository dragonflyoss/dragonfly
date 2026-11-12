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

package job

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mockPlatform = "linux/amd64"
)

func TestPreheat_CreatePreheatRequestsByManifestURL(t *testing.T) {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		t.Skip("DOCKER_USERNAME or DOCKER_PASSWORD is not set, skipping test")
	}

	tests := []struct {
		name   string
		req    *ManifestRequest
		expect func(t *testing.T, layers []*PreheatRequest, err error)
	}{
		{
			name: "get image layers with manifest url",
			req: &ManifestRequest{
				URL:                "https://registry-1.docker.io/v2/dragonflyoss/busybox/manifests/1.35.0",
				Username:           username,
				Password:           password,
				Timeout:            30 * time.Second,
				InsecureSkipVerify: true,
			},
			expect: func(t *testing.T, layers []*PreheatRequest, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(layers, 1)
				assert.Len(layers[0].URLs, 2)
				assert.Equal("registry-1.docker.io/dragonflyoss/busybox:1.35.0", layers[0].Image.GetReference())
				assert.Nil(layers[0].Image.Platform)
			},
		},
		{
			name: "get image layers with multi arch image layers",
			req: &ManifestRequest{
				URL:                "https://registry-1.docker.io/v2/dragonflyoss/scheduler/manifests/v2.1.0",
				Platform:           &mockPlatform,
				Username:           username,
				Password:           password,
				Timeout:            30 * time.Second,
				InsecureSkipVerify: true,
			},
			expect: func(t *testing.T, layers []*PreheatRequest, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(layers, 1)
				assert.Len(layers[0].URLs, 5)
				assert.Equal("registry-1.docker.io/dragonflyoss/scheduler:v2.1.0", layers[0].Image.GetReference())
				assert.Equal(mockPlatform, layers[0].Image.GetPlatform())
			},
		},
		{
			name: "get image layers of all platforms",
			req: &ManifestRequest{
				URL:                "https://registry-1.docker.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Username:           username,
				Password:           password,
				Timeout:            30 * time.Second,
				InsecureSkipVerify: true,
			},
			expect: func(t *testing.T, layers []*PreheatRequest, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(layers, 3)
				assert.Len(layers[0].URLs, 9)
				assert.Equal("linux/amd64", layers[0].Image.GetPlatform())
				assert.Len(layers[1].URLs, 9)
				assert.Equal("linux/arm64", layers[1].Image.GetPlatform())
				assert.Len(layers[2].URLs, 3)
				assert.Equal("unknown/unknown", layers[2].Image.GetPlatform())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			layers, err := NewImage().CreatePreheatRequestsByManifestURL(context.Background(), tc.req)
			tc.expect(t, layers, err)
		})
	}
}
