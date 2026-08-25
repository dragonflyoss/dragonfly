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

package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

func TestGCLogger_Infof(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	original := logger.CoreLogger
	logger.SetCoreLogger(zap.New(core).Sugar())
	defer logger.SetCoreLogger(original)

	gl := &gcLogger{}
	gl.Infof("run %s gc task success, latency: %s", "task", "1s")

	entries := logs.All()
	assert.Len(t, entries, 1)
	assert.Equal(t, "run task gc task success, latency: 1s", entries[0].Message)
}

func TestGCLogger_Errorf(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	original := logger.CoreLogger
	logger.SetCoreLogger(zap.New(core).Sugar())
	defer logger.SetCoreLogger(original)

	gl := &gcLogger{}
	gl.Errorf("run %s gc task failed: %s", "task", "timeout")

	entries := logs.All()
	assert.Len(t, entries, 1)
	assert.Equal(t, "run task gc task failed: timeout", entries[0].Message)
}
