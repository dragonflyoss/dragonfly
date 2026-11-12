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
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGC_Add(t *testing.T) {
	ctl := gomock.NewController(t)
	mockLogger := NewMockLogger(ctl)
	mockRunner := NewMockRunner(ctl)

	tests := []struct {
		name   string
		task   Task
		expect func(t *testing.T, err error)
	}{
		{
			name: "new GC",
			task: Task{
				ID:       "gc",
				Interval: 2 * time.Second,
				Timeout:  1 * time.Second,
				Runner:   mockRunner,
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "add GC task without ID",
			task: Task{
				ID:       "",
				Interval: 2 * time.Second,
				Timeout:  1 * time.Second,
				Runner:   mockRunner,
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "add GC task without interval",
			task: Task{
				ID:       "gc",
				Interval: 0,
				Timeout:  1 * time.Second,
				Runner:   mockRunner,
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "add GC task without timeout",
			task: Task{
				ID:       "gc",
				Interval: 2 * time.Second,
				Timeout:  0,
				Runner:   mockRunner,
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "timeout is greater than interval",
			task: Task{
				ID:       "gc",
				Interval: 1 * time.Second,
				Timeout:  2 * time.Second,
				Runner:   mockRunner,
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "add GC task without runner",
			task: Task{
				ID:       "gc",
				Interval: 2 * time.Second,
				Timeout:  1 * time.Second,
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gc := New(WithLogger(mockLogger))
			tc.expect(t, gc.Add(tc.task))
		})
	}
}

func TestGC_Run(t *testing.T) {
	tests := []struct {
		name   string
		task   Task
		id     string
		mock   func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup)
		expect func(t *testing.T, err error)
	}{
		{
			name: "run task",
			task: Task{
				ID:       "foo",
				Interval: 2 * time.Hour,
				Timeout:  1 * time.Hour,
			},
			id: "foo",
			mock: func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup) {
				wg.Add(3)
				gomock.InOrder(
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
					mr.EXPECT().RunGC(context.Background()).Do(func(_ context.Context) { wg.Done() }).Return(nil).Times(1),
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
				)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "task run GC failed",
			task: Task{
				ID:       "foo",
				Interval: 2 * time.Hour,
				Timeout:  1 * time.Hour,
			},
			id: "foo",
			mock: func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup) {
				wg.Add(4)
				err := errors.New("bar")
				gomock.InOrder(
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
					mr.EXPECT().RunGC(context.Background()).Do(func(_ context.Context) { wg.Done() }).Return(err).Times(1),
					ml.EXPECT().Errorf(gomock.Any(), gomock.Eq("foo"), gomock.Eq(err)).Do(func(template any, args ...any) { wg.Done() }).Times(1),
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
				)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "task load wrong key",
			task: Task{
				ID:       "foo",
				Interval: 2 * time.Hour,
				Timeout:  1 * time.Hour,
			},
			id:   "bar",
			mock: func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup) {},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			mockLogger := NewMockLogger(ctl)
			mockRunner := NewMockRunner(ctl)

			gc := New(WithLogger(mockLogger))
			if err := gc.Add(Task{
				ID:       tc.task.ID,
				Interval: tc.task.Interval,
				Timeout:  tc.task.Timeout,
				Runner:   mockRunner,
			}); err != nil {
				t.Fatal(err)
			}

			var wg sync.WaitGroup
			tc.mock(mockLogger, mockRunner, &wg)
			tc.expect(t, gc.Run(context.Background(), tc.id))
			wg.Wait()
		})
	}
}

func TestGC_RunAll(t *testing.T) {
	tests := []struct {
		name string
		task Task
		mock func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup)
	}{
		{
			name: "run task",
			task: Task{
				ID:       "foo",
				Interval: 2 * time.Hour,
				Timeout:  1 * time.Hour,
			},
			mock: func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup) {
				wg.Add(3)
				gomock.InOrder(
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
					mr.EXPECT().RunGC(context.Background()).Do(func(_ context.Context) { wg.Done() }).Return(nil).Times(1),
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
				)
			},
		},
		{
			name: "task run GC failed",
			task: Task{
				ID:       "foo",
				Interval: 2 * time.Hour,
				Timeout:  1 * time.Hour,
			},
			mock: func(ml *MockLogger, mr *MockRunner, wg *sync.WaitGroup) {
				wg.Add(4)
				err := errors.New("baz")
				gomock.InOrder(
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
					mr.EXPECT().RunGC(context.Background()).Do(func(_ context.Context) { wg.Done() }).Return(err).Times(1),
					ml.EXPECT().Errorf(gomock.Any(), gomock.Eq("foo"), gomock.Eq(err)).Do(func(template any, args ...any) { wg.Done() }).Times(1),
					ml.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template any, args ...any) { wg.Done() }).Times(1),
				)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			mockLogger := NewMockLogger(ctl)
			mockRunner := NewMockRunner(ctl)

			gc := New(WithLogger(mockLogger))
			if err := gc.Add(Task{
				ID:       tc.task.ID,
				Interval: tc.task.Interval,
				Timeout:  tc.task.Timeout,
				Runner:   mockRunner,
			}); err != nil {
				t.Fatal(err)
			}

			var wg sync.WaitGroup
			tc.mock(mockLogger, mockRunner, &wg)
			gc.RunAll(context.Background())
			wg.Wait()
		})
	}
}

func TestGC_Start(t *testing.T) {
	ctl := gomock.NewController(t)
	mockLogger := NewMockLogger(ctl)
	mockRunner := NewMockRunner(ctl)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	gc := New(WithLogger(mockLogger))
	if err := gc.Add(Task{
		ID:       "foo",
		Interval: 2 * time.Hour,
		Timeout:  1 * time.Hour,
		Runner:   mockRunner,
	}); err != nil {
		t.Fatal(err)
	}

	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Eq("foo")).Do(func(template string, args ...any) {
		wg.Done()
	}).Times(1)

	gc.Start(context.Background())
	gc.Stop()
}
