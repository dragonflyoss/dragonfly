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

package config

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	managerv2 "d7y.io/api/v2/pkg/apis/manager/v2"

	"d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/rpc/manager/client/mocks"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
)

var (
	mockIDC      = "foo"
	mockLocation = "bar"

	mockRemoteConfig = &Config{
		DynConfig: DynConfig{
			RefreshInterval: 10 * time.Second,
		},
		Server: ServerConfig{
			Host: "localhost",
		},
		Manager: ManagerConfig{
			Addr:               &mockManagerAddr,
			SchedulerClusterID: 1,
		},
	}
)

func TestRemoteDynconfig_Get(t *testing.T) {
	mockManagerAddr := "localhost"
	mockConfig := &Config{
		DynConfig: DynConfig{},
		Server: ServerConfig{
			Host: "localhost",
		},
		Manager: ManagerConfig{
			Addr:               &mockManagerAddr,
			SchedulerClusterID: 1,
		},
	}

	tests := []struct {
		name            string
		refreshInterval time.Duration
		sleep           func()
		mock            func(m *mocks.MockV2MockRecorder)
		expect          func(t *testing.T, data *RemoteDynconfigData, err error)
	}{
		{
			name:            "get dynconfig success",
			refreshInterval: 10 * time.Second,
			sleep:           func() {},
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					Id:       1,
					Hostname: "foo",
					Idc:      &mockIDC,
					Location: &mockLocation,
					Ip:       "127.0.0.1",
					Port:     8002,
					State:    "active",
					SeedPeers: []*managerv2.SeedPeer{
						{
							Id:           1,
							Hostname:     "bar",
							Type:         pkgtypes.HostTypeSuperSeedName,
							Idc:          &mockIDC,
							Location:     &mockLocation,
							Ip:           "127.0.0.1",
							Port:         8001,
							DownloadPort: 8003,
							SeedPeerCluster: &managerv2.SeedPeerCluster{
								Id:     1,
								Name:   "baz",
								Config: []byte{1},
							},
						},
					},
					SchedulerCluster: &managerv2.SchedulerCluster{
						Id:           1,
						Name:         "bas",
						Config:       []byte{1},
						ClientConfig: []byte{1},
					},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{
					Applications: []*managerv2.Application{
						{
							Id:   1,
							Name: "foo",
							Url:  "example.com",
							Bio:  "bar",
							Priority: &managerv2.ApplicationPriority{
								Value: commonv2.Priority_LEVEL1,
								Urls: []*managerv2.URLPriority{
									{
										Regex: "blobs*",
										Value: commonv2.Priority_LEVEL1,
									},
								},
							},
						},
					},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, data *RemoteDynconfigData, err error) {
				assert := assert.New(t)
				assert.EqualValues(&RemoteDynconfigData{
					Scheduler: &managerv2.Scheduler{
						Id:       1,
						Hostname: "foo",
						Idc:      &mockIDC,
						Location: &mockLocation,
						Ip:       "127.0.0.1",
						Port:     8002,
						State:    "active",
						SeedPeers: []*managerv2.SeedPeer{
							{
								Id:           1,
								Hostname:     "bar",
								Type:         pkgtypes.HostTypeSuperSeedName,
								Idc:          &mockIDC,
								Location:     &mockLocation,
								Ip:           "127.0.0.1",
								Port:         8001,
								DownloadPort: 8003,
								SeedPeerCluster: &managerv2.SeedPeerCluster{
									Id:     1,
									Name:   "baz",
									Config: []byte{1},
								},
							},
						},
						SchedulerCluster: &managerv2.SchedulerCluster{
							Id:           1,
							Name:         "bas",
							Config:       []byte{1},
							ClientConfig: []byte{1},
						},
					},
					Applications: []*managerv2.Application{
						{
							Id:   1,
							Name: "foo",
							Url:  "example.com",
							Bio:  "bar",
							Priority: &managerv2.ApplicationPriority{
								Value: commonv2.Priority_LEVEL1,
								Urls: []*managerv2.URLPriority{
									{
										Regex: "blobs*",
										Value: commonv2.Priority_LEVEL1,
									},
								},
							},
						},
					},
				}, data)
			},
		},
		{
			name:            "get scheduler error",
			refreshInterval: 10 * time.Millisecond,
			sleep: func() {
				time.Sleep(100 * time.Millisecond)
			},
			mock: func(m *mocks.MockV2MockRecorder) {
				gomock.InOrder(
					m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
						Id:       1,
						Hostname: "foo",
						Idc:      &mockIDC,
						Location: &mockLocation,
						Ip:       "127.0.0.1",
						Port:     8002,
						State:    "active",
						SeedPeers: []*managerv2.SeedPeer{
							{
								Id:           1,
								Hostname:     "bar",
								Type:         pkgtypes.HostTypeSuperSeedName,
								Idc:          &mockIDC,
								Location:     &mockLocation,
								Ip:           "127.0.0.1",
								Port:         8001,
								DownloadPort: 8003,
								SeedPeerCluster: &managerv2.SeedPeerCluster{
									Id:     1,
									Name:   "baz",
									Config: []byte{1},
								},
							},
						},
						SchedulerCluster: &managerv2.SchedulerCluster{
							Id:           1,
							Name:         "bas",
							Config:       []byte{1},
							ClientConfig: []byte{1},
						},
					}, nil).Times(1),
					m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{
						Applications: []*managerv2.Application{
							{
								Id:   1,
								Name: "foo",
								Url:  "example.com",
								Bio:  "bar",
								Priority: &managerv2.ApplicationPriority{
									Value: commonv2.Priority_LEVEL1,
									Urls: []*managerv2.URLPriority{
										{
											Regex: "blobs*",
											Value: commonv2.Priority_LEVEL1,
										},
									},
								},
							},
						},
					}, nil).Times(1),
					m.GetScheduler(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, data *RemoteDynconfigData, err error) {
				assert := assert.New(t)
				assert.EqualValues(&RemoteDynconfigData{
					Scheduler: &managerv2.Scheduler{
						Id:       1,
						Hostname: "foo",
						Idc:      &mockIDC,
						Location: &mockLocation,
						Ip:       "127.0.0.1",
						Port:     8002,
						State:    "active",
						SeedPeers: []*managerv2.SeedPeer{
							{
								Id:           1,
								Hostname:     "bar",
								Type:         pkgtypes.HostTypeSuperSeedName,
								Idc:          &mockIDC,
								Location:     &mockLocation,
								Ip:           "127.0.0.1",
								Port:         8001,
								DownloadPort: 8003,
								SeedPeerCluster: &managerv2.SeedPeerCluster{
									Id:     1,
									Name:   "baz",
									Config: []byte{1},
								},
							},
						},
						SchedulerCluster: &managerv2.SchedulerCluster{
							Id:           1,
							Name:         "bas",
							Config:       []byte{1},
							ClientConfig: []byte{1},
						},
					},
					Applications: []*managerv2.Application{
						{
							Id:   1,
							Name: "foo",
							Url:  "example.com",
							Bio:  "bar",
							Priority: &managerv2.ApplicationPriority{
								Value: commonv2.Priority_LEVEL1,
								Urls: []*managerv2.URLPriority{
									{
										Regex: "blobs*",
										Value: commonv2.Priority_LEVEL1,
									},
								},
							},
						},
					},
				}, data)
			},
		},
		{
			name:            "list application error",
			refreshInterval: 10 * time.Millisecond,
			sleep: func() {
				time.Sleep(100 * time.Millisecond)
			},
			mock: func(m *mocks.MockV2MockRecorder) {
				gomock.InOrder(
					m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
						Id:       1,
						Hostname: "foo",
						Idc:      &mockIDC,
						Location: &mockLocation,
						Ip:       "127.0.0.1",
						Port:     8002,
						State:    "active",
						SeedPeers: []*managerv2.SeedPeer{
							{
								Id:           1,
								Hostname:     "bar",
								Type:         pkgtypes.HostTypeSuperSeedName,
								Idc:          &mockIDC,
								Location:     &mockLocation,
								Ip:           "127.0.0.1",
								Port:         8001,
								DownloadPort: 8003,
								SeedPeerCluster: &managerv2.SeedPeerCluster{
									Id:     1,
									Name:   "baz",
									Config: []byte{1},
								},
							},
						},
						SchedulerCluster: &managerv2.SchedulerCluster{
							Id:           1,
							Name:         "bas",
							Config:       []byte{1},
							ClientConfig: []byte{1},
						},
					}, nil).Times(1),
					m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{
						Applications: []*managerv2.Application{
							{
								Id:   1,
								Name: "foo",
								Url:  "example.com",
								Bio:  "bar",
								Priority: &managerv2.ApplicationPriority{
									Value: commonv2.Priority_LEVEL1,
									Urls: []*managerv2.URLPriority{
										{
											Regex: "blobs*",
											Value: commonv2.Priority_LEVEL1,
										},
									},
								},
							},
						},
					}, nil).Times(1),
					m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
						Id:       1,
						Hostname: "foo",
						Idc:      &mockIDC,
						Location: &mockLocation,
						Ip:       "127.0.0.1",
						Port:     8002,
						State:    "active",
						SeedPeers: []*managerv2.SeedPeer{
							{
								Id:           1,
								Hostname:     "bar",
								Type:         pkgtypes.HostTypeSuperSeedName,
								Idc:          &mockIDC,
								Location:     &mockLocation,
								Ip:           "127.0.0.1",
								Port:         8001,
								DownloadPort: 8003,
								SeedPeerCluster: &managerv2.SeedPeerCluster{
									Id:     1,
									Name:   "baz",
									Config: []byte{1},
								},
							},
						},
						SchedulerCluster: &managerv2.SchedulerCluster{
							Id:           1,
							Name:         "bas",
							Config:       []byte{1},
							ClientConfig: []byte{1},
						},
					}, nil).Times(1),
					m.ListApplications(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, data *RemoteDynconfigData, err error) {
				assert := assert.New(t)
				assert.EqualValues(&RemoteDynconfigData{
					Scheduler: &managerv2.Scheduler{
						Id:       1,
						Hostname: "foo",
						Idc:      &mockIDC,
						Location: &mockLocation,
						Ip:       "127.0.0.1",
						Port:     8002,
						State:    "active",
						SeedPeers: []*managerv2.SeedPeer{
							{
								Id:           1,
								Hostname:     "bar",
								Type:         pkgtypes.HostTypeSuperSeedName,
								Idc:          &mockIDC,
								Location:     &mockLocation,
								Ip:           "127.0.0.1",
								Port:         8001,
								DownloadPort: 8003,
								SeedPeerCluster: &managerv2.SeedPeerCluster{
									Id:     1,
									Name:   "baz",
									Config: []byte{1},
								},
							},
						},
						SchedulerCluster: &managerv2.SchedulerCluster{
							Id:           1,
							Name:         "bas",
							Config:       []byte{1},
							ClientConfig: []byte{1},
						},
					},
					Applications: []*managerv2.Application{
						{
							Id:   1,
							Name: "foo",
							Url:  "example.com",
							Bio:  "bar",
							Priority: &managerv2.ApplicationPriority{
								Value: commonv2.Priority_LEVEL1,
								Urls: []*managerv2.URLPriority{
									{
										Regex: "blobs*",
										Value: commonv2.Priority_LEVEL1,
									},
								},
							},
						},
					},
				}, data)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockManagerClient := mocks.NewMockV2(ctl)
			tc.mock(mockManagerClient.EXPECT())

			mockConfig.DynConfig.RefreshInterval = tc.refreshInterval
			d, err := NewDynconfig(mockManagerClient, "", mockConfig)
			if err != nil {
				t.Fatal(err)
			}

			rd, ok := d.(*remoteDynconfig)
			if !ok {
				t.Fatal("invalid remote dynconfig type")
			}

			tc.sleep()
			data, err := rd.Get()
			tc.expect(t, data, err)
		})
	}
}

func TestRemoteDynconfig_GetApplications(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *mocks.MockV2MockRecorder)
		expect func(t *testing.T, applications []*managerv2.Application, err error)
	}{
		{
			name: "applications found",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{
					Applications: []*managerv2.Application{{Id: 1, Name: "foo"}},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, applications []*managerv2.Application, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.EqualValues([]*managerv2.Application{{Id: 1, Name: "foo"}}, applications)
			},
		},
		{
			name: "application not found",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, applications []*managerv2.Application, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockManagerClient := mocks.NewMockV2(ctl)
			tc.mock(mockManagerClient.EXPECT())

			d, err := NewDynconfig(mockManagerClient, "", mockRemoteConfig)
			if err != nil {
				t.Fatal(err)
			}

			applications, err := d.GetApplications()
			tc.expect(t, applications, err)
		})
	}
}

func TestRemoteDynconfig_GetSeedPeerClusterConfig(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *mocks.MockV2MockRecorder)
		expect func(t *testing.T, config types.SeedPeerClusterConfig, err error)
	}{
		{
			name: "scheduler is nil",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SeedPeerClusterConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "seed peer not found",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SeedPeerClusterConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "seed peer cluster config is invalid json",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					SeedPeers: []*managerv2.SeedPeer{{SeedPeerCluster: &managerv2.SeedPeerCluster{Config: []byte("foo")}}},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SeedPeerClusterConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "seed peer cluster config is parsed",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					SeedPeers: []*managerv2.SeedPeer{{SeedPeerCluster: &managerv2.SeedPeerCluster{Config: []byte(`{"load_limit":100}`)}}},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SeedPeerClusterConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.SeedPeerClusterConfig{LoadLimit: 100}, config)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockManagerClient := mocks.NewMockV2(ctl)
			tc.mock(mockManagerClient.EXPECT())

			d, err := NewDynconfig(mockManagerClient, "", mockRemoteConfig)
			if err != nil {
				t.Fatal(err)
			}

			config, err := d.GetSeedPeerClusterConfig()
			tc.expect(t, config, err)
		})
	}
}

func TestRemoteDynconfig_GetSchedulerClusterConfig(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *mocks.MockV2MockRecorder)
		expect func(t *testing.T, config types.SchedulerClusterConfig, err error)
	}{
		{
			name: "scheduler cluster is nil",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SchedulerClusterConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "scheduler cluster config is invalid json",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					SchedulerCluster: &managerv2.SchedulerCluster{Config: []byte("foo")},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SchedulerClusterConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "scheduler cluster config is parsed",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					SchedulerCluster: &managerv2.SchedulerCluster{Config: []byte(`{"candidate_parent_limit":5,"filter_parent_limit":20}`)},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SchedulerClusterConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.SchedulerClusterConfig{CandidateParentLimit: 5, FilterParentLimit: 20}, config)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockManagerClient := mocks.NewMockV2(ctl)
			tc.mock(mockManagerClient.EXPECT())

			d, err := NewDynconfig(mockManagerClient, "", mockRemoteConfig)
			if err != nil {
				t.Fatal(err)
			}

			config, err := d.GetSchedulerClusterConfig()
			tc.expect(t, config, err)
		})
	}
}

func TestRemoteDynconfig_GetSchedulerClusterClientConfig(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *mocks.MockV2MockRecorder)
		expect func(t *testing.T, config types.SchedulerClusterClientConfig, err error)
	}{
		{
			name: "scheduler cluster is nil",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SchedulerClusterClientConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "scheduler cluster client config is invalid json",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					SchedulerCluster: &managerv2.SchedulerCluster{ClientConfig: []byte("foo")},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SchedulerClusterClientConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "scheduler cluster client config is parsed",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{
					SchedulerCluster: &managerv2.SchedulerCluster{ClientConfig: []byte(`{"load_limit":50}`)},
				}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{}, nil).Times(1)
			},
			expect: func(t *testing.T, config types.SchedulerClusterClientConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.SchedulerClusterClientConfig{LoadLimit: 50}, config)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockManagerClient := mocks.NewMockV2(ctl)
			tc.mock(mockManagerClient.EXPECT())

			d, err := NewDynconfig(mockManagerClient, "", mockRemoteConfig)
			if err != nil {
				t.Fatal(err)
			}

			config, err := d.GetSchedulerClusterClientConfig()
			tc.expect(t, config, err)
		})
	}
}

func TestManagerClient_Get(t *testing.T) {
	tests := []struct {
		name   string
		mock   func(m *mocks.MockV2MockRecorder)
		expect func(t *testing.T, data any, err error)
	}{
		{
			name: "get scheduler failed",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, data any, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(data)
			},
		},
		{
			name: "list applications unimplemented on old manager is tolerated",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{Id: 1}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Unimplemented, "foo")).Times(1)
			},
			expect: func(t *testing.T, data any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.EqualValues(RemoteDynconfigData{Scheduler: &managerv2.Scheduler{Id: 1}}, data)
			},
		},
		{
			name: "list applications not found is tolerated",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{Id: 1}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.NotFound, "foo")).Times(1)
			},
			expect: func(t *testing.T, data any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.EqualValues(RemoteDynconfigData{Scheduler: &managerv2.Scheduler{Id: 1}}, data)
			},
		},
		{
			name: "list applications failed with other status code",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{Id: 1}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(nil, status.Error(codes.Internal, "foo")).Times(1)
			},
			expect: func(t *testing.T, data any, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, status.Error(codes.Internal, "foo"))
				assert.Nil(data)
			},
		},
		{
			name: "list applications failed with plain error",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{Id: 1}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, data any, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(data)
			},
		},
		{
			name: "scheduler and applications are returned",
			mock: func(m *mocks.MockV2MockRecorder) {
				m.GetScheduler(gomock.Any(), gomock.Any()).Return(&managerv2.Scheduler{Id: 1}, nil).Times(1)
				m.ListApplications(gomock.Any(), gomock.Any()).Return(&managerv2.ListApplicationsResponse{
					Applications: []*managerv2.Application{{Id: 1, Name: "foo"}},
				}, nil).Times(1)
			},
			expect: func(t *testing.T, data any, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.EqualValues(RemoteDynconfigData{
					Scheduler:    &managerv2.Scheduler{Id: 1},
					Applications: []*managerv2.Application{{Id: 1, Name: "foo"}},
				}, data)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockManagerClient := mocks.NewMockV2(ctl)
			tc.mock(mockManagerClient.EXPECT())

			data, err := newManagerClient(mockManagerClient, mockRemoteConfig).Get()
			tc.expect(t, data, err)
		})
	}
}
