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

package types

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
)

const testPEM = "-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----"

func TestPEMContent_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		expect func(t *testing.T, p PEMContent, err error)
	}{
		{
			name: "inline pem",
			data: `"-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----"`,
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PEMContent(testPEM), p)
			},
		},
		{
			name: "empty string",
			data: `""`,
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PEMContent(""), p)
			},
		},
		{
			name: "number",
			data: `123`,
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p PEMContent
			err := json.Unmarshal([]byte(tc.data), &p)
			tc.expect(t, p, err)
		})
	}
}

func TestPEMContent_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		expect func(t *testing.T, p PEMContent, err error)
	}{
		{
			name: "block scalar pem",
			data: "|\n  -----BEGIN CERTIFICATE-----\n  MIIB\n  -----END CERTIFICATE-----",
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Contains(string(p), "-----BEGIN CERTIFICATE-----")
			},
		},
		{
			name: "sequence",
			data: "[1, 2]",
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p PEMContent
			err := yaml.Unmarshal([]byte(tc.data), &p)
			tc.expect(t, p, err)
		})
	}
}

func TestPEMContent_LoadFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cert.pem")
	if err := os.WriteFile(path, []byte(testPEM+"\n"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		path   string
		expect func(t *testing.T, p PEMContent, err error)
	}{
		{
			name: "existing file",
			path: path,
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(PEMContent(testPEM), p)
			},
		},
		{
			name: "missing file",
			path: filepath.Join(t.TempDir(), "not-exist.pem"),
			expect: func(t *testing.T, p PEMContent, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p PEMContent
			err := json.Unmarshal([]byte(`"`+tc.path+`"`), &p)
			tc.expect(t, p, err)
		})
	}
}

func TestPEMContent_ToBytes(t *testing.T) {
	tests := []struct {
		name   string
		p      PEMContent
		expect func(t *testing.T, b []byte)
	}{
		{
			name: "trims whitespace",
			p:    PEMContent("  " + testPEM + "\n"),
			expect: func(t *testing.T, b []byte) {
				assert := assert.New(t)
				assert.Equal([]byte(testPEM), b)
			},
		},
		{
			name: "empty content",
			p:    PEMContent(""),
			expect: func(t *testing.T, b []byte) {
				assert := assert.New(t)
				assert.Equal([]byte(""), b)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.p.ToBytes())
		})
	}
}

func TestHostType_Name(t *testing.T) {
	tests := []struct {
		name     string
		hostType HostType
		expect   func(t *testing.T, name string)
	}{
		{
			name:     "normal",
			hostType: HostTypeNormal,
			expect: func(t *testing.T, name string) {
				assert := assert.New(t)
				assert.Equal(HostTypeNormalName, name)
			},
		},
		{
			name:     "super seed",
			hostType: HostTypeSuperSeed,
			expect: func(t *testing.T, name string) {
				assert := assert.New(t)
				assert.Equal(HostTypeSuperSeedName, name)
			},
		},
		{
			name:     "unknown falls back to normal",
			hostType: HostType(100),
			expect: func(t *testing.T, name string) {
				assert := assert.New(t)
				assert.Equal(HostTypeNormalName, name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.hostType.Name())
		})
	}
}

func TestParseHostType(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		expect   func(t *testing.T, hostType HostType)
	}{
		{
			name:     "normal",
			typeName: HostTypeNormalName,
			expect: func(t *testing.T, hostType HostType) {
				assert := assert.New(t)
				assert.Equal(HostTypeNormal, hostType)
			},
		},
		{
			name:     "super seed",
			typeName: HostTypeSuperSeedName,
			expect: func(t *testing.T, hostType HostType) {
				assert := assert.New(t)
				assert.Equal(HostTypeSuperSeed, hostType)
			},
		},
		{
			name:     "unknown falls back to normal",
			typeName: "unknown",
			expect: func(t *testing.T, hostType HostType) {
				assert := assert.New(t)
				assert.Equal(HostTypeNormal, hostType)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, ParseHostType(tc.typeName))
		})
	}
}

func TestTaskTypeV1ToV2(t *testing.T) {
	tests := []struct {
		name   string
		typ    commonv1.TaskType
		expect func(t *testing.T, typ commonv2.TaskType)
	}{
		{
			name: "normal",
			typ:  commonv1.TaskType_Normal,
			expect: func(t *testing.T, typ commonv2.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv2.TaskType_STANDARD, typ)
			},
		},
		{
			name: "dfstore",
			typ:  commonv1.TaskType_DfStore,
			expect: func(t *testing.T, typ commonv2.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv2.TaskType_PERSISTENT, typ)
			},
		},
		{
			name: "dfcache",
			typ:  commonv1.TaskType_DfCache,
			expect: func(t *testing.T, typ commonv2.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv2.TaskType_PERSISTENT_CACHE, typ)
			},
		},
		{
			name: "unknown falls back to standard",
			typ:  commonv1.TaskType(100),
			expect: func(t *testing.T, typ commonv2.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv2.TaskType_STANDARD, typ)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, TaskTypeV1ToV2(tc.typ))
		})
	}
}

func TestTaskTypeV2ToV1(t *testing.T) {
	tests := []struct {
		name   string
		typ    commonv2.TaskType
		expect func(t *testing.T, typ commonv1.TaskType)
	}{
		{
			name: "standard",
			typ:  commonv2.TaskType_STANDARD,
			expect: func(t *testing.T, typ commonv1.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv1.TaskType_Normal, typ)
			},
		},
		{
			name: "persistent",
			typ:  commonv2.TaskType_PERSISTENT,
			expect: func(t *testing.T, typ commonv1.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv1.TaskType_DfStore, typ)
			},
		},
		{
			name: "persistent cache",
			typ:  commonv2.TaskType_PERSISTENT_CACHE,
			expect: func(t *testing.T, typ commonv1.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv1.TaskType_DfCache, typ)
			},
		},
		{
			name: "unknown falls back to normal",
			typ:  commonv2.TaskType(100),
			expect: func(t *testing.T, typ commonv1.TaskType) {
				assert := assert.New(t)
				assert.Equal(commonv1.TaskType_Normal, typ)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, TaskTypeV2ToV1(tc.typ))
		})
	}
}

func TestPriorityV1ToV2(t *testing.T) {
	tests := []struct {
		name     string
		priority commonv1.Priority
		expect   func(t *testing.T, priority commonv2.Priority)
	}{
		{
			name:     "level0",
			priority: commonv1.Priority_LEVEL0,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL0, priority)
			},
		},
		{
			name:     "level1",
			priority: commonv1.Priority_LEVEL1,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL1, priority)
			},
		},
		{
			name:     "level2",
			priority: commonv1.Priority_LEVEL2,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL2, priority)
			},
		},
		{
			name:     "level3",
			priority: commonv1.Priority_LEVEL3,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL3, priority)
			},
		},
		{
			name:     "level4",
			priority: commonv1.Priority_LEVEL4,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL4, priority)
			},
		},
		{
			name:     "level5",
			priority: commonv1.Priority_LEVEL5,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL5, priority)
			},
		},
		{
			name:     "level6",
			priority: commonv1.Priority_LEVEL6,
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL6, priority)
			},
		},
		{
			name:     "unknown falls back to level0",
			priority: commonv1.Priority(100),
			expect: func(t *testing.T, priority commonv2.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv2.Priority_LEVEL0, priority)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, PriorityV1ToV2(tc.priority))
		})
	}
}

func TestPriorityV2ToV1(t *testing.T) {
	tests := []struct {
		name     string
		priority commonv2.Priority
		expect   func(t *testing.T, priority commonv1.Priority)
	}{
		{
			name:     "level0",
			priority: commonv2.Priority_LEVEL0,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL0, priority)
			},
		},
		{
			name:     "level1",
			priority: commonv2.Priority_LEVEL1,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL1, priority)
			},
		},
		{
			name:     "level2",
			priority: commonv2.Priority_LEVEL2,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL2, priority)
			},
		},
		{
			name:     "level3",
			priority: commonv2.Priority_LEVEL3,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL3, priority)
			},
		},
		{
			name:     "level4",
			priority: commonv2.Priority_LEVEL4,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL4, priority)
			},
		},
		{
			name:     "level5",
			priority: commonv2.Priority_LEVEL5,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL5, priority)
			},
		},
		{
			name:     "level6",
			priority: commonv2.Priority_LEVEL6,
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL6, priority)
			},
		},
		{
			name:     "unknown falls back to level0",
			priority: commonv2.Priority(100),
			expect: func(t *testing.T, priority commonv1.Priority) {
				assert := assert.New(t)
				assert.Equal(commonv1.Priority_LEVEL0, priority)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, PriorityV2ToV1(tc.priority))
		})
	}
}

func TestSizeScopeV2ToV1(t *testing.T) {
	tests := []struct {
		name      string
		sizeScope commonv2.SizeScope
		expect    func(t *testing.T, sizeScope commonv1.SizeScope)
	}{
		{
			name:      "normal",
			sizeScope: commonv2.SizeScope_NORMAL,
			expect: func(t *testing.T, sizeScope commonv1.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv1.SizeScope_NORMAL, sizeScope)
			},
		},
		{
			name:      "small",
			sizeScope: commonv2.SizeScope_SMALL,
			expect: func(t *testing.T, sizeScope commonv1.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv1.SizeScope_SMALL, sizeScope)
			},
		},
		{
			name:      "tiny",
			sizeScope: commonv2.SizeScope_TINY,
			expect: func(t *testing.T, sizeScope commonv1.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv1.SizeScope_TINY, sizeScope)
			},
		},
		{
			name:      "empty",
			sizeScope: commonv2.SizeScope_EMPTY,
			expect: func(t *testing.T, sizeScope commonv1.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv1.SizeScope_EMPTY, sizeScope)
			},
		},
		{
			name:      "unknow",
			sizeScope: commonv2.SizeScope_UNKNOW,
			expect: func(t *testing.T, sizeScope commonv1.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv1.SizeScope_UNKNOW, sizeScope)
			},
		},
		{
			name:      "unknown falls back to unknow",
			sizeScope: commonv2.SizeScope(100),
			expect: func(t *testing.T, sizeScope commonv1.SizeScope) {
				assert := assert.New(t)
				assert.Equal(commonv1.SizeScope_UNKNOW, sizeScope)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, SizeScopeV2ToV1(tc.sizeScope))
		})
	}
}
