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

package types

import (
	"time"

	"github.com/bits-and-blooms/bitset"
)

type PersistentCacheParams struct {
	// SchedulerClusterID is the scheduler cluster id of the persistent cache.
	SchedulerClusterID uint `uri:"scheduler_cluster_id" binding:"required"`

	// TaskID is the task id of the persistent cache.
	TaskID string `uri:"task_id" binding:"required"`
}

type GetPersistentCachesQuery struct {
	// SchedulerClusterIDs is the scheduler cluster ids of the persistent cache.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`

	// Page is the page number of the persistent cache list.
	Page int `form:"page" binding:"omitempty,gte=1"`

	// PerPage is the item count per page of the persistent cache list.
	PerPage int `form:"per_page" binding:"omitempty,gte=1,lte=10000000"`
}

type GetPersistentCacheResponse struct {
	// TaskID is task id.
	TaskID string `json:"task_id" binding:"omitempty"`

	// PersistentReplicaCount is replica count of the persistent cache task.
	PersistentReplicaCount uint64 `json:"persistent_replica_count" binding:"omitempty"`

	// Tag is used to distinguish different persistent cache tasks.
	Tag string `json:"tag" binding:"omitempty"`

	// Application of persistent cache task.
	Application string `json:"application" binding:"omitempty"`

	// PieceLength is persistent cache task piece length.
	PieceLength uint64 `json:"piece_length" binding:"omitempty"`

	// ContentLength is persistent cache task total content length.
	ContentLength uint64 `json:"content_length" binding:"omitempty"`

	// TotalPieceCount is total piece count.
	TotalPieceCount uint32 `json:"total_piece_count" binding:"omitempty"`

	// State is persistent cache task state.
	State string `json:"state" binding:"omitempty"`

	// TTL is persistent cache task time to live.
	TTL time.Duration `json:"ttl" binding:"omitempty"`

	// CreatedAt is persistent cache task create time.
	CreatedAt time.Time `json:"created_at" binding:"omitempty"`

	// UpdatedAt is persistent cache task update time.
	UpdatedAt time.Time `json:"updated_at" binding:"omitempty"`

	// PersistentCachePeers is the list of persistent cache peers.
	PersistentCachePeers []PersistentCachePeer `json:"persistent_cache_peers" binding:"omitempty"`
}

type PersistentCachePeer struct {
	// ID is persistent cache peer id.
	ID string `json:"id" binding:"omitempty"`

	// Persistent is whether the peer is persistent.
	Persistent bool `json:"persistent" binding:"omitempty"`

	// FinishedPieces is finished pieces bitset.
	FinishedPieces *bitset.BitSet `json:"finished_pieces" binding:"omitempty"`

	// State is persistent cache peer state.
	State string `json:"state" binding:"omitempty"`

	// BlockParents is bad parents ids.
	BlockParents []string `json:"block_parents" binding:"omitempty"`

	// Cost is the cost of downloading.
	Cost time.Duration `json:"cost" binding:"omitempty"`

	// CreatedAt is persistent cache peer create time.
	CreatedAt time.Time `json:"created_at" binding:"omitempty"`

	// UpdatedAt is persistent cache peer update time.
	UpdatedAt time.Time `json:"updated_at" binding:"omitempty"`

	// Host is the peer host.
	Host PersistentCachePeerHost `json:"host" binding:"omitempty"`
}

type PersistentCachePeerHost struct {
	// ID is host id.
	ID string `json:"id" binding:"omitempty"`

	// Type is host type.
	Type string `json:"type" binding:"omitempty"`

	// Hostname is host name.
	Hostname string `json:"hostname" binding:"omitempty"`

	// IP is host ip.
	IP string `json:"ip" binding:"omitempty"`

	// Port is grpc service port.
	Port int32 `json:"port" binding:"omitempty"`

	// DownloadPort is piece downloading port.
	DownloadPort int32 `json:"download_port" binding:"omitempty"`

	// DisableShared is whether the host is disabled for shared with other peers.
	DisableShared bool `json:"disable_shared" binding:"omitempty"`

	// OS is host OS.
	OS string `json:"os" binding:"omitempty"`

	// Platform is host platform.
	Platform string `json:"platform" binding:"omitempty"`

	// PlatformFamily is host platform family.
	PlatformFamily string `json:"platform_family" binding:"omitempty"`

	// PlatformVersion is host platform version.
	PlatformVersion string `json:"platform_version" binding:"omitempty"`

	// KernelVersion is host kernel version.
	KernelVersion string `json:"kernel_version" binding:"omitempty"`

	// CPU contains cpu information.
	CPU struct {
		// LogicalCount is cpu logical count.
		LogicalCount uint32 `json:"logical_count" binding:"omitempty"`

		// PhysicalCount is cpu physical count.
		PhysicalCount uint32 `json:"physical_count" binding:"omitempty"`

		// Percent is cpu usage percent.
		Percent float64 `json:"percent" binding:"omitempty"`

		// ProcessPercent is process cpu usage percent.
		ProcessPercent float64 `json:"process_percent" binding:"omitempty"`

		// Times contains cpu times information.
		Times struct {
			// User is user cpu time.
			User float64 `json:"user" binding:"omitempty"`

			// System is system cpu time.
			System float64 `json:"system" binding:"omitempty"`

			// Idle is idle cpu time.
			Idle float64 `json:"idle" binding:"omitempty"`

			// Nice is nice cpu time.
			Nice float64 `json:"nice" binding:"omitempty"`

			// Iowait is iowait cpu time.
			Iowait float64 `json:"iowait" binding:"omitempty"`

			// Irq is irq cpu time.
			Irq float64 `json:"irq" binding:"omitempty"`

			// Softirq is softirq cpu time.
			Softirq float64 `json:"softirq" binding:"omitempty"`

			// Steal is steal cpu time.
			Steal float64 `json:"steal" binding:"omitempty"`

			// Guest is guest cpu time.
			Guest float64 `json:"guest" binding:"omitempty"`

			// GuestNice is guest nice cpu time.
			GuestNice float64 `json:"guest_nice" binding:"omitempty"`
		} `json:"times" binding:"omitempty"`
	} `json:"cpu" binding:"omitempty"`

	// Memory contains memory information.
	Memory struct {
		// Total is total memory.
		Total uint64 `json:"total" binding:"omitempty"`

		// Available is available memory.
		Available uint64 `json:"available" binding:"omitempty"`

		// Used is used memory.
		Used uint64 `json:"used" binding:"omitempty"`

		// UsedPercent is memory usage percent.
		UsedPercent float64 `json:"used_percent" binding:"omitempty"`

		// ProcessUsedPercent is process memory usage percent.
		ProcessUsedPercent float64 `json:"process_used_percent" binding:"omitempty"`

		// Free is free memory.
		Free uint64 `json:"free" binding:"omitempty"`
	} `json:"memory" binding:"omitempty"`

	// Network contains network information.
	Network struct {
		// TCPConnectionCount is tcp connection count.
		TCPConnectionCount uint32 `json:"tcp_connection_count" binding:"omitempty"`

		// UploadTCPConnectionCount is upload tcp connection count.
		UploadTCPConnectionCount uint32 `json:"upload_tcp_connection_count" binding:"omitempty"`

		// Location is network location.
		Location string `json:"location" binding:"omitempty"`

		// IDC is network idc.
		IDC string `json:"idc" binding:"omitempty"`

		// DownloadRate is download rate.
		DownloadRate uint64 `json:"download_rate" binding:"omitempty"`

		// DownloadRateLimit is download rate limit.
		DownloadRateLimit uint64 `json:"download_rate_limit" binding:"omitempty"`

		// UploadRate is upload rate.
		UploadRate uint64 `json:"upload_rate" binding:"omitempty"`

		// UploadRateLimit is upload rate limit.
		UploadRateLimit uint64 `json:"upload_rate_limit" binding:"omitempty"`
	} `json:"network" binding:"omitempty"`

	// Disk contains disk information.
	Disk struct {
		// Total is total disk space.
		Total uint64 `json:"total" binding:"omitempty"`

		// Free is free disk space.
		Free uint64 `json:"free" binding:"omitempty"`

		// Used is used disk space.
		Used uint64 `json:"used" binding:"omitempty"`

		// UsedPercent is disk usage percent.
		UsedPercent float64 `json:"used_percent" binding:"omitempty"`

		// InodesTotal is total inodes.
		InodesTotal uint64 `json:"inodes_total" binding:"omitempty"`

		// InodesUsed is used inodes.
		InodesUsed uint64 `json:"inodes_used" binding:"omitempty"`

		// InodesFree is free inodes.
		InodesFree uint64 `json:"inodes_free" binding:"omitempty"`

		// InodesUsedPercent is inodes usage percent.
		InodesUsedPercent float64 `json:"inodes_used_percent" binding:"omitempty"`

		// WriteBandwidth is write bandwidth.
		WriteBandwidth uint64 `json:"write_bandwidth" binding:"omitempty"`

		// ReadBandwidth is read bandwidth.
		ReadBandwidth uint64 `json:"read_bandwidth" binding:"omitempty"`
	} `json:"disk" binding:"omitempty"`

	// Build contains build information.
	Build struct {
		// GitVersion is git version.
		GitVersion string `json:"git_version" binding:"omitempty"`

		// GitCommit is git commit.
		GitCommit string `json:"git_commit" binding:"omitempty"`

		// GoVersion is go version.
		GoVersion string `json:"go_version" binding:"omitempty"`

		// Platform is build platform.
		Platform string `json:"platform" binding:"omitempty"`
	} `json:"build" binding:"omitempty"`

	// SchedulerClusterID is the scheduler cluster id matched by scopes.
	SchedulerClusterID uint64 `json:"scheduler_cluster_id" binding:"omitempty"`

	// AnnounceInterval is the interval between host announces to scheduler.
	AnnounceInterval time.Duration `json:"announce_interval" binding:"omitempty"`

	// CreatedAt is host create time.
	CreatedAt time.Time `json:"created_at" binding:"omitempty"`

	// UpdatedAt is host update time.
	UpdatedAt time.Time `json:"updated_at" binding:"omitempty"`
}
