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

package dragonfly

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	managertypes "d7y.io/dragonfly/v2/manager/types"
)

func TestNew(t *testing.T) {
	_, err := New(Options{Endpoint: "", Token: "t"})
	assert.Error(t, err)

	_, err = New(Options{Endpoint: "http://manager", Token: ""})
	assert.Error(t, err)

	_, err = New(Options{Endpoint: "ftp://manager", Token: "t"})
	assert.Error(t, err)

	_, err = New(Options{Endpoint: "https://manager:8080/prefix", Token: "t", InsecureSkipTLSVerify: true})
	assert.NoError(t, err)
}

func TestClient_CreateAndGetJob(t *testing.T) {
	var received struct {
		method string
		path   string
		auth   string
		body   map[string]any
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.method = r.Method
		received.path = r.URL.Path
		received.auth = r.Header.Get("Authorization")

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/prefix/oapi/v1/jobs":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&received.body))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id": 7, "task_id": "group_x", "type": "preheat", "state": "PENDING"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/prefix/oapi/v1/jobs/7":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id": 7, "type": "preheat", "state": "SUCCESS", "result": {"group_uuid": "group_x", "state": "SUCCESS", "job_states": [{"task_name": "preheat", "state": "SUCCESS", "results": [{"success_tasks": [{"url": "https://example.com/a", "hostname": "h", "ip": "1.1.1.1"}], "failure_tasks": [], "scheduler_cluster_id": 1}]}]}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/prefix/oapi/v1/jobs/404":
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message": "record not found"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/prefix/oapi/v1/jobs/500":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"errors": "boom"}`))
		default:
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message": "Unauthorized"}`))
		}
	}))
	defer server.Close()

	c, err := New(Options{Endpoint: server.URL + "/prefix/", Token: "secret"})
	require.NoError(t, err)

	job, err := c.CreatePreheatJob(context.Background(), managertypes.CreatePreheatJobRequest{
		Args: managertypes.PreheatArgs{Type: "file", URLs: []string{"https://example.com/a"}},
	})
	require.NoError(t, err)
	assert.Equal(t, uint(7), job.ID)
	assert.Equal(t, JobStatePending, job.State)
	assert.Equal(t, http.MethodPost, received.method)
	assert.Equal(t, "Bearer secret", received.auth)
	assert.Equal(t, "preheat", received.body["type"], "job type is forced by the client")

	job, err = c.GetJob(context.Background(), 7)
	require.NoError(t, err)
	assert.Equal(t, JobStateSuccess, job.State)
	assert.True(t, IsTerminalState(job.State))

	results, errs, err := PreheatResults(job)
	require.NoError(t, err)
	assert.Empty(t, errs)
	require.Len(t, results, 1)
	require.Len(t, results[0].SuccessTasks, 1)
	assert.Equal(t, "https://example.com/a", results[0].SuccessTasks[0].URL)

	_, err = c.GetJob(context.Background(), 404)
	assert.True(t, errors.Is(err, ErrNotFound))

	_, err = c.GetJob(context.Background(), 500)
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Equal(t, "boom", apiErr.Message)
}

func TestParseGroupJobResult(t *testing.T) {
	group, err := ParseGroupJobResult(&Job{})
	require.NoError(t, err)
	assert.Empty(t, group.JobStates)

	job := &Job{Result: map[string]any{
		"group_uuid": "g",
		"state":      "FAILURE",
		"job_states": []any{
			map[string]any{"task_name": "get_task", "state": "FAILURE", "error": "scheduler down", "results": []any{}},
			map[string]any{"task_name": "get_task", "state": "SUCCESS", "results": []any{
				map[string]any{"peers": []any{map[string]any{"id": "p1", "ip": "1.1.1.1"}}, "scheduler_cluster_id": 2},
			}},
		},
	}}

	results, errs, err := GetTaskResults(job)
	require.NoError(t, err)
	assert.Equal(t, []string{"scheduler down"}, errs)
	require.Len(t, results, 1)
	assert.Len(t, results[0].Peers, 1)
	assert.Equal(t, uint(2), results[0].SchedulerClusterID)
}
