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

// GetPersistentCaches retrieves persistent cache tasks from Redis based on query parameters.
func (s *service) GetPersistentCaches(ctx context.Context, q types.GetPersistentCachesQuery) ([]types.GetPersistentCacheResponse, int64, error) {
	var responses []types.GetPersistentCacheResponse

	// Get all scheduler cluster IDs if none specified.
	schedulerClusterIDs := q.SchedulerClusterIDs
	if len(schedulerClusterIDs) == 0 {
		// Get all available scheduler cluster IDs.
		var cursor uint64
		for {
			prefix := fmt.Sprintf("%s:", pkgredis.MakeNamespaceKeyInScheduler(pkgredis.SchedulerClustersNamespace))
			keys, cursor, err := s.rdb.Scan(ctx, cursor, fmt.Sprintf("%s*", prefix), 10).Result()
			if err != nil {
				logger.Errorf("scan scheduler clusterIDs failed: %v", err)
				return nil, 0, err
			}

			for _, key := range keys {
				// If context is done, return error.
				if err := ctx.Err(); err != nil {
					return nil, 0, err
				}

				// Remove prefix from key.
				suffix := strings.TrimPrefix(key, prefix)
				if suffix == "" {
					logger.Error("invalid key")
					continue
				}

				// Extract scheduler clusterID from suffix.
				if !strings.ContainsRune(suffix, ':') {
					schedulerClusterID, err := strconv.ParseUint(suffix, 10, 32)
					if err != nil {
						logger.Errorf("parse scheduler clusterID failed: %v", err)
						continue
					}
					schedulerClusterIDs = append(schedulerClusterIDs, uint(schedulerClusterID))
				}
			}

			if cursor == 0 {
				break
			}
		}

		if len(schedulerClusterIDs) == 0 {
			return nil, 0, errors.New("no scheduler cluster found")
		}
	}

	// Collect all task IDs from all specified clusters.
	var allTaskKeys []string
	for _, schedulerClusterID := range schedulerClusterIDs {
		// Get all task keys in the cluster.
		taskKeys, err := s.rdb.Keys(ctx, pkgredis.MakePersistentCacheTaskKeyInScheduler(schedulerClusterID, "*")).Result()
		if err != nil {
			logger.Errorf("get task keys failed: %v", err)
			continue
		}
		allTaskKeys = append(allTaskKeys, taskKeys...)
	}

	// Calculate total count and pagination.
	totalCount := int64(len(allTaskKeys))
	startIndex := (q.Page - 1) * q.PerPage
	endIndex := startIndex + q.PerPage

	if startIndex >= int(totalCount) {
		return []types.GetPersistentCacheResponse{}, totalCount, nil
	}

	if endIndex > int(totalCount) {
		endIndex = int(totalCount)
	}

	// Get paginated task IDs and keys.
	taskKeys := allTaskKeys[startIndex:endIndex]

	// Process each task.
	for _, taskKey := range taskKeys {
		schedulerClusterID, err := pkgredis.ExtractSchedulerClusterIDFromPersistentCacheTaskKey(taskKey)
		if err != nil {
			logger.Warnf("extract scheduler cluster ID from persistent cache task key %s failed: %v", taskKey, err)
			continue
		}

		// Get task data from Redis.
		rawTask, err := s.rdb.HGetAll(ctx, taskKey).Result()
		if err != nil {
			logger.Warnf("getting task %s failed from redis: %v", taskKey, err)
			continue
		}

		if len(rawTask) == 0 {
			logger.Warnf("task %s not found in redis", taskKey)
			continue
		}

		// Parse task data.
		response, err := s.parseTaskData(ctx, rawTask)
		if err != nil {
			logger.Warnf("parse task %s data failed: %v", taskKey, err)
			continue
		}

		// Load peers for this task.
		peers, err := s.loadPeersForTask(ctx, taskKey, schedulerClusterID)
		if err != nil {
			logger.Warnf("load peers for task %s failed: %v", taskKey, err)
		}

		response.PersistentCachePeers = peers
		responses = append(responses, response)
	}

	return responses, totalCount, nil
}

// parseTaskData parses raw task data from Redis and returns a GetPersistentCacheResponse.
func (s *service) parseTaskData(ctx context.Context, rawTask map[string]string) (types.GetPersistentCacheResponse, error) {
	var response types.GetPersistentCacheResponse

	// Set task ID.
	response.TaskID = rawTask["id"]

	// Parse PersistentReplicaCount.
	if persistentReplicaCount, err := strconv.ParseUint(rawTask["persistent_replica_count"], 10, 64); err == nil {
		response.PersistentReplicaCount = persistentReplicaCount
	} else {
		logger.Warnf("parsing persistent replica count failed: %v", err)
	}

	// Set Tag and Application.
	response.Tag = rawTask["tag"]
	response.Application = rawTask["application"]

	// Parse PieceLength.
	if pieceLength, err := strconv.ParseUint(rawTask["piece_length"], 10, 64); err == nil {
		response.PieceLength = pieceLength
	} else {
		logger.Warnf("parsing piece length failed: %v", err)
	}

	// Parse ContentLength.
	if contentLength, err := strconv.ParseUint(rawTask["content_length"], 10, 64); err == nil {
		response.ContentLength = contentLength
	} else {
		logger.Warnf("parsing content length failed: %v", err)
	}

	// Parse TotalPieceCount.
	if totalPieceCount, err := strconv.ParseUint(rawTask["total_piece_count"], 10, 32); err == nil {
		response.TotalPieceCount = uint32(totalPieceCount)
	} else {
		logger.Warnf("parsing total piece count failed: %v", err)
	}

	// Set State.
	response.State = rawTask["state"]

	// Parse TTL.
	if ttl, err := strconv.ParseInt(rawTask["ttl"], 10, 64); err == nil {
		response.TTL = time.Duration(ttl)
	} else {
		logger.Warnf("parsing ttl failed: %v", err)
	}

	// Parse CreatedAt.
	if createdAt, err := time.Parse(time.RFC3339, rawTask["created_at"]); err == nil {
		response.CreatedAt = createdAt
	} else {
		logger.Warnf("parsing created at failed: %v", err)
	}

	// Parse UpdatedAt.
	if updatedAt, err := time.Parse(time.RFC3339, rawTask["updated_at"]); err == nil {
		response.UpdatedAt = updatedAt
	} else {
		logger.Warnf("parsing updated at failed: %v", err)
	}

	return response, nil
}

// loadPeersForTask loads peers for a specific task from Redis.
func (s *service) loadPeersForTask(ctx context.Context, taskID string, schedulerClusterID uint) ([]types.PersistentCachePeer, error) {
	var peers []types.PersistentCachePeer

	// Get peer IDs for the task.
	peerIDs, err := s.rdb.SMembers(ctx, pkgredis.MakePersistentCachePeersOfPersistentCacheTaskInScheduler(schedulerClusterID, taskID)).Result()
	if err != nil {
		return nil, err
	}

	// Process each peer.
	for _, peerID := range peerIDs {
		// Get peer data from Redis.
		rawPeer, err := s.rdb.HGetAll(ctx, pkgredis.MakePersistentCachePeerKeyInScheduler(schedulerClusterID, peerID)).Result()
		if err != nil {
			logger.Warnf("getting peer %s failed from redis: %v", peerID, err)
			continue
		}

		if len(rawPeer) == 0 {
			logger.Warnf("peer %s not found in redis", peerID)
			continue
		}

		// Parse peer data.
		peer, err := s.parsePeerData(ctx, rawPeer, schedulerClusterID)
		if err != nil {
			logger.Warnf("parse peer %s data failed: %v", peerID, err)
			continue
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

// parsePeerData parses raw peer data from Redis and returns a PersistentCachePeer.
func (s *service) parsePeerData(ctx context.Context, rawPeer map[string]string, schedulerClusterID uint) (types.PersistentCachePeer, error) {
	var peer types.PersistentCachePeer

	// Set ID.
	peer.ID = rawPeer["id"]

	// Parse Persistent.
	if persistent, err := strconv.ParseBool(rawPeer["persistent"]); err == nil {
		peer.Persistent = persistent
	} else {
		logger.Warnf("parsing persistent failed: %v", err)
	}

	// Parse FinishedPieces.
	finishedPieces := &bitset.BitSet{}
	if err := finishedPieces.UnmarshalBinary([]byte(rawPeer["finished_pieces"])); err == nil {
		peer.FinishedPieces = finishedPieces
	} else {
		logger.Warnf("unmarshal finished pieces failed: %v", err)
	}

	// Set State.
	peer.State = rawPeer["state"]

	// Parse BlockParents.
	blockParents := []string{}
	if err := json.Unmarshal([]byte(rawPeer["block_parents"]), &blockParents); err == nil {
		peer.BlockParents = blockParents
	} else {
		logger.Warnf("unmarshal block parents failed: %v", err)
	}

	// Parse Cost.
	if cost, err := strconv.ParseInt(rawPeer["cost"], 10, 64); err == nil {
		peer.Cost = time.Duration(cost)
	} else {
		logger.Warnf("parsing cost failed: %v", err)
	}

	// Parse CreatedAt.
	if createdAt, err := time.Parse(time.RFC3339, rawPeer["created_at"]); err == nil {
		peer.CreatedAt = createdAt
	} else {
		logger.Warnf("parsing created at failed: %v", err)
	}

	// Parse UpdatedAt.
	if updatedAt, err := time.Parse(time.RFC3339, rawPeer["updated_at"]); err == nil {
		peer.UpdatedAt = updatedAt
	} else {
		logger.Warnf("parsing updated at failed: %v", err)
	}

	// Load host data if host_id is available.
	if hostID, ok := rawPeer["host_id"]; ok && hostID != "" {
		host, err := s.loadHostData(ctx, hostID, schedulerClusterID)
		if err != nil {
			logger.Warnf("load host %s data failed: %v", hostID, err)
		} else {
			peer.Host = host
		}
	}

	return peer, nil
}

// loadHostData loads host data from Redis.
func (s *service) loadHostData(ctx context.Context, hostID string, schedulerClusterID uint) (types.PersistentCachePeerHost, error) {
	var host types.PersistentCachePeerHost

	// Get host data from Redis.
	rawHost, err := s.rdb.HGetAll(ctx, pkgredis.MakePersistentCacheHostKeyInScheduler(schedulerClusterID, hostID)).Result()
	if err != nil {
		return host, err
	}

	if len(rawHost) == 0 {
		return host, nil
	}

	// Set basic host information.
	host.ID = rawHost["id"]
	host.Type = rawHost["type"]
	host.Hostname = rawHost["hostname"]
	host.IP = rawHost["ip"]
	host.OS = rawHost["os"]
	host.Platform = rawHost["platform"]
	host.PlatformFamily = rawHost["platform_family"]
	host.PlatformVersion = rawHost["platform_version"]
	host.KernelVersion = rawHost["kernel_version"]

	// Parse integer fields.
	if port, err := strconv.ParseInt(rawHost["port"], 10, 32); err == nil {
		host.Port = int32(port)
	} else {
		logger.Warnf("parsing port failed: %v", err)
	}

	if downloadPort, err := strconv.ParseInt(rawHost["download_port"], 10, 32); err == nil {
		host.DownloadPort = int32(downloadPort)
	} else {
		logger.Warnf("parsing download port failed: %v", err)
	}

	// Parse boolean fields.
	if disableShared, err := strconv.ParseBool(rawHost["disable_shared"]); err == nil {
		host.DisableShared = disableShared
	} else {
		logger.Warnf("parsing disable shared failed: %v", err)
	}

	// Parse CPU information.
	if cpuLogicalCount, err := strconv.ParseUint(rawHost["cpu_logical_count"], 10, 32); err == nil {
		host.CPU.LogicalCount = uint32(cpuLogicalCount)
	} else {
		logger.Warnf("parsing cpu logical count failed: %v", err)
	}

	if cpuPhysicalCount, err := strconv.ParseUint(rawHost["cpu_physical_count"], 10, 32); err == nil {
		host.CPU.PhysicalCount = uint32(cpuPhysicalCount)
	} else {
		logger.Warnf("parsing cpu physical count failed: %v", err)
	}

	if cpuPercent, err := strconv.ParseFloat(rawHost["cpu_percent"], 64); err == nil {
		host.CPU.Percent = cpuPercent
	} else {
		logger.Warnf("parsing cpu percent failed: %v", err)
	}

	if cpuProcessPercent, err := strconv.ParseFloat(rawHost["cpu_process_percent"], 64); err == nil {
		host.CPU.ProcessPercent = cpuProcessPercent
	} else {
		logger.Warnf("parsing cpu process percent failed: %v", err)
	}

	// Parse CPU Times.
	if cpuTimesUser, err := strconv.ParseFloat(rawHost["cpu_times_user"], 64); err == nil {
		host.CPU.Times.User = cpuTimesUser
	} else {
		logger.Warnf("parsing cpu times user failed: %v", err)
	}

	if cpuTimesSystem, err := strconv.ParseFloat(rawHost["cpu_times_system"], 64); err == nil {
		host.CPU.Times.System = cpuTimesSystem
	} else {
		logger.Warnf("parsing cpu times system failed: %v", err)
	}

	if cpuTimesIdle, err := strconv.ParseFloat(rawHost["cpu_times_idle"], 64); err == nil {
		host.CPU.Times.Idle = cpuTimesIdle
	} else {
		logger.Warnf("parsing cpu times idle failed: %v", err)
	}

	if cpuTimesNice, err := strconv.ParseFloat(rawHost["cpu_times_nice"], 64); err == nil {
		host.CPU.Times.Nice = cpuTimesNice
	} else {
		logger.Warnf("parsing cpu times nice failed: %v", err)
	}

	if cpuTimesIowait, err := strconv.ParseFloat(rawHost["cpu_times_iowait"], 64); err == nil {
		host.CPU.Times.Iowait = cpuTimesIowait
	} else {
		logger.Warnf("parsing cpu times iowait failed: %v", err)
	}

	if cpuTimesIrq, err := strconv.ParseFloat(rawHost["cpu_times_irq"], 64); err == nil {
		host.CPU.Times.Irq = cpuTimesIrq
	} else {
		logger.Warnf("parsing cpu times irq failed: %v", err)
	}

	if cpuTimesSoftirq, err := strconv.ParseFloat(rawHost["cpu_times_softirq"], 64); err == nil {
		host.CPU.Times.Softirq = cpuTimesSoftirq
	} else {
		logger.Warnf("parsing cpu times softirq failed: %v", err)
	}

	if cpuTimesSteal, err := strconv.ParseFloat(rawHost["cpu_times_steal"], 64); err == nil {
		host.CPU.Times.Steal = cpuTimesSteal
	} else {
		logger.Warnf("parsing cpu times steal failed: %v", err)
	}

	if cpuTimesGuest, err := strconv.ParseFloat(rawHost["cpu_times_guest"], 64); err == nil {
		host.CPU.Times.Guest = cpuTimesGuest
	} else {
		logger.Warnf("parsing cpu times guest failed: %v", err)
	}

	if cpuTimesGuestNice, err := strconv.ParseFloat(rawHost["cpu_times_guest_nice"], 64); err == nil {
		host.CPU.Times.GuestNice = cpuTimesGuestNice
	} else {
		logger.Warnf("parsing cpu times guest nice failed: %v", err)
	}

	// Parse Memory information.
	if memoryTotal, err := strconv.ParseUint(rawHost["memory_total"], 10, 64); err == nil {
		host.Memory.Total = memoryTotal
	} else {
		logger.Warnf("parsing memory total failed: %v", err)
	}

	if memoryAvailable, err := strconv.ParseUint(rawHost["memory_available"], 10, 64); err == nil {
		host.Memory.Available = memoryAvailable
	} else {
		logger.Warnf("parsing memory available failed: %v", err)
	}

	if memoryUsed, err := strconv.ParseUint(rawHost["memory_used"], 10, 64); err == nil {
		host.Memory.Used = memoryUsed
	} else {
		logger.Warnf("parsing memory used failed: %v", err)
	}

	if memoryUsedPercent, err := strconv.ParseFloat(rawHost["memory_used_percent"], 64); err == nil {
		host.Memory.UsedPercent = memoryUsedPercent
	} else {
		logger.Warnf("parsing memory used percent failed: %v", err)
	}

	if memoryProcessUsedPercent, err := strconv.ParseFloat(rawHost["memory_process_used_percent"], 64); err == nil {
		host.Memory.ProcessUsedPercent = memoryProcessUsedPercent
	} else {
		logger.Warnf("parsing memory process used percent failed: %v", err)
	}

	if memoryFree, err := strconv.ParseUint(rawHost["memory_free"], 10, 64); err == nil {
		host.Memory.Free = memoryFree
	} else {
		logger.Warnf("parsing memory free failed: %v", err)
	}

	// Parse Network information.
	if tcpConnectionCount, err := strconv.ParseUint(rawHost["tcp_connection_count"], 10, 32); err == nil {
		host.Network.TCPConnectionCount = uint32(tcpConnectionCount)
	} else {
		logger.Warnf("parsing tcp connection count failed: %v", err)
	}

	if uploadTCPConnectionCount, err := strconv.ParseUint(rawHost["upload_tcp_connection_count"], 10, 32); err == nil {
		host.Network.UploadTCPConnectionCount = uint32(uploadTCPConnectionCount)
	} else {
		logger.Warnf("parsing upload tcp connection count failed: %v", err)
	}

	host.Network.Location = rawHost["location"]
	host.Network.IDC = rawHost["idc"]

	if downloadRate, err := strconv.ParseUint(rawHost["download_rate"], 10, 64); err == nil {
		host.Network.DownloadRate = downloadRate
	} else {
		logger.Warnf("parsing download rate failed: %v", err)
	}

	if downloadRateLimit, err := strconv.ParseUint(rawHost["download_rate_limit"], 10, 64); err == nil {
		host.Network.DownloadRateLimit = downloadRateLimit
	} else {
		logger.Warnf("parsing download rate limit failed: %v", err)
	}

	if uploadRate, err := strconv.ParseUint(rawHost["upload_rate"], 10, 64); err == nil {
		host.Network.UploadRate = uploadRate
	} else {
		logger.Warnf("parsing upload rate failed: %v", err)
	}

	if uploadRateLimit, err := strconv.ParseUint(rawHost["upload_rate_limit"], 10, 64); err == nil {
		host.Network.UploadRateLimit = uploadRateLimit
	} else {
		logger.Warnf("parsing upload rate limit failed: %v", err)
	}

	// Parse Disk information.
	if diskTotal, err := strconv.ParseUint(rawHost["disk_total"], 10, 64); err == nil {
		host.Disk.Total = diskTotal
	} else {
		logger.Warnf("parsing disk total failed: %v", err)
	}

	if diskFree, err := strconv.ParseUint(rawHost["disk_free"], 10, 64); err == nil {
		host.Disk.Free = diskFree
	} else {
		logger.Warnf("parsing disk free failed: %v", err)
	}

	if diskUsed, err := strconv.ParseUint(rawHost["disk_used"], 10, 64); err == nil {
		host.Disk.Used = diskUsed
	} else {
		logger.Warnf("parsing disk used failed: %v", err)
	}

	if diskUsedPercent, err := strconv.ParseFloat(rawHost["disk_used_percent"], 64); err == nil {
		host.Disk.UsedPercent = diskUsedPercent
	} else {
		logger.Warnf("parsing disk used percent failed: %v", err)
	}

	if diskInodesTotal, err := strconv.ParseUint(rawHost["disk_inodes_total"], 10, 64); err == nil {
		host.Disk.InodesTotal = diskInodesTotal
	} else {
		logger.Warnf("parsing disk inodes total failed: %v", err)
	}

	if diskInodesUsed, err := strconv.ParseUint(rawHost["disk_inodes_used"], 10, 64); err == nil {
		host.Disk.InodesUsed = diskInodesUsed
	} else {
		logger.Warnf("parsing disk inodes used failed: %v", err)
	}

	if diskInodesFree, err := strconv.ParseUint(rawHost["disk_inodes_free"], 10, 64); err == nil {
		host.Disk.InodesFree = diskInodesFree
	} else {
		logger.Warnf("parsing disk inodes free failed: %v", err)
	}

	if diskInodesUsedPercent, err := strconv.ParseFloat(rawHost["disk_inodes_used_percent"], 64); err == nil {
		host.Disk.InodesUsedPercent = diskInodesUsedPercent
	} else {
		logger.Warnf("parsing disk inodes used percent failed: %v", err)
	}

	if diskWriteBandwidth, err := strconv.ParseUint(rawHost["disk_write_bandwidth"], 10, 64); err == nil {
		host.Disk.WriteBandwidth = diskWriteBandwidth
	} else {
		logger.Warnf("parsing disk write bandwidth failed: %v", err)
	}

	if diskReadBandwidth, err := strconv.ParseUint(rawHost["disk_read_bandwidth"], 10, 64); err == nil {
		host.Disk.ReadBandwidth = diskReadBandwidth
	} else {
		logger.Warnf("parsing disk read bandwidth failed: %v", err)
	}

	// Parse Build information.
	host.Build.GitVersion = rawHost["build_git_version"]
	host.Build.GitCommit = rawHost["build_git_commit"]
	host.Build.GoVersion = rawHost["build_go_version"]
	host.Build.Platform = rawHost["build_platform"]

	// Parse SchedulerClusterID.
	if schedulerClusterID, err := strconv.ParseUint(rawHost["scheduler_cluster_id"], 10, 64); err == nil {
		host.SchedulerClusterID = schedulerClusterID
	} else {
		logger.Warnf("parsing scheduler cluster id failed: %v", err)
	}

	// Parse AnnounceInterval.
	if announceInterval, err := strconv.ParseInt(rawHost["announce_interval"], 10, 64); err == nil {
		host.AnnounceInterval = time.Duration(announceInterval)
	} else {
		logger.Warnf("parsing announce interval failed: %v", err)
	}

	// Parse CreatedAt.
	if createdAt, err := time.Parse(time.RFC3339, rawHost["created_at"]); err == nil {
		host.CreatedAt = createdAt
	} else {
		logger.Warnf("parsing created at failed: %v", err)
	}

	// Parse UpdatedAt.
	if updatedAt, err := time.Parse(time.RFC3339, rawHost["updated_at"]); err == nil {
		host.UpdatedAt = updatedAt
	} else {
		logger.Warnf("parsing updated at failed: %v", err)
	}

	return host, nil
}
