/*
 *     Copyright 2026 The Dragonfly Authors
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

package jwt

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.Equal(t, ModeDisabled, config.Mode)
	assert.Equal(t, DefaultIssuer, config.JWT.Issuer)
	assert.Equal(t, DefaultTokenTTL, config.JWT.TokenTTL)
	assert.Equal(t, DefaultMaxTokenTTL, config.JWT.MaxTokenTTL)
	assert.Equal(t, DefaultClockSkew, config.JWT.ClockSkew)
	assert.Equal(t, DefaultRefreshBefore, config.JWT.RefreshBefore)
}

func TestConfigEffectiveMode(t *testing.T) {
	assert.Equal(t, ModeDisabled, Config{}.EffectiveMode())
	assert.False(t, Config{}.Enabled())
	assert.True(t, Config{Mode: ModeRequired}.Enabled())
}

func TestValidateConfig(t *testing.T) {
	keyFile := writeTestKey(t, strings.Repeat("k", MinimumKeySize))
	valid := validTestConfig(keyFile)

	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{
			name:    "unsupported mode",
			mutate:  func(config *Config) { config.Mode = "unknown" },
			wantErr: `grpc auth has unsupported mode "unknown"`,
		},
		{
			name:    "missing issuer",
			mutate:  func(config *Config) { config.JWT.Issuer = "" },
			wantErr: "grpc auth jwt requires parameter issuer",
		},
		{
			name:    "token ttl below one second",
			mutate:  func(config *Config) { config.JWT.TokenTTL = 0 },
			wantErr: "grpc auth jwt tokenTTL must be at least one second",
		},
		{
			name:    "max token ttl below one second",
			mutate:  func(config *Config) { config.JWT.MaxTokenTTL = 0 },
			wantErr: "grpc auth jwt maxTokenTTL must be at least one second",
		},
		{
			name:    "duration with fractional second",
			mutate:  func(config *Config) { config.JWT.TokenTTL += time.Millisecond },
			wantErr: "grpc auth jwt durations must use whole seconds",
		},
		{
			name:    "token ttl exceeds maximum",
			mutate:  func(config *Config) { config.JWT.TokenTTL = config.JWT.MaxTokenTTL + time.Second },
			wantErr: "grpc auth jwt tokenTTL must not exceed maxTokenTTL",
		},
		{
			name:    "negative clock skew",
			mutate:  func(config *Config) { config.JWT.ClockSkew = -time.Second },
			wantErr: "grpc auth jwt clockSkew must not be negative",
		},
		{
			name:    "invalid refresh before",
			mutate:  func(config *Config) { config.JWT.RefreshBefore = config.JWT.TokenTTL },
			wantErr: "grpc auth jwt refreshBefore must be positive and less than tokenTTL",
		},
		{
			name:    "missing active key id",
			mutate:  func(config *Config) { config.JWT.ActiveKeyID = "" },
			wantErr: "grpc auth jwt requires parameter activeKeyID",
		},
		{
			name:    "missing keys",
			mutate:  func(config *Config) { config.JWT.Keys = nil },
			wantErr: "grpc auth jwt requires at least one key",
		},
		{
			name: "duplicate key id",
			mutate: func(config *Config) {
				config.JWT.Keys = append(config.JWT.Keys, config.JWT.Keys[0])
			},
			wantErr: `grpc auth jwt key id "test-key" is duplicated`,
		},
		{
			name:    "untrusted active key",
			mutate:  func(config *Config) { config.JWT.ActiveKeyID = "unknown" },
			wantErr: `grpc auth jwt active key id "unknown" is not trusted`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := valid
			config.JWT.Keys = append([]KeyConfig(nil), valid.JWT.Keys...)
			tc.mutate(&config)
			assert.EqualError(t, ValidateConfig(config), tc.wantErr)
		})
	}

	assert.NoError(t, ValidateConfig(valid))
	assert.NoError(t, ValidateConfig(Config{}))
}

func TestValidateConfigRejectsInvalidKeyFiles(t *testing.T) {
	shortKeyFile := writeTestKey(t, "short")
	invalidBase64File := filepath.Join(t.TempDir(), "invalid")
	require.NoError(t, os.WriteFile(invalidBase64File, []byte("not-base64!"), 0o600))

	directory := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{
			name:    "missing file",
			path:    filepath.Join(t.TempDir(), "missing"),
			wantErr: "stat secret file",
		},
		{
			name:    "directory",
			path:    directory,
			wantErr: "secret file is not a regular file",
		},
		{
			name:    "invalid base64",
			path:    invalidBase64File,
			wantErr: "secret file is not valid standard Base64",
		},
		{
			name:    "short key",
			path:    shortKeyFile,
			wantErr: "decoded secret must contain at least 32 bytes",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := validTestConfig(tc.path)
			assert.ErrorContains(t, ValidateConfig(config), tc.wantErr)
		})
	}
}

func validTestConfig(keyFile string) Config {
	config := DefaultConfig()
	config.Mode = ModeRequired
	config.JWT.ActiveKeyID = "test-key"
	config.JWT.Keys = []KeyConfig{{ID: "test-key", SecretFile: keyFile}}
	return config
}

func writeTestKey(t *testing.T, secret string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "key")
	encoded := base64.StdEncoding.EncodeToString([]byte(secret)) + "\n"
	require.NoError(t, os.WriteFile(path, []byte(encoded), 0o600))
	return path
}
