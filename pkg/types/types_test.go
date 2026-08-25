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
	assert := assert.New(t)

	var p PEMContent
	assert.NoError(json.Unmarshal([]byte(`"`+`-----BEGIN CERTIFICATE-----\nMIIB\n-----END CERTIFICATE-----`+`"`), &p))
	assert.Equal(PEMContent(testPEM), p)

	var empty PEMContent
	assert.NoError(json.Unmarshal([]byte(`""`), &empty))
	assert.Equal(PEMContent(""), empty)

	var invalid PEMContent
	assert.Error(json.Unmarshal([]byte(`123`), &invalid))
}

func TestPEMContent_UnmarshalYAML(t *testing.T) {
	assert := assert.New(t)

	var p PEMContent
	assert.NoError(yaml.Unmarshal([]byte("|\n  -----BEGIN CERTIFICATE-----\n  MIIB\n  -----END CERTIFICATE-----"), &p))
	assert.Contains(string(p), "-----BEGIN CERTIFICATE-----")

	var invalid PEMContent
	assert.Error(yaml.Unmarshal([]byte("[1, 2]"), &invalid))
}

func TestPEMContent_LoadFromFile(t *testing.T) {
	assert := assert.New(t)

	path := filepath.Join(t.TempDir(), "cert.pem")
	assert.NoError(os.WriteFile(path, []byte(testPEM+"\n"), 0600))

	var p PEMContent
	assert.NoError(json.Unmarshal([]byte(`"`+path+`"`), &p))
	assert.Equal(PEMContent(testPEM), p)

	var missing PEMContent
	assert.Error(json.Unmarshal([]byte(`"`+filepath.Join(t.TempDir(), "not-exist.pem")+`"`), &missing))
}

func TestPEMContent_ToBytes(t *testing.T) {
	assert := assert.New(t)
	assert.Equal([]byte(testPEM), PEMContent("  "+testPEM+"\n").ToBytes())
	assert.Equal([]byte(""), PEMContent("").ToBytes())
}

func TestHostType_Name(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(HostTypeNormalName, HostTypeNormal.Name())
	assert.Equal(HostTypeSuperSeedName, HostTypeSuperSeed.Name())
	assert.Equal(HostTypeNormalName, HostType(100).Name())
}

func TestParseHostType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(HostTypeNormal, ParseHostType(HostTypeNormalName))
	assert.Equal(HostTypeSuperSeed, ParseHostType(HostTypeSuperSeedName))
	assert.Equal(HostTypeNormal, ParseHostType("unknown"))
}

func TestTaskTypeV1ToV2(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(commonv2.TaskType_STANDARD, TaskTypeV1ToV2(commonv1.TaskType_Normal))
	assert.Equal(commonv2.TaskType_PERSISTENT, TaskTypeV1ToV2(commonv1.TaskType_DfStore))
	assert.Equal(commonv2.TaskType_PERSISTENT_CACHE, TaskTypeV1ToV2(commonv1.TaskType_DfCache))
	assert.Equal(commonv2.TaskType_STANDARD, TaskTypeV1ToV2(commonv1.TaskType(100)))
}

func TestTaskTypeV2ToV1(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(commonv1.TaskType_Normal, TaskTypeV2ToV1(commonv2.TaskType_STANDARD))
	assert.Equal(commonv1.TaskType_DfStore, TaskTypeV2ToV1(commonv2.TaskType_PERSISTENT))
	assert.Equal(commonv1.TaskType_DfCache, TaskTypeV2ToV1(commonv2.TaskType_PERSISTENT_CACHE))
	assert.Equal(commonv1.TaskType_Normal, TaskTypeV2ToV1(commonv2.TaskType(100)))
}

func TestPriorityV1ToV2(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		v1 commonv1.Priority
		v2 commonv2.Priority
	}{
		{commonv1.Priority_LEVEL0, commonv2.Priority_LEVEL0},
		{commonv1.Priority_LEVEL1, commonv2.Priority_LEVEL1},
		{commonv1.Priority_LEVEL2, commonv2.Priority_LEVEL2},
		{commonv1.Priority_LEVEL3, commonv2.Priority_LEVEL3},
		{commonv1.Priority_LEVEL4, commonv2.Priority_LEVEL4},
		{commonv1.Priority_LEVEL5, commonv2.Priority_LEVEL5},
		{commonv1.Priority_LEVEL6, commonv2.Priority_LEVEL6},
	}
	for _, tt := range tests {
		assert.Equal(tt.v2, PriorityV1ToV2(tt.v1))
	}
	assert.Equal(commonv2.Priority_LEVEL0, PriorityV1ToV2(commonv1.Priority(100)))
}

func TestPriorityV2ToV1(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		v2 commonv2.Priority
		v1 commonv1.Priority
	}{
		{commonv2.Priority_LEVEL0, commonv1.Priority_LEVEL0},
		{commonv2.Priority_LEVEL1, commonv1.Priority_LEVEL1},
		{commonv2.Priority_LEVEL2, commonv1.Priority_LEVEL2},
		{commonv2.Priority_LEVEL3, commonv1.Priority_LEVEL3},
		{commonv2.Priority_LEVEL4, commonv1.Priority_LEVEL4},
		{commonv2.Priority_LEVEL5, commonv1.Priority_LEVEL5},
		{commonv2.Priority_LEVEL6, commonv1.Priority_LEVEL6},
	}
	for _, tt := range tests {
		assert.Equal(tt.v1, PriorityV2ToV1(tt.v2))
	}
	assert.Equal(commonv1.Priority_LEVEL0, PriorityV2ToV1(commonv2.Priority(100)))
}

func TestSizeScopeV2ToV1(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		v2 commonv2.SizeScope
		v1 commonv1.SizeScope
	}{
		{commonv2.SizeScope_NORMAL, commonv1.SizeScope_NORMAL},
		{commonv2.SizeScope_SMALL, commonv1.SizeScope_SMALL},
		{commonv2.SizeScope_TINY, commonv1.SizeScope_TINY},
		{commonv2.SizeScope_EMPTY, commonv1.SizeScope_EMPTY},
		{commonv2.SizeScope_UNKNOW, commonv1.SizeScope_UNKNOW},
	}
	for _, tt := range tests {
		assert.Equal(tt.v1, SizeScopeV2ToV1(tt.v2))
	}
	assert.Equal(commonv1.SizeScope_UNKNOW, SizeScopeV2ToV1(commonv2.SizeScope(100)))
}
