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

package evaluator

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	"d7y.io/dragonfly/v2/pkg/digest"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
)

var (
	mockRawHost = standard.Host{
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
		CreatedAt:       atomic.NewTime(time.Now()),
		UpdatedAt:       atomic.NewTime(time.Now()),
	}

	mockRawSeedHost = standard.Host{
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
		CreatedAt:       atomic.NewTime(time.Now()),
		UpdatedAt:       atomic.NewTime(time.Now()),
	}

	mockCPU = standard.CPU{
		LogicalCount:   4,
		PhysicalCount:  2,
		Percent:        1,
		ProcessPercent: 0.5,
		Times: standard.CPUTimes{
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

	mockMemory = standard.Memory{
		Total:              17179869184,
		Available:          5962813440,
		Used:               11217055744,
		UsedPercent:        65.291858,
		ProcessUsedPercent: 41.525125,
		Free:               2749598908,
	}

	mockNetwork = standard.Network{
		TCPConnectionCount:       10,
		UploadTCPConnectionCount: 1,
		Location:                 mockHostLocation,
		IDC:                      mockHostIDC,
		RxBandwidth:              100,
		MaxRxBandwidth:           200,
		TxBandwidth:              100,
		MaxTxBandwidth:           200,
	}

	mockDisk = standard.Disk{
		Total:             499963174912,
		Free:              37226479616,
		Used:              423809622016,
		UsedPercent:       91.92547406065952,
		InodesTotal:       4882452880,
		InodesUsed:        7835772,
		InodesFree:        4874617108,
		InodesUsedPercent: 0.1604884305611568,
		WriteBandwidth:    1,
		ReadBandwidth:     1,
	}

	mockBuild = standard.Build{
		GitVersion: "v1.0.0",
		GitCommit:  "221176b117c6d59366d68f2b34d38be50c935883",
		GoVersion:  "1.18",
		Platform:   "darwin",
	}

	mockTaskBackToSourceLimit   int32  = 200
	mockTaskURL                        = "http://example.com/foo"
	mockTaskPieceLength         uint64 = 2048
	mockTaskID                         = idgen.TaskIDV2ByURLBased(mockTaskURL, &mockTaskPieceLength, mockTaskTag, mockTaskApplication, mockTaskFilteredQueryParams)
	mockTaskDigest                     = digest.New(digest.AlgorithmSHA256, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4")
	mockTaskTag                        = "d7y"
	mockTaskApplication                = "foo"
	mockTaskFilteredQueryParams        = []string{"bar"}
	mockTaskHeader                     = map[string]string{"content-length": "100"}
	mockHostID                         = idgen.HostIDV2("127.0.0.1", "foo", false)
	mockSeedHostID                     = idgen.HostIDV2("127.0.0.1", "bar", true)
	mockHostLocation                   = "bas"
	mockHostIDC                        = "baz"
	mockPeerID                         = idgen.PeerIDV2()
)

func TestEvaluatorDefault_newEvaluatorDefault(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, e any)
	}{
		{
			name: "new evaluator commonv1",
			expect: func(t *testing.T, e any) {
				assert := assert.New(t)
				assert.Equal(reflect.TypeOf(e).Elem().Name(), "evaluatorDefault")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, newEvaluatorDefault())
		})
	}
}

func TestEvaluatorDefault_EvaluateParents(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(parents []*standard.Peer, child *standard.Peer)
		expect func(t *testing.T, parents []*standard.Peer)
	}{
		{
			name: "sort parents by score in descending order",
			mock: func(parents []*standard.Peer, child *standard.Peer) {
				parents[0].Host.Network.MaxTxBandwidth = 1000
				parents[0].Host.TxBandwidth.Store(900)
				parents[0].Host.Network.IDC = mockHostIDC
				parents[0].Host.Network.Location = mockHostLocation

				parents[1].Host.Network.MaxTxBandwidth = 1000
				parents[1].Host.TxBandwidth.Store(100)
				parents[1].Host.Network.IDC = mockHostIDC
				parents[1].Host.Network.Location = mockHostLocation

				parents[2].Host.Network.MaxTxBandwidth = 1000
				parents[2].Host.TxBandwidth.Store(500)
				parents[2].Host.Network.IDC = mockHostIDC
				parents[2].Host.Network.Location = mockHostLocation

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = mockHostLocation
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal("127.0.0.2-host2", parents[0].Host.ID)
				assert.Equal("127.0.0.3-host3", parents[1].Host.ID)
				assert.Equal("127.0.0.1-host1", parents[2].Host.ID)
			},
		},
		{
			name: "sort parents with IDC affinity considered",
			mock: func(parents []*standard.Peer, child *standard.Peer) {
				parents[0].Host.Network.MaxTxBandwidth = 1000
				parents[0].Host.TxBandwidth.Store(500)
				parents[0].Host.Network.IDC = mockHostIDC
				parents[0].Host.Network.Location = mockHostLocation

				parents[1].Host.Network.MaxTxBandwidth = 1000
				parents[1].Host.TxBandwidth.Store(100)
				parents[1].Host.Network.IDC = "idc-2"
				parents[1].Host.Network.Location = mockHostLocation

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = mockHostLocation
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal("127.0.0.1-host1", parents[0].Host.ID)
				assert.Equal("127.0.0.2-host2", parents[1].Host.ID)
				assert.Equal("127.0.0.3-host3", parents[2].Host.ID)
			},
		},
		{
			name: "sort parents with location affinity considered",
			mock: func(parents []*standard.Peer, child *standard.Peer) {
				parents[0].Host.Network.MaxTxBandwidth = 1000
				parents[0].Host.TxBandwidth.Store(500)
				parents[0].Host.Network.IDC = mockHostIDC
				parents[0].Host.Network.Location = "country-3|province-3"

				parents[1].Host.Network.MaxTxBandwidth = 1000
				parents[1].Host.TxBandwidth.Store(500)
				parents[1].Host.Network.IDC = mockHostIDC
				parents[1].Host.Network.Location = "country-3|province-3|city-3"

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = "country-3|province-3|city-3"
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal("127.0.0.2-host2", parents[0].Host.ID)
				assert.Equal("127.0.0.1-host1", parents[1].Host.ID)
				assert.Equal("127.0.0.3-host3", parents[2].Host.ID)
			},
		},
		{
			name: "sort parents with host type considered",
			mock: func(parents []*standard.Peer, child *standard.Peer) {
				parents[0].Host.Type = types.HostTypeNormal
				parents[1].Host.Type = types.HostTypeSuperSeed
				parents[2].Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, parents []*standard.Peer) {
				assert := assert.New(t)
				assert.Equal("127.0.0.1-host1", parents[0].Host.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			parents := []*standard.Peer{
				standard.NewPeer(idgen.PeerIDV2(), mockTask, standard.NewHost(
					idgen.HostIDV2("127.0.0.1", "host1", false), "127.0.0.1", "host1",
					8003, 8001, 8004, types.HostTypeNormal)),
				standard.NewPeer(idgen.PeerIDV2(), mockTask, standard.NewHost(
					idgen.HostIDV2("127.0.0.2", "host2", false), "127.0.0.2", "host2",
					8003, 8001, 8004, types.HostTypeNormal)),
				standard.NewPeer(idgen.PeerIDV2(), mockTask, standard.NewHost(
					idgen.HostIDV2("127.0.0.3", "host3", false), "127.0.0.3", "host3",
					8003, 8001, 8004, types.HostTypeNormal)),
			}

			child := standard.NewPeer(mockPeerID, mockTask, standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type))

			e := newEvaluatorDefault()
			tc.mock(parents, child)
			result := e.EvaluateParents(parents, child)
			tc.expect(t, result)
		})
	}
}

func TestEvaluatorDefault_evaluateParents(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(parent *standard.Peer, child *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "perfect score with all optimal conditions",
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.Host.Network.MaxTxBandwidth = 4294967296
				parent.Host.TxBandwidth.Store(0)
				parent.Host.UploadContentLength.Store(0)
				parent.Host.ConcurrentUploadPieceCount.Store(0)
				parent.Host.Network.IDC = mockHostIDC
				parent.Host.Network.Location = "country-2|province-2|city-2"
				parent.Host.Type = types.HostTypeSuperSeed
				parent.FSM.SetState(standard.PeerStateRunning)

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = "country-2|province-2|city-2"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(1.0, score)
			},
		},
		{
			name: "zero score with worst conditions",
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.Host.Network.MaxTxBandwidth = 0
				parent.Host.TxBandwidth.Store(0)
				parent.Host.UploadContentLength.Store(0)
				parent.Host.ConcurrentUploadPieceCount.Store(0)
				parent.Host.Network.IDC = ""
				parent.Host.Network.Location = ""
				parent.Host.Type = types.HostTypeSuperSeed
				parent.FSM.SetState(standard.PeerStateSucceeded)

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = mockHostLocation
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "mixed score with partial matches",
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.Host.Network.MaxTxBandwidth = 4294967296
				parent.Host.TxBandwidth.Store(500)
				parent.Host.UploadContentLength.Store(8589934592)
				parent.Host.ConcurrentUploadPieceCount.Store(2)
				parent.Host.Network.IDC = mockHostIDC
				parent.Host.Network.Location = "country-1|province-1"
				parent.Host.Type = types.HostTypeNormal

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = "country-1|province-1|city-1"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.InDelta(0.84, score, 0.01)
			},
		},
		{
			name: "seed peer in running state",
			mock: func(parent *standard.Peer, child *standard.Peer) {
				parent.Host.Network.MaxTxBandwidth = 4294967296
				parent.Host.TxBandwidth.Store(4294967296)
				parent.Host.UploadContentLength.Store(12884901888)
				parent.Host.ConcurrentUploadPieceCount.Store(16)
				parent.Host.Network.IDC = mockHostIDC
				parent.Host.Network.Location = mockHostLocation
				parent.Host.Type = types.HostTypeSuperSeed
				parent.FSM.SetState(standard.PeerStateRunning)

				child.Host.Network.IDC = mockHostIDC
				child.Host.Network.Location = mockHostLocation
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.628, score)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))

			parent := standard.NewPeer(idgen.PeerIDV2(), mockTask, standard.NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type))

			child := standard.NewPeer(mockPeerID, mockTask, standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type))

			e := newEvaluatorDefault()
			tc.mock(parent, child)
			tc.expect(t, e.(*evaluatorDefault).evaluateParents(parent, child))
		})
	}
}

func TestEvaluatorDefault_calculateLoadQualityScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "perfect load quality score",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 4294967296
				peer.Host.TxBandwidth.Store(0)
				peer.Host.UploadContentLength.Store(0)
				peer.Host.ConcurrentUploadPieceCount.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(1.0, score)
			},
		},
		{
			name: "zero load quality score",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 0
				peer.Host.TxBandwidth.Store(0)
				peer.Host.UploadContentLength.Store(0)
				peer.Host.ConcurrentUploadPieceCount.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "high bandwidth usage",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 4294967296
				peer.Host.TxBandwidth.Store(8589934592)
				peer.Host.UploadContentLength.Store(0)
				peer.Host.ConcurrentUploadPieceCount.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.5, score)
			},
		},
		{
			name: "high upload content length",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 4294967296
				peer.Host.TxBandwidth.Store(0)
				peer.Host.UploadContentLength.Store(8589934592)
				peer.Host.ConcurrentUploadPieceCount.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.InDelta(0.91, score, 0.01)
			},
		},
		{
			name: "high concurrent upload count",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 4294967296
				peer.Host.TxBandwidth.Store(0)
				peer.Host.UploadContentLength.Store(0)
				peer.Host.ConcurrentUploadPieceCount.Store(200)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.InDelta(0.83, score, 0.01)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type))

			e := newEvaluatorDefault()
			tc.mock(peer)
			tc.expect(t, e.(*evaluatorDefault).calculateLoadQualityScore(peer))
		})
	}
}

func TestEvaluatorDefault_calculatePeakBandwidthUsageScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "no bandwidth usage",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000
				peer.Host.TxBandwidth.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(1.0, score)
			},
		},
		{
			name: "half bandwidth usage",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000
				peer.Host.TxBandwidth.Store(500)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.5, score)
			},
		},
		{
			name: "full bandwidth usage",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000
				peer.Host.TxBandwidth.Store(1000)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "over bandwidth capacity",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000
				peer.Host.TxBandwidth.Store(1500)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "zero max bandwidth",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 0
				peer.Host.TxBandwidth.Store(100)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "low bandwidth usage",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000
				peer.Host.TxBandwidth.Store(100)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.9, score)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type))

			e := newEvaluatorDefault()
			tc.mock(peer)
			tc.expect(t, e.(*evaluatorDefault).calculatePeakBandwidthUsageScore(peer))
		})
	}
}

func TestEvaluatorDefault_calculateBandwidthDurationScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "no upload content",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000
				peer.Host.UploadContentLength.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(1.0, score)
			},
		},
		{
			name: "zero max bandwidth",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 0
				peer.Host.UploadContentLength.Store(1000)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "low upload content length",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000000
				peer.Host.UploadContentLength.Store(1000)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Greater(score, 0.99)
			},
		},
		{
			name: "high upload content length",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000000
				peer.Host.UploadContentLength.Store(10000000)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "medium upload content length",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 1000000
				peer.Host.UploadContentLength.Store(3750000)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.5, score)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type))

			e := newEvaluatorDefault()
			tc.mock(peer)
			tc.expect(t, e.(*evaluatorDefault).calculateBandwidthDurationScore(peer))
		})
	}
}

func TestEvaluatorDefault_calculateConcurrencyScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "no concurrent uploads",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 400000000
				peer.Host.ConcurrentUploadPieceCount.Store(0)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(1.0, score)
			},
		},
		{
			name: "zero max bandwidth",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 0
				peer.Host.ConcurrentUploadPieceCount.Store(10)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "low bandwidth with concurrent uploads",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 100000
				peer.Host.ConcurrentUploadPieceCount.Store(5)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(0.0, score)
			},
		},
		{
			name: "high bandwidth with low concurrent uploads",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 400000000
				peer.Host.ConcurrentUploadPieceCount.Store(2)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(1.0, score)
			},
		},
		{
			name: "high bandwidth with high concurrent uploads",
			mock: func(peer *standard.Peer) {
				peer.Host.Network.MaxTxBandwidth = 400000000
				peer.Host.ConcurrentUploadPieceCount.Store(10)
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.InDelta(0.29, score, 0.01)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type))

			e := newEvaluatorDefault()
			tc.mock(peer)
			tc.expect(t, e.(*evaluatorDefault).calculateConcurrencyScore(peer))
		})
	}
}

func TestEvaluatorDefault_calculateHostTypeScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(peer *standard.Peer)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "host is normal peer",
			mock: func(peer *standard.Peer) {},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.5))
			},
		},
		{
			name: "host is seed peer but peer state is PeerStateSucceeded",
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateSucceeded)
				peer.Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "host is seed peer but peer state is PeerStateRunning",
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.Host.Type = types.HostTypeSuperSeed
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))
			peer := standard.NewPeer(mockPeerID, mockTask, mockHost)
			e := newEvaluatorDefault()
			tc.mock(peer)
			tc.expect(t, e.(*evaluatorDefault).calculateHostTypeScore(peer))
		})
	}
}

func TestEvaluatorDefault_calculateIDCAffinityScore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(dstHost *standard.Host, srcHost *standard.Host)
		expect func(t *testing.T, score float64)
	}{
		{
			name: "idc is empty",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = ""
				srcHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst host idc is empty",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "src host idc is empty",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				srcHost.Network.IDC = ""
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "idc is not the same",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = "foo"
				srcHost.Network.IDC = "bar"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "idc is the same",
			mock: func(dstHost *standard.Host, srcHost *standard.Host) {
				dstHost.Network.IDC = "example"
				srcHost.Network.IDC = "example"
			},
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dstHost := standard.NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			srcHost := standard.NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
			e := newEvaluatorDefault()
			tc.mock(dstHost, srcHost)
			tc.expect(t, e.(*evaluatorDefault).calculateIDCAffinityScore(dstHost.Network.IDC, srcHost.Network.IDC))
		})
	}
}

func TestEvaluatorDefault_calculateLocationAffinityScore(t *testing.T) {
	tests := []struct {
		name   string
		dst    string
		src    string
		expect func(t *testing.T, score float64)
	}{
		{
			name: "dst is empty and src is empty",
			dst:  "",
			src:  "",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst is empty",
			dst:  "",
			src:  "baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "src is empty",
			dst:  "bar",
			src:  "",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "has only one element and matches",
			dst:  "foo",
			src:  "foo",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "has only one element and does not match",
			dst:  "foo",
			src:  "bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "has multi element and match",
			dst:  "foo|bar",
			src:  "foo|bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
		{
			name: "has multi element and does not match",
			dst:  "foo|bar",
			src:  "bar|foo",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0))
			},
		},
		{
			name: "dst length is greater than src",
			dst:  "foo|bar|baz",
			src:  "foo|bar",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.4))
			},
		},
		{
			name: "src length is greater than dst",
			dst:  "foo|bar",
			src:  "foo|bar|baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.4))
			},
		},
		{
			name: "dst exceeds maximum length",
			dst:  "foo|bar|baz|bac|bae|baf",
			src:  "foo|bar|baz",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.6))
			},
		},
		{
			name: "src exceeds maximum length",
			dst:  "foo|bar|baz",
			src:  "foo|bar|baz|bac|bae|baf",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(0.6))
			},
		},
		{
			name: "dst and src both exceeds maximum length",
			dst:  "foo|bar|baz|bac|bae|baf",
			src:  "foo|bar|baz|bac|bae|bai",
			expect: func(t *testing.T, score float64) {
				assert := assert.New(t)
				assert.Equal(score, float64(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorDefault()
			tc.expect(t, e.(*evaluatorDefault).calculateLocationAffinityScore(tc.dst, tc.src))
		})
	}
}

func TestEvaluatorDefault_IsBadParent(t *testing.T) {
	mockHost := standard.NewHost(
		mockRawHost.ID, mockRawHost.IP, mockRawHost.Hostname,
		mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
	mockTask := standard.NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit, standard.WithDigest(mockTaskDigest))

	tests := []struct {
		name            string
		peer            *standard.Peer
		totalPieceCount int32
		mock            func(peer *standard.Peer)
		expect          func(t *testing.T, isBadParent bool)
	}{
		{
			name:            "peer state is PeerStateFailed",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateFailed)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateLeave",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateLeave)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStatePending",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStatePending)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateReceivedTiny",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateReceivedTiny)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateReceivedSmall",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateReceivedSmall)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "peer state is PeerStateReceivedNormal",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateReceivedNormal)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "download costs does not meet the normal distribution and last cost is twenty times more than mean",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.AppendPieceCost(10)
				peer.AppendPieceCost(201)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "download costs does not meet the normal distribution and last cost is twenty times lower than mean",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				peer.AppendPieceCost(10)
				peer.AppendPieceCost(200)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is too long",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				for i := range 30 {
					peer.AppendPieceCost(time.Duration(i))
				}
				peer.AppendPieceCost(50)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is normal",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				for i := range 30 {
					peer.AppendPieceCost(time.Duration(i))
				}
				peer.AppendPieceCost(18)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
		{
			name:            "download costs meet the normal distribution and last cost is too short",
			peer:            standard.NewPeer(mockPeerID, mockTask, mockHost),
			totalPieceCount: 1,
			mock: func(peer *standard.Peer) {
				peer.FSM.SetState(standard.PeerStateRunning)
				for i := 20; i < 50; i++ {
					peer.AppendPieceCost(time.Duration(i))
				}
				peer.AppendPieceCost(0)
			},
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorDefault()
			tc.mock(tc.peer)
			tc.expect(t, e.IsBadParent(tc.peer))
		})
	}
}
