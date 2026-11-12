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

package job

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/dragonflyoss/machinery/v1"
	"github.com/dragonflyoss/machinery/v1/backends/eager"
	backendsiface "github.com/dragonflyoss/machinery/v1/backends/iface"
	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/stretchr/testify/assert"
)

func TestJob_MarshalResponse(t *testing.T) {
	tests := []struct {
		name   string
		v      any
		expect func(t *testing.T, arg string, err error)
	}{
		{
			name: "marshal common struct",
			v: struct {
				I int64   `json:"i" binding:"required"`
				F float64 `json:"f" binding:"required"`
				S string  `json:"s" binding:"required"`
			}{
				I: 1,
				F: 1.1,
				S: "foo",
			},
			expect: func(t *testing.T, arg string, err error) {
				assert := assert.New(t)
				assert.Equal("{\"i\":1,\"f\":1.1,\"s\":\"foo\"}", arg)
			},
		},
		{
			name: "marshal empty struct",
			v: struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, arg string, err error) {
				assert := assert.New(t)
				assert.Equal("{\"i\":0,\"f\":0,\"s\":\"\"}", arg)
			},
		},
		{
			name: "marshal struct with slice",
			v: struct {
				S []string `json:"s" binding:"required"`
			}{
				S: []string{},
			},
			expect: func(t *testing.T, arg string, err error) {
				assert := assert.New(t)
				assert.Equal("{\"s\":[]}", arg)
			},
		},
		{
			name: "marshal struct with nil slice",
			v: struct {
				S []string `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, arg string, err error) {
				assert := assert.New(t)
				assert.Equal("{\"s\":null}", arg)
			},
		},
		{
			name: "marshal nil",
			v:    nil,
			expect: func(t *testing.T, arg string, err error) {
				assert := assert.New(t)
				assert.Equal("null", arg)
			},
		},
		{
			name: "marshal unsupported type",
			v: struct {
				C chan struct{} `json:"c" binding:"required"`
			}{
				C: make(chan struct{}),
			},
			expect: func(t *testing.T, arg string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arg, err := MarshalResponse(tc.v)
			tc.expect(t, arg, err)
		})
	}
}

func TestJob_MarshalRequest(t *testing.T) {
	tests := []struct {
		name   string
		value  any
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
				assert.Equal([]machineryv1tasks.Arg{{Type: "string", Value: "{\"i\":1,\"f\":1.1,\"s\":\"foo\"}"}}, result)
			},
		},
		{
			name: "marshal empty struct",
			value: struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name: "", Type: "string", Value: "{\"i\":0,\"f\":0,\"s\":\"\"}"}}, result)
			},
		},
		{
			name: "marshal struct with slice",
			value: struct {
				S []string `json:"s" binding:"required"`
			}{
				S: []string{},
			},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name: "", Type: "string", Value: "{\"s\":[]}"}}, result)
			},
		},
		{
			name: "marshal struct with nil slice",
			value: struct {
				S []string `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name: "", Type: "string", Value: "{\"s\":null}"}}, result)
			},
		},
		{
			name:  "marshal nil",
			value: nil,
			expect: func(t *testing.T, result []machineryv1tasks.Arg, err error) {
				assert := assert.New(t)
				assert.Equal([]machineryv1tasks.Arg{{Name: "", Type: "string", Value: "null"}}, result)
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
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			arg, err := MarshalRequest(tc.value)
			tc.expect(t, arg, err)
		})
	}
}

func TestJob_UnmarshalResponse(t *testing.T) {
	tests := []struct {
		name   string
		data   []reflect.Value
		value  any
		expect func(t *testing.T, result any, err error)
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
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					I int64   `json:"i" binding:"omitempty"`
					F float64 `json:"f" binding:"omitempty"`
					S string  `json:"s" binding:"omitempty"`
				}{1, 1.1, "foo"}, result)
			},
		},
		{
			name: "unmarshal struct lack of parameters",
			data: []reflect.Value{
				reflect.ValueOf("{}"),
			},
			value: &struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					I int64   `json:"i" binding:"omitempty"`
					F float64 `json:"f" binding:"omitempty"`
					S string  `json:"s" binding:"omitempty"`
				}{0, 0, ""}, result)
			},
		},
		{
			name: "unmarshal struct with slice",
			data: []reflect.Value{
				reflect.ValueOf("{\"s\":[]}"),
			},
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					S []string `json:"s" binding:"required"`
				}{S: []string{}}, result)
			},
		},
		{
			name: "unmarshal struct with nil slice",
			data: []reflect.Value{
				reflect.ValueOf("{\"s\":null}"),
			},
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					S []string `json:"s" binding:"required"`
				}{S: nil}, result)
			},
		},
		{
			name: "unmarshal nil data",
			data: []reflect.Value{},
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "unmarshal invalid json",
			data: []reflect.Value{
				reflect.ValueOf("{\"s\":"),
			},
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := UnmarshalResponse(tc.data, tc.value)
			tc.expect(t, tc.value, err)
		})
	}
}

func TestJob_UnmarshalRequest(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		value  any
		expect func(t *testing.T, result any, err error)
	}{
		{
			name: "unmarshal common struct",
			data: "{\"i\":1,\"f\":1.1,\"s\":\"foo\"}",
			value: &struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					I int64   `json:"i" binding:"omitempty"`
					F float64 `json:"f" binding:"omitempty"`
					S string  `json:"s" binding:"omitempty"`
				}{1, 1.1, "foo"}, result)
			},
		},
		{
			name: "unmarshal empty struct",
			data: "{}",
			value: &struct {
				I int64   `json:"i" binding:"omitempty"`
				F float64 `json:"f" binding:"omitempty"`
				S string  `json:"s" binding:"omitempty"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					I int64   `json:"i" binding:"omitempty"`
					F float64 `json:"f" binding:"omitempty"`
					S string  `json:"s" binding:"omitempty"`
				}{0, 0, ""}, result)
			},
		},
		{
			name: "unmarshal struct with slice",
			data: "{\"s\":[]}",
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					S []string `json:"s" binding:"required"`
				}{S: []string{}}, result)
			},
		},
		{
			name: "unmarshal struct with nil slice",
			data: "{\"s\":null}",
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Equal(&struct {
					S []string `json:"s" binding:"required"`
				}{S: nil}, result)
			},
		},
		{
			name: "unmarshal nil data",
			data: "",
			value: &struct {
				S []string `json:"s" binding:"required"`
			}{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := UnmarshalRequest(tc.data, tc.value)
			tc.expect(t, tc.value, err)
		})
	}
}

func TestJob_UnmarshalTaskResult(t *testing.T) {
	tests := []struct {
		name   string
		data   any
		value  any
		expect func(t *testing.T, result any, err error)
	}{
		{
			name:  "unmarshal common struct",
			data:  "{\"scheduler_cluster_id\":1,\"peers\":[]}",
			value: &GetTaskResponse{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&GetTaskResponse{Peers: []*Peer{}, SchedulerClusterID: 1}, result)
			},
		},
		{
			name:  "unmarshal data that is not a string",
			data:  1,
			value: &GetTaskResponse{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Equal(&GetTaskResponse{}, result)
			},
		},
		{
			name:  "unmarshal nil data",
			data:  nil,
			value: &GetTaskResponse{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:  "unmarshal invalid json",
			data:  "{\"scheduler_cluster_id\":",
			value: &GetTaskResponse{},
			expect: func(t *testing.T, result any, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := UnmarshalTaskResult(tc.data, tc.value)
			tc.expect(t, tc.value, err)
		})
	}
}

func TestJob_RequestRoundTrip(t *testing.T) {
	timeout := 30 * time.Second
	pieceLength := uint64(4194304)
	percentage := uint32(50)

	tests := []struct {
		name   string
		value  any
		target any
		expect func(t *testing.T, target any, err error)
	}{
		{
			name: "get task request",
			value: GetTaskRequest{
				TaskID:              "foo",
				Timeout:             timeout,
				GroupUUID:           "group_1",
				TaskUUID:            "task_1",
				ConcurrentPeerCount: 3,
				Scope:               "all",
			},
			target: &GetTaskRequest{},
			expect: func(t *testing.T, target any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&GetTaskRequest{
					TaskID:              "foo",
					Timeout:             timeout,
					GroupUUID:           "group_1",
					TaskUUID:            "task_1",
					ConcurrentPeerCount: 3,
					Scope:               "all",
				}, target)
			},
		},
		{
			name: "delete task request",
			value: DeleteTaskRequest{
				TaskID:    "bar",
				Timeout:   timeout,
				GroupUUID: "group_2",
				TaskUUID:  "task_2",
			},
			target: &DeleteTaskRequest{},
			expect: func(t *testing.T, target any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&DeleteTaskRequest{
					TaskID:    "bar",
					Timeout:   timeout,
					GroupUUID: "group_2",
					TaskUUID:  "task_2",
				}, target)
			},
		},
		{
			name: "preheat request with optional pointers and maps",
			value: PreheatRequest{
				URLs:                []string{"https://example.com/a", "https://example.com/b"},
				PieceLength:         &pieceLength,
				Tag:                 "tag",
				FilteredQueryParams: "a&b",
				Headers:             map[string]string{"Authorization": "Bearer foo"},
				Application:         "app",
				Priority:            5,
				Scope:               "single_seed_peer",
				IPs:                 []string{"127.0.0.1"},
				Percentage:          &percentage,
			},
			target: &PreheatRequest{},
			expect: func(t *testing.T, target any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&PreheatRequest{
					URLs:                []string{"https://example.com/a", "https://example.com/b"},
					PieceLength:         &pieceLength,
					Tag:                 "tag",
					FilteredQueryParams: "a&b",
					Headers:             map[string]string{"Authorization": "Bearer foo"},
					Application:         "app",
					Priority:            5,
					Scope:               "single_seed_peer",
					IPs:                 []string{"127.0.0.1"},
					Percentage:          &percentage,
				}, target)
			},
		},
		{
			name:   "empty preheat request keeps optional pointers nil",
			value:  PreheatRequest{},
			target: &PreheatRequest{},
			expect: func(t *testing.T, target any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(&PreheatRequest{}, target)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args, err := MarshalRequest(tc.value)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, tc.target, UnmarshalRequest(args[0].Value.(string), tc.target))
		})
	}
}

func TestJob_GetGroupJobState(t *testing.T) {
	tests := []struct {
		name      string
		jobName   string
		groupUUID string
		setup     func(backend backendsiface.Backend)
		expect    func(t *testing.T, state *GroupJobState, err error)
	}{
		{
			name:      "group not found",
			jobName:   PreheatJob,
			groupUUID: "group_1",
			setup:     func(backend backendsiface.Backend) {},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(state)
			},
		},
		{
			name:      "group without tasks",
			jobName:   PreheatJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(state)
			},
		},
		{
			name:      "unsupported job name",
			jobName:   "foo",
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{}"}})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(state)
			},
		},
		{
			name:      "invalid task result",
			jobName:   PreheatJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{"}})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				var syntaxErr *json.SyntaxError
				assert.ErrorAs(err, &syntaxErr)
				assert.Nil(state)
			},
		},
		{
			name:      "any failed task fails the group",
			jobName:   PreheatJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1", "task_2", "task_3"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{\"scheduler_cluster_id\":1}"}})
				_ = backend.SetStateFailure(&machineryv1tasks.Signature{UUID: "task_2"}, "boom")
				_ = backend.SetStatePending(&machineryv1tasks.Signature{UUID: "task_3", Name: PreheatJob})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("group_1", state.GroupUUID)
				assert.Equal(machineryv1tasks.StateFailure, state.State)
				assert.Len(state.JobStates, 3)
				assert.WithinDuration(time.Now(), state.UpdatedAt, time.Minute)
			},
		},
		{
			name:      "any unfinished task keeps the group pending",
			jobName:   PreheatJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1", "task_2"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{\"scheduler_cluster_id\":1}"}})
				_ = backend.SetStateStarted(&machineryv1tasks.Signature{UUID: "task_2"})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(machineryv1tasks.StatePending, state.State)
				assert.Len(state.JobStates, 2)
			},
		},
		{
			name:      "all preheat tasks succeeded",
			jobName:   PreheatJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1", "task_2"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{\"scheduler_cluster_id\":1}"}})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_2"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{\"scheduler_cluster_id\":2,\"failure_tasks\":[{\"url\":\"https://example.com\",\"hostname\":\"foo\",\"ip\":\"127.0.0.1\",\"description\":\"boom\"}]}"}})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(machineryv1tasks.StateSuccess, state.State)
				assert.Len(state.JobStates, 2)

				results := make(map[string]any, len(state.JobStates))
				for _, jobState := range state.JobStates {
					assert.Equal(machineryv1tasks.StateSuccess, jobState.State)
					assert.Len(jobState.Results, 1)
					results[jobState.TaskUUID] = jobState.Results[0]
				}

				assert.Equal(PreheatResponse{SchedulerClusterID: 1}, results["task_1"])
				assert.Equal(PreheatResponse{
					SchedulerClusterID: 2,
					FailureTasks: []*PreheatFailureTask{{
						URL:         "https://example.com",
						Hostname:    "foo",
						IP:          "127.0.0.1",
						Description: "boom",
					}},
				}, results["task_2"])
			},
		},
		{
			name:      "all get task tasks succeeded",
			jobName:   GetTaskJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{\"scheduler_cluster_id\":3,\"peers\":[]}"}})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(machineryv1tasks.StateSuccess, state.State)
				assert.Len(state.JobStates, 1)
				assert.Equal([]any{GetTaskResponse{Peers: []*Peer{}, SchedulerClusterID: 3}}, state.JobStates[0].Results)
			},
		},
		{
			name:      "all delete task tasks succeeded",
			jobName:   DeleteTaskJob,
			groupUUID: "group_1",
			setup: func(backend backendsiface.Backend) {
				_ = backend.InitGroup("group_1", []string{"task_1"})
				_ = backend.SetStateSuccess(&machineryv1tasks.Signature{UUID: "task_1"}, []*machineryv1tasks.TaskResult{{Type: "string", Value: "{\"scheduler_cluster_id\":4}"}})
			},
			expect: func(t *testing.T, state *GroupJobState, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(machineryv1tasks.StateSuccess, state.State)
				assert.Len(state.JobStates, 1)
				assert.Equal([]any{DeleteTaskResponse{SchedulerClusterID: 4}}, state.JobStates[0].Results)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			backend := eager.New()
			tc.setup(backend)

			server := &machinery.Server{}
			server.SetBackend(backend)

			state, err := (&Job{Server: server}).GetGroupJobState(tc.jobName, tc.groupUUID)
			tc.expect(t, state, err)
		})
	}
}

func TestConstructGroupJobState(t *testing.T) {
	createdAt := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	jobStates := []jobState{{TaskUUID: "task_1", State: machineryv1tasks.StateSuccess}}

	tests := []struct {
		name      string
		state     string
		jobStates []jobState
		expect    func(t *testing.T, state *GroupJobState)
	}{
		{
			name:      "success state with job states",
			state:     machineryv1tasks.StateSuccess,
			jobStates: jobStates,
			expect: func(t *testing.T, state *GroupJobState) {
				assert := assert.New(t)
				assert.Equal("group_1", state.GroupUUID)
				assert.Equal(machineryv1tasks.StateSuccess, state.State)
				assert.Equal(createdAt, state.CreatedAt)
				assert.Equal(jobStates, state.JobStates)
				assert.WithinDuration(time.Now(), state.UpdatedAt, time.Minute)
			},
		},
		{
			name:      "pending state without job states",
			state:     machineryv1tasks.StatePending,
			jobStates: nil,
			expect: func(t *testing.T, state *GroupJobState) {
				assert := assert.New(t)
				assert.Equal(machineryv1tasks.StatePending, state.State)
				assert.Nil(state.JobStates)
				assert.True(state.UpdatedAt.After(state.CreatedAt))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, constructGroupJobState("group_1", tc.state, createdAt, tc.jobStates))
		})
	}
}
