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

//go:generate mockgen -destination host_manager_mock.go -source host_manager.go -package persistent

package persistent

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/container/set"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
)

const (
	// GC persistent host id.
	GCHostID = "persistent-host"
)

// HostManager is the interface used for host manager.
type HostManager interface {
	// Load returns host by a key.
	Load(context.Context, string) (*Host, bool)

	// Store sets host.
	Store(context.Context, *Host) error

	// Delete deletes host by a key.
	Delete(context.Context, string) error

	// LoadAll returns all hosts.
	LoadAll(context.Context) ([]*Host, error)

	// LoadRandom loads host randomly through the set of redis.
	LoadRandom(context.Context, int, set.SafeSet[string]) ([]*Host, error)

	// RunGC runs garbage collection.
	RunGC(context.Context) error
}

// hostManager contains content for host manager.
type hostManager struct {
	// Config is scheduler config.
	config *config.Config

	// Redis universal client interface.
	rdb redis.UniversalClient
}

// New host manager interface.
func newHostManager(cfg *config.Config, gc pkggc.GC, rdb redis.UniversalClient) (HostManager, error) {
	h := &hostManager{config: cfg, rdb: rdb}

	if err := gc.Add(pkggc.Task{
		ID:       GCHostID,
		Interval: cfg.Scheduler.GC.HostGCInterval,
		Timeout:  cfg.Scheduler.GC.HostGCInterval,
		Runner:   h,
	}); err != nil {
		return nil, err
	}

	return h, nil
}

// Load returns host by a key.
func (h *hostManager) Load(ctx context.Context, hostID string) (*Host, bool) {
	log := logger.WithHostID(hostID)
	rawHost, err := h.rdb.HGetAll(ctx, pkgredis.MakePersistentCacheHostKeyInScheduler(h.config.Manager.SchedulerClusterID, hostID)).Result()
	if err != nil {
		log.Errorf("getting host failed from redis: %v", err)
		return nil, false
	}

	if len(rawHost) == 0 {
		return nil, false
	}

	// Set integer fields from raw host.
	port, err := strconv.ParseInt(rawHost["port"], 10, 32)
	if err != nil {
		log.Errorf("parsing port failed: %v", err)
		return nil, false
	}

	downloadPort, err := strconv.ParseInt(rawHost["download_port"], 10, 32)
	if err != nil {
		log.Errorf("parsing download port failed: %v", err)
		return nil, false
	}

	proxyPort, err := strconv.ParseInt(rawHost["proxy_port"], 10, 32)
	if err != nil {
		log.Errorf("parsing proxy port failed: %v", err)
		return nil, false
	}

	// Set cpu fields from raw host.
	schedulerClusterID, err := strconv.ParseUint(rawHost["scheduler_cluster_id"], 10, 64)
	if err != nil {
		log.Errorf("parsing scheduler cluster id failed: %v", err)
		return nil, false
	}

	// Set boolean fields from raw host.
	disableShared, err := strconv.ParseBool(rawHost["disable_shared"])
	if err != nil {
		log.Errorf("parsing disable shared failed: %v", err)
		return nil, false
	}

	// Set cpu fields from raw host.
	cpuLogicalCount, err := strconv.ParseUint(rawHost["cpu_logical_count"], 10, 32)
	if err != nil {
		log.Errorf("parsing cpu logical count failed: %v", err)
		return nil, false
	}

	cpuPhysicalCount, err := strconv.ParseUint(rawHost["cpu_physical_count"], 10, 32)
	if err != nil {
		log.Errorf("parsing cpu physical count failed: %v", err)
		return nil, false
	}

	cpuPercent, err := strconv.ParseFloat(rawHost["cpu_percent"], 64)
	if err != nil {
		log.Errorf("parsing cpu percent failed: %v", err)
		return nil, false
	}

	cpuProcessPercent, err := strconv.ParseFloat(rawHost["cpu_processe_percent"], 64)
	if err != nil {
		log.Errorf("parsing cpu process percent failed: %v", err)
		return nil, false
	}

	cpuTimesUser, err := strconv.ParseFloat(rawHost["cpu_times_user"], 64)
	if err != nil {
		log.Errorf("parsing cpu times user failed: %v", err)
		return nil, false
	}

	cpuTimesSystem, err := strconv.ParseFloat(rawHost["cpu_times_system"], 64)
	if err != nil {
		log.Errorf("parsing cpu times system failed: %v", err)
		return nil, false
	}

	cpuTimesIdle, err := strconv.ParseFloat(rawHost["cpu_times_idle"], 64)
	if err != nil {
		log.Errorf("parsing cpu times idle failed: %v", err)
		return nil, false
	}

	cpuTimesNice, err := strconv.ParseFloat(rawHost["cpu_times_nice"], 64)
	if err != nil {
		log.Errorf("parsing cpu times nice failed: %v", err)
		return nil, false
	}

	cpuTimesIowait, err := strconv.ParseFloat(rawHost["cpu_times_iowait"], 64)
	if err != nil {
		log.Errorf("parsing cpu times iowait failed: %v", err)
		return nil, false
	}

	cpuTimesIrq, err := strconv.ParseFloat(rawHost["cpu_times_irq"], 64)
	if err != nil {
		log.Errorf("parsing cpu times irq failed: %v", err)
		return nil, false
	}

	cpuTimesSoftirq, err := strconv.ParseFloat(rawHost["cpu_times_softirq"], 64)
	if err != nil {
		log.Errorf("parsing cpu times softirq failed: %v", err)
		return nil, false
	}

	cpuTimesSteal, err := strconv.ParseFloat(rawHost["cpu_times_steal"], 64)
	if err != nil {
		log.Errorf("parsing cpu times steal failed: %v", err)
		return nil, false
	}

	cpuTimesGuest, err := strconv.ParseFloat(rawHost["cpu_times_guest"], 64)
	if err != nil {
		log.Errorf("parsing cpu times guest failed: %v", err)
		return nil, false
	}

	cpuTimesGuestNice, err := strconv.ParseFloat(rawHost["cpu_times_guest_nice"], 64)
	if err != nil {
		log.Errorf("parsing cpu times guest nice failed: %v", err)
		return nil, false
	}

	cpu := CPU{
		LogicalCount:   uint32(cpuLogicalCount),
		PhysicalCount:  uint32(cpuPhysicalCount),
		Percent:        cpuPercent,
		ProcessPercent: cpuProcessPercent,
		Times: CPUTimes{
			User:      cpuTimesUser,
			System:    cpuTimesSystem,
			Idle:      cpuTimesIdle,
			Nice:      cpuTimesNice,
			Iowait:    cpuTimesIowait,
			Irq:       cpuTimesIrq,
			Softirq:   cpuTimesSoftirq,
			Steal:     cpuTimesSteal,
			Guest:     cpuTimesGuest,
			GuestNice: cpuTimesGuestNice,
		},
	}

	// Set memory fields from raw host.
	memoryTotal, err := strconv.ParseUint(rawHost["memory_total"], 10, 64)
	if err != nil {
		log.Errorf("parsing memory total failed: %v", err)
		return nil, false
	}

	memoryAvailable, err := strconv.ParseUint(rawHost["memory_available"], 10, 64)
	if err != nil {
		log.Errorf("parsing memory available failed: %v", err)
		return nil, false
	}

	memoryUsed, err := strconv.ParseUint(rawHost["memory_used"], 10, 64)
	if err != nil {
		log.Errorf("parsing memory used failed: %v", err)
		return nil, false
	}

	memoryUsedPercent, err := strconv.ParseFloat(rawHost["memory_used_percent"], 64)
	if err != nil {
		log.Errorf("parsing memory used percent failed: %v", err)
		return nil, false
	}

	memoryProcessUsedPercent, err := strconv.ParseFloat(rawHost["memory_processe_used_percent"], 64)
	if err != nil {
		log.Errorf("parsing memory process used percent failed: %v", err)
		return nil, false
	}

	memoryFree, err := strconv.ParseUint(rawHost["memory_free"], 10, 64)
	if err != nil {
		log.Errorf("parsing memory free failed: %v", err)
		return nil, false
	}

	memory := Memory{
		Total:              memoryTotal,
		Available:          memoryAvailable,
		Used:               memoryUsed,
		UsedPercent:        memoryUsedPercent,
		ProcessUsedPercent: memoryProcessUsedPercent,
		Free:               memoryFree,
	}

	// Set network fields from raw host.
	networkTCPConnectionCount, err := strconv.ParseUint(rawHost["network_tcp_connection_count"], 10, 32)
	if err != nil {
		log.Errorf("parsing network tcp connection count failed: %v", err)
		return nil, false
	}

	networkUploadTCPConnectionCount, err := strconv.ParseUint(rawHost["network_upload_tcp_connection_count"], 10, 32)
	if err != nil {
		log.Errorf("parsing network upload tcp connection count failed: %v", err)
		return nil, false
	}

	rxBandwidth, err := strconv.ParseUint(rawHost["network_rx_bandwidth"], 10, 64)
	if err != nil {
		log.Errorf("parsing download rate failed: %v", err)
		return nil, false
	}

	maxRxBandwidth, err := strconv.ParseUint(rawHost["network_max_rx_bandwidth"], 10, 64)
	if err != nil {
		log.Errorf("parsing download rate limit failed: %v", err)
		return nil, false
	}

	txBandwidth, err := strconv.ParseUint(rawHost["network_tx_bandwidth"], 10, 64)
	if err != nil {
		log.Errorf("parsing upload rate failed: %v", err)
		return nil, false
	}

	maxTxBandwidth, err := strconv.ParseUint(rawHost["network_max_tx_bandwidth"], 10, 64)
	if err != nil {
		log.Errorf("parsing upload rate limit failed: %v", err)
		return nil, false
	}

	network := Network{
		TCPConnectionCount:       uint32(networkTCPConnectionCount),
		UploadTCPConnectionCount: uint32(networkUploadTCPConnectionCount),
		Location:                 rawHost["network_location"],
		IDC:                      rawHost["network_idc"],
		RxBandwidth:              rxBandwidth,
		MaxRxBandwidth:           maxRxBandwidth,
		TxBandwidth:              txBandwidth,
		MaxTxBandwidth:           maxTxBandwidth,
	}

	// Set disk fields from raw host.
	diskTotal, err := strconv.ParseUint(rawHost["disk_total"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk total failed: %v", err)
		return nil, false
	}

	diskFree, err := strconv.ParseUint(rawHost["disk_free"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk free failed: %v", err)
		return nil, false
	}

	diskUsed, err := strconv.ParseUint(rawHost["disk_used"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk used failed: %v", err)
		return nil, false
	}

	diskUsedPercent, err := strconv.ParseFloat(rawHost["disk_used_percent"], 64)
	if err != nil {
		log.Errorf("parsing disk used percent failed: %v", err)
		return nil, false
	}

	diskInodesTotal, err := strconv.ParseUint(rawHost["disk_inodes_total"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk inodes total failed: %v", err)
		return nil, false
	}

	diskInodesUsed, err := strconv.ParseUint(rawHost["disk_inodes_used"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk inodes used failed: %v", err)
		return nil, false
	}

	diskInodesFree, err := strconv.ParseUint(rawHost["disk_inodes_free"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk inodes free failed: %v", err)
		return nil, false
	}

	diskInodesUsedPercent, err := strconv.ParseFloat(rawHost["disk_inodes_used_percent"], 64)
	if err != nil {
		log.Errorf("parsing disk inodes used percent failed: %v", err)
		return nil, false
	}

	diskWriteBandwidth, err := strconv.ParseUint(rawHost["disk_write_bandwidth"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk write bandwidth failed: %v", err)
		return nil, false
	}

	diskReadBandwidth, err := strconv.ParseUint(rawHost["disk_read_bandwidth"], 10, 64)
	if err != nil {
		log.Errorf("parsing disk read bandwidth failed: %v", err)
		return nil, false
	}

	disk := Disk{
		Total:             diskTotal,
		Free:              diskFree,
		Used:              diskUsed,
		UsedPercent:       diskUsedPercent,
		InodesTotal:       diskInodesTotal,
		InodesUsed:        diskInodesUsed,
		InodesFree:        diskInodesFree,
		InodesUsedPercent: diskInodesUsedPercent,
		WriteBandwidth:    diskWriteBandwidth,
		ReadBandwidth:     diskReadBandwidth,
	}

	build := Build{
		GitVersion: rawHost["build_git_version"],
		GitCommit:  rawHost["build_git_commit"],
		GoVersion:  rawHost["build_go_version"],
		Platform:   rawHost["build_platform"],
	}

	// Set time fields from raw host.
	announceInterval, err := strconv.ParseUint(rawHost["announce_interval"], 10, 64)
	if err != nil {
		log.Errorf("parsing announce interval failed: %v", err)
		return nil, false
	}

	createdAt, err := time.Parse(time.RFC3339, rawHost["created_at"])
	if err != nil {
		log.Errorf("parsing created at failed: %v", err)
		return nil, false
	}

	updatedAt, err := time.Parse(time.RFC3339, rawHost["updated_at"])
	if err != nil {
		log.Errorf("parsing updated at failed: %v", err)
		return nil, false
	}

	return NewHost(
		rawHost["id"],
		rawHost["hostname"],
		rawHost["ip"],
		rawHost["os"],
		rawHost["platform"],
		rawHost["platform_family"],
		rawHost["platform_version"],
		rawHost["kernel_version"],
		int32(port),
		int32(downloadPort),
		int32(proxyPort),
		uint64(schedulerClusterID),
		disableShared,
		pkgtypes.ParseHostType(rawHost["type"]),
		cpu,
		memory,
		network,
		disk,
		build,
		time.Duration(announceInterval),
		createdAt,
		updatedAt,
		logger.WithHost(rawHost["id"], rawHost["hostname"], rawHost["ip"]),
	), true
}

// Store sets host.
func (h *hostManager) Store(ctx context.Context, host *Host) error {
	// Define the Lua script as a string.
	const storeHostScript = `
-- Extract keys and arguments
local host_key = KEYS[1]  -- Key for the host hash
local hosts_set_key = KEYS[2]  -- Key for the set of hosts

-- Extract host fields from arguments
local host_id = ARGV[1]
local host_type = ARGV[2]
local hostname = ARGV[3]
local ip = ARGV[4]
local port = ARGV[5]
local download_port = ARGV[6]
local proxy_port = ARGV[7]
local disable_shared = tonumber(ARGV[8])
local os = ARGV[9]
local platform = ARGV[10]
local platform_family = ARGV[11]
local platform_version = ARGV[12]
local kernel_version = ARGV[13]
local cpu_logical_count = ARGV[14]
local cpu_physical_count = ARGV[15]
local cpu_percent = ARGV[16]
local cpu_process_percent = ARGV[17]
local cpu_times_user = ARGV[18]
local cpu_times_system = ARGV[19]
local cpu_times_idle = ARGV[20]
local cpu_times_nice = ARGV[21]
local cpu_times_iowait = ARGV[22]
local cpu_times_irq = ARGV[23]
local cpu_times_softirq = ARGV[24]
local cpu_times_steal = ARGV[25]
local cpu_times_guest = ARGV[26]
local cpu_times_guest_nice = ARGV[27]
local memory_total = ARGV[28]
local memory_available = ARGV[29]
local memory_used = ARGV[30]
local memory_used_percent = ARGV[31]
local memory_process_used_percent = ARGV[32]
local memory_free = ARGV[33]
local network_tcp_connection_count = ARGV[34]
local network_upload_tcp_connection_count = ARGV[35]
local network_location = ARGV[36]
local network_idc = ARGV[37]
local network_rx_bandwidth = ARGV[38]
local network_max_rx_bandwidth = ARGV[39]
local network_tx_bandwidth = ARGV[40]
local network_max_tx_bandwidth = ARGV[41]
local disk_total = ARGV[42]
local disk_free = ARGV[43]
local disk_used = ARGV[44]
local disk_used_percent = ARGV[45]
local disk_inodes_total = ARGV[46]
local disk_inodes_used = ARGV[47]
local disk_inodes_free = ARGV[48]
local disk_inodes_used_percent = ARGV[49]
local disk_write_bandwidth = ARGV[50]
local disk_read_bandwidth = ARGV[51]
local build_git_version = ARGV[52]
local build_git_commit = ARGV[53]
local build_go_version = ARGV[54]
local build_platform = ARGV[55]
local scheduler_cluster_id = ARGV[56]
local announce_interval = ARGV[57]
local created_at = ARGV[58]
local updated_at = ARGV[59]

-- Perform HSET operation
redis.call("HSET", host_key,
    "id", host_id,
    "type", host_type,
    "hostname", hostname,
    "ip", ip,
    "port", port,
    "download_port", download_port,
    "proxy_port", proxy_port,
    "disable_shared", disable_shared,
    "os", os,
    "platform", platform,
    "platform_family", platform_family,
    "platform_version", platform_version,
    "kernel_version", kernel_version,
    "cpu_logical_count", cpu_logical_count,
    "cpu_physical_count", cpu_physical_count,
    "cpu_percent", cpu_percent,
    "cpu_processe_percent", cpu_process_percent,
    "cpu_times_user", cpu_times_user,
    "cpu_times_system", cpu_times_system,
    "cpu_times_idle", cpu_times_idle,
    "cpu_times_nice", cpu_times_nice,
    "cpu_times_iowait", cpu_times_iowait,
    "cpu_times_irq", cpu_times_irq,
    "cpu_times_softirq", cpu_times_softirq,
    "cpu_times_steal", cpu_times_steal,
    "cpu_times_guest", cpu_times_guest,
    "cpu_times_guest_nice", cpu_times_guest_nice,
    "memory_total", memory_total,
    "memory_available", memory_available,
    "memory_used", memory_used,
    "memory_used_percent", memory_used_percent,
    "memory_processe_used_percent", memory_process_used_percent,
    "memory_free", memory_free,
    "network_tcp_connection_count", network_tcp_connection_count,
    "network_upload_tcp_connection_count", network_upload_tcp_connection_count,
    "network_location", network_location,
    "network_idc", network_idc,
    "network_rx_bandwidth", network_rx_bandwidth,
    "network_max_rx_bandwidth", network_max_rx_bandwidth,
    "network_tx_bandwidth", network_tx_bandwidth,
    "network_max_tx_bandwidth", network_max_tx_bandwidth,
    "disk_total", disk_total,
    "disk_free", disk_free,
    "disk_used", disk_used,
    "disk_used_percent", disk_used_percent,
    "disk_inodes_total", disk_inodes_total,
    "disk_inodes_used", disk_inodes_used,
    "disk_inodes_free", disk_inodes_free,
    "disk_inodes_used_percent", disk_inodes_used_percent,
    "disk_write_bandwidth", disk_write_bandwidth,
    "disk_read_bandwidth", disk_read_bandwidth,
    "build_git_version", build_git_version,
    "build_git_commit", build_git_commit,
    "build_go_version", build_go_version,
    "build_platform", build_platform,
    "scheduler_cluster_id", scheduler_cluster_id,
    "announce_interval", announce_interval,
    "created_at", created_at,
    "updated_at", updated_at)

-- Perform SADD operation
redis.call("SADD", hosts_set_key, host_id)

return true
`

	// Create a new Redis script.
	script := redis.NewScript(storeHostScript)

	// Prepare keys.
	keys := []string{
		pkgredis.MakePersistentCacheHostKeyInScheduler(h.config.Manager.SchedulerClusterID, host.ID),
		pkgredis.MakePersistentCacheHostsInScheduler(h.config.Manager.SchedulerClusterID),
	}

	// Prepare arguments.
	args := []any{
		host.ID,
		host.Type.Name(),
		host.Hostname,
		host.IP,
		host.Port,
		host.DownloadPort,
		host.ProxyPort,
		host.DisableShared,
		host.OS,
		host.Platform,
		host.PlatformFamily,
		host.PlatformVersion,
		host.KernelVersion,
		host.CPU.LogicalCount,
		host.CPU.PhysicalCount,
		host.CPU.Percent,
		host.CPU.ProcessPercent,
		host.CPU.Times.User,
		host.CPU.Times.System,
		host.CPU.Times.Idle,
		host.CPU.Times.Nice,
		host.CPU.Times.Iowait,
		host.CPU.Times.Irq,
		host.CPU.Times.Softirq,
		host.CPU.Times.Steal,
		host.CPU.Times.Guest,
		host.CPU.Times.GuestNice,
		host.Memory.Total,
		host.Memory.Available,
		host.Memory.Used,
		host.Memory.UsedPercent,
		host.Memory.ProcessUsedPercent,
		host.Memory.Free,
		host.Network.TCPConnectionCount,
		host.Network.UploadTCPConnectionCount,
		host.Network.Location,
		host.Network.IDC,
		host.Network.RxBandwidth,
		host.Network.MaxRxBandwidth,
		host.Network.TxBandwidth,
		host.Network.MaxTxBandwidth,
		host.Disk.Total,
		host.Disk.Free,
		host.Disk.Used,
		host.Disk.UsedPercent,
		host.Disk.InodesTotal,
		host.Disk.InodesUsed,
		host.Disk.InodesFree,
		host.Disk.InodesUsedPercent,
		host.Disk.WriteBandwidth,
		host.Disk.ReadBandwidth,
		host.Build.GitVersion,
		host.Build.GitCommit,
		host.Build.GoVersion,
		host.Build.Platform,
		host.SchedulerClusterID,
		host.AnnounceInterval.Nanoseconds(),
		host.CreatedAt.Format(time.RFC3339),
		host.UpdatedAt.Format(time.RFC3339),
	}

	// Execute the script.
	if err := script.Run(ctx, h.rdb, keys, args...).Err(); err != nil {
		host.Log.Errorf("store host failed: %v", err)
		return err
	}

	return nil
}

// Delete deletes host by a key.
func (h *hostManager) Delete(ctx context.Context, hostID string) error {
	// Define the Lua script as a string.
	const deleteHostScript = `
-- Extract keys
local host_key = KEYS[1]  -- Key for the host hash
local hosts_set_key = KEYS[2]  -- Key for the set of hosts

-- Extract arguments
local host_id = ARGV[1]

-- Perform DEL operation to delete the host hash
redis.call("DEL", host_key)

-- Perform SREM operation to remove the host ID from the set
redis.call("SREM", hosts_set_key, host_id)

return true
`

	log := logger.WithHostID(hostID)

	// Create a new Redis script.
	script := redis.NewScript(deleteHostScript)

	// Prepare keys.
	keys := []string{
		pkgredis.MakePersistentCacheHostKeyInScheduler(h.config.Manager.SchedulerClusterID, hostID),
		pkgredis.MakePersistentCacheHostsInScheduler(h.config.Manager.SchedulerClusterID),
	}

	// Prepare arguments.
	args := []any{
		hostID,
	}

	// Execute the script.
	err := script.Run(ctx, h.rdb, keys, args...).Err()
	if err != nil {
		log.Errorf("delete host failed: %v", err)
		return err
	}

	return nil
}

// LoadAll returns all hosts.
func (h *hostManager) LoadAll(ctx context.Context) ([]*Host, error) {
	var (
		hosts  []*Host
		cursor uint64
	)

	for {
		var (
			hostKeys []string
			err      error
		)

		hostKeys, cursor, err = h.rdb.SScan(ctx, pkgredis.MakePersistentCacheHostsInScheduler(h.config.Manager.SchedulerClusterID), cursor, "*", 10).Result()
		if err != nil {
			logger.Error("scan hosts failed")
			return nil, err
		}

		for _, hostKey := range hostKeys {
			host, loaded := h.Load(ctx, hostKey)
			if !loaded {
				logger.WithHostID(hostKey).Error("load host failed")
				continue
			}

			hosts = append(hosts, host)
		}

		if cursor == 0 {
			break
		}
	}

	return hosts, nil
}

// LoadRandom loads host randomly through the set of redis.
func (h *hostManager) LoadRandom(ctx context.Context, n int, blocklist set.SafeSet[string]) ([]*Host, error) {
	hostKeys, err := h.rdb.SMembers(ctx, pkgredis.MakePersistentCacheHostsInScheduler(h.config.Manager.SchedulerClusterID)).Result()
	if err != nil {
		logger.Error("smembers hosts failed")
		return nil, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(hostKeys), func(i, j int) {
		hostKeys[i], hostKeys[j] = hostKeys[j], hostKeys[i]
	})

	hosts := make([]*Host, 0, n)
	for _, hostKey := range hostKeys {
		if len(hosts) >= n {
			break
		}

		if blocklist.Contains(hostKey) {
			continue
		}

		host, loaded := h.Load(ctx, hostKey)
		if !loaded {
			logger.WithHostID(hostKey).Error("load host failed")
			continue
		}

		hosts = append(hosts, host)
	}

	return hosts, nil
}

// RunGC runs garbage collection.
func (h *hostManager) RunGC(ctx context.Context) error {
	hosts, err := h.LoadAll(ctx)
	if err != nil {
		logger.Error("load all hosts failed")
		return err
	}

	for _, host := range hosts {
		// If the host's elapsed exceeds twice the announcing interval,
		// then leave peers in host.
		elapsed := time.Since(host.UpdatedAt)
		if host.AnnounceInterval > 0 && elapsed > host.AnnounceInterval*2 {
			host.Log.Info("host has been reclaimed")
			if err := h.Delete(ctx, host.ID); err != nil {
				host.Log.Errorf("delete host failed: %v", err)
			}
		}
	}

	return nil
}
