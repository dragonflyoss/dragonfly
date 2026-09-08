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

package standard

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
)

var (
	mockHostGCConfig = &config.GCConfig{
		HostGCInterval: 1 * time.Second,
	}
)

func TestHostManager_newHostManager(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hostManager HostManager, err error)
	}{
		{
			name: "new host manager",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("hostManager", reflect.TypeOf(hostManager).Elem().Name())
			},
		},
		{
			name: "new host manager failed because of gc error",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())
			hostManager, err := newHostManager(mockHostGCConfig, gc)

			tc.expect(t, hostManager, err)
		})
	}
}

func TestHostManager_Load(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hostManager HostManager, mockHost *Host)
	}{
		{
			name: "load host",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				hostManager.Store(mockHost)
				host, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
		{
			name: "host does not exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				_, loaded := hostManager.Load(mockHost.ID)
				assert.False(loaded)
			},
		},
		{
			name: "load key is empty",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				mockHost.ID = ""
				hostManager.Store(mockHost)
				host, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, hostManager, mockHost)
		})
	}
}

func TestHostManager_Store(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hostManager HostManager, mockHost *Host)
	}{
		{
			name: "store host",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				hostManager.Store(mockHost)
				host, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
		{
			name: "store key is empty",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				mockHost.ID = ""
				hostManager.Store(mockHost)
				host, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, hostManager, mockHost)
		})
	}
}

func TestHostManager_LoadOrStore(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hostManager HostManager, mockHost *Host)
	}{
		{
			name: "load host exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				hostManager.Store(mockHost)
				host, loaded := hostManager.LoadOrStore(mockHost)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
		{
			name: "load host does not exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				host, loaded := hostManager.LoadOrStore(mockHost)
				assert.False(loaded)
				assert.Equal(mockHost.ID, host.ID)
				assert.Len(hostManager.LoadAllNormals(), 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, hostManager, mockHost)
		})
	}
}

func TestHostManager_Delete(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hostManager HostManager, mockHost *Host)
	}{
		{
			name: "delete host",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				hostManager.Store(mockHost)
				hostManager.Delete(mockHost.ID)
				_, loaded := hostManager.Load(mockHost.ID)
				assert.False(loaded)
				assert.Empty(hostManager.LoadAllNormals())
			},
		},
		{
			name: "delete key does not exist",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host) {
				assert := assert.New(t)
				mockHost.ID = ""
				hostManager.Store(mockHost)
				hostManager.Delete(mockHost.ID)
				_, loaded := hostManager.Load(mockHost.ID)
				assert.False(loaded)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, hostManager, mockHost)
		})
	}
}

func TestHostManager_Range(t *testing.T) {
	tests := []struct {
		name   string
		expect func(t *testing.T, hostManager HostManager, mockHost *Host, mockSeedHost *Host)
	}{
		{
			name: "range visits every host",
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockSeedHost *Host) {
				assert := assert.New(t)
				var ids []string
				hostManager.Range(func(_, value any) bool {
					ids = append(ids, value.(*Host).ID)
					return true
				})

				assert.ElementsMatch([]string{mockHost.ID, mockSeedHost.ID}, ids)
			},
		},
		{
			name: "range stops when f returns false",
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockSeedHost *Host) {
				assert := assert.New(t)
				var ids []string
				hostManager.Range(func(_, value any) bool {
					ids = append(ids, value.(*Host).ID)
					return false
				})

				assert.Len(ids, 1)
			},
		},
		{
			name: "range normals visits only normal hosts",
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockSeedHost *Host) {
				assert := assert.New(t)
				var ids []string
				hostManager.RangeNormals(func(_, value any) bool {
					ids = append(ids, value.(*Host).ID)
					return true
				})

				assert.Equal([]string{mockHost.ID}, ids)
			},
		},
		{
			name: "range seeds visits only seed hosts",
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockSeedHost *Host) {
				assert := assert.New(t)
				var ids []string
				hostManager.RangeSeeds(func(_, value any) bool {
					ids = append(ids, value.(*Host).ID)
					return true
				})

				assert.Equal([]string{mockSeedHost.ID}, ids)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			gc.EXPECT().Add(gomock.Any()).Return(nil).Times(1)

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockSeedHost := NewHost(
				mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
				mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			hostManager.Store(mockHost)
			hostManager.Store(mockSeedHost)
			tc.expect(t, hostManager, mockHost, mockSeedHost)
		})
	}
}

func TestHostManager_LoadAll(t *testing.T) {
	tests := []struct {
		name   string
		hosts  []*Host
		expect func(t *testing.T, hostManager HostManager, hosts []*Host)
	}{
		{
			name:  "map is empty",
			hosts: []*Host{},
			expect: func(t *testing.T, hostManager HostManager, hosts []*Host) {
				assert := assert.New(t)
				assert.Equal(0, hostManager.Len())
				assert.Empty(hostManager.LoadAll())
				assert.Empty(hostManager.LoadAllNormals())
				assert.Empty(hostManager.LoadAllSeeds())
			},
		},
		{
			name: "normal and seed hosts are grouped by type",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type),
				NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type),
			},
			expect: func(t *testing.T, hostManager HostManager, hosts []*Host) {
				assert := assert.New(t)
				assert.Equal(2, hostManager.Len())
				assert.ElementsMatch(hosts, hostManager.LoadAll())
				assert.Equal([]*Host{hosts[0]}, hostManager.LoadAllNormals())
				assert.Equal([]*Host{hosts[1]}, hostManager.LoadAllSeeds())
			},
		},
		{
			name: "host with unknown type is stored only in the all map",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, types.HostType(100)),
			},
			expect: func(t *testing.T, hostManager HostManager, hosts []*Host) {
				assert := assert.New(t)
				assert.Equal(1, hostManager.Len())
				assert.Equal(hosts, hostManager.LoadAll())
				assert.Empty(hostManager.LoadAllNormals())
				assert.Empty(hostManager.LoadAllSeeds())
			},
		},
		{
			name: "deleted host is removed from every map",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type),
				NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type),
			},
			expect: func(t *testing.T, hostManager HostManager, hosts []*Host) {
				assert := assert.New(t)
				hostManager.Delete(hosts[1].ID)
				assert.Equal(1, hostManager.Len())
				assert.Equal([]*Host{hosts[0]}, hostManager.LoadAll())
				assert.Equal([]*Host{hosts[0]}, hostManager.LoadAllNormals())
				assert.Empty(hostManager.LoadAllSeeds())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			gc.EXPECT().Add(gomock.Any()).Return(nil).Times(1)

			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			for _, host := range tc.hosts {
				hostManager.Store(host)
			}

			tc.expect(t, hostManager, tc.hosts)
		})
	}
}

func TestHostManager_LoadRandom(t *testing.T) {
	tests := []struct {
		name   string
		hosts  []*Host
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hm HostManager, hosts []*Host)
	}{
		{
			name: "load random hosts",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type),
				NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type),
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hm HostManager, hosts []*Host) {
				assert := assert.New(t)
				for _, host := range hosts {
					hm.Store(host)
				}

				blocklist := set.NewSafeSet[string]()
				blocklist.Add(mockRawSeedHost.ID)
				h := hm.LoadRandom(2, blocklist)
				assert.Len(h, 1)
				assert.Equal(mockRawHost.ID, h[0].ID)
			},
		},
		{
			name: "load random hosts when the load number is 0",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type),
				NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type),
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hm HostManager, hosts []*Host) {
				assert := assert.New(t)
				for _, host := range hosts {
					hm.Store(host)
				}

				blocklist := set.NewSafeSet[string]()
				blocklist.Add(mockRawSeedHost.ID)
				assert.Empty(hm.LoadRandom(0, blocklist))
			},
		},
		{
			name:  "map is empty",
			hosts: []*Host{},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hm HostManager, hosts []*Host) {
				assert := assert.New(t)
				for _, host := range hosts {
					hm.Store(host)
				}

				blocklist := set.NewSafeSet[string]()
				blocklist.Add(mockRawSeedHost.ID)
				assert.Empty(hm.LoadRandom(1, blocklist))
			},
		},
		{
			name: "the number of hosts in the map is insufficient",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type),
				NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type),
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hm HostManager, hosts []*Host) {
				assert := assert.New(t)
				for _, host := range hosts {
					hm.Store(host)
				}

				blocklist := set.NewSafeSet[string]()
				blocklist.Add(mockRawSeedHost.ID)
				assert.Len(hm.LoadRandom(3, blocklist), 1)
			},
		},
		{
			name: "load random hosts stops at the requested number",
			hosts: []*Host{
				NewHost(
					mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
					mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type),
				NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type),
			},
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hm HostManager, hosts []*Host) {
				assert := assert.New(t)
				for _, host := range hosts {
					hm.Store(host)
				}

				assert.Len(hm.LoadRandom(1, set.NewSafeSet[string]()), 1)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			hm, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, hm, tc.hosts)
		})
	}
}

func TestHostManager_RunGC(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *gc.MockGCMockRecorder)
		expect func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer)
	}{
		{
			name: "host has peers",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				hostManager.Store(mockHost)
				mockHost.StorePeer(mockPeer)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				host, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
		{
			name: "host has upload peers",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				hostManager.Store(mockHost)
				mockHost.StorePeer(mockPeer)
				mockHost.PeerCount.Add(0)
				mockHost.ConcurrentUploadCount.Add(1)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				host, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Equal(mockHost.ID, host.ID)
			},
		},
		{
			name: "host is seed peer",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				mockSeedHost := NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type)
				hostManager.Store(mockSeedHost)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				host, loaded := hostManager.Load(mockSeedHost.ID)
				assert.True(loaded)
				assert.Equal(mockSeedHost.ID, host.ID)
			},
		},
		{
			name: "host with zero announce interval is never reclaimed",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				mockHost.UpdatedAt.Store(time.Now().Add(-time.Hour))
				hostManager.Store(mockHost)
				mockHost.StorePeer(mockPeer)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.True(mockPeer.FSM.Is(PeerStatePending))
			},
		},
		{
			name: "host elapsed is within twice the announce interval",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				mockHost.AnnounceInterval = time.Hour
				mockHost.UpdatedAt.Store(time.Now().Add(-time.Hour))
				hostManager.Store(mockHost)
				mockHost.StorePeer(mockPeer)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				_, loaded := hostManager.Load(mockHost.ID)
				assert.True(loaded)
				assert.Len(hostManager.LoadAllNormals(), 1)
				assert.True(mockPeer.FSM.Is(PeerStatePending))
			},
		},
		{
			name: "host elapsed exceeds twice the announce interval",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				mockHost.AnnounceInterval = 1 * time.Microsecond
				hostManager.Store(mockHost)
				mockHost.StorePeer(mockPeer)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				mockHost.Peers.Range(func(_, value any) bool {
					peer := value.(*Peer)
					assert.True(peer.FSM.Is(PeerStateLeave))
					return true
				})

				_, loaded := hostManager.Load(mockHost.ID)
				assert.False(loaded)
				assert.Empty(hostManager.LoadAllNormals())
			},
		},
		{
			name: "seed host elapsed exceeds twice the announce interval",
			mock: func(m *gc.MockGCMockRecorder) {
				m.Add(gomock.Any()).Return(nil).Times(1)
			},
			expect: func(t *testing.T, hostManager HostManager, mockHost *Host, mockPeer *Peer) {
				assert := assert.New(t)
				mockSeedHost := NewHost(
					mockRawSeedHost.ID, mockRawSeedHost.IP, mockRawSeedHost.Name, mockRawSeedHost.Hostname,
					mockRawSeedHost.Port, mockRawSeedHost.DownloadPort, mockRawSeedHost.ProxyPort, mockRawSeedHost.Type,
					WithAnnounceInterval(1*time.Microsecond))
				mockSeedPeer := NewPeer(mockSeedPeerID, mockPeer.Task, mockSeedHost)
				hostManager.Store(mockSeedHost)
				mockSeedHost.StorePeer(mockSeedPeer)
				err := hostManager.RunGC(context.Background())
				assert.NoError(err)

				assert.True(mockSeedPeer.FSM.Is(PeerStateLeave))
				_, loaded := hostManager.Load(mockSeedHost.ID)
				assert.False(loaded)
				assert.Empty(hostManager.LoadAllSeeds())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			mockHost := NewHost(
				mockRawHost.ID, mockRawHost.IP, mockRawHost.Name, mockRawHost.Hostname,
				mockRawHost.Port, mockRawHost.DownloadPort, mockRawHost.ProxyPort, mockRawHost.Type)
			mockTask := NewTask(mockTaskID, mockTaskURL, mockTaskTag, mockTaskApplication, commonv2.TaskType_STANDARD, mockTaskFilteredQueryParams, mockTaskHeader, mockTaskBackToSourceLimit)
			mockPeer := NewPeer(mockPeerID, mockTask, mockHost)
			hostManager, err := newHostManager(mockHostGCConfig, gc)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, hostManager, mockHost, mockPeer)
		})
	}
}
