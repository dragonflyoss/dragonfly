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

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bits-and-blooms/bitset"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/types"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
)

// DestroyPersistentCacheTask deletes a persistent cache task from Redis based on query parameters.
func (s *service) DestroyPersistentCacheTask(ctx context.Context, schedulerClusterID uint, id string) error {
	return s.rdb.Del(ctx, pkgredis.MakePersistentCacheTaskKeyInScheduler(schedulerClusterID, id)).Err()
}

// GetPersistentCacheTask retrieves a persistent cache task from Redis based on query parameters.
func (s *service) GetPersistentCacheTask(ctx context.Context, schedulerClusterID uint, id string) (types.PersistentCacheTask, error) {
	return s.loadTask(ctx, schedulerClusterID, id)
}

// GetPersistentCacheTasks retrieves persistent cache tasks from Redis based on query parameters.
func (s *service) GetPersistentCacheTasks(ctx context.Context, q types.GetPersistentCacheTasksQuery) ([]types.PersistentCacheTask, int64, error) {
	tasks, err := s.loadAllTasks(ctx, q.SchedulerClusterID)
	if err != nil {
		return nil, 0, err
	}

	if len(tasks) == 0 {
		return []types.PersistentCacheTask{}, 0, nil
	}

	return tasks, int64(len(tasks)), nil
}

// loadAllTasks loads all persistent cache tasks from Redis based on the provided scheduler cluster ID.
func (s *service) loadAllTasks(ctx context.Context, schedulerClusterID uint) ([]types.PersistentCacheTask, error) {
	var (
		tasks  []types.PersistentCacheTask
		cursor uint64
	)

	for {
		var (
			taskKeys []string
			err      error
		)

		// Example: If {prefix} is "scheduler:scheduler-clusters:1:persistent-cache-tasks:".
		//
		// Keys include:
		// - "{prefix}{taskID}:persistent-cache-peers"
		// - "{prefix}{taskID}:persistent-peers"
		// - "{prefix}{taskID}"
		//
		// Purpose: Scan all keys starting with {prefix}.
		prefix := fmt.Sprintf("%s:", pkgredis.MakePersistentCacheTasksInScheduler(schedulerClusterID))
		taskKeys, cursor, err = s.rdb.Scan(ctx, cursor, fmt.Sprintf("%s*", prefix), 10).Result()
		if err != nil {
			logger.Error("scan tasks failed")
			return nil, err
		}

		taskIDs := make(map[string]struct{})
		for _, taskKey := range taskKeys {
			// If context is done, return error.
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			// Remove prefix from task key.
			suffix := strings.TrimPrefix(taskKey, prefix)
			if suffix == "" {
				logger.Error("invalid task key")
				continue
			}

			// Suffix is a non-empty string.
			//
			// Keys include:
			// - "{taskID}:persistent-cache-peers"
			// - "{taskID}:persistent-peers"
			// - "{taskID}"
			//
			// Purpose: Extract taskID from suffix and ensure no duplicate taskIDs.
			taskID := strings.Split(suffix, ":")[0]
			if _, ok := taskIDs[taskID]; ok {
				continue
			}
			taskIDs[taskID] = struct{}{}

			task, err := s.loadTask(ctx, schedulerClusterID, taskID)
			if err != nil {
				logger.WithTaskID(taskID).Error("load task failed")
				continue
			}

			tasks = append(tasks, task)
		}

		if cursor == 0 {
			break
		}
	}

	return tasks, nil
}

// loadTask loads a task from Redis based on the provided key.
func (s *service) loadTask(ctx context.Context, schedulerClusterID uint, id string) (types.PersistentCacheTask, error) {
	taskKey := pkgredis.MakePersistentCacheTaskKeyInScheduler(schedulerClusterID, id)
	rawTask, err := s.rdb.HGetAll(ctx, taskKey).Result()
	if err != nil {
		return types.PersistentCacheTask{}, err
	}

	if len(rawTask) == 0 {
		return types.PersistentCacheTask{}, errors.New("task not found")
	}

	task := types.PersistentCacheTask{
		ID:          rawTask["id"],
		Tag:         rawTask["tag"],
		Application: rawTask["application"],
		State:       rawTask["state"],
	}

	// Parse PersistentReplicaCount.
	persistentReplicaCount, err := strconv.ParseUint(rawTask["persistent_replica_count"], 10, 64)
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.PersistentReplicaCount = persistentReplicaCount

	// Parse PieceLength.
	pieceLength, err := strconv.ParseUint(rawTask["piece_length"], 10, 64)
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.PieceLength = pieceLength

	// Parse ContentLength.
	contentLength, err := strconv.ParseUint(rawTask["content_length"], 10, 64)
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.ContentLength = contentLength

	// Parse TotalPieceCount.
	totalPieceCount, err := strconv.ParseUint(rawTask["total_piece_count"], 10, 32)
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.TotalPieceCount = uint32(totalPieceCount)

	// Parse TTL.
	ttl, err := strconv.ParseInt(rawTask["ttl"], 10, 64)
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.TTL = time.Duration(ttl)

	// Parse CreatedAt.
	createdAt, err := time.Parse(time.RFC3339, rawTask["created_at"])
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.CreatedAt = createdAt

	// Parse UpdatedAt.
	updatedAt, err := time.Parse(time.RFC3339, rawTask["updated_at"])
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.UpdatedAt = updatedAt

	peers, err := s.loadAllPeersByTaskID(ctx, schedulerClusterID, task.ID)
	if err != nil {
		return types.PersistentCacheTask{}, err
	}
	task.Peers = peers
	return task, nil
}

// loadAllPeersByTaskID loads all peers associated with a task ID from Redis.
func (s *service) loadAllPeersByTaskID(ctx context.Context, schedulerClusterID uint, taskID string) ([]*types.PersistentCachePeer, error) {
	var peers []*types.PersistentCachePeer
	peerIDs, err := s.rdb.SMembers(ctx, pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(schedulerClusterID, taskID)).Result()
	if err != nil {
		return nil, err
	}

	for _, peerID := range peerIDs {
		peer, err := s.loadPeer(ctx, schedulerClusterID, peerID)
		if err != nil {
			logger.Warnf("load peer %s failed: %v", peerID, err)
			continue
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

// loadPeer loads a peer from Redis based on the provided key.
func (s *service) loadPeer(ctx context.Context, schedulerClusterID uint, id string) (*types.PersistentCachePeer, error) {
	peerKey := pkgredis.MakePersistentCachePeerKeyInScheduler(schedulerClusterID, id)
	rawPeer, err := s.rdb.HGetAll(ctx, peerKey).Result()
	if err != nil {
		return nil, err
	}

	if len(rawPeer) == 0 {
		return nil, errors.New("peer not found")
	}

	peer := &types.PersistentCachePeer{
		ID:    rawPeer["id"],
		State: rawPeer["state"],
	}

	// Parse Persistent.
	persistent, err := strconv.ParseBool(rawPeer["persistent"])
	if err != nil {
		return nil, err
	}
	peer.Persistent = persistent

	// Parse FinishedPieces.
	finishedPieces := &bitset.BitSet{}
	if err := finishedPieces.UnmarshalBinary([]byte(rawPeer["finished_pieces"])); err != nil {
		return nil, err
	}
	peer.FinishedPieces = finishedPieces

	// Parse BlockParents.
	blockParents := []string{}
	if err := json.Unmarshal([]byte(rawPeer["block_parents"]), &blockParents); err != nil {
		return nil, err
	}
	peer.BlockParents = blockParents

	// Parse Cost.
	cost, err := strconv.ParseInt(rawPeer["cost"], 10, 64)
	if err != nil {
		return nil, err
	}
	peer.Cost = time.Duration(cost)

	// Parse CreatedAt.
	createdAt, err := time.Parse(time.RFC3339, rawPeer["created_at"])
	if err != nil {
		return nil, err
	}
	peer.CreatedAt = createdAt

	// Parse UpdatedAt.
	updatedAt, err := time.Parse(time.RFC3339, rawPeer["updated_at"])
	if err != nil {
		return nil, err
	}
	peer.UpdatedAt = updatedAt

	host, err := s.loadHost(ctx, schedulerClusterID, rawPeer["host_id"])
	if err != nil {
		return nil, err
	}
	peer.Host = host
	return peer, nil
}

// loadHost loads a host from Redis based on the provided key.
func (s *service) loadHost(ctx context.Context, schedulerClusterID uint, id string) (*types.PersistentCacheHost, error) {
	hostKey := pkgredis.MakePersistentCacheHostKeyInScheduler(schedulerClusterID, id)
	rawHost, err := s.rdb.HGetAll(ctx, hostKey).Result()
	if err != nil {
		return nil, err
	}

	if len(rawHost) == 0 {
		return nil, errors.New("host not found")
	}

	host := &types.PersistentCacheHost{
		ID:                 rawHost["id"],
		Type:               rawHost["type"],
		Hostname:           rawHost["hostname"],
		IP:                 rawHost["ip"],
		OS:                 rawHost["os"],
		Platform:           rawHost["platform"],
		PlatformFamily:     rawHost["platform_family"],
		PlatformVersion:    rawHost["platform_version"],
		KernelVersion:      rawHost["kernel_version"],
		SchedulerClusterID: schedulerClusterID,
	}

	// Parse integer fields.
	port, err := strconv.ParseInt(rawHost["port"], 10, 32)
	if err != nil {
		return nil, err
	}
	host.Port = int32(port)

	// Parse download port.
	downloadPort, err := strconv.ParseInt(rawHost["download_port"], 10, 32)
	if err != nil {
		return nil, err
	}
	host.DownloadPort = int32(downloadPort)

	// Parse boolean fields.
	disableShared, err := strconv.ParseBool(rawHost["disable_shared"])
	if err != nil {
		return nil, err
	}
	host.DisableShared = disableShared

	// Parse CPU information.
	cpuLogicalCount, err := strconv.ParseUint(rawHost["cpu_logical_count"], 10, 32)
	if err != nil {
		return nil, err
	}
	host.CPU.LogicalCount = uint32(cpuLogicalCount)

	cpuPhysicalCount, err := strconv.ParseUint(rawHost["cpu_physical_count"], 10, 32)
	if err != nil {
		return nil, err
	}
	host.CPU.PhysicalCount = uint32(cpuPhysicalCount)

	cpuPercent, err := strconv.ParseFloat(rawHost["cpu_percent"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Percent = cpuPercent

	cpuProcessPercent, err := strconv.ParseFloat(rawHost["cpu_process_percent"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.ProcessPercent = cpuProcessPercent

	// Parse CPU Times.
	cpuTimesUser, err := strconv.ParseFloat(rawHost["cpu_times_user"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.User = cpuTimesUser

	cpuTimesSystem, err := strconv.ParseFloat(rawHost["cpu_times_system"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.System = cpuTimesSystem

	cpuTimesIdle, err := strconv.ParseFloat(rawHost["cpu_times_idle"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Idle = cpuTimesIdle

	cpuTimesNice, err := strconv.ParseFloat(rawHost["cpu_times_nice"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Nice = cpuTimesNice

	cpuTimesIowait, err := strconv.ParseFloat(rawHost["cpu_times_iowait"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Iowait = cpuTimesIowait

	cpuTimesIrq, err := strconv.ParseFloat(rawHost["cpu_times_irq"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Irq = cpuTimesIrq

	cpuTimesSoftirq, err := strconv.ParseFloat(rawHost["cpu_times_softirq"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Softirq = cpuTimesSoftirq

	cpuTimesSteal, err := strconv.ParseFloat(rawHost["cpu_times_steal"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Steal = cpuTimesSteal

	cpuTimesGuest, err := strconv.ParseFloat(rawHost["cpu_times_guest"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.Guest = cpuTimesGuest

	cpuTimesGuestNice, err := strconv.ParseFloat(rawHost["cpu_times_guest_nice"], 64)
	if err != nil {
		return nil, err
	}
	host.CPU.Times.GuestNice = cpuTimesGuestNice

	// Parse Memory information.
	memoryTotal, err := strconv.ParseUint(rawHost["memory_total"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Memory.Total = memoryTotal

	memoryAvailable, err := strconv.ParseUint(rawHost["memory_available"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Memory.Available = memoryAvailable

	memoryUsed, err := strconv.ParseUint(rawHost["memory_used"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Memory.Used = memoryUsed

	memoryUsedPercent, err := strconv.ParseFloat(rawHost["memory_used_percent"], 64)
	if err != nil {
		return nil, err
	}
	host.Memory.UsedPercent = memoryUsedPercent

	memoryProcessUsedPercent, err := strconv.ParseFloat(rawHost["memory_process_used_percent"], 64)
	if err != nil {
		return nil, err
	}
	host.Memory.ProcessUsedPercent = memoryProcessUsedPercent

	memoryFree, err := strconv.ParseUint(rawHost["memory_free"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Memory.Free = memoryFree

	// Parse Network information.
	tcpConnectionCount, err := strconv.ParseUint(rawHost["tcp_connection_count"], 10, 32)
	if err != nil {
		return nil, err
	}
	host.Network.TCPConnectionCount = uint32(tcpConnectionCount)

	uploadTCPConnectionCount, err := strconv.ParseUint(rawHost["upload_tcp_connection_count"], 10, 32)
	if err != nil {
		return nil, err
	}
	host.Network.UploadTCPConnectionCount = uint32(uploadTCPConnectionCount)
	host.Network.Location = rawHost["location"]
	host.Network.IDC = rawHost["idc"]

	rxBandwidth, err := strconv.ParseUint(rawHost["rx_bandwidth"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Network.RxBandwidth = rxBandwidth

	maxRxBandwidth, err := strconv.ParseUint(rawHost["max_rx_bandwidth"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Network.MaxRxBandwidth = maxRxBandwidth

	txBandwidth, err := strconv.ParseUint(rawHost["tx_bandwidth"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Network.TxBandwidth = txBandwidth

	maxTxBandwidth, err := strconv.ParseUint(rawHost["max_tx_bandwidth"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Network.MaxTxBandwidth = maxTxBandwidth

	// Parse Disk information.
	diskTotal, err := strconv.ParseUint(rawHost["disk_total"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.Total = diskTotal

	diskFree, err := strconv.ParseUint(rawHost["disk_free"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.Free = diskFree

	diskUsed, err := strconv.ParseUint(rawHost["disk_used"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.Used = diskUsed

	diskUsedPercent, err := strconv.ParseFloat(rawHost["disk_used_percent"], 64)
	if err != nil {
		return nil, err
	}
	host.Disk.UsedPercent = diskUsedPercent

	diskInodesTotal, err := strconv.ParseUint(rawHost["disk_inodes_total"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.InodesTotal = diskInodesTotal

	diskInodesUsed, err := strconv.ParseUint(rawHost["disk_inodes_used"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.InodesUsed = diskInodesUsed

	diskInodesFree, err := strconv.ParseUint(rawHost["disk_inodes_free"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.InodesFree = diskInodesFree

	diskInodesUsedPercent, err := strconv.ParseFloat(rawHost["disk_inodes_used_percent"], 64)
	if err != nil {
		return nil, err
	}
	host.Disk.InodesUsedPercent = diskInodesUsedPercent

	diskWriteBandwidth, err := strconv.ParseUint(rawHost["disk_write_bandwidth"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.WriteBandwidth = diskWriteBandwidth

	diskReadBandwidth, err := strconv.ParseUint(rawHost["disk_read_bandwidth"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.Disk.ReadBandwidth = diskReadBandwidth

	// Parse Build information.
	host.Build.GitVersion = rawHost["build_git_version"]
	host.Build.GitCommit = rawHost["build_git_commit"]
	host.Build.GoVersion = rawHost["build_go_version"]
	host.Build.Platform = rawHost["build_platform"]

	// Parse AnnounceInterval.
	announceInterval, err := strconv.ParseInt(rawHost["announce_interval"], 10, 64)
	if err != nil {
		return nil, err
	}
	host.AnnounceInterval = time.Duration(announceInterval)

	// Parse CreatedAt.
	createdAt, err := time.Parse(time.RFC3339, rawHost["created_at"])
	if err != nil {
		return nil, err
	}
	host.CreatedAt = createdAt

	// Parse UpdatedAt.
	updatedAt, err := time.Parse(time.RFC3339, rawHost["updated_at"])
	if err != nil {
		return nil, err
	}
	host.UpdatedAt = updatedAt
	return host, nil
}
