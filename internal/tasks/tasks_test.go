package tasks

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	machineryv1tasks "github.com/RichardKnop/machinery/v1/tasks"
)

func TestTaskMarshal(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		expect func(t *testing.T, result []machineryv1tasks.Arg, err error)
	}{
		{
			name: "marshal common struct",
			value: struct {
				I int64   `json:"i" binding:"required"`
				F float64 `json:"f" binding:"required"`
				S string  `json:"s" binding:"required"`
			}{
				I: 1,
				F: 1.1,
				S: "foo",
			},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Type:"string", Value:"{\"i\":1,\"f\":1.1,\"s\":\"foo\"}"}}, result)
			},
		},
		{
			name: "marshal struct with zero value",
			value: struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name:"", Type:"string", Value:"{\"i\":0,\"f\":0,\"s\":\"\"}"}}, result)
			},
		},
		{
			name: "marshal struct with slice",
			value: struct {
				S []string  `json:"s" binding:"required"`
			}{
				S: []string{},
			},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name:"", Type:"string", Value:"{\"s\":[]}"}}, result)
			},
		},
		{
			name: "marshal struct with nil slice",
			value: struct {
				S []string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name:"", Type:"string", Value:"{\"s\":null}"}}, result)
			},
		},
		{
			name: "marshal nil",
			value: nil,
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name:"", Type:"string", Value:"null"}}, result)
			},
		},
		{
			name: "marshal unsupported type",
			value: struct {
				C chan struct{} `json:"c" binding:"required"`
			}{
				C: make(chan struct{}),
			},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal("json: unsupported type: chan struct {}", err.Error())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arg, err := Marshal(tc.value)
			tc.expect(t, arg, err)
		})
	}
}

func TestTaskUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		data   []reflect.Value
		value  interface{}
		expect func(t *testing.T, result interface{}, err error)
	}{
		{
			name: "unmarshal common struct",
			data: []reflect.Value{
				reflect.ValueOf("{\"i\":1,\"f\":1.1,\"s\":\"foo\"}"),
			},
			value: &struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result interface{}, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					I int64   `json:"i" binding:"omitempty"`
					F float64 `json:"f" binding:"omitempty"`
					S string  `json:"s" binding:"omitempty"`
				}{1, 1.1, "foo"}, result)
			},
		},
		//TODO: add more test cases
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Unmarshal(tc.data, tc.value)
			tc.expect(t, tc.value, err)
		})
	}
}
