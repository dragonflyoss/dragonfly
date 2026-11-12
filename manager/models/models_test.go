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

package models

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils/tests"
)

type mockDialector struct {
	tests.DummyDialector
	name string
}

func (d mockDialector) Name() string {
	return d.name
}

func TestPaginate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{Logger: gormlogger.Discard})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		page    int
		perPage int
		expect  func(t *testing.T, sql string)
	}{
		{
			name:    "first page has no offset",
			page:    1,
			perPage: 10,
			expect: func(t *testing.T, sql string) {
				assert := assert.New(t)
				assert.Contains(sql, "LIMIT 10")
				assert.NotContains(sql, "OFFSET")
			},
		},
		{
			name:    "later page offsets by the preceding pages",
			page:    3,
			perPage: 25,
			expect: func(t *testing.T, sql string) {
				assert := assert.New(t)
				assert.Contains(sql, "LIMIT 25 OFFSET 50")
			},
		},
		{
			name:    "second page with a single item per page",
			page:    2,
			perPage: 1,
			expect: func(t *testing.T, sql string) {
				assert := assert.New(t)
				assert.Contains(sql, "LIMIT 1 OFFSET 1")
			},
		},
		{
			name:    "page zero yields a negative offset that is not emitted",
			page:    0,
			perPage: 10,
			expect: func(t *testing.T, sql string) {
				assert := assert.New(t)
				assert.Contains(sql, "LIMIT 10")
				assert.NotContains(sql, "OFFSET")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				return tx.Scopes(Paginate(tc.page, tc.perPage)).Find(&[]Config{})
			}))
		})
	}
}

func TestJSONMap_Value(t *testing.T) {
	tests := []struct {
		name   string
		m      JSONMap
		expect func(t *testing.T, value driver.Value, err error)
	}{
		{
			name: "nil map",
			m:    nil,
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("{}", value)
			},
		},
		{
			name: "empty map",
			m:    JSONMap{},
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("{}", value)
			},
		},
		{
			name: "map with nested values",
			m:    JSONMap{"foo": "bar", "nested": map[string]any{"n": float64(1)}, "list": []any{"a"}},
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.JSONEq(`{"foo":"bar","nested":{"n":1},"list":["a"]}`, value.(string))
			},
		},
		{
			name: "map with unsupported value",
			m:    JSONMap{"ch": make(chan struct{})},
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				var unsupportedErr *json.UnsupportedTypeError
				assert.ErrorAs(err, &unsupportedErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.m.Value()
			tc.expect(t, value, err)
		})
	}
}

func TestJSONMap_Scan(t *testing.T) {
	tests := []struct {
		name   string
		val    any
		expect func(t *testing.T, m JSONMap, err error)
	}{
		{
			name: "nil value",
			val:  nil,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{}, m)
			},
		},
		{
			name: "bytes value",
			val:  []byte(`{"foo":"bar","n":1}`),
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{"foo": "bar", "n": float64(1)}, m)
			},
		},
		{
			name: "string value",
			val:  `{"foo":"bar"}`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{"foo": "bar"}, m)
			},
		},
		{
			name: "empty bytes value",
			val:  []byte{},
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{}, m)
			},
		},
		{
			name: "empty string value",
			val:  "",
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{}, m)
			},
		},
		{
			name: "null literal",
			val:  "null",
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{}, m)
			},
		},
		{
			name: "unsupported type",
			val:  1,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(m)
			},
		},
		{
			name: "invalid json",
			val:  `{"foo":`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
		{
			name: "json array is not an object",
			val:  `["foo"]`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				var typeErr *json.UnmarshalTypeError
				assert.ErrorAs(err, &typeErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var m JSONMap
			err := m.Scan(tc.val)
			tc.expect(t, m, err)
		})
	}
}

func TestJSONMap_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		m      JSONMap
		expect func(t *testing.T, b []byte, err error)
	}{
		{
			name: "nil map",
			m:    nil,
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("null", string(b))
			},
		},
		{
			name: "empty map",
			m:    JSONMap{},
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("{}", string(b))
			},
		},
		{
			name: "map with values",
			m:    JSONMap{"foo": "bar", "n": float64(1)},
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.JSONEq(`{"foo":"bar","n":1}`, string(b))
			},
		},
		{
			name: "nested map marshals through the struct field",
			m:    JSONMap{"inner": JSONMap{"k": "v"}},
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.JSONEq(`{"inner":{"k":"v"}}`, string(b))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.m)
			tc.expect(t, b, err)
		})
	}
}

func TestJSONMap_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		expect func(t *testing.T, m JSONMap, err error)
	}{
		{
			name: "object",
			data: `{"foo":"bar","n":1}`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{"foo": "bar", "n": float64(1)}, m)
			},
		},
		{
			name: "empty object",
			data: `{}`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(JSONMap{}, m)
			},
		},
		{
			name: "invalid json",
			data: `{"foo":`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
		{
			name: "json string is not an object",
			data: `"foo"`,
			expect: func(t *testing.T, m JSONMap, err error) {
				assert := assert.New(t)
				var typeErr *json.UnmarshalTypeError
				assert.ErrorAs(err, &typeErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var m JSONMap
			err := json.Unmarshal([]byte(tc.data), &m)
			tc.expect(t, m, err)
		})
	}
}

func TestJSONMap_GormDBDataType(t *testing.T) {
	tests := []struct {
		name    string
		dialect string
		expect  func(t *testing.T, dataType string)
	}{
		{
			name:    "postgres",
			dialect: "postgres",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("text", dataType)
			},
		},
		{
			name:    "sqlite",
			dialect: "sqlite",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("TEXT", dataType)
			},
		},
		{
			name:    "sqlserver",
			dialect: "sqlserver",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("NVARCHAR(MAX)", dataType)
			},
		},
		{
			name:    "mysql",
			dialect: "mysql",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("longtext", dataType)
			},
		},
		{
			name:    "unknown dialect falls back to longtext",
			dialect: "dummy",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("longtext", dataType)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := &gorm.DB{Config: &gorm.Config{Dialector: mockDialector{name: tc.dialect}}}
			tc.expect(t, JSONMap{}.GormDBDataType(db, nil))
		})
	}
}

func TestArray_Value(t *testing.T) {
	tests := []struct {
		name   string
		a      Array
		expect func(t *testing.T, value driver.Value, err error)
	}{
		{
			name: "nil array",
			a:    nil,
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Nil(value)
			},
		},
		{
			name: "empty array",
			a:    Array{},
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("[]", value)
			},
		},
		{
			name: "array with values",
			a:    Array{"job", "cluster"},
			expect: func(t *testing.T, value driver.Value, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(`["job","cluster"]`, value)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.a.Value()
			tc.expect(t, value, err)
		})
	}
}

func TestArray_Scan(t *testing.T) {
	tests := []struct {
		name   string
		val    any
		expect func(t *testing.T, a Array, err error)
	}{
		{
			name: "nil value",
			val:  nil,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{}, a)
			},
		},
		{
			name: "bytes value",
			val:  []byte(`["job","cluster"]`),
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{"job", "cluster"}, a)
			},
		},
		{
			name: "string value",
			val:  `["job"]`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{"job"}, a)
			},
		},
		{
			name: "empty bytes value",
			val:  []byte{},
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{}, a)
			},
		},
		{
			name: "null literal",
			val:  "null",
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{}, a)
			},
		},
		{
			name: "unsupported type",
			val:  1,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(a)
			},
		},
		{
			name: "invalid json",
			val:  `["job"`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
		{
			name: "json object is not an array",
			val:  `{"foo":"bar"}`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				var typeErr *json.UnmarshalTypeError
				assert.ErrorAs(err, &typeErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var a Array
			err := a.Scan(tc.val)
			tc.expect(t, a, err)
		})
	}
}

func TestArray_MarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		a      Array
		expect func(t *testing.T, b []byte, err error)
	}{
		{
			name: "nil array",
			a:    nil,
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("null", string(b))
			},
		},
		{
			name: "empty array",
			a:    Array{},
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("[]", string(b))
			},
		},
		{
			name: "array with values",
			a:    Array{"job", "cluster"},
			expect: func(t *testing.T, b []byte, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(`["job","cluster"]`, string(b))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.a)
			tc.expect(t, b, err)
		})
	}
}

func TestArray_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		expect func(t *testing.T, a Array, err error)
	}{
		{
			name: "array",
			data: `["job","cluster"]`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{"job", "cluster"}, a)
			},
		},
		{
			name: "empty array",
			data: `[]`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(Array{}, a)
			},
		},
		{
			name: "invalid json",
			data: `["job"`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
		{
			name: "array with non string element",
			data: `[1]`,
			expect: func(t *testing.T, a Array, err error) {
				assert := assert.New(t)
				var typeErr *json.UnmarshalTypeError
				assert.ErrorAs(err, &typeErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var a Array
			err := json.Unmarshal([]byte(tc.data), &a)
			tc.expect(t, a, err)
		})
	}
}

func TestArray_GormDBDataType(t *testing.T) {
	tests := []struct {
		name    string
		dialect string
		expect  func(t *testing.T, dataType string)
	}{
		{
			name:    "postgres",
			dialect: "postgres",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("text", dataType)
			},
		},
		{
			name:    "sqlite",
			dialect: "sqlite",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("TEXT", dataType)
			},
		},
		{
			name:    "sqlserver",
			dialect: "sqlserver",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("NVARCHAR(MAX)", dataType)
			},
		},
		{
			name:    "mysql",
			dialect: "mysql",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("longtext", dataType)
			},
		},
		{
			name:    "unknown dialect falls back to longtext",
			dialect: "dummy",
			expect: func(t *testing.T, dataType string) {
				assert := assert.New(t)
				assert.Equal("longtext", dataType)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := &gorm.DB{Config: &gorm.Config{Dialector: mockDialector{name: tc.dialect}}}
			tc.expect(t, Array{}.GormDBDataType(db, nil))
		})
	}
}
