/*
 *     Copyright 2024 The Dragonfly Authors
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

package persistentcache

import (
	"time"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/types"
)

// Host contains content for host.
type Host struct {
	// ID is host id.
	ID string

	// Type is host type.
	Type types.HostType

	// Hostname is host name.
	Hostname string

	// IP is host ip.
	IP string

	// Port is grpc service port.
	Port int32

	// Port is tcp service port.
	TCPPort int32

	// Port is quic service port.
	QUICPort int32

	// DownloadPort is piece downloading port.
	DownloadPort int32

	// ProxyPort is the proxy server port.
	ProxyPort int32

	// DisableShared is whether the host is disabled for
	// shared with other peers.
	DisableShared bool

	// Host OS.
	OS string

	// Host platform.
	Platform string

	// Host platform family.
	PlatformFamily string

	// Host platform version.
	PlatformVersion string

	// Host kernel version.
	KernelVersion string

	// CPU Stat.
	CPU CPU

	// Memory Stat.
	Memory Memory

	// Network Stat.
	Network Network

	// Dist Stat.
	Disk Disk

	// Build information.
	Build Build

	// SchedulerClusterID is the scheduler cluster id matched by scopes.
	SchedulerClusterID uint64

	// AnnounceInterval is the interval between host announces to scheduler.
	AnnounceInterval time.Duration

	// CreatedAt is host create time.
	CreatedAt time.Time

	// UpdatedAt is host update time.
	UpdatedAt time.Time

	// Host log.
	Log *logger.SugaredLoggerOnWith
}

// CPU contains content for cpu.
type CPU struct {
	// Number of logical cores in the system.
	LogicalCount uint32

	// Number of physical cores in the system.
	PhysicalCount uint32

	// Percent calculates the percentage of cpu used.
	Percent float64

	// Calculates the percentage of cpu used by process.
	ProcessPercent float64

	// Times contains the amounts of time the CPU has spent performing different kinds of work.
	Times CPUTimes
}

// CPUTimes contains content for cpu times.
type CPUTimes struct {
	// CPU time of user.
	User float64

	// CPU time of system.
	System float64

	// CPU time of idle.
	Idle float64

	// CPU time of nice.
	Nice float64

	// CPU time of iowait.
	Iowait float64

	// CPU time of irq.
	Irq float64

	// CPU time of softirq.
	Softirq float64

	// CPU time of steal.
	Steal float64

	// CPU time of guest.
	Guest float64

	// CPU time of guest nice.
	GuestNice float64
}

// Memory contains content for memory.
type Memory struct {
	// Total amount of RAM on this system.
	Total uint64

	// RAM available for programs to allocate.
	Available uint64

	// RAM used by programs.
	Used uint64

	// Percentage of RAM used by programs.
	UsedPercent float64

	// Calculates the percentage of memory used by process.
	ProcessUsedPercent float64

	// This is the kernel's notion of free memory.
	Free uint64
}

// Network contains content for network.
type Network struct {
	// Return count of tcp connections opened and status is ESTABLISHED.
	TCPConnectionCount uint32

	// Return count of upload tcp connections opened and status is ESTABLISHED.
	UploadTCPConnectionCount uint32

	// Location path(area|country|province|city|...).
	Location string

	// IDC where the peer host is located
	IDC string

	// RxBandwidth of the host, unit is byte/s.
	RxBandwidth uint64

	// MaxRxBandwidth of the host, unit is byte/s.
	MaxRxBandwidth uint64

	// TxBandwidth of the host, unit is byte/s.
	TxBandwidth uint64

	// MaxTxBandwidth of the host, unit is byte/s.
	MaxTxBandwidth uint64
}

// Build contains content for build.
type Build struct {
	// Git version.
	GitVersion string

	// Git commit.
	GitCommit string

	// Golang version.
	GoVersion string

	// Rust version.
	RustVersion string

	// Build platform.
	Platform string
}

// Disk contains content for disk.
type Disk struct {
	// Total amount of disk on the data path of dragonfly.
	Total uint64

	// Free amount of disk on the data path of dragonfly.
	Free uint64

	// Used amount of disk on the data path of dragonfly.
	Used uint64

	// Used percent of disk on the data path of dragonfly directory.
	UsedPercent float64

	// Total amount of indoes on the data path of dragonfly directory.
	InodesTotal uint64

	// Used amount of indoes on the data path of dragonfly directory.
	InodesUsed uint64

	// Free amount of indoes on the data path of dragonfly directory.
	InodesFree uint64

	// Used percent of indoes on the data path of dragonfly directory.
	InodesUsedPercent float64

	// Disk write bandwidth, unit is byte/s.
	WriteBandwidth uint64

	// Disk read bandwidth, unit is byte/s.
	ReadBandwidth uint64
}

// NewHost returns a new host.
func NewHost(
	id, hostname, ip, os, platform, platformFamily, platformVersion, kernelVersion string, port, downloadPort, proxyPort int32,
	schedulerClusterId uint64, disableShared bool, typ types.HostType, cpu CPU, memory Memory, network Network, disk Disk,
	build Build, announceInterval time.Duration, createdAt, updatedAt time.Time, log *logger.SugaredLoggerOnWith,
) *Host {
	return &Host{
		ID:                 id,
		Type:               types.HostType(typ),
		Hostname:           hostname,
		IP:                 ip,
		Port:               port,
		DownloadPort:       downloadPort,
		ProxyPort:          proxyPort,
		DisableShared:      disableShared,
		OS:                 os,
		Platform:           platform,
		PlatformFamily:     platformFamily,
		PlatformVersion:    platformVersion,
		KernelVersion:      kernelVersion,
		CPU:                cpu,
		Memory:             memory,
		Network:            network,
		Disk:               disk,
		Build:              build,
		SchedulerClusterID: schedulerClusterId,
		AnnounceInterval:   announceInterval,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		Log:                logger.WithHost(id, hostname, ip),
	}
}
