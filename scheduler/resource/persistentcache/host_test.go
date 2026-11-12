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

package persistentcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/types"
)

var (
	mockRawHost = Host{
		ID:              mockHostID,
		Type:            types.HostTypeNormal,
		Hostname:        "foo",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		ProxyPort:       8004,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Log:             logger.With("host", "foo"),
	}

	mockRawSeedHost = Host{
		ID:              mockSeedHostID,
		Type:            types.HostTypeSuperSeed,
		Hostname:        "bar",
		IP:              "127.0.0.1",
		Port:            8003,
		DownloadPort:    8001,
		ProxyPort:       8004,
		OS:              "darwin",
		Platform:        "darwin",
		PlatformFamily:  "Standalone Workstation",
		PlatformVersion: "11.1",
		KernelVersion:   "20.2.0",
		CPU:             mockCPU,
		Memory:          mockMemory,
		Network:         mockNetwork,
		Disk:            mockDisk,
		Build:           mockBuild,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Log:             logger.With("host", "foo"),
	}

	mockCPU = CPU{
		LogicalCount:   4,
		PhysicalCount:  2,
		Percent:        1,
		ProcessPercent: 0.5,
		Times: CPUTimes{
			User:      240662.2,
			System:    317950.1,
			Idle:      3393691.3,
			Nice:      0,
			Iowait:    0,
			Irq:       0,
			Softirq:   0,
			Steal:     0,
			Guest:     0,
			GuestNice: 0,
		},
	}

	mockMemory = Memory{
		Total:              17179869184,
		Available:          5962813440,
		Used:               11217055744,
		UsedPercent:        65.291858,
		ProcessUsedPercent: 41.525125,
		Free:               2749598908,
	}

	mockNetwork = Network{
		TCPConnectionCount:       10,
		UploadTCPConnectionCount: 1,
		Location:                 mockHostLocation,
		IDC:                      mockHostIDC,
		RxBandwidth:              100,
		MaxRxBandwidth:           200,
		TxBandwidth:              100,
		MaxTxBandwidth:           200,
	}

	mockDisk = Disk{
		Total:             499963174912,
		Free:              37226479616,
		Used:              423809622016,
		UsedPercent:       91.92547406065952,
		InodesTotal:       4882452880,
		InodesUsed:        7835772,
		InodesFree:        4874617108,
		InodesUsedPercent: 0.1604884305611568,
	}

	mockBuild = Build{
		GitVersion: "v1.0.0",
		GitCommit:  "221176b117c6d59366d68f2b34d38be50c935883",
		GoVersion:  "1.18",
		Platform:   "darwin",
	}

	mockAnnounceInterval = 5 * time.Minute

	mockHostID       = idgen.HostID("127.0.0.1", "foo", false)
	mockSeedHostID   = idgen.HostID("127.0.0.1", "bar", true)
	mockHostLocation = "baz"
	mockHostIDC      = "bas"
)

func TestNewHost(t *testing.T) {
	cpu := CPU{
		LogicalCount:   4,
		PhysicalCount:  2,
		Percent:        50.0,
		ProcessPercent: 25.0,
		Times: CPUTimes{
			User:      100.0,
			System:    50.0,
			Idle:      200.0,
			Nice:      10.0,
			Iowait:    5.0,
			Irq:       1.0,
			Softirq:   2.0,
			Steal:     0.0,
			Guest:     0.0,
			GuestNice: 0.0,
		},
	}
	memory := Memory{
		Total:              16 * 1024 * 1024 * 1024,
		Available:          8 * 1024 * 1024 * 1024,
		Used:               8 * 1024 * 1024 * 1024,
		UsedPercent:        50.0,
		ProcessUsedPercent: 25.0,
		Free:               4 * 1024 * 1024 * 1024,
	}
	network := Network{
		TCPConnectionCount:       100,
		UploadTCPConnectionCount: 50,
		Location:                 "us-west",
		IDC:                      "test-idc",
		RxBandwidth:              1024 * 1024,
		MaxRxBandwidth:           2 * 1024 * 1024,
		TxBandwidth:              512 * 1024,
		MaxTxBandwidth:           1024 * 1024,
	}
	disk := Disk{
		Total:             1000 * 1024 * 1024 * 1024,
		Free:              500 * 1024 * 1024 * 1024,
		Used:              500 * 1024 * 1024 * 1024,
		UsedPercent:       50.0,
		InodesTotal:       1000000,
		InodesUsed:        500000,
		InodesFree:        500000,
		InodesUsedPercent: 50.0,
		WriteBandwidth:    100 * 1024 * 1024,
		ReadBandwidth:     200 * 1024 * 1024,
	}
	build := Build{
		GitVersion:  "v1.0.0",
		GitCommit:   "abc123",
		GoVersion:   "go1.17",
		RustVersion: "1.57",
		Platform:    "linux/amd64",
	}

	host := NewHost(
		"test-id", "test-name", "test-host", "127.0.0.1", "linux", "amd64", "debian", "11", "5.10.0", 8002, 8001, 8004,
		1, false, types.HostTypeNormal, cpu, memory, network, disk, build, 30*time.Second,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), nil,
	)

	assert := assert.New(t)
	assert.Equal("test-id", host.ID)
	assert.Equal("test-name", host.Name)
	assert.Equal("test-host", host.Hostname)
	assert.Equal("127.0.0.1", host.IP)
	assert.Equal("linux", host.OS)
	assert.Equal("amd64", host.Platform)
	assert.Equal("debian", host.PlatformFamily)
	assert.Equal("11", host.PlatformVersion)
	assert.Equal("5.10.0", host.KernelVersion)
	assert.Equal(int32(8002), host.Port)
	assert.Equal(int32(8001), host.DownloadPort)
	assert.Equal(int32(8004), host.ProxyPort)
	assert.Equal(uint64(1), host.SchedulerClusterID)
	assert.False(host.DisableShared)
	assert.Equal(types.HostTypeNormal, host.Type)
	assert.Equal(cpu, host.CPU)
	assert.Equal(memory, host.Memory)
	assert.Equal(network, host.Network)
	assert.Equal(disk, host.Disk)
	assert.Equal(build, host.Build)
	assert.Equal(30*time.Second, host.AnnounceInterval)
	assert.Equal(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), host.CreatedAt)
	assert.Equal(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), host.UpdatedAt)
	assert.NotNil(host.Log)
}
