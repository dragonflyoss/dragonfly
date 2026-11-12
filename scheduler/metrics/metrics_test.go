/*
 *     Copyright 2020 The Dragonfly Authors
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

package metrics

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"d7y.io/dragonfly/v2/scheduler/config"
)

func TestNew(t *testing.T) {
	cfg := &config.MetricsConfig{
		Addr: "localhost:8080",
	}
	svr := grpc.NewServer()
	server := New(cfg, svr)

	assert := assert.New(t)
	assert.Equal(cfg.Addr, server.Addr)
	assert.IsType(&http.ServeMux{}, server.Handler)
}

func TestCalculateSizeLevel(t *testing.T) {
	tests := []struct {
		name   string
		size   int64
		expect func(t *testing.T, level TaskSizeLevel)
	}{
		{
			name: "negative size is unknown",
			size: -1,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel0, level)
			},
		},
		{
			name: "zero size is unknown",
			size: 0,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel0, level)
			},
		},
		{
			name: "one byte is level 1",
			size: 1,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel1, level)
			},
		},
		{
			name: "just below 1MB is level 1",
			size: Size1MB - 1,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel1, level)
			},
		},
		{
			name: "exactly 1MB is level 2",
			size: Size1MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel2, level)
			},
		},
		{
			name: "exactly 4MB is level 3",
			size: Size4MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel3, level)
			},
		},
		{
			name: "exactly 8MB is level 4",
			size: Size8MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel4, level)
			},
		},
		{
			name: "exactly 16MB is level 5",
			size: Size16MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel5, level)
			},
		},
		{
			name: "exactly 32MB is level 6",
			size: Size32MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel6, level)
			},
		},
		{
			name: "exactly 64MB is level 7",
			size: Size64MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel7, level)
			},
		},
		{
			name: "exactly 128MB is level 8",
			size: Size128MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel8, level)
			},
		},
		{
			name: "exactly 256MB is level 9",
			size: Size256MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel9, level)
			},
		},
		{
			name: "exactly 512MB is level 10",
			size: Size512MB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel10, level)
			},
		},
		{
			name: "exactly 1GB is level 11",
			size: Size1GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel11, level)
			},
		},
		{
			name: "exactly 4GB is level 12",
			size: Size4GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel12, level)
			},
		},
		{
			name: "exactly 8GB is level 13",
			size: Size8GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel13, level)
			},
		},
		{
			name: "exactly 16GB is level 14",
			size: Size16GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel14, level)
			},
		},
		{
			name: "exactly 32GB is level 15",
			size: Size32GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel15, level)
			},
		},
		{
			name: "exactly 64GB is level 16",
			size: Size64GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel16, level)
			},
		},
		{
			name: "exactly 128GB is level 17",
			size: Size128GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel17, level)
			},
		},
		{
			name: "exactly 256GB is level 18",
			size: Size256GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel18, level)
			},
		},
		{
			name: "exactly 512GB is level 19",
			size: Size512GB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel19, level)
			},
		},
		{
			name: "exactly 1TB is level 20",
			size: Size1TB,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel20, level)
			},
		},
		{
			name: "above 1TB is level 20",
			size: Size1TB * 8,
			expect: func(t *testing.T, level TaskSizeLevel) {
				assert := assert.New(t)
				assert.Equal(TaskSizeLevel20, level)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, CalculateSizeLevel(tc.size))
		})
	}
}

func TestTaskSizeLevel_String(t *testing.T) {
	tests := []struct {
		name   string
		level  TaskSizeLevel
		expect func(t *testing.T, s string)
	}{
		{
			name:  "level 0",
			level: TaskSizeLevel0,
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("0", s)
			},
		},
		{
			name:  "level 1",
			level: TaskSizeLevel1,
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("1", s)
			},
		},
		{
			name:  "level 20",
			level: TaskSizeLevel20,
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("20", s)
			},
		},
		{
			name:  "unknown level falls back to 0",
			level: TaskSizeLevel(100),
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("0", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.level.String())
		})
	}
}
