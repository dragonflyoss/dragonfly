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

package dependency

import (
	"net"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/cmd/dependency/base"
	managerconfig "d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/pkg/dfnet"
	"d7y.io/dragonfly/v2/pkg/types"
	schedulerconfig "d7y.io/dragonfly/v2/scheduler/config"
)

type mockDecoderConfig struct {
	Addr dfnet.NetAddr    `mapstructure:"addr"`
	Cert types.PEMContent `mapstructure:"cert"`
	IP   net.IP           `mapstructure:"ip"`
}

type mockTLSConfig struct {
	CACert string `yaml:"caCert" mapstructure:"caCert"`
}

type mockPortRange struct {
	Start int
	End   int
}

type mockConfig struct {
	base.Options `yaml:",inline" mapstructure:",squash"`

	Name string         `yaml:"name" mapstructure:"name"`
	TLS  *mockTLSConfig `yaml:"tls" mapstructure:"tls"`
	Port mockPortRange  `yaml:"port" mapstructure:"port"`
}

func mockViper(prefix string) {
	viper.Reset()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func TestBindEnvsFromConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		envs       map[string]string
		expect     func(t *testing.T, cfg *mockConfig, err error)
	}{
		{
			name: "no env keeps defaults",
			expect: func(t *testing.T, cfg *mockConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("default-name", cfg.Name)
				assert.Nil(cfg.TLS)
			},
		},
		{
			name: "env overrides top-level default",
			envs: map[string]string{"TEST_NAME": "from-env"},
			expect: func(t *testing.T, cfg *mockConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("from-env", cfg.Name)
			},
		},
		{
			name: "env materializes section behind nil pointer",
			envs: map[string]string{"TEST_TLS_CACERT": "/etc/ssl/ca.crt"},
			expect: func(t *testing.T, cfg *mockConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				if assert.NotNil(cfg.TLS) {
					assert.Equal("/etc/ssl/ca.crt", cfg.TLS.CACert)
				}
			},
		},
		{
			name: "env binds squashed embedded section at top level",
			envs: map[string]string{"TEST_CONSOLE": "true"},
			expect: func(t *testing.T, cfg *mockConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(cfg.Console)
			},
		},
		{
			name: "env binds untagged field by field name",
			envs: map[string]string{"TEST_PORT_START": "65003"},
			expect: func(t *testing.T, cfg *mockConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(65003, cfg.Port.Start)
			},
		},
		{
			name:       "env overrides config file value and keeps file-only values",
			configFile: "name: from-file\nport:\n  start: 7000\n",
			envs:       map[string]string{"TEST_NAME": "from-env"},
			expect: func(t *testing.T, cfg *mockConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("from-env", cfg.Name)
				assert.Equal(7000, cfg.Port.Start)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockViper("test")
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			cfg := &mockConfig{Name: "default-name"}
			bindEnvsFromConfig(cfg)

			if tc.configFile != "" {
				viper.SetConfigType("yaml")
				if err := viper.ReadConfig(strings.NewReader(tc.configFile)); err != nil {
					t.Fatal(err)
				}
			}

			tc.expect(t, cfg, viper.Unmarshal(cfg, initDecoderConfig))
		})
	}
}

func TestBindEnvsFromConfig_RealSchedulerConfig(t *testing.T) {
	assert := assert.New(t)
	mockViper("scheduler")
	t.Setenv("SCHEDULER_SERVER_HOST", "override-host")
	t.Setenv("SCHEDULER_SERVER_ADVERTISEIP", "192.0.2.1")
	t.Setenv("SCHEDULER_SERVER_TLS_CACERT", "/etc/ssl/ca.crt")

	cfg := schedulerconfig.New()
	assert.Nil(cfg.Server.TLS)

	bindEnvsFromConfig(cfg)
	assert.NoError(viper.Unmarshal(cfg, initDecoderConfig))
	assert.Equal("override-host", cfg.Server.Host)
	assert.True(cfg.Server.AdvertiseIP.Equal(net.ParseIP("192.0.2.1")))
	if assert.NotNil(cfg.Server.TLS) {
		assert.Equal("/etc/ssl/ca.crt", cfg.Server.TLS.CACert)
	}
}

func TestBindEnvsFromConfig_RealManagerConfig(t *testing.T) {
	assert := assert.New(t)
	mockViper("manager")
	t.Setenv("MANAGER_SERVER_GRPC_PORT_START", "65003")

	cfg := managerconfig.New()
	bindEnvsFromConfig(cfg)
	assert.NoError(viper.Unmarshal(cfg, initDecoderConfig))
	assert.Equal(65003, cfg.Server.GRPC.Port.Start)
}

func TestInitDecoderConfig(t *testing.T) {
	tests := []struct {
		name       string
		configFile string
		expect     func(t *testing.T, cfg *mockDecoderConfig, err error)
	}{
		{
			name:       "scalar net addr decodes as tcp",
			configFile: "addr: 127.0.0.1:8002\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(dfnet.NetAddr{Type: dfnet.TCP, Addr: "127.0.0.1:8002"}, cfg.Addr)
			},
		},
		{
			name:       "mapping net addr keeps its type",
			configFile: "addr:\n  type: unix\n  addr: /var/run/dfdaemon.sock\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(dfnet.NetAddr{Type: dfnet.UNIX, Addr: "/var/run/dfdaemon.sock"}, cfg.Addr)
			},
		},
		{
			name:       "inline pem content is trimmed and kept",
			configFile: "cert: |\n  -----BEGIN CERTIFICATE-----\n  Zm9v\n  -----END CERTIFICATE-----\n\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.PEMContent("-----BEGIN CERTIFICATE-----\nZm9v\n-----END CERTIFICATE-----"), cfg.Cert)
			},
		},
		{
			name:       "empty pem content stays empty",
			configFile: "cert: \"\"\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.PEMContent(""), cfg.Cert)
			},
		},
		{
			name:       "pem path that does not exist fails",
			configFile: "cert: /nonexistent/ca.crt\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:       "ip string is parsed",
			configFile: "ip: 192.0.2.1\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(net.ParseIP("192.0.2.1").Equal(cfg.IP))
			},
		},
		{
			name:       "invalid ip fails",
			configFile: "ip: not-an-ip\n",
			expect: func(t *testing.T, cfg *mockDecoderConfig, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockViper("test")
			viper.SetConfigType("yaml")
			if err := viper.ReadConfig(strings.NewReader(tc.configFile)); err != nil {
				t.Fatal(err)
			}

			cfg := &mockDecoderConfig{}
			tc.expect(t, cfg, viper.Unmarshal(cfg, initDecoderConfig))
		})
	}
}
