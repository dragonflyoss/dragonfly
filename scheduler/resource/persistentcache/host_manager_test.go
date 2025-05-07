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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/container/set"
	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
)

func TestHostManager_Load(t *testing.T) {
	tests := []struct {
		name           string
		hostID         string
		mockRedis      func(mock redismock.ClientMock)
		expectedHost   *Host
		expectedLoaded bool
		expectedError  bool
		errorMsg       string
	}{
		{
			name:   "host exists in Redis",
			hostID: "host1",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":                                  "host1",
					"type":                                "normal",
					"hostname":                            "hostname1",
					"ip":                                  "127.0.0.1",
					"port":                                "8080",
					"download_port":                       "8081",
					"disable_shared":                      "false",
					"os":                                  "linux",
					"platform":                            "x86_64",
					"platform_family":                     "debian",
					"platform_version":                    "11",
					"kernel_version":                      "5.10",
					"cpu_logical_count":                   "4",
					"cpu_physical_count":                  "2",
					"cpu_percent":                         "50.0",
					"cpu_processe_percent":                "25.0",
					"cpu_times_user":                      "10.0",
					"cpu_times_system":                    "5.0",
					"cpu_times_idle":                      "100.0",
					"cpu_times_nice":                      "0.0",
					"cpu_times_iowait":                    "1.0",
					"cpu_times_irq":                       "0.5",
					"cpu_times_softirq":                   "0.2",
					"cpu_times_steal":                     "0.1",
					"cpu_times_guest":                     "0.0",
					"cpu_times_guest_nice":                "0.0",
					"memory_total":                        "8000000000",
					"memory_available":                    "4000000000",
					"memory_used":                         "4000000000",
					"memory_used_percent":                 "50.0",
					"memory_processe_used_percent":        "25.0",
					"memory_free":                         "2000000000",
					"network_tcp_connection_count":        "100",
					"network_upload_tcp_connection_count": "50",
					"network_location":                    "location1",
					"network_idc":                         "idc1",
					"network_download_rate":               "1000000",
					"network_download_rate_limit":         "2000000",
					"network_upload_rate":                 "500000",
					"network_upload_rate_limit":           "1000000",
					"disk_total":                          "100000000000",
					"disk_free":                           "50000000000",
					"disk_used":                           "50000000000",
					"disk_used_percent":                   "50.0",
					"disk_inodes_total":                   "100000",
					"disk_inodes_used":                    "50000",
					"disk_inodes_free":                    "50000",
					"disk_inodes_used_percent":            "50.0",
					"disk_write_bandwidth":                "10000000",
					"disk_read_bandwidth":                 "20000000",
					"build_git_version":                   "v1.0.0",
					"build_git_commit":                    "commit1",
					"build_go_version":                    "1.16",
					"build_platform":                      "linux/amd64",
					"scheduler_cluster_id":                "1",
					"announce_interval":                   "300",
					"created_at":                          time.Now().Format(time.RFC3339),
					"updated_at":                          time.Now().Format(time.RFC3339),
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host1")).SetVal(mockData)
			},
			expectedHost: &Host{
				ID:                 "host1",
				Type:               pkgtypes.HostTypeNormal,
				Hostname:           "hostname1",
				IP:                 "127.0.0.1",
				Port:               8080,
				DownloadPort:       8081,
				DisableShared:      false,
				OS:                 "linux",
				Platform:           "x86_64",
				PlatformFamily:     "debian",
				PlatformVersion:    "11",
				KernelVersion:      "5.10",
				CPU:                CPU{LogicalCount: 4, PhysicalCount: 2, Percent: 50.0, ProcessPercent: 25.0, Times: CPUTimes{User: 10.0, System: 5.0, Idle: 100.0, Nice: 0.0, Iowait: 1.0, Irq: 0.5, Softirq: 0.2, Steal: 0.1, Guest: 0.0, GuestNice: 0.0}},
				Memory:             Memory{Total: 8000000000, Available: 4000000000, Used: 4000000000, UsedPercent: 50.0, ProcessUsedPercent: 25.0, Free: 2000000000},
				Network:            Network{TCPConnectionCount: 100, UploadTCPConnectionCount: 50, Location: "location1", IDC: "idc1", DownloadRate: 1000000, DownloadRateLimit: 2000000, UploadRate: 500000, UploadRateLimit: 1000000},
				Disk:               Disk{Total: 100000000000, Free: 50000000000, Used: 50000000000, UsedPercent: 50.0, InodesTotal: 100000, InodesUsed: 50000, InodesFree: 50000, InodesUsedPercent: 50.0, WriteBandwidth: 10000000, ReadBandwidth: 20000000},
				Build:              Build{GitVersion: "v1.0.0", GitCommit: "commit1", GoVersion: "1.16", Platform: "linux/amd64"},
				SchedulerClusterID: 1,
				AnnounceInterval:   time.Duration(300),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				Log:                logger.WithHost("host1", "hostname1", "127.0.0.1"),
			},
			expectedLoaded: true,
		},
		{
			name:   "host does not exist in Redis",
			hostID: "host2",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host2")).SetVal(map[string]string{})
			},
			expectedHost:   nil,
			expectedLoaded: false,
		},
		{
			name:   "redis returns an error",
			hostID: "host3",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host3")).SetErr(fmt.Errorf("Redis error"))
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid port value",
			hostID: "host4",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":   "host4",
					"port": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host4")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid scheduler_cluster_id value",
			hostID: "host5",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":                   "host5",
					"scheduler_cluster_id": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host5")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid disable_shared value",
			hostID: "host6",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":             "host6",
					"disable_shared": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host6")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid cpu_logical_count value",
			hostID: "host7",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":                "host7",
					"cpu_logical_count": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host7")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid cpu_percent value",
			hostID: "host8",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":          "host8",
					"cpu_percent": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host8")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid memory_total value",
			hostID: "host9",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":           "host9",
					"memory_total": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host9")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid network_tcp_connection_count value",
			hostID: "host10",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":                           "host10",
					"network_tcp_connection_count": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host10")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid disk_total value",
			hostID: "host11",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":         "host11",
					"disk_total": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host11")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid announce_interval value",
			hostID: "host12",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":                "host12",
					"announce_interval": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host12")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid created_at value",
			hostID: "host13",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":         "host13",
					"created_at": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host13")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
		{
			name:   "invalid updated_at value",
			hostID: "host14",
			mockRedis: func(mock redismock.ClientMock) {
				mockData := map[string]string{
					"id":         "host14",
					"updated_at": "invalid",
				}
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host14")).SetVal(mockData)
			},
			expectedHost:   nil,
			expectedLoaded: false,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			h := &hostManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 1,
					},
				},
				rdb: rdb,
			}

			host, loaded := h.Load(context.Background(), tt.hostID)
			if tt.expectedError {
				assert.Error(t, nil)
				assert.Contains(t, "", tt.errorMsg)
			} else {
				assert.NoError(t, nil)
			}

			if tt.expectedLoaded {
				assert.NotNil(t, host)
				assert.Equal(t, tt.expectedHost.ID, host.ID)
				assert.Equal(t, tt.expectedHost.Type, host.Type)
				assert.Equal(t, tt.expectedHost.Hostname, host.Hostname)
				assert.Equal(t, tt.expectedHost.IP, host.IP)
				assert.Equal(t, tt.expectedHost.Port, host.Port)
				assert.Equal(t, tt.expectedHost.DownloadPort, host.DownloadPort)
				assert.Equal(t, tt.expectedHost.DisableShared, host.DisableShared)
				assert.Equal(t, tt.expectedHost.OS, host.OS)
				assert.Equal(t, tt.expectedHost.Platform, host.Platform)
				assert.Equal(t, tt.expectedHost.PlatformFamily, host.PlatformFamily)
				assert.Equal(t, tt.expectedHost.PlatformVersion, host.PlatformVersion)
				assert.Equal(t, tt.expectedHost.KernelVersion, host.KernelVersion)
				assert.Equal(t, tt.expectedHost.CPU, host.CPU)
				assert.Equal(t, tt.expectedHost.Memory, host.Memory)
				assert.Equal(t, tt.expectedHost.Network, host.Network)
				assert.Equal(t, tt.expectedHost.Disk, host.Disk)
				assert.Equal(t, tt.expectedHost.Build, host.Build)
				assert.Equal(t, tt.expectedHost.SchedulerClusterID, host.SchedulerClusterID)
				assert.Equal(t, tt.expectedHost.AnnounceInterval.Abs(), host.AnnounceInterval.Abs())
				assert.Equal(t, tt.expectedHost.CreatedAt.Format(time.RFC3339), host.CreatedAt.Format(time.RFC3339))
				assert.Equal(t, tt.expectedHost.UpdatedAt.Format(time.RFC3339), host.UpdatedAt.Format(time.RFC3339))
			}

			assert.Equal(t, tt.expectedLoaded, loaded)
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestHostManager_LoadAll(t *testing.T) {
	tests := []struct {
		name          string
		mockRedis     func(mock redismock.ClientMock)
		expectedError bool
		expectedHosts int
	}{
		{
			name: "scan fails",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSScan(pkgredis.MakePersistentCacheHostsInScheduler(1), 0, "*", 10).
					SetErr(fmt.Errorf("redis scan error"))
			},
			expectedError: true,
			expectedHosts: 0,
		},
		{
			name: "some hosts fail to load",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSScan(pkgredis.MakePersistentCacheHostsInScheduler(1), 0, "*", 10).
					SetVal([]string{"host1", "host2"}, 0)
				// host1 loaded successfully
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host1")).
					SetVal(map[string]string{
						"id":                                  "host1",
						"type":                                "normal",
						"hostname":                            "hostname1",
						"ip":                                  "127.0.0.1",
						"port":                                "8080",
						"download_port":                       "8081",
						"disable_shared":                      "false",
						"os":                                  "linux",
						"platform":                            "x86_64",
						"platform_family":                     "debian",
						"platform_version":                    "11",
						"kernel_version":                      "5.10",
						"cpu_logical_count":                   "4",
						"cpu_physical_count":                  "2",
						"cpu_percent":                         "50.0",
						"cpu_processe_percent":                "25.0",
						"cpu_times_user":                      "10.0",
						"cpu_times_system":                    "5.0",
						"cpu_times_idle":                      "100.0",
						"cpu_times_nice":                      "0.0",
						"cpu_times_iowait":                    "1.0",
						"cpu_times_irq":                       "0.5",
						"cpu_times_softirq":                   "0.2",
						"cpu_times_steal":                     "0.1",
						"cpu_times_guest":                     "0.0",
						"cpu_times_guest_nice":                "0.0",
						"memory_total":                        "8000000000",
						"memory_available":                    "4000000000",
						"memory_used":                         "4000000000",
						"memory_used_percent":                 "50.0",
						"memory_processe_used_percent":        "25.0",
						"memory_free":                         "2000000000",
						"network_tcp_connection_count":        "100",
						"network_upload_tcp_connection_count": "50",
						"network_location":                    "location1",
						"network_idc":                         "idc1",
						"network_download_rate":               "1000000",
						"network_download_rate_limit":         "2000000",
						"network_upload_rate":                 "500000",
						"network_upload_rate_limit":           "1000000",
						"disk_total":                          "100000000000",
						"disk_free":                           "50000000000",
						"disk_used":                           "50000000000",
						"disk_used_percent":                   "50.0",
						"disk_inodes_total":                   "100000",
						"disk_inodes_used":                    "50000",
						"disk_inodes_free":                    "50000",
						"disk_inodes_used_percent":            "50.0",
						"disk_write_bandwidth":                "10000000",
						"disk_read_bandwidth":                 "20000000",
						"build_git_version":                   "v1.0.0",
						"build_git_commit":                    "commit1",
						"build_go_version":                    "1.16",
						"build_platform":                      "linux/amd64",
						"scheduler_cluster_id":                "1",
						"announce_interval":                   "0",
						"updated_at":                          time.Now().Format(time.RFC3339),
						"created_at":                          time.Now().Format(time.RFC3339),
					})
				// host2 load fails
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host2")).
					SetErr(fmt.Errorf("redis hgetall error"))
			},
			expectedError: false,
			expectedHosts: 1,
		},
		{
			name: "multiple scans, all loaded successfully",
			mockRedis: func(mock redismock.ClientMock) {
				// First scan
				mock.ExpectSScan(pkgredis.MakePersistentCacheHostsInScheduler(1), 0, "*", 10).
					SetVal([]string{"host3"}, 123)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host3")).
					SetVal(map[string]string{
						"id":                                  "host3",
						"type":                                "normal",
						"hostname":                            "hostname1",
						"ip":                                  "127.0.0.1",
						"port":                                "8080",
						"download_port":                       "8081",
						"disable_shared":                      "false",
						"os":                                  "linux",
						"platform":                            "x86_64",
						"platform_family":                     "debian",
						"platform_version":                    "11",
						"kernel_version":                      "5.10",
						"cpu_logical_count":                   "4",
						"cpu_physical_count":                  "2",
						"cpu_percent":                         "50.0",
						"cpu_processe_percent":                "25.0",
						"cpu_times_user":                      "10.0",
						"cpu_times_system":                    "5.0",
						"cpu_times_idle":                      "100.0",
						"cpu_times_nice":                      "0.0",
						"cpu_times_iowait":                    "1.0",
						"cpu_times_irq":                       "0.5",
						"cpu_times_softirq":                   "0.2",
						"cpu_times_steal":                     "0.1",
						"cpu_times_guest":                     "0.0",
						"cpu_times_guest_nice":                "0.0",
						"memory_total":                        "8000000000",
						"memory_available":                    "4000000000",
						"memory_used":                         "4000000000",
						"memory_used_percent":                 "50.0",
						"memory_processe_used_percent":        "25.0",
						"memory_free":                         "2000000000",
						"network_tcp_connection_count":        "100",
						"network_upload_tcp_connection_count": "50",
						"network_location":                    "location1",
						"network_idc":                         "idc1",
						"network_download_rate":               "1000000",
						"network_download_rate_limit":         "2000000",
						"network_upload_rate":                 "500000",
						"network_upload_rate_limit":           "1000000",
						"disk_total":                          "100000000000",
						"disk_free":                           "50000000000",
						"disk_used":                           "50000000000",
						"disk_used_percent":                   "50.0",
						"disk_inodes_total":                   "100000",
						"disk_inodes_used":                    "50000",
						"disk_inodes_free":                    "50000",
						"disk_inodes_used_percent":            "50.0",
						"disk_write_bandwidth":                "10000000",
						"disk_read_bandwidth":                 "20000000",
						"build_git_version":                   "v1.0.0",
						"build_git_commit":                    "commit1",
						"build_go_version":                    "1.16",
						"build_platform":                      "linux/amd64",
						"scheduler_cluster_id":                "1",
						"announce_interval":                   "0",
						"updated_at":                          time.Now().Format(time.RFC3339),
						"created_at":                          time.Now().Format(time.RFC3339),
					})
				// Second scan
				mock.ExpectSScan(pkgredis.MakePersistentCacheHostsInScheduler(1), 123, "*", 10).
					SetVal([]string{"host4"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host4")).
					SetVal(map[string]string{
						"id":                                  "host4",
						"type":                                "normal",
						"hostname":                            "hostname1",
						"ip":                                  "127.0.0.1",
						"port":                                "8080",
						"download_port":                       "8081",
						"disable_shared":                      "false",
						"os":                                  "linux",
						"platform":                            "x86_64",
						"platform_family":                     "debian",
						"platform_version":                    "11",
						"kernel_version":                      "5.10",
						"cpu_logical_count":                   "4",
						"cpu_physical_count":                  "2",
						"cpu_percent":                         "50.0",
						"cpu_processe_percent":                "25.0",
						"cpu_times_user":                      "10.0",
						"cpu_times_system":                    "5.0",
						"cpu_times_idle":                      "100.0",
						"cpu_times_nice":                      "0.0",
						"cpu_times_iowait":                    "1.0",
						"cpu_times_irq":                       "0.5",
						"cpu_times_softirq":                   "0.2",
						"cpu_times_steal":                     "0.1",
						"cpu_times_guest":                     "0.0",
						"cpu_times_guest_nice":                "0.0",
						"memory_total":                        "8000000000",
						"memory_available":                    "4000000000",
						"memory_used":                         "4000000000",
						"memory_used_percent":                 "50.0",
						"memory_processe_used_percent":        "25.0",
						"memory_free":                         "2000000000",
						"network_tcp_connection_count":        "100",
						"network_upload_tcp_connection_count": "50",
						"network_location":                    "location1",
						"network_idc":                         "idc1",
						"network_download_rate":               "1000000",
						"network_download_rate_limit":         "2000000",
						"network_upload_rate":                 "500000",
						"network_upload_rate_limit":           "1000000",
						"disk_total":                          "100000000000",
						"disk_free":                           "50000000000",
						"disk_used":                           "50000000000",
						"disk_used_percent":                   "50.0",
						"disk_inodes_total":                   "100000",
						"disk_inodes_used":                    "50000",
						"disk_inodes_free":                    "50000",
						"disk_inodes_used_percent":            "50.0",
						"disk_write_bandwidth":                "10000000",
						"disk_read_bandwidth":                 "20000000",
						"build_git_version":                   "v1.0.0",
						"build_git_commit":                    "commit1",
						"build_go_version":                    "1.16",
						"build_platform":                      "linux/amd64",
						"scheduler_cluster_id":                "1",
						"announce_interval":                   "0",
						"updated_at":                          time.Now().Format(time.RFC3339),
						"created_at":                          time.Now().Format(time.RFC3339),
					})
			},
			expectedError: false,
			expectedHosts: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			h := &hostManager{
				config: &config.Config{
					Manager: config.ManagerConfig{SchedulerClusterID: 1},
				},
				rdb: rdb,
			}

			hosts, err := h.LoadAll(context.Background())
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedHosts, len(hosts))

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("redis expectations were not met: %v", err)
			}
		})
	}
}

func TestHostManager_LoadRandom(t *testing.T) {
	tests := []struct {
		name              string
		n                 int
		blocklist         set.SafeSet[string]
		mockRedis         func(mock redismock.ClientMock)
		expectedErr       bool
		expectedHostCount int
	}{
		{
			name:      "smembers fails",
			n:         2,
			blocklist: set.NewSafeSet[string](),
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(pkgredis.MakePersistentCacheHostsInScheduler(1)).
					SetErr(fmt.Errorf("redis error"))
			},
			expectedErr:       true,
			expectedHostCount: 0,
		},
		{
			name: "some hosts in blocklist",
			n:    3,
			blocklist: func() set.SafeSet[string] {
				s := set.NewSafeSet[string]()
				s.Add("host1")
				s.Add("host2")
				s.Add("host3")
				return s
			}(),
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSMembers(pkgredis.MakePersistentCacheHostsInScheduler(1)).
					SetVal([]string{"host1", "host2", "host3"})
			},
			expectedErr:       false,
			expectedHostCount: 0, // host2 is skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			h := &hostManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 1,
					},
				},
				rdb: rdb,
			}

			hosts, err := h.LoadRandom(context.Background(), tt.n, tt.blocklist)
			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, hosts, tt.expectedHostCount)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}

func TestHostManager_RunGC(t *testing.T) {
	tests := []struct {
		name      string
		mockRedis func(mock redismock.ClientMock)
		expectErr bool
	}{
		{
			name: "loadAll fails",
			mockRedis: func(mock redismock.ClientMock) {
				mock.ExpectSScan(pkgredis.MakePersistentCacheHostsInScheduler(1), 0, "*", 10).
					SetErr(fmt.Errorf("redis error"))
			},
			expectErr: true,
		},
		{
			name: "hosts found, none older than 2 intervals => no delete",
			mockRedis: func(mock redismock.ClientMock) {
				// Return a single host
				mock.ExpectSScan(pkgredis.MakePersistentCacheHostsInScheduler(1), 0, "*", 10).
					SetVal([]string{"host1"}, 0)
				mock.ExpectHGetAll(pkgredis.MakePersistentCacheHostKeyInScheduler(1, "host1")).
					SetVal(map[string]string{
						"id":                                  "host1",
						"type":                                "normal",
						"hostname":                            "hostname1",
						"ip":                                  "127.0.0.1",
						"port":                                "8080",
						"download_port":                       "8081",
						"disable_shared":                      "false",
						"os":                                  "linux",
						"platform":                            "x86_64",
						"platform_family":                     "debian",
						"platform_version":                    "11",
						"kernel_version":                      "5.10",
						"cpu_logical_count":                   "4",
						"cpu_physical_count":                  "2",
						"cpu_percent":                         "50.0",
						"cpu_processe_percent":                "25.0",
						"cpu_times_user":                      "10.0",
						"cpu_times_system":                    "5.0",
						"cpu_times_idle":                      "100.0",
						"cpu_times_nice":                      "0.0",
						"cpu_times_iowait":                    "1.0",
						"cpu_times_irq":                       "0.5",
						"cpu_times_softirq":                   "0.2",
						"cpu_times_steal":                     "0.1",
						"cpu_times_guest":                     "0.0",
						"cpu_times_guest_nice":                "0.0",
						"memory_total":                        "8000000000",
						"memory_available":                    "4000000000",
						"memory_used":                         "4000000000",
						"memory_used_percent":                 "50.0",
						"memory_processe_used_percent":        "25.0",
						"memory_free":                         "2000000000",
						"network_tcp_connection_count":        "100",
						"network_upload_tcp_connection_count": "50",
						"network_location":                    "location1",
						"network_idc":                         "idc1",
						"network_download_rate":               "1000000",
						"network_download_rate_limit":         "2000000",
						"network_upload_rate":                 "500000",
						"network_upload_rate_limit":           "1000000",
						"disk_total":                          "100000000000",
						"disk_free":                           "50000000000",
						"disk_used":                           "50000000000",
						"disk_used_percent":                   "50.0",
						"disk_inodes_total":                   "100000",
						"disk_inodes_used":                    "50000",
						"disk_inodes_free":                    "50000",
						"disk_inodes_used_percent":            "50.0",
						"disk_write_bandwidth":                "10000000",
						"disk_read_bandwidth":                 "20000000",
						"build_git_version":                   "v1.0.0",
						"build_git_commit":                    "commit1",
						"build_go_version":                    "1.16",
						"build_platform":                      "linux/amd64",
						"scheduler_cluster_id":                "1",
						"announce_interval":                   "0",
						"created_at":                          time.Now().Format(time.RFC3339),
						"updated_at":                          time.Now().Format(time.RFC3339),
					})
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rdb, mock := redismock.NewClientMock()
			tt.mockRedis(mock)

			h := &hostManager{
				config: &config.Config{
					Manager: config.ManagerConfig{
						SchedulerClusterID: 1,
					},
				},
				rdb: rdb,
			}

			err := h.RunGC(context.Background())
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %v", err)
			}
		})
	}
}
