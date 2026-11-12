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

// Package dragonfly is a small client for the job endpoints of the Dragonfly
// manager open API, authenticated with a personal access token.
package dragonfly

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	managertypes "d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/version"
)

// Job states reported by the manager, mirroring the machinery task states.
const (
	JobStatePending  = "PENDING"
	JobStateReceived = "RECEIVED"
	JobStateStarted  = "STARTED"
	JobStateRetry    = "RETRY"
	JobStateSuccess  = "SUCCESS"
	JobStateFailure  = "FAILURE"
)

// Job types created by the controller. They mirror the names in internal/job,
// which is not imported to keep the machinery dependency out of the controller.
const (
	JobTypePreheat    = "preheat"
	JobTypeGetTask    = "get_task"
	JobTypeDeleteTask = "delete_task"
)

const (
	// DefaultRequestTimeout bounds every request to the manager.
	DefaultRequestTimeout = 30 * time.Second

	jobsPath = "/oapi/v1/jobs"
)

// ErrNotFound is returned when the manager does not know the job.
var ErrNotFound = errors.New("job not found")

// IsTerminalState reports whether a job state is final.
func IsTerminalState(state string) bool {
	return state == JobStateSuccess || state == JobStateFailure
}

// Job is the manager representation of a job.
type Job struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	TaskID    string         `json:"task_id"`
	BIO       string         `json:"bio"`
	Type      string         `json:"type"`
	State     string         `json:"state"`
	Args      map[string]any `json:"args"`
	Result    map[string]any `json:"result"`
}

// GroupJobResult is the aggregated result of a job group as stored by the manager.
type GroupJobResult struct {
	GroupUUID string     `json:"group_uuid"`
	State     string     `json:"state"`
	JobStates []JobState `json:"job_states"`
}

// JobState is the result of one task of a job group, one per scheduler.
type JobState struct {
	TaskUUID string            `json:"task_uuid"`
	TaskName string            `json:"task_name"`
	State    string            `json:"state"`
	Results  []json.RawMessage `json:"results"`
	Error    string            `json:"error"`
}

// PreheatResult is the result of a preheat task on one scheduler.
type PreheatResult struct {
	SuccessTasks       []PreheatSuccessTask `json:"success_tasks"`
	FailureTasks       []PreheatFailureTask `json:"failure_tasks"`
	SchedulerClusterID uint                 `json:"scheduler_cluster_id"`
}

// PreheatSuccessTask is a task preheated on one peer.
type PreheatSuccessTask struct {
	URL      string `json:"url"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

// PreheatFailureTask is a task that failed to preheat on one peer.
type PreheatFailureTask struct {
	URL         string `json:"url"`
	Hostname    string `json:"hostname"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

// GetTaskResult is the result of a get_task task on one scheduler.
type GetTaskResult struct {
	Peers              []TaskPeer `json:"peers"`
	SchedulerClusterID uint       `json:"scheduler_cluster_id"`
}

// TaskPeer is a peer holding a task.
type TaskPeer struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	HostType string `json:"host_type"`
}

// DeleteTaskResult is the result of a delete_task task on one scheduler.
type DeleteTaskResult struct {
	SuccessTasks       []DeleteTaskPeer `json:"success_tasks"`
	FailureTasks       []DeleteTaskPeer `json:"failure_tasks"`
	SchedulerClusterID uint             `json:"scheduler_cluster_id"`
}

// DeleteTaskPeer is a peer a task was deleted from, or failed to be deleted from.
type DeleteTaskPeer struct {
	Hostname    string `json:"hostname"`
	IP          string `json:"ip"`
	HostType    string `json:"host_type"`
	Description string `json:"description"`
}

// Client talks to the job API of a Dragonfly manager.
type Client interface {
	// CreatePreheatJob creates a preheat job.
	CreatePreheatJob(ctx context.Context, req managertypes.CreatePreheatJobRequest) (*Job, error)

	// CreateGetTaskJob creates a get_task job.
	CreateGetTaskJob(ctx context.Context, req managertypes.CreateGetTaskJobRequest) (*Job, error)

	// CreateDeleteTaskJob creates a delete_task job.
	CreateDeleteTaskJob(ctx context.Context, req managertypes.CreateDeleteTaskJobRequest) (*Job, error)

	// GetJob returns a job by ID. It returns ErrNotFound when the manager does not know the job.
	GetJob(ctx context.Context, id uint) (*Job, error)
}

// Options configure a Client.
type Options struct {
	// Endpoint is the base URL of the manager, for example http://dragonfly-manager:8080.
	Endpoint string

	// Token is a personal access token with the job scope.
	Token string

	// InsecureSkipTLSVerify disables TLS verification.
	InsecureSkipTLSVerify bool

	// Timeout bounds every request. Zero means DefaultRequestTimeout.
	Timeout time.Duration

	// HTTPClient overrides the HTTP client, mainly for tests.
	HTTPClient *http.Client
}

// APIError is a non 2xx response from the manager.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("manager api error: status %d: %s", e.StatusCode, e.Message)
}

type client struct {
	baseURL    *url.URL
	token      string
	httpClient *http.Client
}

// New creates a Client.
func New(opts Options) (Client, error) {
	if opts.Endpoint == "" {
		return nil, errors.New("manager endpoint is required")
	}

	if opts.Token == "" {
		return nil, errors.New("manager token is required")
	}

	baseURL, err := url.Parse(opts.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse manager endpoint %q: %w", opts.Endpoint, err)
	}

	if baseURL.Scheme != "http" && baseURL.Scheme != "https" {
		return nil, fmt.Errorf("manager endpoint %q must use http or https", opts.Endpoint)
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		timeout := opts.Timeout
		if timeout == 0 {
			timeout = DefaultRequestTimeout
		}

		transport := http.DefaultTransport.(*http.Transport).Clone()
		if opts.InsecureSkipTLSVerify {
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // nolint:gosec
		}

		httpClient = &http.Client{Timeout: timeout, Transport: transport}
	}

	return &client{
		baseURL:    baseURL,
		token:      opts.Token,
		httpClient: httpClient,
	}, nil
}

func (c *client) CreatePreheatJob(ctx context.Context, req managertypes.CreatePreheatJobRequest) (*Job, error) {
	req.Type = JobTypePreheat
	return c.createJob(ctx, req)
}

func (c *client) CreateGetTaskJob(ctx context.Context, req managertypes.CreateGetTaskJobRequest) (*Job, error) {
	req.Type = JobTypeGetTask
	return c.createJob(ctx, req)
}

func (c *client) CreateDeleteTaskJob(ctx context.Context, req managertypes.CreateDeleteTaskJobRequest) (*Job, error) {
	req.Type = JobTypeDeleteTask
	return c.createJob(ctx, req)
}

func (c *client) GetJob(ctx context.Context, id uint) (*Job, error) {
	var job Job
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("%s/%d", jobsPath, id), nil, &job); err != nil {
		var apiErr *APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("job %d: %w", id, ErrNotFound)
		}

		return nil, err
	}

	return &job, nil
}

func (c *client) createJob(ctx context.Context, body any) (*Job, error) {
	var job Job
	if err := c.do(ctx, http.MethodPost, jobsPath, body, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *client) do(ctx context.Context, method, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}

		reader = bytes.NewReader(payload)
	}

	u := *c.baseURL
	u.Path = strings.TrimSuffix(u.Path, "/") + path

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "dragonfly-data-controller/"+version.GitVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, u.Redacted(), err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Message: errorMessage(data)}
	}

	if out == nil || len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// errorMessage extracts the message of an error response, which the manager
// returns either as {"message": ...} or {"errors": ...}.
func errorMessage(data []byte) string {
	var body struct {
		Message string `json:"message"`
		Errors  any    `json:"errors"`
	}

	if err := json.Unmarshal(data, &body); err == nil {
		if body.Message != "" {
			return body.Message
		}

		if body.Errors != nil {
			return fmt.Sprint(body.Errors)
		}
	}

	msg := strings.TrimSpace(string(data))
	if len(msg) > 512 {
		msg = msg[:512]
	}

	return msg
}

// ParseGroupJobResult decodes the aggregated result stored on a job.
func ParseGroupJobResult(j *Job) (*GroupJobResult, error) {
	if j == nil || len(j.Result) == 0 {
		return &GroupJobResult{}, nil
	}

	data, err := json.Marshal(j.Result)
	if err != nil {
		return nil, fmt.Errorf("encode job result: %w", err)
	}

	var result GroupJobResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode job result: %w", err)
	}

	return &result, nil
}

// decodeResults decodes every result of every task in a group into T.
func decodeResults[T any](j *Job) ([]T, []string, error) {
	group, err := ParseGroupJobResult(j)
	if err != nil {
		return nil, nil, err
	}

	var (
		results []T
		errs    []string
	)

	for _, state := range group.JobStates {
		if state.Error != "" {
			errs = append(errs, state.Error)
		}

		for _, raw := range state.Results {
			var result T
			if err := json.Unmarshal(raw, &result); err != nil {
				return nil, nil, fmt.Errorf("decode task result: %w", err)
			}

			results = append(results, result)
		}
	}

	return results, errs, nil
}

// PreheatResults decodes the results of a preheat job. The second value lists the
// errors reported by the schedulers.
func PreheatResults(j *Job) ([]PreheatResult, []string, error) {
	return decodeResults[PreheatResult](j)
}

// GetTaskResults decodes the results of a get_task job.
func GetTaskResults(j *Job) ([]GetTaskResult, []string, error) {
	return decodeResults[GetTaskResult](j)
}

// DeleteTaskResults decodes the results of a delete_task job.
func DeleteTaskResults(j *Job) ([]DeleteTaskResult, []string, error) {
	return decodeResults[DeleteTaskResult](j)
}
