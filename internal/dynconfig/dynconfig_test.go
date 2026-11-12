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

package dynconfig

import (
	"errors"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"d7y.io/dragonfly/v2/internal/dynconfig/mocks"
)

type mockDynconfig struct {
	Scheduler mockSchedulerOption
}

type mockSchedulerOption struct {
	Name string
}

func TestDynconfig_Get(t *testing.T) {
	schedulerName := "scheduler"

	tests := []struct {
		name   string
		expire time.Duration
		sleep  func()
		mock   func(m *mocks.MockClientMockRecorder)
		expect func(t *testing.T, d Dynconfig[mockDynconfig])
	}{
		{
			name:   "get config success",
			expire: 1 * time.Millisecond,
			sleep:  func() {},
			mock: func(m *mocks.MockClientMockRecorder) {
				var d map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, &d); err != nil {
					t.Error(err)
				}

				m.Get().Return(d, nil).AnyTimes()
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)
			},
		},
		{
			name:   "get expired config",
			expire: 1 * time.Millisecond,
			sleep: func() {
				time.Sleep(30 * time.Millisecond)
			},
			mock: func(m *mocks.MockClientMockRecorder) {
				var d map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: "foo",
					},
				}, &d); err != nil {
					t.Error(err)
				}

				m.Get().Return(d, nil).Times(2)
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: "foo",
					},
				}, data)
			},
		},
		{
			name:   "get config failed",
			expire: 1 * time.Millisecond,
			sleep: func() {
				time.Sleep(30 * time.Millisecond)
			},
			mock: func(m *mocks.MockClientMockRecorder) {
				var d map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, &d); err != nil {
					t.Error(err)
				}

				gomock.InOrder(
					m.Get().Return(d, nil).Times(1),
					m.Get().Return(nil, errors.New("manager service error")).Times(1),
				)
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)
			},
		},
		{
			name:   "get config with cache",
			expire: 10 * time.Second,
			sleep:  func() {},
			mock: func(m *mocks.MockClientMockRecorder) {
				var d map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, &d); err != nil {
					t.Error(err)
				}

				m.Get().Return(d, nil).Times(1)
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)

				data, err = d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)
			},
		},
		{
			name:   "config is changed",
			expire: 20 * time.Millisecond,
			sleep: func() {
				time.Sleep(30 * time.Millisecond)
			},
			mock: func(m *mocks.MockClientMockRecorder) {
				var df map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, &df); err != nil {
					t.Error(err)
				}

				var ds map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: "foo",
					},
				}, &ds); err != nil {
					t.Error(err)
				}

				gomock.InOrder(
					m.Get().Return(df, nil).Times(2),
					m.Get().Return(ds, nil).Times(1),
				)
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)

				data, err = d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)

				time.Sleep(30 * time.Millisecond)
				data, err = d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: "foo",
					},
				}, data)
			},
		},
		{
			name:   "config is not changed",
			expire: 1 * time.Millisecond,
			sleep: func() {
				time.Sleep(30 * time.Millisecond)
			},
			mock: func(m *mocks.MockClientMockRecorder) {
				var df map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, &df); err != nil {
					t.Error(err)
				}

				m.Get().Return(df, nil).Times(3)
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)

				time.Sleep(30 * time.Millisecond)
				data, err = d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)
			},
		},
		{
			name:   "config is not changed and cache is not expired",
			expire: 10 * time.Second,
			sleep:  func() {},
			mock: func(m *mocks.MockClientMockRecorder) {
				var df map[string]any
				if err := mapstructure.Decode(mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, &df); err != nil {
					t.Error(err)
				}

				m.Get().Return(df, nil).Times(1)
			},
			expect: func(t *testing.T, d Dynconfig[mockDynconfig]) {
				assert := assert.New(t)
				data, err := d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)

				data, err = d.Get()
				assert.NoError(err)
				assert.EqualValues(&mockDynconfig{
					Scheduler: mockSchedulerOption{
						Name: schedulerName,
					},
				}, data)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()
			mockClient := mocks.NewMockClient(ctl)
			tc.mock(mockClient.EXPECT())

			d, err := New[mockDynconfig](mockClient, tc.expire)
			if err != nil {
				t.Fatal(err)
			}

			tc.sleep()
			tc.expect(t, d)
		})
	}
}
