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

package redis

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-redis/redismock/v9"
	goredis "github.com/redis/go-redis/v9"
)

func TestSchedulerClusterHashTags(t *testing.T) {
	keys := []string{
		MakePersistentHostKeyInScheduler(1, "host"), MakePersistentHostsInScheduler(1),
		MakePersistentTaskKeyInScheduler(1, "task"), MakePersistentTasksInScheduler(1),
		MakePersistentCachePeerKeyForPersistentTaskInScheduler(1, "peer"),
		MakePersistentCachePeersOfPersistentTaskInScheduler(1, "task"),
		MakePersistentPeersOfPersistentTaskInScheduler(1, "task"),
		MakePersistentCachePeersOfPersistentHostInScheduler(1, "host"),
		MakePersistentCachePeersForPersistentTaskInScheduler(1),
		MakePersistentCacheHostKeyInScheduler(1, "host"), MakePersistentCacheHostsInScheduler(1),
		MakePersistentCacheTaskKeyInScheduler(1, "task"), MakePersistentCacheTasksInScheduler(1),
		MakePersistentCachePeerKeyForPersistentCacheTaskInScheduler(1, "peer"),
		MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task"),
		MakePersistentPeersOfPersistentCacheTaskInScheduler(1, "task"),
		MakePersistentCachePeersOfPersistentCacheHostInScheduler(1, "host"),
		MakePersistentCachePeersForPersistentCacheTaskInScheduler(1),
	}
	for _, key := range keys {
		if !strings.HasPrefix(key, "scheduler:scheduler-clusters:1:") || SchedulerKey(key, false) != key {
			t.Fatalf("Default format changed: %s", key)
		}
		key = SchedulerKey(key, true)
		if !strings.HasPrefix(key, "scheduler:scheduler-clusters:{1}:") {
			t.Fatalf("Key does not share the cluster hash tag: %s", key)
		}
	}
	if !strings.Contains(SchedulerKey(MakePersistentHostKeyInScheduler(2, "host"), true), "{2}") {
		t.Fatal("Different scheduler clusters must have different tags")
	}
	if MakeKeyInScheduler("other", "id") != "scheduler:other:id" {
		t.Fatal("Unrelated namespaces must remain unchanged")
	}
}

func TestRejectClusterDatabase(t *testing.T) {
	if _, err := NewRedis(&goredis.UniversalOptions{Addrs: []string{"host1:6379", "host2:6379"}, DB: 10}); err == nil {
		t.Fatal("Nonzero Cluster database accepted")
	}
}

func TestSchedulerKeyPreservesOtherNamespaces(t *testing.T) {
	for _, key := range []string{"scheduler:other:id", "manager:scheduler-clusters:1:schedulers", "scheduler:scheduler-clusters:{1}:persistent-tasks:task"} {
		if got := SchedulerKey(key, true); got != key {
			t.Fatalf("SchedulerKey(%q) = %q", key, got)
		}
	}
}

func TestScanSchedulerKeysUsesOwningPrimary(t *testing.T) {
	primary, primaryMock := redismock.NewClientMock()
	replica, replicaMock := redismock.NewClientMock()
	cluster := goredis.NewClusterClient(&goredis.ClusterOptions{
		Addrs: []string{"primary:6379"},
		ClusterSlots: func(context.Context) ([]goredis.ClusterSlot, error) {
			return []goredis.ClusterSlot{{Start: 0, End: 16383, Nodes: []goredis.ClusterNode{
				{Addr: "primary:6379"}, {Addr: "replica:6379"},
			}}}, nil
		},
		NewClient: func(options *goredis.Options) *goredis.Client {
			if options.Addr == "primary:6379" {
				return primary
			}
			return replica
		},
	})
	cluster.AddHook(rejectClusterScanHook{})
	t.Cleanup(func() { _ = cluster.Close() })
	pattern := "scheduler:scheduler-clusters:{1}:persistent-tasks:*"
	key := "scheduler:scheduler-clusters:{1}:persistent-tasks:task"
	primaryMock.ExpectScan(5, pattern, 10).SetVal([]string{key}, 7)
	keys, cursor, err := ScanSchedulerKeys(context.Background(), cluster, 5, pattern, 10)
	if err != nil || cursor != 7 || len(keys) != 1 || keys[0] != key {
		t.Fatalf("ScanSchedulerKeys = %v, %d, %v", keys, cursor, err)
	}
	if err := primaryMock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
	if err := replicaMock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

type rejectClusterScanHook struct{}

func (rejectClusterScanHook) DialHook(next goredis.DialHook) goredis.DialHook {
	return next
}

func (rejectClusterScanHook) ProcessHook(next goredis.ProcessHook) goredis.ProcessHook {
	return func(ctx context.Context, cmd goredis.Cmder) error {
		if cmd.Name() == "scan" {
			return errors.New("SCAN must use the owning primary client")
		}
		return next(ctx, cmd)
	}
}

func (rejectClusterScanHook) ProcessPipelineHook(next goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return next
}
