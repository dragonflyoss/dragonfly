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

package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bits-and-blooms/bitset"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/manager/types"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
)

var mockPersistentCacheTime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func mockRawPersistentCacheTask(id string) map[string]string {
	return map[string]string{
		"id":                       id,
		"tag":                      "tag",
		"application":              "app",
		"state":                    "Succeeded",
		"persistent_replica_count": "3",
		"piece_length":             "4194304",
		"content_length":           "8388608",
		"total_piece_count":        "2",
		"ttl":                      fmt.Sprint(int64(time.Hour)),
		"created_at":               mockPersistentCacheTime.Format(time.RFC3339),
		"updated_at":               mockPersistentCacheTime.Add(time.Minute).Format(time.RFC3339),
	}
}

func mockRawPersistentCachePeer(t *testing.T, id, hostID string) map[string]string {
	assert := assert.New(t)
	finishedPieces := bitset.New(2)
	finishedPieces.Set(0)
	finishedPieces.Set(1)
	data, err := finishedPieces.MarshalBinary()
	assert.NoError(err)

	return map[string]string{
		"id":              id,
		"state":           "Succeeded",
		"persistent":      "true",
		"finished_pieces": string(data),
		"block_parents":   `["parent-1"]`,
		"cost":            fmt.Sprint(int64(time.Second)),
		"host_id":         hostID,
		"created_at":      mockPersistentCacheTime.Format(time.RFC3339),
		"updated_at":      mockPersistentCacheTime.Format(time.RFC3339),
	}
}

func mockRawPersistentCacheHost(id string) map[string]string {
	return map[string]string{
		"id":                                  id,
		"type":                                "normal",
		"hostname":                            "host-1",
		"ip":                                  "10.0.0.1",
		"port":                                "4000",
		"download_port":                       "4001",
		"disable_shared":                      "false",
		"os":                                  "linux",
		"platform":                            "ubuntu",
		"platform_family":                     "debian",
		"platform_version":                    "22.04",
		"kernel_version":                      "5.15",
		"cpu_logical_count":                   "8",
		"cpu_physical_count":                  "4",
		"cpu_percent":                         "12.5",
		"cpu_process_percent":                 "1.5",
		"cpu_times_user":                      "1",
		"cpu_times_system":                    "2",
		"cpu_times_idle":                      "3",
		"cpu_times_nice":                      "4",
		"cpu_times_iowait":                    "5",
		"cpu_times_irq":                       "6",
		"cpu_times_softirq":                   "7",
		"cpu_times_steal":                     "8",
		"cpu_times_guest":                     "9",
		"cpu_times_guest_nice":                "10",
		"memory_total":                        "1000",
		"memory_available":                    "600",
		"memory_used":                         "400",
		"memory_used_percent":                 "40",
		"memory_process_used_percent":         "4",
		"memory_free":                         "600",
		"network_tcp_connection_count":        "20",
		"network_upload_tcp_connection_count": "10",
		"network_location":                    "hangzhou",
		"network_idc":                         "idc-1",
		"network_rx_bandwidth":                "100",
		"network_max_rx_bandwidth":            "1000",
		"network_tx_bandwidth":                "200",
		"network_max_tx_bandwidth":            "2000",
		"disk_total":                          "5000",
		"disk_free":                           "3000",
		"disk_used":                           "2000",
		"disk_used_percent":                   "40",
		"disk_inodes_total":                   "500",
		"disk_inodes_used":                    "100",
		"disk_inodes_free":                    "400",
		"disk_inodes_used_percent":            "20",
		"disk_write_bandwidth":                "30",
		"disk_read_bandwidth":                 "40",
		"build_git_version":                   "v2.3.0",
		"build_git_commit":                    "abcdef",
		"build_go_version":                    "1.25",
		"build_platform":                      "linux/amd64",
		"announce_interval":                   fmt.Sprint(int64(5 * time.Minute)),
		"created_at":                          mockPersistentCacheTime.Format(time.RFC3339),
		"updated_at":                          mockPersistentCacheTime.Format(time.RFC3339),
	}
}

func TestService_DestroyPersistentCacheTask(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, err error)
	}{
		{
			name: "deletes the task key",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectDel(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(1)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "redis error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectDel(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			s := &service{rdb: rdb}
			tc.expect(t, s.DestroyPersistentCacheTask(context.Background(), 1, "task-1"))
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetPersistentCacheTask(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(t *testing.T, mock redismock.ClientMock)
		expect func(t *testing.T, task types.PersistentCacheTask, err error)
	}{
		{
			name: "redis error",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetErr(errors.New("redis error"))
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(task.ID)
			},
		},
		{
			name: "task not found",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(map[string]string{})
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "malformed replica count",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				rawTask := mockRawPersistentCacheTask("task-1")
				rawTask["persistent_replica_count"] = "many"
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(rawTask)
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(task.ID)
			},
		},
		{
			name: "malformed created_at",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				rawTask := mockRawPersistentCacheTask("task-1")
				rawTask["created_at"] = "yesterday"
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(rawTask)
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "peers lookup fails",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetErr(errors.New("smembers error"))
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Empty(task.ID)
			},
		},
		{
			name: "task without peers",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{})
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("task-1", task.ID)
				assert.Equal("tag", task.Tag)
				assert.Equal("app", task.Application)
				assert.Equal("Succeeded", task.State)
				assert.Equal(uint64(3), task.PersistentReplicaCount)
				assert.Equal(uint64(4194304), task.PieceLength)
				assert.Equal(uint64(8388608), task.ContentLength)
				assert.Equal(uint32(2), task.TotalPieceCount)
				assert.Equal(time.Hour, task.TTL)
				assert.True(mockPersistentCacheTime.Equal(task.CreatedAt))
				assert.True(mockPersistentCacheTime.Add(time.Minute).Equal(task.UpdatedAt))
				assert.Empty(task.Peers)
			},
		},
		{
			name: "task with peer and host",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{"peer-1"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentCacheTaskInScheduler(1, "peer-1")).SetVal(mockRawPersistentCachePeer(t, "peer-1", "host-1"))
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host-1")).SetVal(mockRawPersistentCacheHost("host-1"))
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(task.Peers, 1)

				peer := task.Peers[0]
				assert.Equal("peer-1", peer.ID)
				assert.Equal("Succeeded", peer.State)
				assert.True(peer.Persistent)
				assert.Equal(uint(2), peer.FinishedPieces.Count())
				assert.Equal([]string{"parent-1"}, peer.BlockParents)
				assert.Equal(time.Second, peer.Cost)
				assert.True(mockPersistentCacheTime.Equal(peer.CreatedAt))

				host := peer.Host
				assert.Equal("host-1", host.ID)
				assert.Equal("normal", host.Type)
				assert.Equal("host-1", host.Hostname)
				assert.Equal("10.0.0.1", host.IP)
				assert.Equal(int32(4000), host.Port)
				assert.Equal(int32(4001), host.DownloadPort)
				assert.False(host.DisableShared)
				assert.Equal("linux", host.OS)
				assert.Equal(uint32(8), host.CPU.LogicalCount)
				assert.Equal(uint32(4), host.CPU.PhysicalCount)
				assert.Equal(12.5, host.CPU.Percent)
				assert.Equal(float64(10), host.CPU.Times.GuestNice)
				assert.Equal(uint64(1000), host.Memory.Total)
				assert.Equal(float64(40), host.Memory.UsedPercent)
				assert.Equal(uint32(20), host.Network.TCPConnectionCount)
				assert.Equal("hangzhou", host.Network.Location)
				assert.Equal("idc-1", host.Network.IDC)
				assert.Equal(uint64(2000), host.Network.MaxTxBandwidth)
				assert.Equal(uint64(5000), host.Disk.Total)
				assert.Equal(float64(20), host.Disk.InodesUsedPercent)
				assert.Equal(uint64(40), host.Disk.ReadBandwidth)
				assert.Equal("v2.3.0", host.Build.GitVersion)
				assert.Equal("linux/amd64", host.Build.Platform)
				assert.Equal(uint(1), host.SchedulerClusterID)
				assert.Equal(5*time.Minute, host.AnnounceInterval)
				assert.True(mockPersistentCacheTime.Equal(host.CreatedAt))
			},
		},
		{
			name: "peer that fails to load is skipped",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{"peer-1", "peer-2"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentCacheTaskInScheduler(1, "peer-1")).SetVal(map[string]string{})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentCacheTaskInScheduler(1, "peer-2")).SetVal(mockRawPersistentCachePeer(t, "peer-2", "host-1"))
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host-1")).SetVal(mockRawPersistentCacheHost("host-1"))
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(task.Peers, 1)
				assert.Equal("peer-2", task.Peers[0].ID)
			},
		},
		{
			name: "peer whose host is missing is skipped",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{"peer-1"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentCacheTaskInScheduler(1, "peer-1")).SetVal(mockRawPersistentCachePeer(t, "peer-1", "host-1"))
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host-1")).SetVal(map[string]string{})
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("task-1", task.ID)
				assert.Empty(task.Peers)
			},
		},
		{
			name: "peer with malformed block parents is skipped",
			mock: func(t *testing.T, mock redismock.ClientMock) {
				rawPeer := mockRawPersistentCachePeer(t, "peer-1", "host-1")
				rawPeer["block_parents"] = "not-json"
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{"peer-1"})
				mock.ExpectHGetAll(pkgredis.MakePersistentCachePeerKeyForPersistentCacheTaskInScheduler(1, "peer-1")).SetVal(rawPeer)
			},
			expect: func(t *testing.T, task types.PersistentCacheTask, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(task.Peers)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(t, mock)

			s := &service{rdb: rdb}
			task, err := s.GetPersistentCacheTask(context.Background(), 1, "task-1")
			tc.expect(t, task, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestService_GetPersistentCacheTasks(t *testing.T) {
	prefix := fmt.Sprintf("%s:", pkgredis.MakePersistentCacheTasksInScheduler(1))
	tests := []struct {
		name   string
		mock   func(mock redismock.ClientMock)
		expect func(t *testing.T, tasks []types.PersistentCacheTask, count int64, err error)
	}{
		{
			name: "scan error",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetErr(errors.New("scan error"))
			},
			expect: func(t *testing.T, tasks []types.PersistentCacheTask, count int64, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(tasks)
				assert.Equal(int64(0), count)
			},
		},
		{
			name: "no tasks",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{}, 0)
			},
			expect: func(t *testing.T, tasks []types.PersistentCacheTask, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]types.PersistentCacheTask{}, tasks)
				assert.Equal(int64(0), count)
			},
		},
		{
			name: "deduplicates task ids and skips the bare prefix key",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{
					prefix,
					prefix + "task-1",
					prefix + "task-1:persistent-cache-peers",
					prefix + "task-2:persistent-peers",
				}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{})
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-2")).SetVal(mockRawPersistentCacheTask("task-2"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-2")).SetVal([]string{})
			},
			expect: func(t *testing.T, tasks []types.PersistentCacheTask, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(tasks, 2)
				assert.Equal("task-1", tasks[0].ID)
				assert.Equal("task-2", tasks[1].ID)
				assert.Equal(int64(2), count)
			},
		},
		{
			name: "task that fails to load is skipped",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "task-1", prefix + "task-2"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetErr(errors.New("load error"))
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-2")).SetVal(mockRawPersistentCacheTask("task-2"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-2")).SetVal([]string{})
			},
			expect: func(t *testing.T, tasks []types.PersistentCacheTask, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(tasks, 1)
				assert.Equal("task-2", tasks[0].ID)
				assert.Equal(int64(1), count)
			},
		},
		{
			name: "follows the scan cursor across pages",
			mock: func(mock redismock.ClientMock) {
				mock.ExpectScan(0, prefix+"*", 10).SetVal([]string{prefix + "task-1"}, 7)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-1")).SetVal(mockRawPersistentCacheTask("task-1"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-1")).SetVal([]string{})
				mock.ExpectScan(7, prefix+"*", 10).SetVal([]string{prefix + "task-2"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheTaskKeyInScheduler(1, "task-2")).SetVal(mockRawPersistentCacheTask("task-2"))
				mock.ExpectSMembers(pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(1, "task-2")).SetVal([]string{})
			},
			expect: func(t *testing.T, tasks []types.PersistentCacheTask, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(tasks, 2)
				assert.Equal(int64(2), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			rdb, mock := redismock.NewClientMock()
			tc.mock(mock)

			s := &service{rdb: rdb}
			tasks, count, err := s.GetPersistentCacheTasks(context.Background(), types.GetPersistentCacheTasksQuery{SchedulerClusterID: 1, Page: 1, PerPage: 10})
			tc.expect(t, tasks, count, err)
			assert.NoError(mock.ExpectationsWereMet())
		})
	}
}
