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

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dragonflyoss/machinery/v1"
	"github.com/dragonflyoss/machinery/v1/backends/eager"
	backendsiface "github.com/dragonflyoss/machinery/v1/backends/iface"
	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"

	internaljob "d7y.io/dragonfly/v2/internal/job"
	managerjob "d7y.io/dragonfly/v2/manager/job"
	"d7y.io/dragonfly/v2/manager/job/mocks"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/net/http"
)

func newTestInternalJob(backend backendsiface.Backend) *internaljob.Job {
	server := &machinery.Server{}
	server.SetBackend(backend)

	return &internaljob.Job{Server: server}
}

func mockJobSchedulers(t *testing.T, db *gorm.DB) models.User {
	full := mockSchedulerCluster(t, db, "full")
	inactive := mockSchedulerCluster(t, db, "inactive")
	scheduleOnly := mockSchedulerCluster(t, db, "schedule-only")
	mockSchedulerCluster(t, db, "empty")
	mockScheduler(t, db, full.ID, "full-preheat", models.SchedulerStateActive, types.DefaultSchedulerFeatures)
	mockScheduler(t, db, full.ID, "full-schedule", models.SchedulerStateActive, []string{types.SchedulerFeatureSchedule})
	mockScheduler(t, db, inactive.ID, "inactive-preheat", models.SchedulerStateInactive, types.DefaultSchedulerFeatures)
	mockScheduler(t, db, scheduleOnly.ID, "schedule-only", models.SchedulerStateActive, []string{types.SchedulerFeatureSchedule})

	return mockUser(t, db, "foo")
}

func TestService_CreatePreheatJob(t *testing.T) {
	s := newTestService(t)
	user := mockJobSchedulers(t, s.db)
	ctrl := gomock.NewController(t)
	mockPreheat := mocks.NewMockPreheat(ctrl)
	s.job = &managerjob.Job{Job: newTestInternalJob(eager.New()), Preheat: mockPreheat}

	tests := []struct {
		name   string
		req    types.CreatePreheatJobRequest
		mock   func(m *mocks.MockPreheatMockRecorder)
		expect func(t *testing.T, s *service, job *models.Job, err error)
	}{
		{
			name: "scheduler cluster not found",
			req:  types.CreatePreheatJobRequest{Type: internaljob.PreheatJob, Args: types.PreheatArgs{Type: "file", URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{99}},
			mock: func(m *mocks.MockPreheatMockRecorder) {},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(job)
			},
		},
		{
			name: "cluster without preheat capable scheduler",
			req:  types.CreatePreheatJobRequest{Type: internaljob.PreheatJob, Args: types.PreheatArgs{Type: "file", URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{3}},
			mock: func(m *mocks.MockPreheatMockRecorder) {},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(job)
			},
		},
		{
			name: "job creation fails",
			req:  types.CreatePreheatJobRequest{Type: internaljob.PreheatJob, Args: types.PreheatArgs{Type: "file", URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{1}},
			mock: func(m *mocks.MockPreheatMockRecorder) {
				m.CreatePreheat(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("create preheat failed"))
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(job)

				var count int64
				assert.NoError(s.db.Model(&models.Job{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
		{
			name: "applies defaults and picks the preheat capable scheduler",
			req: types.CreatePreheatJobRequest{
				BIO:                 "bio",
				Type:                internaljob.PreheatJob,
				Args:                types.PreheatArgs{Type: "file", URL: "https://example.com/file"},
				UserID:              user.ID,
				SchedulerClusterIDs: []uint{1},
			},
			mock: func(m *mocks.MockPreheatMockRecorder) {
				m.CreatePreheat(gomock.Any(), gomock.Cond(func(schedulers []models.Scheduler) bool {
					return len(schedulers) == 1 && schedulers[0].Hostname == "full-preheat" && schedulers[0].SchedulerCluster.ID == 1
				}), gomock.Any()).Return(&internaljob.GroupJobState{GroupUUID: "group-1", State: machineryv1tasks.StatePending}, nil)
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("group-1", job.TaskID)
				assert.Equal(machineryv1tasks.StatePending, job.State)
				assert.Equal(internaljob.PreheatJob, job.Type)
				assert.Equal("bio", job.BIO)
				assert.Equal(user.ID, job.UserID)
				assert.Len(job.SchedulerClusters, 1)
				assert.Equal("full", job.SchedulerClusters[0].Name)
				assert.Equal("https://example.com/file", job.Args["url"])
				assert.Equal(types.SingleSeedPeerScope, job.Args["scope"])
				assert.Equal(float64(types.DefaultPreheatConcurrentTaskCount), job.Args["concurrent_task_count"])
				assert.Equal(float64(types.DefaultPreheatConcurrentPeerCount), job.Args["concurrent_peer_count"])
				assert.Equal(float64(types.DefaultJobTimeout), job.Args["timeout"])
				assert.Equal(http.RawDefaultFilteredQueryParams, job.Args["filtered_query_params"])

				stored := models.Job{}
				assert.NoError(s.db.Preload("SchedulerClusters").First(&stored, job.ID).Error)
				assert.Equal("group-1", stored.TaskID)
				assert.Len(stored.SchedulerClusters, 1)
			},
		},
		{
			name: "keeps explicit args",
			req: types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type:                "image",
					URL:                 "https://example.com/image",
					Scope:               types.AllPeersScope,
					ConcurrentTaskCount: 2,
					ConcurrentPeerCount: 10,
					Timeout:             5 * time.Minute,
					FilteredQueryParams: "foo&bar",
				},
				SchedulerClusterIDs: []uint{1},
			},
			mock: func(m *mocks.MockPreheatMockRecorder) {
				m.CreatePreheat(gomock.Any(), gomock.Len(1), gomock.Any()).Return(&internaljob.GroupJobState{GroupUUID: "group-2", State: machineryv1tasks.StatePending}, nil)
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(types.AllPeersScope, job.Args["scope"])
				assert.Equal(float64(2), job.Args["concurrent_task_count"])
				assert.Equal(float64(10), job.Args["concurrent_peer_count"])
				assert.Equal(float64(5*time.Minute), job.Args["timeout"])
				assert.Equal("foo&bar", job.Args["filtered_query_params"])
			},
		},
		{
			name: "without cluster ids only clusters with a capable active scheduler are used",
			req:  types.CreatePreheatJobRequest{Type: internaljob.PreheatJob, Args: types.PreheatArgs{Type: "file", URL: "https://example.com/file"}},
			mock: func(m *mocks.MockPreheatMockRecorder) {
				m.CreatePreheat(gomock.Any(), gomock.Cond(func(schedulers []models.Scheduler) bool {
					return len(schedulers) == 1 && schedulers[0].Hostname == "full-preheat"
				}), gomock.Any()).Return(&internaljob.GroupJobState{GroupUUID: "group-3", State: machineryv1tasks.StatePending}, nil)
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(job.SchedulerClusters, 1)
				assert.Equal("full", job.SchedulerClusters[0].Name)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(mockPreheat.EXPECT())

			job, err := s.CreatePreheatJob(context.Background(), tc.req)
			tc.expect(t, s, job, err)
		})
	}
}

func TestService_CreateGetTaskJob(t *testing.T) {
	s := newTestService(t)
	mockJobSchedulers(t, s.db)
	ctrl := gomock.NewController(t)
	mockTask := mocks.NewMockTask(ctrl)
	s.job = &managerjob.Job{Job: newTestInternalJob(eager.New()), Task: mockTask}

	tests := []struct {
		name   string
		req    types.CreateGetTaskJobRequest
		mock   func(m *mocks.MockTaskMockRecorder)
		expect func(t *testing.T, s *service, job *models.Job, err error)
	}{
		{
			name: "scheduler cluster not found",
			req:  types.CreateGetTaskJobRequest{Type: internaljob.GetTaskJob, Args: types.GetTaskArgs{URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{99}},
			mock: func(m *mocks.MockTaskMockRecorder) {},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(job)
			},
		},
		{
			name: "cluster with only inactive schedulers",
			req:  types.CreateGetTaskJobRequest{Type: internaljob.GetTaskJob, Args: types.GetTaskArgs{URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{2}},
			mock: func(m *mocks.MockTaskMockRecorder) {},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(job)
			},
		},
		{
			name: "job creation fails",
			req:  types.CreateGetTaskJobRequest{Type: internaljob.GetTaskJob, Args: types.GetTaskArgs{URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{1}},
			mock: func(m *mocks.MockTaskMockRecorder) {
				m.CreateGetTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("create get task failed"))
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(job)
			},
		},
		{
			name: "applies defaults and uses every active scheduler of the cluster",
			req:  types.CreateGetTaskJobRequest{BIO: "bio", Type: internaljob.GetTaskJob, Args: types.GetTaskArgs{URL: "https://example.com/file", Tag: "tag"}, SchedulerClusterIDs: []uint{1}},
			mock: func(m *mocks.MockTaskMockRecorder) {
				m.CreateGetTask(gomock.Any(), gomock.Cond(func(schedulers []models.Scheduler) bool {
					return len(schedulers) == 2 && schedulers[0].SchedulerCluster.ID == 1 && schedulers[1].SchedulerCluster.ID == 1
				}), gomock.Any()).Return(&internaljob.GroupJobState{GroupUUID: "group-1", State: machineryv1tasks.StatePending}, nil)
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("group-1", job.TaskID)
				assert.Equal(internaljob.GetTaskJob, job.Type)
				assert.Equal("tag", job.Args["tag"])
				assert.Equal(float64(types.DefaultGetTaskConcurrentPeerCount), job.Args["concurrent_peer_count"])
				assert.Equal(float64(types.DefaultJobTimeout), job.Args["timeout"])
				assert.Equal(http.RawDefaultFilteredQueryParams, job.Args["filtered_query_params"])

				stored := models.Job{}
				assert.NoError(s.db.Preload("SchedulerClusters").First(&stored, job.ID).Error)
				assert.Len(stored.SchedulerClusters, 1)
				assert.Equal("full", stored.SchedulerClusters[0].Name)
			},
		},
		{
			name: "without cluster ids uses active schedulers of all clusters",
			req:  types.CreateGetTaskJobRequest{Type: internaljob.GetTaskJob, Args: types.GetTaskArgs{URL: "https://example.com/file"}},
			mock: func(m *mocks.MockTaskMockRecorder) {
				m.CreateGetTask(gomock.Any(), gomock.Len(3), gomock.Any()).Return(&internaljob.GroupJobState{GroupUUID: "group-2", State: machineryv1tasks.StatePending}, nil)
			},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				stored := models.Job{}
				assert.NoError(s.db.Preload("SchedulerClusters").First(&stored, job.ID).Error)
				assert.Len(stored.SchedulerClusters, 2)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(mockTask.EXPECT())

			job, err := s.CreateGetTaskJob(context.Background(), tc.req)
			tc.expect(t, s, job, err)
		})
	}
}

func TestService_CreateDeleteTaskJob(t *testing.T) {
	s := newTestService(t)
	mockJobSchedulers(t, s.db)
	ctrl := gomock.NewController(t)
	mockTask := mocks.NewMockTask(ctrl)
	s.job = &managerjob.Job{Job: newTestInternalJob(eager.New()), Task: mockTask}

	tests := []struct {
		name   string
		req    types.CreateDeleteTaskJobRequest
		mock   func(m *mocks.MockTaskMockRecorder)
		expect func(t *testing.T, job *models.Job, err error)
	}{
		{
			name: "scheduler cluster not found",
			req:  types.CreateDeleteTaskJobRequest{Type: internaljob.DeleteTaskJob, Args: types.DeleteTaskArgs{URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{99}},
			mock: func(m *mocks.MockTaskMockRecorder) {},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(job)
			},
		},
		{
			name: "job creation fails",
			req:  types.CreateDeleteTaskJobRequest{Type: internaljob.DeleteTaskJob, Args: types.DeleteTaskArgs{URL: "https://example.com/file"}, SchedulerClusterIDs: []uint{1}},
			mock: func(m *mocks.MockTaskMockRecorder) {
				m.CreateDeleteTask(gomock.Any(), gomock.Len(2), gomock.Any()).Return(nil, errors.New("create delete task failed"))
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(job)
			},
		},
		{
			name: "applies defaults",
			req:  types.CreateDeleteTaskJobRequest{Type: internaljob.DeleteTaskJob, Args: types.DeleteTaskArgs{TaskID: "task-1"}, SchedulerClusterIDs: []uint{1, 3}},
			mock: func(m *mocks.MockTaskMockRecorder) {
				m.CreateDeleteTask(gomock.Any(), gomock.Len(3), gomock.Any()).Return(&internaljob.GroupJobState{GroupUUID: "group-1", State: machineryv1tasks.StatePending}, nil)
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("group-1", job.TaskID)
				assert.Equal(internaljob.DeleteTaskJob, job.Type)
				assert.Equal("task-1", job.Args["task_id"])
				assert.Equal(float64(types.DefaultJobTimeout), job.Args["timeout"])
				assert.Equal(http.RawDefaultFilteredQueryParams, job.Args["filtered_query_params"])
				assert.Len(job.SchedulerClusters, 3)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(mockTask.EXPECT())

			job, err := s.CreateDeleteTaskJob(context.Background(), tc.req)
			tc.expect(t, job, err)
		})
	}
}

func TestService_CreateSyncPeersJob(t *testing.T) {
	s := newTestService(t)
	mockJobSchedulers(t, s.db)
	ctrl := gomock.NewController(t)
	mockSyncPeers := mocks.NewMockSyncPeers(ctrl)
	s.job = &managerjob.Job{Job: newTestInternalJob(eager.New()), SyncPeers: mockSyncPeers}

	tests := []struct {
		name   string
		req    types.CreateSyncPeersJobRequest
		mock   func(m *mocks.MockSyncPeersMockRecorder)
		expect func(t *testing.T, err error)
	}{
		{
			name: "scheduler cluster not found",
			req:  types.CreateSyncPeersJobRequest{Type: internaljob.SyncPeersJob, SchedulerClusterIDs: []uint{99}},
			mock: func(m *mocks.MockSyncPeersMockRecorder) {},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "cluster without active scheduler",
			req:  types.CreateSyncPeersJobRequest{Type: internaljob.SyncPeersJob, SchedulerClusterIDs: []uint{2}},
			mock: func(m *mocks.MockSyncPeersMockRecorder) {},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "one active scheduler per requested cluster",
			req:  types.CreateSyncPeersJobRequest{Type: internaljob.SyncPeersJob, SchedulerClusterIDs: []uint{1, 3}},
			mock: func(m *mocks.MockSyncPeersMockRecorder) {
				m.CreateSyncPeers(gomock.Any(), gomock.Cond(func(schedulers []models.Scheduler) bool {
					return len(schedulers) == 2 && schedulers[0].SchedulerCluster.ID == 1 && schedulers[1].SchedulerCluster.ID == 3
				})).Return(nil)
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.NoError(err)
			},
		},
		{
			name: "without cluster ids skips clusters without active scheduler",
			req:  types.CreateSyncPeersJobRequest{Type: internaljob.SyncPeersJob},
			mock: func(m *mocks.MockSyncPeersMockRecorder) {
				m.CreateSyncPeers(gomock.Any(), gomock.Len(2)).Return(errors.New("sync peers failed"))
			},
			expect: func(t *testing.T, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mock(mockSyncPeers.EXPECT())

			tc.expect(t, s.CreateSyncPeersJob(context.Background(), tc.req))
		})
	}
}

func TestService_CreateGCJob(t *testing.T) {
	tests := []struct {
		name   string
		ctx    func() context.Context
		req    types.CreateGCJobRequest
		mock   func(m *pkggc.MockGCMockRecorder)
		expect func(t *testing.T, job *models.Job, err error)
	}{
		{
			name: "gc run fails",
			ctx: func() context.Context {
				return context.Background()
			},
			req: types.CreateGCJobRequest{Type: internaljob.GCJob, Args: types.GCArgs{Type: models.ConfigGC}, UserID: 7},
			mock: func(m *pkggc.MockGCMockRecorder) {
				m.Run(gomock.Any(), models.ConfigGC).Return(errors.New("gc failed"))
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(job)
			},
		},
		{
			name: "polling stops when the context is cancelled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			req: types.CreateGCJobRequest{Type: internaljob.GCJob, Args: types.GCArgs{Type: "audit"}, UserID: 7},
			mock: func(m *pkggc.MockGCMockRecorder) {
				m.Run(gomock.Cond(func(ctx context.Context) bool {
					taskID, ok := ctx.Value(pkggc.ContextKeyTaskID).(string)
					return ok && taskID != "" && ctx.Value(pkggc.ContextKeyUserID) == uint(7)
				}), "audit").Return(nil)
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, context.Canceled)
				assert.Nil(job)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockGC := pkggc.NewMockGC(ctrl)
			tc.mock(mockGC.EXPECT())

			s := newTestService(t)
			s.gc = mockGC

			job, err := s.CreateGCJob(tc.ctx(), tc.req)
			tc.expect(t, job, err)
		})
	}
}

func TestService_findSchedulerInClusters(t *testing.T) {
	s := newTestService(t)
	mockJobSchedulers(t, s.db)

	tests := []struct {
		name                string
		schedulerClusterIDs []uint
		expect              func(t *testing.T, schedulers []models.Scheduler, err error)
	}{
		{
			name:                "scheduler cluster not found",
			schedulerClusterIDs: []uint{99},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Error(err)
				assert.Nil(schedulers)
			},
		},
		{
			name:                "cluster without active scheduler",
			schedulerClusterIDs: []uint{2},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(schedulers)
			},
		},
		{
			name:                "one active scheduler per requested cluster",
			schedulerClusterIDs: []uint{1, 3},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 2)
				assert.Equal("full", schedulers[0].SchedulerCluster.Name)
				assert.Equal("schedule-only", schedulers[1].SchedulerCluster.Name)
			},
		},
		{
			name:                "all clusters with an active scheduler",
			schedulerClusterIDs: nil,
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 2)
				assert.Equal(uint(1), schedulers[0].SchedulerClusterID)
				assert.Equal(uint(3), schedulers[1].SchedulerClusterID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schedulers, err := s.findSchedulerInClusters(context.Background(), tc.schedulerClusterIDs)
			tc.expect(t, schedulers, err)
		})
	}
}

func TestService_findAllSchedulersInClusters(t *testing.T) {
	s := newTestService(t)
	mockJobSchedulers(t, s.db)

	tests := []struct {
		name                string
		schedulerClusterIDs []uint
		expect              func(t *testing.T, schedulers []models.Scheduler, err error)
	}{
		{
			name:                "scheduler cluster not found",
			schedulerClusterIDs: []uint{1, 99},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(schedulers)
			},
		},
		{
			name:                "cluster with only inactive schedulers",
			schedulerClusterIDs: []uint{2},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(schedulers)
			},
		},
		{
			name:                "every active scheduler of the requested cluster",
			schedulerClusterIDs: []uint{1},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 2)
				assert.Equal("full", schedulers[0].SchedulerCluster.Name)
				assert.Equal("full", schedulers[1].SchedulerCluster.Name)
			},
		},
		{
			name:                "every active scheduler of all clusters",
			schedulerClusterIDs: nil,
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 3)
				for _, scheduler := range schedulers {
					assert.Equal(models.SchedulerStateActive, scheduler.State)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schedulers, err := s.findAllSchedulersInClusters(context.Background(), tc.schedulerClusterIDs)
			tc.expect(t, schedulers, err)
		})
	}
}

func TestService_findAllCandidateSchedulersInClusters(t *testing.T) {
	s := newTestService(t)
	mockJobSchedulers(t, s.db)

	tests := []struct {
		name                string
		schedulerClusterIDs []uint
		features            []string
		expect              func(t *testing.T, schedulers []models.Scheduler, err error)
	}{
		{
			name:                "scheduler cluster not found",
			schedulerClusterIDs: []uint{99},
			features:            []string{types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(schedulers)
			},
		},
		{
			name:                "first scheduler supporting every feature",
			schedulerClusterIDs: []uint{1},
			features:            []string{types.SchedulerFeatureSchedule, types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 1)
				assert.Equal("full-preheat", schedulers[0].Hostname)
				assert.Equal("full", schedulers[0].SchedulerCluster.Name)
			},
		},
		{
			name:                "cluster whose schedulers lack the feature",
			schedulerClusterIDs: []uint{3},
			features:            []string{types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.Error(err)
				assert.Nil(schedulers)
			},
		},
		{
			name:                "inactive schedulers are ignored",
			schedulerClusterIDs: []uint{2},
			features:            []string{types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name:                "no features picks the first active scheduler",
			schedulerClusterIDs: []uint{3},
			features:            nil,
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 1)
				assert.Equal("schedule-only", schedulers[0].Hostname)
			},
		},
		{
			name:                "all clusters with feature",
			schedulerClusterIDs: nil,
			features:            []string{types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 1)
				assert.Equal("full-preheat", schedulers[0].Hostname)
			},
		},
		{
			name:                "all clusters without features picks one scheduler per active cluster",
			schedulerClusterIDs: nil,
			features:            nil,
			expect: func(t *testing.T, schedulers []models.Scheduler, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(schedulers, 2)
				assert.Equal("full-preheat", schedulers[0].Hostname)
				assert.Equal("schedule-only", schedulers[1].Hostname)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schedulers, err := s.findAllCandidateSchedulersInClusters(context.Background(), tc.schedulerClusterIDs, tc.features)
			tc.expect(t, schedulers, err)
		})
	}
}

func TestContainsAllFeatures(t *testing.T) {
	tests := []struct {
		name              string
		schedulerFeatures []string
		features          []string
		expect            func(t *testing.T, ok bool)
	}{
		{
			name:              "no required features",
			schedulerFeatures: nil,
			features:          nil,
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:              "scheduler supports a superset",
			schedulerFeatures: []string{types.SchedulerFeatureSchedule, types.SchedulerFeaturePreheat},
			features:          []string{types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:              "scheduler misses one feature",
			schedulerFeatures: []string{types.SchedulerFeatureSchedule},
			features:          []string{types.SchedulerFeatureSchedule, types.SchedulerFeaturePreheat},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name:              "scheduler without features",
			schedulerFeatures: nil,
			features:          []string{types.SchedulerFeatureSchedule},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, containsAllFeatures(tc.schedulerFeatures, tc.features))
		})
	}
}

func TestService_pollingJob(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, backend backendsiface.Backend, signature *machineryv1tasks.Signature)
		expect func(t *testing.T, job models.Job)
	}{
		{
			name: "group succeeded",
			setup: func(t *testing.T, backend backendsiface.Backend, signature *machineryv1tasks.Signature) {
				assert := assert.New(t)
				assert.NoError(backend.InitGroup(signature.GroupUUID, []string{signature.UUID}))
				assert.NoError(backend.SetStateSuccess(signature, nil))
			},
			expect: func(t *testing.T, job models.Job) {
				assert := assert.New(t)
				assert.Equal(machineryv1tasks.StateSuccess, job.State)
				assert.Equal(machineryv1tasks.StateSuccess, job.Result["state"])
				assert.Equal("group-1", job.Result["group_uuid"])
				assert.Len(job.Result["job_states"], 1)
			},
		},
		{
			name: "group failed",
			setup: func(t *testing.T, backend backendsiface.Backend, signature *machineryv1tasks.Signature) {
				assert := assert.New(t)
				assert.NoError(backend.InitGroup(signature.GroupUUID, []string{signature.UUID}))
				assert.NoError(backend.SetStateFailure(signature, "boom"))
			},
			expect: func(t *testing.T, job models.Job) {
				assert := assert.New(t)
				assert.Equal(machineryv1tasks.StateFailure, job.State)
				assert.Equal(machineryv1tasks.StateFailure, job.Result["state"])
			},
		},
		{
			name: "group still running after the last attempt is marked failed",
			setup: func(t *testing.T, backend backendsiface.Backend, signature *machineryv1tasks.Signature) {
				assert := assert.New(t)
				assert.NoError(backend.InitGroup(signature.GroupUUID, []string{signature.UUID}))
				assert.NoError(backend.SetStateStarted(signature))
			},
			expect: func(t *testing.T, job models.Job) {
				assert := assert.New(t)
				assert.Equal(machineryv1tasks.StateFailure, job.State)
				assert.Equal(machineryv1tasks.StatePending, job.Result["state"])
			},
		},
		{
			name:  "group missing from backend is marked failed without result",
			setup: func(t *testing.T, backend backendsiface.Backend, signature *machineryv1tasks.Signature) {},
			expect: func(t *testing.T, job models.Job) {
				assert := assert.New(t)
				assert.Equal(machineryv1tasks.StateFailure, job.State)
				assert.Empty(job.Result)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			s := newTestService(t)
			backend := eager.New()
			s.job = &managerjob.Job{Job: newTestInternalJob(backend)}
			signature := &machineryv1tasks.Signature{UUID: "task-1", Name: internaljob.PreheatJob, GroupUUID: "group-1"}
			tc.setup(t, backend, signature)

			job := models.Job{TaskID: signature.GroupUUID, Type: internaljob.PreheatJob, State: machineryv1tasks.StatePending, Args: models.JSONMap{}}
			assert.NoError(s.db.Create(&job).Error)

			s.pollingJob(context.Background(), internaljob.PreheatJob, job.ID, signature.GroupUUID, time.Millisecond, time.Millisecond, 2)

			stored := models.Job{}
			assert.NoError(s.db.First(&stored, job.ID).Error)
			tc.expect(t, stored)
		})
	}
}

func mockJobResult(peers ...map[string]any) models.JSONMap {
	anyPeers := make([]any, 0, len(peers))
	for _, peer := range peers {
		anyPeers = append(anyPeers, peer)
	}

	return models.JSONMap{
		"job_states": []any{
			map[string]any{
				"results": []any{
					map[string]any{
						"scheduler_cluster_id": float64(1),
						"peers":                anyPeers,
					},
				},
			},
		},
	}
}

func TestService_extractPeersFromJobs(t *testing.T) {
	tests := []struct {
		name   string
		jobs   []*models.Job
		expect func(t *testing.T, peers []types.Peer, err error)
	}{
		{
			name: "collects normal peers and skips seed peers",
			jobs: []*models.Job{
				{
					State: machineryv1tasks.StateSuccess,
					Args:  models.JSONMap{"url": "https://example.com/layer-1"},
					Result: mockJobResult(map[string]any{"id": "peer-1", "host_type": "normal", "ip": "10.0.0.1", "hostname": "host-1"},
						map[string]any{"id": "peer-2", "host_type": "super", "ip": "10.0.0.2", "hostname": "host-2"},
					),
				},
			},
			expect: func(t *testing.T, peers []types.Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal([]types.Peer{{
					IP:                 "10.0.0.1",
					Hostname:           "host-1",
					CachedLayers:       []types.Layer{{URL: "https://example.com/layer-1"}},
					SchedulerClusterID: 1,
				}}, peers)
			},
		},
		{
			name: "merges layers of the same host across jobs and deduplicates urls",
			jobs: []*models.Job{
				{
					State:  machineryv1tasks.StateSuccess,
					Args:   models.JSONMap{"url": "https://example.com/layer-1"},
					Result: mockJobResult(map[string]any{"id": "peer-1", "host_type": "normal", "ip": "10.0.0.1", "hostname": "host-1"}),
				},
				{
					State: machineryv1tasks.StateSuccess,
					Args:  models.JSONMap{"url": "https://example.com/layer-2"},
					Result: mockJobResult(map[string]any{"id": "peer-1", "host_type": "normal", "ip": "10.0.0.1", "hostname": "host-1"},
						map[string]any{"id": "peer-1-again", "host_type": "normal", "ip": "10.0.0.1", "hostname": "host-1"},
						map[string]any{"id": "peer-3", "host_type": "normal", "ip": "10.0.0.3", "hostname": "host-3"},
					),
				},
			},
			expect: func(t *testing.T, peers []types.Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.ElementsMatch([]types.Peer{
					{
						IP:                 "10.0.0.1",
						Hostname:           "host-1",
						CachedLayers:       []types.Layer{{URL: "https://example.com/layer-1"}, {URL: "https://example.com/layer-2"}},
						SchedulerClusterID: 1,
					},
					{
						IP:                 "10.0.0.3",
						Hostname:           "host-3",
						CachedLayers:       []types.Layer{{URL: "https://example.com/layer-2"}},
						SchedulerClusterID: 1,
					},
				}, peers)
			},
		},
		{
			name: "skips unsuccessful jobs and jobs with malformed results",
			jobs: []*models.Job{
				{
					State:  machineryv1tasks.StateFailure,
					Args:   models.JSONMap{"url": "https://example.com/layer-1"},
					Result: mockJobResult(map[string]any{"id": "peer-1", "host_type": "normal", "ip": "10.0.0.1", "hostname": "host-1"}),
				},
				{
					State:  machineryv1tasks.StateSuccess,
					Args:   models.JSONMap{"url": "https://example.com/layer-1"},
					Result: models.JSONMap{"job_states": "invalid"},
				},
				{
					State:  machineryv1tasks.StateSuccess,
					Args:   models.JSONMap{},
					Result: mockJobResult(map[string]any{"id": "peer-1", "host_type": "normal", "ip": "10.0.0.1", "hostname": "host-1"}),
				},
				{
					State:  machineryv1tasks.StateSuccess,
					Args:   models.JSONMap{"url": "https://example.com/layer-1"},
					Result: mockJobResult(map[string]any{"id": "peer-1", "host_type": "normal", "ip": "10.0.0.1"}),
				},
			},
			expect: func(t *testing.T, peers []types.Peer, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Empty(peers)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := &service{}
			peers, err := s.extractPeersFromJobs(tc.jobs)
			tc.expect(t, peers, err)
		})
	}
}

func TestService_GetJob(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, job *models.Job, err error)
	}{
		{
			name: "job not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(job)
			},
		},
		{
			name: "loads clusters and owner",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				user := mockUser(t, s.db, "foo")
				schedulerCluster := mockSchedulerCluster(t, s.db, "cluster")
				seedPeerCluster := mockSeedPeerCluster(t, s.db, "seed")
				job := models.Job{
					TaskID:            "group-1",
					Type:              internaljob.PreheatJob,
					State:             machineryv1tasks.StateSuccess,
					Args:              models.JSONMap{},
					UserID:            user.ID,
					SchedulerClusters: []models.SchedulerCluster{schedulerCluster},
					SeedPeerClusters:  []models.SeedPeerCluster{seedPeerCluster},
				}
				assert.NoError(s.db.Create(&job).Error)
				return job.ID
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal("group-1", job.TaskID)
				assert.Equal("foo", job.User.Name)
				assert.Len(job.SchedulerClusters, 1)
				assert.Equal("cluster", job.SchedulerClusters[0].Name)
				assert.Len(job.SeedPeerClusters, 1)
				assert.Equal("seed", job.SeedPeerClusters[0].Name)
			},
		},
		{
			name: "owner missing",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				job := models.Job{TaskID: "group-1", Type: internaljob.PreheatJob, Args: models.JSONMap{}, UserID: 99}
				assert.NoError(s.db.Create(&job).Error)
				return job.ID
			},
			expect: func(t *testing.T, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(job)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			job, err := s.GetJob(context.Background(), id)
			tc.expect(t, job, err)
		})
	}
}

func TestService_UpdateJob(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		req    types.UpdateJobRequest
		expect func(t *testing.T, s *service, job *models.Job, err error)
	}{
		{
			name: "job not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			req: types.UpdateJobRequest{BIO: "bio"},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
				assert.Nil(job)
			},
		},
		{
			name: "updates bio and owner only",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				mockUser(t, s.db, "foo")
				schedulerCluster := mockSchedulerCluster(t, s.db, "cluster")
				job := models.Job{
					TaskID:            "group-1",
					Type:              internaljob.PreheatJob,
					State:             machineryv1tasks.StateSuccess,
					Args:              models.JSONMap{"url": "https://example.com/file"},
					SchedulerClusters: []models.SchedulerCluster{schedulerCluster},
				}
				assert.NoError(s.db.Create(&job).Error)
				return job.ID
			},
			req: types.UpdateJobRequest{BIO: "bio", UserID: 1},
			expect: func(t *testing.T, s *service, job *models.Job, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(job.SchedulerClusters, 1)

				stored := models.Job{}
				assert.NoError(s.db.First(&stored, job.ID).Error)
				assert.Equal("bio", stored.BIO)
				assert.Equal(uint(1), stored.UserID)
				assert.Equal("foo", stored.User.Name)
				assert.Equal("group-1", stored.TaskID)
				assert.Equal(machineryv1tasks.StateSuccess, stored.State)
				assert.Equal("https://example.com/file", stored.Args["url"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			job, err := s.UpdateJob(context.Background(), id, tc.req)
			tc.expect(t, s, job, err)
		})
	}
}

func TestService_DestroyJob(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(t *testing.T, s *service) uint
		expect func(t *testing.T, s *service, err error)
	}{
		{
			name: "job not found",
			setup: func(t *testing.T, s *service) uint {
				return 99
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.ErrorIs(err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "hard deletes",
			setup: func(t *testing.T, s *service) uint {
				assert := assert.New(t)
				job := models.Job{TaskID: "group-1", Type: internaljob.PreheatJob, Args: models.JSONMap{}}
				assert.NoError(s.db.Create(&job).Error)
				return job.ID
			},
			expect: func(t *testing.T, s *service, err error) {
				assert := assert.New(t)
				assert.NoError(err)

				var count int64
				assert.NoError(s.db.Unscoped().Model(&models.Job{}).Count(&count).Error)
				assert.Equal(int64(0), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestService(t)
			id := tc.setup(t, s)

			tc.expect(t, s, s.DestroyJob(context.Background(), id))
		})
	}
}

func TestService_GetJobs(t *testing.T) {
	s := newTestService(t)
	user := mockUser(t, s.db, "foo")
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	for i, job := range []models.Job{
		{TaskID: "group-1", Type: internaljob.PreheatJob, State: machineryv1tasks.StateSuccess, Args: models.JSONMap{}, UserID: user.ID},
		{TaskID: "group-2", Type: internaljob.PreheatJob, State: machineryv1tasks.StateFailure, Args: models.JSONMap{}},
		{TaskID: "group-3", Type: internaljob.GetTaskJob, State: machineryv1tasks.StateSuccess, Args: models.JSONMap{}, UserID: user.ID},
	} {
		job.CreatedAt = now.Add(time.Duration(i) * time.Minute)
		if err := s.db.Create(&job).Error; err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name   string
		query  types.GetJobsQuery
		expect func(t *testing.T, jobs []models.Job, count int64, err error)
	}{
		{
			name:  "newest first with pagination and owners",
			query: types.GetJobsQuery{Page: 1, PerPage: 2},
			expect: func(t *testing.T, jobs []models.Job, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(jobs, 2)
				assert.Equal(int64(3), count)
				assert.Equal("group-3", jobs[0].TaskID)
				assert.Equal("foo", jobs[0].User.Name)
				assert.Equal("group-2", jobs[1].TaskID)
				assert.Empty(jobs[1].User.Name)
			},
		},
		{
			name:  "last page",
			query: types.GetJobsQuery{Page: 2, PerPage: 2},
			expect: func(t *testing.T, jobs []models.Job, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(jobs, 1)
				assert.Equal("group-1", jobs[0].TaskID)
				assert.Equal(int64(3), count)
			},
		},
		{
			name:  "filter by type",
			query: types.GetJobsQuery{Type: internaljob.GetTaskJob, Page: 1, PerPage: 10},
			expect: func(t *testing.T, jobs []models.Job, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(jobs, 1)
				assert.Equal("group-3", jobs[0].TaskID)
				assert.Equal(int64(1), count)
			},
		},
		{
			name:  "filter by state",
			query: types.GetJobsQuery{State: machineryv1tasks.StateFailure, Page: 1, PerPage: 10},
			expect: func(t *testing.T, jobs []models.Job, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(jobs, 1)
				assert.Equal("group-2", jobs[0].TaskID)
			},
		},
		{
			name:  "filter by type and user",
			query: types.GetJobsQuery{Type: internaljob.PreheatJob, UserID: user.ID, Page: 1, PerPage: 10},
			expect: func(t *testing.T, jobs []models.Job, count int64, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Len(jobs, 1)
				assert.Equal("group-1", jobs[0].TaskID)
				assert.Equal(int64(1), count)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jobs, count, err := s.GetJobs(context.Background(), tc.query)
			tc.expect(t, jobs, count, err)
		})
	}
}
