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

package structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToMap(t *testing.T) {
	tests := []struct {
		name   string
		s      any
		expect func(t *testing.T, m map[string]any, err error)
	}{
		{
			name: "conver struct to map",
			s: struct {
				Name string
				Age  float64
			}{
				Name: "foo",
				Age:  18,
			},
			expect: func(t *testing.T, m map[string]any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(map[string]any{
					"Name": "foo",
					"Age":  float64(18),
				}, m)
			},
		},
		{
			name: "conver string to map failed",
			s:    "foo",
			expect: func(t *testing.T, m map[string]any, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(m)
			},
		},
		{
			name: "conver number to map failed",
			s:    1,
			expect: func(t *testing.T, m map[string]any, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(m)
			},
		},
		{
			name: "conver nil to map",
			s:    nil,
			expect: func(t *testing.T, m map[string]any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(map[string]any(nil), m)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := StructToMap(tc.s)
			tc.expect(t, m, err)
		})
	}
}

func TestMapToStruct(t *testing.T) {
	type person struct {
		Name string
		Age  float64
	}

	tests := []struct {
		name   string
		m      map[string]any
		expect func(t *testing.T, s *person, err error)
	}{
		{
			name: "conver map to struct",
			m: map[string]any{
				"Name": "foo",
				"Age":  float64(18),
			},
			expect: func(t *testing.T, s *person, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&person{
					Name: "foo",
					Age:  18,
				}, s)
			},
		},
		{
			name: "conver nil map to struct",
			m:    nil,
			expect: func(t *testing.T, s *person, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&person{Name: "", Age: 0}, s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &person{}
			err := MapToStruct(tc.m, s)
			tc.expect(t, s, err)
		})
	}
}
