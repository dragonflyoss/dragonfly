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

package evaluator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/scheduler/resource/persistent"
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
)

func TestEvaluator_New(t *testing.T) {
	pluginDir := "."
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	tests := []struct {
		name      string
		algorithm string
		expect    func(t *testing.T, e any)
	}{
		{
			name:      "new evaluator with default algorithm",
			algorithm: "default",
			expect: func(t *testing.T, e any) {
				assert := assert.New(t)
				assert.Equal("evaluatorDefault", reflect.TypeOf(e).Elem().Name())
			},
		},
		{
			name:      "new evaluator with plugin",
			algorithm: "plugin",
			expect: func(t *testing.T, e any) {
				assert := assert.New(t)
				assert.Equal("evaluatorDefault", reflect.TypeOf(e).Elem().Name())
			},
		},
		{
			name:      "new evaluator with empty string",
			algorithm: "",
			expect: func(t *testing.T, e any) {
				assert := assert.New(t)
				assert.Equal("evaluatorDefault", reflect.TypeOf(e).Elem().Name())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, New(tc.algorithm, pluginDir))
		})
	}
}

func TestEvaluator_IsBadPersistentParent(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		expect func(t *testing.T, isBadParent bool)
	}{
		{
			name:  "peer state is PeerStatePending",
			state: persistent.PeerStatePending,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateUploading",
			state: persistent.PeerStateUploading,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateReceivedEmpty",
			state: persistent.PeerStateReceivedEmpty,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateReceivedNormal",
			state: persistent.PeerStateReceivedNormal,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateFailed",
			state: persistent.PeerStateFailed,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateRunning",
			state: persistent.PeerStateRunning,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateSucceeded",
			state: persistent.PeerStateSucceeded,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorDefault()
			tc.expect(t, e.IsBadPersistentParent(mockPersistentPeer("parent", tc.state, mockHostIDC, mockHostLocation)))
		})
	}
}

func TestEvaluator_IsBadPersistentCacheParent(t *testing.T) {
	tests := []struct {
		name   string
		state  string
		expect func(t *testing.T, isBadParent bool)
	}{
		{
			name:  "peer state is PeerStatePending",
			state: persistentcache.PeerStatePending,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateUploading",
			state: persistentcache.PeerStateUploading,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateReceivedEmpty",
			state: persistentcache.PeerStateReceivedEmpty,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateReceivedNormal",
			state: persistentcache.PeerStateReceivedNormal,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateFailed",
			state: persistentcache.PeerStateFailed,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.True(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateRunning",
			state: persistentcache.PeerStateRunning,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
		{
			name:  "peer state is PeerStateSucceeded",
			state: persistentcache.PeerStateSucceeded,
			expect: func(t *testing.T, isBadParent bool) {
				assert := assert.New(t)
				assert.False(isBadParent)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := newEvaluatorDefault()
			tc.expect(t, e.IsBadPersistentCacheParent(mockPersistentCachePeer("parent", tc.state, mockHostIDC, mockHostLocation)))
		})
	}
}
