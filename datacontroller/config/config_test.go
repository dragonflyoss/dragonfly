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

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	t.Run("defaults are valid", func(t *testing.T) {
		assert.NoError(t, New().Validate())
	})

	t.Run("manager endpoint needs a token", func(t *testing.T) {
		cfg := New()
		cfg.Manager.Endpoint = "http://manager:8080"
		assert.Error(t, cfg.Validate())

		cfg.Manager.Token = "t"
		assert.NoError(t, cfg.Validate())
	})

	t.Run("manager endpoint must be http or https", func(t *testing.T) {
		cfg := New()
		cfg.Manager.Endpoint = "manager:8080"
		cfg.Manager.Token = "t"
		assert.Error(t, cfg.Validate())
	})

	t.Run("controller parameters", func(t *testing.T) {
		cfg := New()
		cfg.Controller.MaxConcurrentReconciles = 0
		assert.Error(t, cfg.Validate())

		cfg = New()
		cfg.Controller.PollInterval = 0
		assert.Error(t, cfg.Validate())

		cfg = New()
		cfg.Controller.MaxRetryBackoff = cfg.Controller.RetryBackoff / 2
		assert.Error(t, cfg.Validate())

		cfg = New()
		cfg.Controller.LeaderElection = true
		cfg.Controller.LeaderElectionID = ""
		assert.Error(t, cfg.Validate())

		cfg = New()
		cfg.Server.LogLevel = "loud"
		assert.Error(t, cfg.Validate())
	})
}

func TestConfig_Convert(t *testing.T) {
	dir := t.TempDir()
	tokenFile := filepath.Join(dir, "token")
	require.NoError(t, os.WriteFile(tokenFile, []byte("  secret\n"), 0o600))

	cfg := New()
	cfg.Manager.TokenFile = tokenFile
	require.NoError(t, cfg.Convert())
	assert.Equal(t, "secret", cfg.Manager.Token)

	cfg = New()
	cfg.Manager.Token = "inline"
	cfg.Manager.TokenFile = tokenFile
	require.NoError(t, cfg.Convert())
	assert.Equal(t, "inline", cfg.Manager.Token, "an explicit token wins over the file")

	cfg = New()
	cfg.Manager.TokenFile = filepath.Join(dir, "missing")
	assert.Error(t, cfg.Convert())

	empty := filepath.Join(dir, "empty")
	require.NoError(t, os.WriteFile(empty, []byte(" \n"), 0o600))
	cfg = New()
	cfg.Manager.TokenFile = empty
	assert.Error(t, cfg.Convert())
}
