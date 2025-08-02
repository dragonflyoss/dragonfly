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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/dragonflyoss/machinery/v1"
	machineryv1config "github.com/dragonflyoss/machinery/v1/config"
	machineryv1log "github.com/dragonflyoss/machinery/v1/log"
	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

// tracer is a global tracer for job.
var tracer = otel.Tracer("job")

type Config struct {
	Addrs            []string
	MasterName       string
	Username         string
	Password         string
	SentinelUsername string
	SentinelPassword string
	BrokerDB         int
	BackendDB        int
}

type Job struct {
	Server *machinery.Server
	Worker *machinery.Worker
	Queue  Queue
}

func New(cfg *Config, queue Queue) (*Job, error) {
	// Set logger
	machineryv1log.Set(&MachineryLogger{})

	if err := ping(&redis.UniversalOptions{
		Addrs:            cfg.Addrs,
		MasterName:       cfg.MasterName,
		Username:         cfg.Username,
		Password:         cfg.Password,
		SentinelUsername: cfg.SentinelUsername,
		SentinelPassword: cfg.SentinelPassword,
		DB:               cfg.BackendDB,
	}); err != nil {
		return nil, err
	}

	var broker string
	if cfg.Username != "" {
		broker = fmt.Sprintf("redis://%s:%s@%s/%d", url.QueryEscape(cfg.Username), url.QueryEscape(cfg.Password), strings.Join(cfg.Addrs, ","), cfg.BrokerDB)
	} else {
		broker = fmt.Sprintf("redis://%s@%s/%d", url.QueryEscape(cfg.Password), strings.Join(cfg.Addrs, ","), cfg.BrokerDB)
	}

	var backend string
	if cfg.Username != "" {
		backend = fmt.Sprintf("redis://%s:%s@%s/%d", url.QueryEscape(cfg.Username), url.QueryEscape(cfg.Password), strings.Join(cfg.Addrs, ","), cfg.BackendDB)
	} else {
		backend = fmt.Sprintf("redis://%s@%s/%d", url.QueryEscape(cfg.Password), strings.Join(cfg.Addrs, ","), cfg.BackendDB)
	}

	server, err := machinery.NewServer(&machineryv1config.Config{
		Broker:          broker,
		DefaultQueue:    queue.String(),
		ResultBackend:   backend,
		ResultsExpireIn: DefaultResultsExpireIn,
		Redis: &machineryv1config.RedisConfig{
			MasterName:             cfg.MasterName,
			MaxIdle:                DefaultRedisMaxIdle,
			MaxActive:              DefaultRedisMaxActive,
			IdleTimeout:            DefaultRedisIdleTimeout,
			ReadTimeout:            DefaultRedisReadTimeout,
			WriteTimeout:           DefaultRedisWriteTimeout,
			ConnectTimeout:         DefaultRedisConnectTimeout,
			NormalTasksPollPeriod:  DefaultRedisNormalTasksPollPeriod,
			DelayedTasksPollPeriod: DefaultRedisDelayedTasksPollPeriod,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Job{
		Server: server,
		Queue:  queue,
	}, nil
}

func ping(options *redis.UniversalOptions) error {
	return redis.NewUniversalClient(options).Ping(context.Background()).Err()
}

func (t *Job) RegisterJob(namedJobFuncs map[string]any) error {
	return t.Server.RegisterTasks(namedJobFuncs)
}

func (t *Job) LaunchWorker(consumerTag string, concurrency int) error {
	t.Worker = t.Server.NewWorker(consumerTag, concurrency)
	return t.Worker.Launch()
}

type GroupJobState struct {
	GroupUUID string     `json:"group_uuid"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	JobStates []jobState `json:"job_states"`
}

type jobState struct {
	TaskUUID  string    `json:"task_uuid"`
	TaskName  string    `json:"task_name"`
	State     string    `json:"state"`
	Results   []any     `json:"results"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	TTL       int64     `json:"ttl"`
}

// newGroupJobState 构建 GroupJobState 实例，便于主流程复用
func newGroupJobState(groupUUID, state string, createdAt time.Time, jobStates []jobState) *GroupJobState {
	return &GroupJobState{
		GroupUUID: groupUUID,
		State:     state,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
		JobStates: jobStates,
	}
}

func (t *Job) GetGroupJobState(name string, groupUUID string) (*GroupJobState, error) {
	taskStates, err := t.Server.GetBackend().GroupTaskStates(groupUUID, 0)
	if err != nil {
		return nil, err
	}

	if len(taskStates) == 0 {
		return nil, errors.New("empty group")
	}

	jobStates := make([]jobState, len(taskStates))
	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
	)

	for i, taskState := range taskStates {
		wg.Add(1)
		go func(i int, taskState *machineryv1tasks.TaskState) {
			defer wg.Done()
			var results []any
			for _, result := range taskState.Results {
				var err error
				switch name {
				case PreheatJob:
					var resp PreheatResponse
					err = UnmarshalTaskResult(result.Value, &resp)
					results = append(results, resp)
				case GetTaskJob:
					var resp GetTaskResponse
					err = UnmarshalTaskResult(result.Value, &resp)
					results = append(results, resp)
				case DeleteTaskJob:
					var resp DeleteTaskResponse
					err = UnmarshalTaskResult(result.Value, &resp)
					results = append(results, resp)
				default:
					err = errors.New("unsupported unmarshal task result")
				}
				if err != nil {
					once.Do(func() { firstErr = err })
					return
				}
			}
			jobStates[i] = jobState{
				TaskUUID:  taskState.TaskUUID,
				TaskName:  taskState.TaskName,
				State:     taskState.State,
				Results:   results,
				Error:     taskState.Error,
				CreatedAt: taskState.CreatedAt,
				TTL:       taskState.TTL,
			}
		}(i, taskState)
	}
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}

	for _, taskState := range taskStates {
		if taskState.IsFailure() {
			logger.WithGroupAndTaskUUID(groupUUID, taskState.TaskUUID).Errorf("task is failed: %#v", taskState)
			return newGroupJobState(groupUUID, machineryv1tasks.StateFailure, taskState.CreatedAt, jobStates), nil
		}
	}

	for _, taskState := range taskStates {
		if !taskState.IsSuccess() {
			logger.WithGroupAndTaskUUID(groupUUID, taskState.TaskUUID).Infof("task is not succeeded: %#v", taskState)
			return newGroupJobState(groupUUID, machineryv1tasks.StatePending, taskState.CreatedAt, jobStates), nil
		}
	}

	return newGroupJobState(groupUUID, machineryv1tasks.StateSuccess, taskStates[0].CreatedAt, jobStates), nil
}

func MarshalResponse(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func MarshalRequest(v any) ([]machineryv1tasks.Arg, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return []machineryv1tasks.Arg{{
		Type:  "string",
		Value: string(b),
	}}, nil
}

func UnmarshalTaskResult(data any, v any) error {
	s, ok := data.(string)
	if !ok {
		return errors.New("invalid task result")
	}

	return json.Unmarshal([]byte(s), v)
}

func UnmarshalResponse(data []reflect.Value, v any) error {
	if len(data) == 0 {
		return errors.New("empty data is not specified")
	}

	if err := json.Unmarshal([]byte(data[0].String()), v); err != nil {
		return err
	}

	return nil
}

func UnmarshalRequest(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}
