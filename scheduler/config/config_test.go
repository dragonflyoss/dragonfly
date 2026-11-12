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
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"d7y.io/dragonfly/v2/pkg/net/ip"
)

var (
	mockManagerAddr = "localhost"

	mockManagerConfig = ManagerConfig{
		Addr:               &mockManagerAddr,
		SchedulerClusterID: DefaultManagerSchedulerClusterID,
		KeepAlive: KeepAliveConfig{
			Interval: DefaultManagerKeepAliveInterval,
		},
	}

	mockJobConfig = JobConfig{
		Enable:             true,
		GlobalWorkerNum:    DefaultJobGlobalWorkerNum,
		SchedulerWorkerNum: DefaultJobSchedulerWorkerNum,
		LocalWorkerNum:     DefaultJobLocalWorkerNum,
	}

	mockMetricsConfig = MetricsConfig{
		Enable: true,
		Addr:   DefaultMetricsAddr,
	}

	mockRedisConfig = RedisConfig{
		Addrs:      []string{"127.0.0.0:6379"},
		MasterName: "master",
		Username:   "baz",
		Password:   "bax",
		BrokerDB:   DefaultRedisBrokerDB,
		BackendDB:  DefaultRedisBackendDB,
	}
)

func TestConfig_Load(t *testing.T) {
	mockManagerLoadAddr := "127.0.0.1:65003"
	config := &Config{
		Scheduler: SchedulerConfig{
			Algorithm:              "default",
			BackToSourceCount:      3,
			RetryBackToSourceLimit: 2,
			RetryLimit:             10,
			RetryInterval:          10 * time.Second,
			GC: GCConfig{
				PieceDownloadTimeout: 5 * time.Second,
				PeerGCInterval:       10 * time.Second,
				PeerTTL:              1 * time.Minute,
				TaskGCInterval:       30 * time.Second,
				HostGCInterval:       1 * time.Minute,
				HostTTL:              1 * time.Minute,
			},
		},
		Server: ServerConfig{
			AdvertiseIP:      net.ParseIP("127.0.0.1"),
			AdvertisePort:    8004,
			ListenIP:         net.ParseIP("0.0.0.0"),
			Port:             8002,
			Host:             "foo",
			RequestRateLimit: 100,
			TLS: &GRPCTLSServerConfig{
				CACert: "foo",
				Cert:   "foo",
				Key:    "foo",
			},
			LogDir:        "foo",
			LogLevel:      "debug",
			LogMaxSize:    512,
			LogMaxAge:     5,
			LogMaxBackups: 3,
			PluginDir:     "foo",
		},
		Database: DatabaseConfig{
			Redis: RedisConfig{
				Host:        "127.0.0.1",
				Password:    "foo",
				Addrs:       []string{"foo", "bar"},
				MasterName:  "baz",
				Port:        6379,
				BrokerDB:    DefaultRedisBrokerDB,
				BackendDB:   DefaultRedisBackendDB,
				PoolSize:    10,
				PoolTimeout: 1 * time.Second,
			},
		},
		DynConfig: DynConfig{
			RefreshInterval: 10 * time.Second,
		},
		Manager: ManagerConfig{
			Addr: &mockManagerLoadAddr,
			TLS: &GRPCTLSClientConfig{
				CACert: "foo",
				Cert:   "foo",
				Key:    "foo",
			},
			SchedulerClusterID: 1,
			KeepAlive: KeepAliveConfig{
				Interval: 5 * time.Second,
			},
		},
		SeedPeer: SeedPeerConfig{
			TLS: &GRPCTLSClientConfig{
				CACert: "foo",
				Cert:   "foo",
				Key:    "foo",
			},
			TaskDownloadTimeout: 12 * time.Hour,
		},
		Host: HostConfig{
			IDC:      "foo",
			Location: "baz",
		},
		Job: JobConfig{
			Enable:             true,
			GlobalWorkerNum:    1,
			SchedulerWorkerNum: 1,
			LocalWorkerNum:     5,
		},
		Metrics: MetricsConfig{
			Enable:     false,
			Addr:       ":8000",
			EnableHost: true,
		},
		Network: NetworkConfig{
			EnableIPv6: true,
		},
	}

	schedulerConfigYAML := &Config{}
	contentYAML, _ := os.ReadFile("./testdata/scheduler.yaml")
	if err := yaml.Unmarshal(contentYAML, &schedulerConfigYAML); err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)
	assert.EqualValues(config, schedulerConfigYAML)
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		mock   func(cfg *Config)
		expect func(t *testing.T, err error)
	}{
		{
			name:   "valid config",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:   "valid config without manager",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager.Addr = nil
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:   "server requires parameter advertiseIP",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Job = mockJobConfig
				cfg.Server.AdvertiseIP = nil
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server requires parameter advertisePort",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Job = mockJobConfig
				cfg.Server.AdvertisePort = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server requires parameter listenIP",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Job = mockJobConfig
				cfg.Server.ListenIP = nil
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server requires parameter port",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Job = mockJobConfig
				cfg.Server.Port = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server requires parameter host",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Job = mockJobConfig
				cfg.Server.Host = ""
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server tls requires parameter caCert",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Server.TLS = &GRPCTLSServerConfig{
					CACert: "",
					Cert:   "foo",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server tls requires parameter cert",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Server.TLS = &GRPCTLSServerConfig{
					CACert: "foo",
					Cert:   "",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server tls requires parameter key",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Server.TLS = &GRPCTLSServerConfig{
					CACert: "foo",
					Cert:   "foo",
					Key:    "",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "server requires parameter requestRateLimit",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Job = mockJobConfig
				cfg.Server.RequestRateLimit = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "redis requires parameter brokerDB",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.BrokerDB = -1
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "redis requires parameter backendDB",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.BackendDB = -1
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "redis tls allows insecureSkipVerify only",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.TLS = &RedisTLSClientConfig{
					InsecureSkipVerify: true,
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:   "redis tls allows caCert only",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.TLS = &RedisTLSClientConfig{
					CACert: "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:   "redis tls allows mutual tls",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.TLS = &RedisTLSClientConfig{
					CACert: "foo",
					Cert:   "foo",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name:   "redis tls requires caCert or insecureSkipVerify",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.TLS = &RedisTLSClientConfig{}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "redis tls cert and key must be paired",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.TLS = &RedisTLSClientConfig{
					CACert: "foo",
					Cert:   "foo",
					Key:    "",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "redis tls cert and key must be paired (reverse)",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Database.Redis.TLS = &RedisTLSClientConfig{
					CACert: "foo",
					Cert:   "",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter algorithm",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.Algorithm = ""
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter backToSourceCount",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.BackToSourceCount = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter retryBackToSourceLimit",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.RetryBackToSourceLimit = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter retryLimit",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.RetryLimit = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter retryInterval",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.RetryInterval = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter pieceDownloadTimeout",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.GC.PieceDownloadTimeout = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter peerTTL",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.GC.PeerTTL = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter peerGCInterval",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.GC.PeerGCInterval = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter taskGCInterval",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.GC.TaskGCInterval = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter hostGCInterval",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.GC.HostGCInterval = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "scheduler requires parameter hostTTL",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Scheduler.GC.HostTTL = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "dynconfig requires parameter refreshInterval",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.DynConfig.RefreshInterval = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "manager requires parameter addr",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				emptyAddr := ""
				cfg.Manager.Addr = &emptyAddr
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "manager requires parameter schedulerClusterID",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Manager.SchedulerClusterID = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "manager requires parameter keepAlive interval",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Manager.KeepAlive.Interval = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "manager tls requires parameter caCert",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Manager.TLS = &GRPCTLSClientConfig{
					CACert: "",
					Cert:   "foo",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "manager tls requires parameter cert",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Manager.TLS = &GRPCTLSClientConfig{
					CACert: "foo",
					Cert:   "",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "manager tls requires parameter key",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Manager.TLS = &GRPCTLSClientConfig{
					CACert: "foo",
					Cert:   "foo",
					Key:    "",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "seedPeer requires parameter taskDownloadTimeout",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.SeedPeer.TaskDownloadTimeout = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "seedPeer tls requires parameter caCert",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.SeedPeer.TLS = &GRPCTLSClientConfig{
					CACert: "",
					Cert:   "foo",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "seedPeer tls requires parameter cert",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.SeedPeer.TLS = &GRPCTLSClientConfig{
					CACert: "foo",
					Cert:   "",
					Key:    "foo",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "seedPeer tls requires parameter key",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.SeedPeer.TLS = &GRPCTLSClientConfig{
					CACert: "foo",
					Cert:   "foo",
					Key:    "",
				}
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "job requires parameter globalWorkerNum",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Job.GlobalWorkerNum = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "job requires parameter schedulerWorkerNum",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Job.SchedulerWorkerNum = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "job requires parameter localWorkerNum",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Job.LocalWorkerNum = 0
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "metrics requires parameter addr",
			config: New(),
			mock: func(cfg *Config) {
				cfg.Manager = mockManagerConfig
				cfg.Database.Redis = mockRedisConfig
				cfg.Job = mockJobConfig
				cfg.Metrics = mockMetricsConfig
				cfg.Metrics.Addr = ""
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.config.Convert(); err != nil {
				t.Fatal(err)
			}

			tc.mock(tc.config)
			tc.expect(t, tc.config.Validate())
		})
	}
}

func TestConfig_Convert(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		expect func(t *testing.T, cfg *Config)
	}{
		{
			name: "deprecated job redis addrs fill database redis addrs",
			config: &Config{
				Job: JobConfig{Redis: RedisConfig{Addrs: []string{"foo:6379"}}},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal([]string{"foo:6379"}, cfg.Database.Redis.Addrs)
			},
		},
		{
			name: "deprecated job redis host and port fill database redis addrs",
			config: &Config{
				Job: JobConfig{Redis: RedisConfig{Host: "127.0.0.1", Port: 6379}},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal([]string{"127.0.0.1:6379"}, cfg.Database.Redis.Addrs)
			},
		},
		{
			name: "deprecated job redis host without port is ignored",
			config: &Config{
				Job: JobConfig{Redis: RedisConfig{Host: "127.0.0.1"}},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Empty(cfg.Database.Redis.Addrs)
			},
		},
		{
			name: "database redis addrs take precedence over deprecated job redis",
			config: &Config{
				Database: DatabaseConfig{Redis: RedisConfig{Addrs: []string{"foo:6379"}}},
				Job:      JobConfig{Redis: RedisConfig{Addrs: []string{"bar:6379"}, Host: "127.0.0.1", Port: 6379}},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal([]string{"foo:6379"}, cfg.Database.Redis.Addrs)
			},
		},
		{
			name: "deprecated job redis credentials and databases fill database redis",
			config: &Config{
				Job: JobConfig{Redis: RedisConfig{MasterName: "master", Username: "foo", Password: "bar", BrokerDB: 3, BackendDB: 4}},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal("master", cfg.Database.Redis.MasterName)
				assert.Equal("foo", cfg.Database.Redis.Username)
				assert.Equal("bar", cfg.Database.Redis.Password)
				assert.Equal(3, cfg.Database.Redis.BrokerDB)
				assert.Equal(4, cfg.Database.Redis.BackendDB)
			},
		},
		{
			name: "database redis credentials and databases take precedence over deprecated job redis",
			config: &Config{
				Database: DatabaseConfig{Redis: RedisConfig{MasterName: "master", Username: "foo", Password: "bar", BrokerDB: 1, BackendDB: 2}},
				Job:      JobConfig{Redis: RedisConfig{MasterName: "baz", Username: "bax", Password: "bac", BrokerDB: 3, BackendDB: 4}},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal("master", cfg.Database.Redis.MasterName)
				assert.Equal("foo", cfg.Database.Redis.Username)
				assert.Equal("bar", cfg.Database.Redis.Password)
				assert.Equal(1, cfg.Database.Redis.BrokerDB)
				assert.Equal(2, cfg.Database.Redis.BackendDB)
			},
		},
		{
			name:   "advertiseIP and listenIP default to ipv4",
			config: &Config{},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal(ip.IPv4, cfg.Server.AdvertiseIP)
				assert.Equal(net.IPv4zero, cfg.Server.ListenIP)
			},
		},
		{
			name: "advertiseIP and listenIP default to ipv6 when ipv6 is enabled",
			config: &Config{
				Network: NetworkConfig{EnableIPv6: true},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal(ip.IPv6, cfg.Server.AdvertiseIP)
				assert.Equal(net.IPv6zero, cfg.Server.ListenIP)
			},
		},
		{
			name: "explicit advertiseIP and listenIP are kept",
			config: &Config{
				Server: ServerConfig{AdvertiseIP: net.ParseIP("10.0.0.1"), ListenIP: net.ParseIP("10.0.0.2")},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal(net.ParseIP("10.0.0.1"), cfg.Server.AdvertiseIP)
				assert.Equal(net.ParseIP("10.0.0.2"), cfg.Server.ListenIP)
			},
		},
		{
			name: "advertisePort defaults to port",
			config: &Config{
				Server: ServerConfig{Port: 8002},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal(8002, cfg.Server.AdvertisePort)
			},
		},
		{
			name: "explicit advertisePort is kept",
			config: &Config{
				Server: ServerConfig{Port: 8002, AdvertisePort: 8004},
			},
			expect: func(t *testing.T, cfg *Config) {
				assert := assert.New(t)
				assert.Equal(8004, cfg.Server.AdvertisePort)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.config.Convert(); err != nil {
				t.Fatal(err)
			}

			tc.expect(t, tc.config)
		})
	}
}
