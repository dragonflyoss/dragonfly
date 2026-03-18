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
	"time"

	"github.com/RichardKnop/machinery/v1"
	machineryv1config "github.com/RichardKnop/machinery/v1/config"
	machineryv1log "github.com/RichardKnop/machinery/v1/log"
	machineryv1tasks "github.com/RichardKnop/machinery/v1/tasks"
	"github.com/redis/go-redis/v9"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
)

var (
	ErrRedisNotInitialized = errors.New("redis client not initialized")
)

type Config struct {
	Addrs      []string
	MasterName string
	Username   string
	Password   string
	BrokerDB   int
	BackendDB  int
}

type Job struct {
	Server *machinery.Server
	Worker *machinery.Worker
	Queue  Queue
	rdb    redis.UniversalClient
}

func New(cfg *Config, queue Queue) (*Job, error) {
	// Set logger
	machineryv1log.Set(&MachineryLogger{})

	if err := ping(&redis.UniversalOptions{
		Addrs:      cfg.Addrs,
		MasterName: cfg.MasterName,
		Username:   cfg.Username,
		Password:   cfg.Password,
		DB:         cfg.BackendDB,
	}); err != nil {
		return nil, err
	}

	server, err := machinery.NewServer(&machineryv1config.Config{
		Broker:          fmt.Sprintf("redis://%s@%s/%d", url.QueryEscape(cfg.Password), strings.Join(cfg.Addrs, ","), cfg.BrokerDB),
		DefaultQueue:    queue.String(),
		ResultBackend:   fmt.Sprintf("redis://%s@%s/%d", url.QueryEscape(cfg.Password), strings.Join(cfg.Addrs, ","), cfg.BackendDB),
		ResultsExpireIn: DefaultResultsExpireIn,
		Redis: &machineryv1config.RedisConfig{
			MasterName:     cfg.MasterName,
			MaxIdle:        DefaultRedisMaxIdle,
			IdleTimeout:    DefaultRedisIdleTimeout,
			ReadTimeout:    DefaultRedisReadTimeout,
			WriteTimeout:   DefaultRedisWriteTimeout,
			ConnectTimeout: DefaultRedisConnectTimeout,
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

// InitRdb initializes the redis client for job.
func (t *Job) InitRdb(cfg *Config) error {
	if t.rdb != nil {
		return nil
	}
	// Initialize redis client.
	var rdb redis.UniversalClient
	rdb, err := pkgredis.NewRedis(&redis.UniversalOptions{
		Addrs:      cfg.Addrs,
		MasterName: cfg.MasterName,
		Username:   cfg.Username,
		Password:   cfg.Password,
		DB:         cfg.BackendDB,
	})
	if err != nil {
		return err
	}
	t.rdb = rdb
	return nil
}

type GroupJobState struct {
	GroupUUID string
	State     string
	CreatedAt time.Time
	JobStates []*machineryv1tasks.TaskState
}

func (t *Job) GetGroupJobState(groupID string) (*GroupJobState, error) {
	taskStates, err := t.Server.GetBackend().GroupTaskStates(groupID, 0)
	if err != nil {
		return nil, err
	}

	if len(taskStates) == 0 {
		return nil, errors.New("empty group")
	}

	for _, taskState := range taskStates {
		if taskState.IsFailure() {
			logger.WithGroupAndTaskID(groupID, taskState.TaskUUID).Errorf("task is failed: %#v", taskState)
			return &GroupJobState{
				GroupUUID: groupID,
				State:     machineryv1tasks.StateFailure,
				CreatedAt: taskState.CreatedAt,
				JobStates: taskStates,
			}, nil
		}
	}

	for _, taskState := range taskStates {
		if !taskState.IsSuccess() {
			logger.WithGroupAndTaskID(groupID, taskState.TaskUUID).Infof("task is not succeeded: %#v", taskState)
			return &GroupJobState{
				GroupUUID: groupID,
				State:     machineryv1tasks.StatePending,
				CreatedAt: taskState.CreatedAt,
				JobStates: taskStates,
			}, nil
		}
	}

	return &GroupJobState{
		GroupUUID: groupID,
		State:     machineryv1tasks.StateSuccess,
		CreatedAt: taskStates[0].CreatedAt,
		JobStates: taskStates,
	}, nil
}

// SetTaskResults sets task results to redis by key.
func (t *Job) SetTaskResults(data []string, jobName string) (string, error) {
	if t.rdb == nil {
		return "", ErrRedisNotInitialized
	}

	if len(data) == 0 {
		return "", nil
	}

	key := fmt.Sprintf("%s_results:%s:%d", jobName, t.Queue.String(), time.Now().Unix())

	ctx := context.Background()

	// save results to redis set
	for _, val := range data {
		if err := t.rdb.SAdd(ctx, key, val).Err(); err != nil {
			logger.Errorf("Failed to SAdd: %v", err)
			return "", err
		}
	}

	// set expire time for the results
	if err := t.rdb.Expire(ctx, key, time.Duration(DefaultResultsExpireIn)*time.Second).Err(); err != nil {
		logger.Errorf("Failed to set expire: %v", err)
		return "", err
	}

	return key, nil
}

// GetTaskResults gets task results from redis by key.
// results is not sorted, caller should sort it if needed.
func (t *Job) GetTaskResults(key string) ([]string, error) {
	if t.rdb == nil {
		return nil, ErrRedisNotInitialized
	}
	var cursor uint64 = 0
	var results []string
	for {
		items, nextCursor, err := t.rdb.SScan(context.Background(), key, cursor, "", 50).Result()
		if err != nil {
			break
		}
		for _, item := range items {
			results = append(results, item)
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return results, nil
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
