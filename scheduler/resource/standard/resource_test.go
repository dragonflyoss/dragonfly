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

package standard

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/rpc"
	"d7y.io/dragonfly/v2/scheduler/config"
)

func TestResource_New(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
		mock   func(mg *gc.MockGCMockRecorder)
		expect func(t *testing.T, resource Resource, err error)
	}{
		{
			name:   "new resource",
			config: config.New(),
			mock: func(mg *gc.MockGCMockRecorder) {
				gomock.InOrder(
					mg.Add(gomock.Any()).Return(nil).Times(3),
				)
			},
			expect: func(t *testing.T, resource Resource, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("resource", reflect.TypeOf(resource).Elem().Name())
				assert.NotNil(resource.SeedPeer())
				assert.NotNil(resource.HostManager())
				assert.NotNil(resource.PeerManager())
				assert.NotNil(resource.TaskManager())
				assert.NotNil(resource.PeerClientPool())
			},
		},
		{
			name:   "new resource failed because of host manager error",
			config: config.New(),
			mock: func(mg *gc.MockGCMockRecorder) {
				mg.Add(gomock.Any()).Return(errors.New("foo")).Times(1)
			},
			expect: func(t *testing.T, resource Resource, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "new resource failed because of task manager error",
			config: config.New(),
			mock: func(mg *gc.MockGCMockRecorder) {
				gomock.InOrder(
					mg.Add(gomock.Any()).Return(nil).Times(1),
					mg.Add(gomock.Any()).Return(errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, resource Resource, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "new resource failed because of peer manager error",
			config: config.New(),
			mock: func(mg *gc.MockGCMockRecorder) {
				gomock.InOrder(
					mg.Add(gomock.Any()).Return(nil).Times(2),
					mg.Add(gomock.Any()).Return(errors.New("foo")).Times(1),
				)
			},
			expect: func(t *testing.T, resource Resource, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:   "new resource with empty seed peer list",
			config: config.New(),
			mock: func(mg *gc.MockGCMockRecorder) {
				gomock.InOrder(
					mg.Add(gomock.Any()).Return(nil).Times(3),
				)
			},
			expect: func(t *testing.T, resource Resource, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "new resource without seed peer",
			config: &config.Config{
				Scheduler: config.SchedulerConfig{
					GC: config.GCConfig{
						PeerGCInterval: 100,
						PeerTTL:        1000,
						TaskGCInterval: 100,
						HostGCInterval: 100,
					},
				},
			},
			mock: func(mg *gc.MockGCMockRecorder) {
				mg.Add(gomock.Any()).Return(nil).Times(3)
			},
			expect: func(t *testing.T, resource Resource, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("resource", reflect.TypeOf(resource).Elem().Name())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			gc := gc.NewMockGC(ctl)
			tc.mock(gc.EXPECT())

			resource, err := New(tc.config, gc, rpc.NewInsecureCredentials())
			tc.expect(t, resource, err)
		})
	}
}

func TestResource_Serve(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	gc := gc.NewMockGC(ctl)
	gc.EXPECT().Add(gomock.Any()).Return(nil).Times(3)

	resource, err := New(config.New(), gc, rpc.NewInsecureCredentials())
	if err != nil {
		t.Fatal(err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- resource.Serve()
	}()
	resource.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	assert := assert.New(t)
	select {
	case err := <-errCh:
		assert.NoError(err)
	case <-ctx.Done():
		assert.NoError(ctx.Err())
	}
}
