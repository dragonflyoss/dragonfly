/*
 *     Copyright 2022 The Dragonfly Authors
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

package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	mockVertexID    = "foo"
	mockVertexValue = "bar"
)

func TestVertex_New(t *testing.T) {
	v := NewVertex(mockVertexID, mockVertexValue)
	assert := assert.New(t)
	assert.Equal(mockVertexID, v.ID)
	assert.Equal(mockVertexValue, v.Value)
	assert.Equal(uint(0), v.Parents.Len())
	assert.Equal(uint(0), v.Children.Len())
}

func TestVertex_Degree(t *testing.T) {
	v := NewVertex(mockVertexID, mockVertexValue)
	assert := assert.New(t)
	assert.Equal(mockVertexID, v.ID)
	assert.Equal(mockVertexValue, v.Value)
	assert.Equal(0, v.Degree())

	v.Parents.Add(v)
	assert.Equal(1, v.Degree())

	v.Children.Add(v)
	assert.Equal(2, v.Degree())

	v.Parents.Delete(v)
	assert.Equal(1, v.Degree())

	v.Children.Delete(v)
	assert.Equal(0, v.Degree())
}

func TestVertex_InDegree(t *testing.T) {
	v := NewVertex(mockVertexID, mockVertexValue)
	assert := assert.New(t)
	assert.Equal(mockVertexID, v.ID)
	assert.Equal(mockVertexValue, v.Value)
	assert.Equal(0, v.InDegree())

	v.Parents.Add(v)
	assert.Equal(1, v.InDegree())

	v.Children.Add(v)
	assert.Equal(1, v.InDegree())

	v.Parents.Delete(v)
	assert.Equal(0, v.InDegree())

	v.Children.Delete(v)
	assert.Equal(0, v.InDegree())
}

func TestVertex_OutDegree(t *testing.T) {
	v := NewVertex(mockVertexID, mockVertexValue)
	assert := assert.New(t)
	assert.Equal(mockVertexID, v.ID)
	assert.Equal(mockVertexValue, v.Value)
	assert.Equal(0, v.OutDegree())

	v.Parents.Add(v)
	assert.Equal(0, v.OutDegree())

	v.Children.Add(v)
	assert.Equal(1, v.OutDegree())

	v.Parents.Delete(v)
	assert.Equal(1, v.OutDegree())

	v.Children.Delete(v)
	assert.Equal(0, v.OutDegree())
}
