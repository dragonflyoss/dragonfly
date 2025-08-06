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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPreheat_CreatePreheatRequestsByManifestURL(t *testing.T) {
	tests := []struct {
		name   string
		req    *ManifestRequest
		expect func(t *testing.T, layers []*PreheatRequest)
	}{
		{
			name: "get image layers with manifest url",
			req: &ManifestRequest{
				URL:                "https://registry-1.docker.io/v2/dragonflyoss/busybox/manifests/1.35.0",
				Timeout:            30 * time.Second,
				InsecureSkipVerify: true,
			},
			expect: func(t *testing.T, layers []*PreheatRequest) {
				assert := assert.New(t)
				assert.Equal(2, len(layers[0].URLs))
			},
		},
		{
			name: "get image layers with multi arch image layers",
			req: &ManifestRequest{
				URL:                "https://registry-1.docker.io/v2/dragonflyoss/scheduler/manifests/v2.1.0",
				Platform:           "linux/amd64",
				Timeout:            30 * time.Second,
				InsecureSkipVerify: true,
			},
			expect: func(t *testing.T, layers []*PreheatRequest) {
				assert := assert.New(t)
				assert.Equal(5, len(layers[0].URLs))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			layers, err := NewImage().CreatePreheatRequestsByManifestURL(context.Background(), tc.req)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, layers)
		})
	}
}
