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
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	schedulerconfig "d7y.io/dragonfly/v2/scheduler/config"
)

// tlsConfig mirrors the optional, pointer-typed TLS sections of the real
// configs (e.g. scheduler/config.Config.Server.TLS), which are left nil by the
// default config and therefore never appear in a marshalled value snapshot.
type tlsConfig struct {
	CACert string `yaml:"caCert" mapstructure:"caCert"`
}

// testConfig mirrors the real config shape: populated scalar defaults plus an
// optional nested section behind a nil pointer.
type testConfig struct {
	Name string     `yaml:"name" mapstructure:"name"`
	TLS  *tlsConfig `yaml:"tls" mapstructure:"tls"`
}

// newTestConfig returns the default config: scalar defaults set, optional TLS
// section left nil (as config.New() does for the real configs).
func newTestConfig() *testConfig {
	return &testConfig{Name: "default-name"}
}

// setupViper resets and configures the package-global viper the same way
// InitCommandAndConfig does, with the given env prefix.
func setupViper(prefix string) {
	viper.Reset()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func TestBindEnvsFromConfig_TopLevelOverride(t *testing.T) {
	setupViper("test")
	t.Setenv("TEST_NAME", "from-env")

	cfg := newTestConfig()
	bindEnvsFromConfig(cfg)
	require.NoError(t, viper.Unmarshal(cfg, initDecoderConfig))

	assert.Equal(t, "from-env", cfg.Name)
}

// TestBindEnvsFromConfig_NestedNilPointerOverride is the regression guard: the
// TLS section is nil in the default config, so the previous marshal-of-value
// implementation never bound TEST_TLS_CACERT and this override was silently
// dropped. Enumerating keys from the type closes that gap.
func TestBindEnvsFromConfig_NestedNilPointerOverride(t *testing.T) {
	setupViper("test")
	t.Setenv("TEST_TLS_CACERT", "/etc/ssl/ca.crt")

	cfg := newTestConfig()
	require.Nil(t, cfg.TLS, "TLS should start nil to model the default config")

	bindEnvsFromConfig(cfg)
	require.NoError(t, viper.Unmarshal(cfg, initDecoderConfig))

	require.NotNil(t, cfg.TLS, "nested env override should materialize the TLS section")
	assert.Equal(t, "/etc/ssl/ca.crt", cfg.TLS.CACert)
}

func TestBindEnvsFromConfig_NoEnvKeepsDefaults(t *testing.T) {
	setupViper("test")

	cfg := newTestConfig()
	bindEnvsFromConfig(cfg)
	require.NoError(t, viper.Unmarshal(cfg, initDecoderConfig))

	assert.Equal(t, "default-name", cfg.Name)
	assert.Nil(t, cfg.TLS, "no env set should leave the optional section nil")
}

// TestBindEnvsFromConfig_RealSchedulerConfig exercises the actual production
// config type end-to-end, including a nested key (Server.TLS.CACert) that lives
// under a pointer left nil by config.New() — the regression the fix targets.
func TestBindEnvsFromConfig_RealSchedulerConfig(t *testing.T) {
	setupViper("scheduler")
	t.Setenv("SCHEDULER_SERVER_HOST", "override-host")
	t.Setenv("SCHEDULER_SERVER_TLS_CACERT", "/etc/ssl/ca.crt")

	cfg := schedulerconfig.New()
	require.Nil(t, cfg.Server.TLS, "default scheduler config leaves Server.TLS nil")

	bindEnvsFromConfig(cfg)
	require.NoError(t, viper.Unmarshal(cfg, initDecoderConfig))

	assert.Equal(t, "override-host", cfg.Server.Host, "env should override a populated default")

	require.NotNil(t, cfg.Server.TLS, "nested env override should materialize the TLS section")
	assert.Equal(t, "/etc/ssl/ca.crt", cfg.Server.TLS.CACert)
}
