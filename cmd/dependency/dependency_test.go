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
	schedulerconfig "d7y.io/dragonfly/v2/scheduler/config"
)

type tlsConfig struct {
	CACert string `yaml:"caCert" mapstructure:"caCert"`
}

type portRange struct {
	Start int
	End   int
}

type testConfig struct {
	base.Options `yaml:",inline" mapstructure:",squash"`

	Name string     `yaml:"name" mapstructure:"name"`
	TLS  *tlsConfig `yaml:"tls" mapstructure:"tls"`
	Port portRange  `yaml:"port" mapstructure:"port"`
}

func newTestConfig() *testConfig {
	return &testConfig{Name: "default-name"}
}

func setupViper(prefix string) {
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
		expect     func(t *testing.T, cfg *testConfig, err error)
	}{
		{
			name: "no env keeps defaults",
			expect: func(t *testing.T, cfg *testConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("default-name", cfg.Name)
				assert.Nil(cfg.TLS)
			},
		},
		{
			name: "env overrides top-level default",
			envs: map[string]string{"TEST_NAME": "from-env"},
			expect: func(t *testing.T, cfg *testConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("from-env", cfg.Name)
			},
		},
		{
			name: "env materializes section behind nil pointer",
			envs: map[string]string{"TEST_TLS_CACERT": "/etc/ssl/ca.crt"},
			expect: func(t *testing.T, cfg *testConfig, err error) {
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
			expect: func(t *testing.T, cfg *testConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.True(cfg.Console)
			},
		},
		{
			name: "env binds untagged field by field name",
			envs: map[string]string{"TEST_PORT_START": "65003"},
			expect: func(t *testing.T, cfg *testConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(65003, cfg.Port.Start)
			},
		},
		{
			name:       "env overrides config file value and keeps file-only values",
			configFile: "name: from-file\nport:\n  start: 7000\n",
			envs:       map[string]string{"TEST_NAME": "from-env"},
			expect: func(t *testing.T, cfg *testConfig, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("from-env", cfg.Name)
				assert.Equal(7000, cfg.Port.Start)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			setupViper("test")
			for k, v := range tc.envs {
				t.Setenv(k, v)
			}

			cfg := newTestConfig()
			bindEnvsFromConfig(cfg)

			if tc.configFile != "" {
				viper.SetConfigType("yaml")
				assert.NoError(viper.ReadConfig(strings.NewReader(tc.configFile)))
			}

			tc.expect(t, cfg, viper.Unmarshal(cfg, initDecoderConfig))
		})
	}
}

func TestBindEnvsFromConfig_RealSchedulerConfig(t *testing.T) {
	assert := assert.New(t)
	setupViper("scheduler")
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
	setupViper("manager")
	t.Setenv("MANAGER_SERVER_GRPC_PORT_START", "65003")

	cfg := managerconfig.New()
	bindEnvsFromConfig(cfg)
	assert.NoError(viper.Unmarshal(cfg, initDecoderConfig))
	assert.Equal(65003, cfg.Server.GRPC.Port.Start)
}
