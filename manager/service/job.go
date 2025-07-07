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

package service

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"sync"
	"time"

	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/job"
	"d7y.io/dragonfly/v2/manager/metrics"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
	pkggc "d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/idgen"
	"d7y.io/dragonfly/v2/pkg/net/http"
	"d7y.io/dragonfly/v2/pkg/retry"
	"d7y.io/dragonfly/v2/pkg/slices"
	"d7y.io/dragonfly/v2/pkg/structure"
	pkgtypes "d7y.io/dragonfly/v2/pkg/types"
)

const (
	// DefaultGCJobPollingTimeout is the default timeout for polling GC job.
	DefaultGCJobPollingTimeout = 30 * time.Minute

	// DefaultGCJobPollingInterval is the default interval for polling GC job.
	DefaultGCJobPollingInterval = 5 * time.Second

	// DefaultGetImageDistributionJobConcurrentCount is the default concurrent count for getting image distribution job.
	DefaultGetImageDistributionJobConcurrentCount = 30
)

func (s *service) CreateGCJob(ctx context.Context, json types.CreateGCJobRequest) (*models.Job, error) {
	taskID := uuid.NewString()
	ctx = context.WithValue(ctx, pkggc.ContextKeyTaskID, taskID)
	ctx = context.WithValue(ctx, pkggc.ContextKeyUserID, json.UserID)

	// This is a non-block function to run the gc task, which will run the task asynchronously in the backend.
	if err := s.gc.Run(ctx, json.Args.Type); err != nil {
		logger.Errorf("run gc job failed: %w", err)
		return nil, err
	}

	return s.pollingGCJob(ctx, json.Type, json.UserID, taskID)
}

func (s *service) pollingGCJob(ctx context.Context, jobType string, userID uint, taskID string) (*models.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultGCJobPollingTimeout)
	defer cancel()

	ticker := time.NewTicker(DefaultGCJobPollingInterval)
	defer ticker.Stop()

	job := models.Job{}

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context done: %w", ctx.Err())

		case <-ticker.C:
			if err := s.db.WithContext(ctx).First(&job, models.Job{
				Type:   jobType,
				UserID: userID,
				TaskID: taskID,
			}).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					continue
				}

				return nil, err
			}

			// Return the job if the job is in success or failure state, otherwise continue polling.
			if job.State == machineryv1tasks.StateSuccess || job.State == machineryv1tasks.StateFailure {
				return &job, nil
			}
		}
	}
}

func (s *service) CreateSyncPeersJob(ctx context.Context, json types.CreateSyncPeersJobRequest) error {
	schedulers, err := s.findSchedulerInClusters(ctx, json.SchedulerClusterIDs)
	if err != nil {
		logger.Errorf("find scheduler in clusters failed: %w", err)
		return err
	}

	return s.job.SyncPeers.CreateSyncPeers(ctx, schedulers)
}

func (s *service) CreatePreheatJob(ctx context.Context, json types.CreatePreheatJobRequest) (*models.Job, error) {
	if json.Args.Scope == "" {
		json.Args.Scope = types.SingleSeedPeerScope
	}

	if json.Args.ConcurrentCount == 0 {
		json.Args.ConcurrentCount = types.DefaultPreheatConcurrentCount
	}

	if json.Args.Timeout == 0 {
		json.Args.Timeout = types.DefaultJobTimeout
	}

	if json.Args.FilteredQueryParams == "" {
		json.Args.FilteredQueryParams = http.RawDefaultFilteredQueryParams
	}

	args, err := structure.StructToMap(json.Args)
	if err != nil {
		logger.Errorf("convert preheat args to map failed: %w", err)
		return nil, err
	}

	candidateSchedulers, err := s.findAllCandidateSchedulersInClusters(ctx, json.SchedulerClusterIDs, []string{types.SchedulerFeaturePreheat})
	if err != nil {
		logger.Errorf("find candidate schedulers in clusters failed: %w", err)
		return nil, err
	}

	groupJobState, err := s.job.CreatePreheat(ctx, candidateSchedulers, json.Args)
	if err != nil {
		logger.Errorf("create preheat job failed: %w", err)
		return nil, err
	}

	var candidateSchedulerClusters []models.SchedulerCluster
	for _, candidateScheduler := range candidateSchedulers {
		candidateSchedulerClusters = append(candidateSchedulerClusters, candidateScheduler.SchedulerCluster)
	}

	job := models.Job{
		TaskID:            groupJobState.GroupUUID,
		BIO:               json.BIO,
		Type:              json.Type,
		State:             groupJobState.State,
		Args:              args,
		UserID:            json.UserID,
		SchedulerClusters: candidateSchedulerClusters,
	}

	if err := s.db.WithContext(ctx).Create(&job).Error; err != nil {
		logger.Errorf("create preheat job failed: %w", err)
		return nil, err
	}

	go s.pollingJob(context.Background(), internaljob.PreheatJob, job.ID, job.TaskID, 30, 300, 16)
	return &job, nil
}

func (s *service) CreateGetTaskJob(ctx context.Context, json types.CreateGetTaskJobRequest) (*models.Job, error) {
	if json.Args.FilteredQueryParams == "" {
		json.Args.FilteredQueryParams = http.RawDefaultFilteredQueryParams
	}

	args, err := structure.StructToMap(json.Args)
	if err != nil {
		logger.Errorf("convert get task args to map failed: %w", err)
		return nil, err
	}

	schedulers, err := s.findAllSchedulersInClusters(ctx, json.SchedulerClusterIDs)
	if err != nil {
		logger.Errorf("find schedulers in clusters failed: %w", err)
		return nil, err
	}

	groupJobState, err := s.job.CreateGetTask(ctx, schedulers, json.Args)
	if err != nil {
		logger.Errorf("create get task job failed: %w", err)
		return nil, err
	}

	var schedulerClusters []models.SchedulerCluster
	for _, scheduler := range schedulers {
		schedulerClusters = append(schedulerClusters, scheduler.SchedulerCluster)
	}

	job := models.Job{
		TaskID:            groupJobState.GroupUUID,
		BIO:               json.BIO,
		Type:              json.Type,
		State:             groupJobState.State,
		Args:              args,
		UserID:            json.UserID,
		SchedulerClusters: schedulerClusters,
	}

	if err := s.db.WithContext(ctx).Create(&job).Error; err != nil {
		logger.Errorf("create get task job failed: %w", err)
		return nil, err
	}

	go s.pollingJob(context.Background(), internaljob.GetTaskJob, job.ID, job.TaskID, 30, 300, 16)
	logger.Infof("create get task job %s for task %s in scheduler clusters %v", job.ID, job.TaskID, json.SchedulerClusterIDs)
	return &job, nil
}

func (s *service) CreateGetImageDistributionJob(ctx context.Context, json types.CreateGetImageDistributionJobRequest) (*types.CreateGetImageDistributionJobResponse, error) {
	imageLayers, err := s.getImageLayers(ctx, json)
	if err != nil {
		err = fmt.Errorf("get image layers failed: %w", err)
		logger.Error(err)
		return nil, err
	}

	var layers []internaljob.PreheatRequest
	for _, imageLayer := range imageLayers {
		for _, url := range imageLayer.URLs {
			layers = append(layers, internaljob.PreheatRequest{
				URL:                 url,
				PieceLength:         imageLayer.PieceLength,
				Tag:                 imageLayer.Tag,
				Application:         imageLayer.Application,
				FilteredQueryParams: imageLayer.FilteredQueryParams,
			})
		}
	}

	if len(layers) == 0 {
		err = errors.New("no valid image layers found")
		logger.Error(err)
		return nil, err
	}

	image := types.Image{Layers: make([]types.Layer, 0, len(layers))}
	for _, file := range layers {
		image.Layers = append(image.Layers, types.Layer{URL: file.URL})
	}

	schedulers, err := s.findAllSchedulersInClusters(ctx, json.SchedulerClusterIDs)
	if err != nil {
		err = fmt.Errorf("find schedulers in clusters failed: %w", err)
		logger.Error(err)
		return nil, err
	}

	// Create multiple get task jobs synchronously for each layer and
	// extract peers from the jobs.
	jobs := s.createGetTaskJobsSync(ctx, layers, json, schedulers)
	peers, err := s.extractPeersFromJobs(jobs)
	if err != nil {
		err = fmt.Errorf("extract peers from jobs failed: %w", err)
		logger.Error(err)
		return nil, err
	}

	return &types.CreateGetImageDistributionJobResponse{
		Image: image,
		Peers: peers,
	}, nil
}

func (s *service) getImageLayers(ctx context.Context, json types.CreateGetImageDistributionJobRequest) ([]internaljob.PreheatRequest, error) {
	var certPool *x509.CertPool
	if len(s.config.Job.Preheat.TLS.CACert) != 0 {
		certPool = x509.NewCertPool()
		if !certPool.AppendCertsFromPEM([]byte(s.config.Job.Preheat.TLS.CACert)) {
			return nil, errors.New("invalid CA Cert")
		}
	}

	layers, err := job.GetImageLayers(ctx, types.PreheatArgs{
		Type:                job.PreheatImageType.String(),
		URL:                 json.Args.URL,
		PieceLength:         json.Args.PieceLength,
		Tag:                 json.Args.Tag,
		Application:         json.Args.Application,
		FilteredQueryParams: json.Args.FilteredQueryParams,
		Headers:             json.Args.Headers,
		Username:            json.Args.Username,
		Password:            json.Args.Password,
		Platform:            json.Args.Platform,
	}, s.config.Job.Preheat.RegistryTimeout, certPool, s.config.Job.Preheat.TLS.InsecureSkipVerify)
	if err != nil {
		return nil, fmt.Errorf("get image layers failed: %w", err)
	}

	if len(layers) == 0 {
		return nil, errors.New("no valid image layers found")
	}

	return layers, nil
}

func (s *service) createGetTaskJobsSync(ctx context.Context, layers []internaljob.PreheatRequest, json types.CreateGetImageDistributionJobRequest, schedulers []models.Scheduler) []*models.Job {
	var mu sync.Mutex
	jobs := make([]*models.Job, 0, len(layers))
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(DefaultGetImageDistributionJobConcurrentCount)
	for _, file := range layers {
		eg.Go(func() error {
			job, err := s.createGetTaskJobSync(ctx, types.CreateGetTaskJobRequest{
				BIO:  json.BIO,
				Type: internaljob.GetTaskJob,
				Args: types.GetTaskArgs{
					URL:                 file.URL,
					PieceLength:         file.PieceLength,
					Tag:                 file.Tag,
					Application:         file.Application,
					FilteredQueryParams: file.FilteredQueryParams,
				},
				SchedulerClusterIDs: json.SchedulerClusterIDs,
			}, schedulers)
			if err != nil {
				logger.Warnf("failed to create task job for image layer %s: %w", file.URL, err)
				return nil
			}

			mu.Lock()
			jobs = append(jobs, job)
			mu.Unlock()
			return nil
		})
	}

	// If any of the goroutines return an error, ignore it and continue processing.
	if err := eg.Wait(); err != nil {
		logger.Errorf("failed to create get task jobs: %w", err)
	}

	return jobs
}

func (s *service) createGetTaskJobSync(ctx context.Context, json types.CreateGetTaskJobRequest, schedulers []models.Scheduler) (*models.Job, error) {
	if json.Args.FilteredQueryParams == "" {
		json.Args.FilteredQueryParams = http.RawDefaultFilteredQueryParams
	}

	args, err := structure.StructToMap(json.Args)
	if err != nil {
		return nil, err
	}

	groupJobState, err := s.job.CreateGetTask(ctx, schedulers, json.Args)
	if err != nil {
		return nil, err
	}

	var schedulerClusters []models.SchedulerCluster
	for _, scheduler := range schedulers {
		schedulerClusters = append(schedulerClusters, scheduler.SchedulerCluster)
	}

	job := models.Job{
		TaskID:            groupJobState.GroupUUID,
		BIO:               json.BIO,
		Type:              json.Type,
		State:             groupJobState.State,
		Args:              args,
		UserID:            json.UserID,
		SchedulerClusters: schedulerClusters,
	}

	if err := s.db.WithContext(ctx).Create(&job).Error; err != nil {
		return nil, err
	}

	s.pollingJob(context.Background(), internaljob.GetTaskJob, job.ID, job.TaskID, 3, 5, 4)
	if err := s.db.WithContext(ctx).Preload("SeedPeerClusters").Preload("SchedulerClusters").First(&job, job.ID).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *service) extractPeersFromJobs(jobs []*models.Job) ([]types.Peer, error) {
	m := make(map[string]*types.Peer, len(jobs))
	for _, job := range jobs {
		if job.State != machineryv1tasks.StateSuccess {
			continue
		}

		jobStates, ok := job.Result["job_states"].([]any)
		if !ok {
			logger.Warnf("job %s has no job_states in result", job.ID)
			continue
		}

		url, ok := job.Args["url"].(string)
		if !ok {
			logger.Warnf("job %s has no url in args", job.ID)
			continue
		}

		for _, jobState := range jobStates {
			jobState, ok := jobState.(map[string]any)
			if !ok {
				logger.Warnf("job %s has invalid job_state in result", job.ID)
				continue
			}

			results, ok := jobState["results"].([]any)
			if !ok {
				logger.Warnf("job %s has no results in job_state", job.ID)
				continue
			}

			for _, result := range results {
				result, ok := result.(map[string]any)
				if !ok {
					logger.Warnf("job %s has invalid result in job_state", job.ID)
					continue
				}

				schedulerClusterID, ok := result["scheduler_cluster_id"].(float64)
				if !ok {
					logger.Warnf("job %s has no scheduler_cluster_id in result", job.ID)
					continue
				}

				peers, ok := result["peers"].([]any)
				if !ok {
					logger.Warnf("job %s has no peers in result", job.ID)
					continue
				}

				for _, peer := range peers {
					peer, ok := peer.(map[string]any)
					if !ok {
						logger.Warnf("job %s has invalid peer in result", job.ID)
						continue
					}

					id, ok := peer["id"].(string)
					if !ok {
						logger.Warnf("job %s has invalid peer id in result", job.ID)
						continue
					}

					hostType, ok := peer["host_type"].(string)
					if !ok {
						logger.Warnf("job %s has no host_type in result for peer %s", job.ID, id)
						continue
					}

					// Only collect normal peers and skip seed peers.
					if hostType != pkgtypes.HostTypeNormalName {
						continue
					}

					ip, ok := peer["ip"].(string)
					if !ok {
						logger.Warnf("job %s has no ip in result for peer %s", job.ID, id)
						continue
					}

					hostname, ok := peer["hostname"].(string)
					if !ok {
						logger.Warnf("job %s has no hostname in result for peer %s", job.ID, id)
						continue
					}

					hostID := idgen.HostIDV2(ip, hostname, false)
					p, found := m[hostID]
					if !found {
						m[hostID] = &types.Peer{
							IP:                 ip,
							Hostname:           hostname,
							CachedLayers:       []types.Layer{{URL: url}},
							SchedulerClusterID: uint(schedulerClusterID),
						}
					} else {
						if slices.Contains(p.CachedLayers, types.Layer{URL: url}) {
							continue
						}

						p.CachedLayers = append(p.CachedLayers, types.Layer{URL: url})
					}
				}
			}
		}
	}

	peers := make([]types.Peer, 0, len(m))
	for _, peer := range m {
		peers = append(peers, *peer)
	}

	return peers, nil
}

func (s *service) CreateDeleteTaskJob(ctx context.Context, json types.CreateDeleteTaskJobRequest) (*models.Job, error) {
	if json.Args.Timeout == 0 {
		json.Args.Timeout = types.DefaultJobTimeout
	}

	if json.Args.FilteredQueryParams == "" {
		json.Args.FilteredQueryParams = http.RawDefaultFilteredQueryParams
	}

	args, err := structure.StructToMap(json.Args)
	if err != nil {
		logger.Errorf("convert delete task args to map failed: %w", err)
		return nil, err
	}

	schedulers, err := s.findAllSchedulersInClusters(ctx, json.SchedulerClusterIDs)
	if err != nil {
		logger.Errorf("find schedulers in clusters failed: %w", err)
		return nil, err
	}

	groupJobState, err := s.job.CreateDeleteTask(ctx, schedulers, json.Args)
	if err != nil {
		logger.Errorf("create delete task job failed: %w", err)
		return nil, err
	}

	var schedulerClusters []models.SchedulerCluster
	for _, scheduler := range schedulers {
		schedulerClusters = append(schedulerClusters, scheduler.SchedulerCluster)
	}

	job := models.Job{
		TaskID:            groupJobState.GroupUUID,
		BIO:               json.BIO,
		Type:              json.Type,
		State:             groupJobState.State,
		Args:              args,
		UserID:            json.UserID,
		SchedulerClusters: schedulerClusters,
	}

	if err := s.db.WithContext(ctx).Create(&job).Error; err != nil {
		logger.Errorf("create delete task job failed: %w", err)
		return nil, err
	}

	go s.pollingJob(context.Background(), internaljob.DeleteTaskJob, job.ID, job.TaskID, 30, 300, 16)
	return &job, nil
}

func (s *service) findSchedulerInClusters(ctx context.Context, schedulerClusterIDs []uint) ([]models.Scheduler, error) {
	var activeSchedulers []models.Scheduler
	if len(schedulerClusterIDs) != 0 {
		// Find the scheduler clusters by request.
		for _, schedulerClusterID := range schedulerClusterIDs {
			schedulerCluster := models.SchedulerCluster{}
			if err := s.db.WithContext(ctx).First(&schedulerCluster, schedulerClusterID).Error; err != nil {
				return nil, fmt.Errorf("scheduler cluster id %d: %w", schedulerClusterID, err)
			}

			scheduler := models.Scheduler{}
			if err := s.db.WithContext(ctx).Preload("SchedulerCluster").First(&scheduler, models.Scheduler{
				SchedulerClusterID: schedulerCluster.ID,
				State:              models.SchedulerStateActive,
			}).Error; err != nil {
				return nil, fmt.Errorf("scheduler cluster id %d: %w", schedulerClusterID, err)
			}

			activeSchedulers = append(activeSchedulers, scheduler)
		}
	} else {
		// Find all of the scheduler clusters that has active scheduler.
		var schedulerClusters []models.SchedulerCluster
		if err := s.db.WithContext(ctx).Find(&schedulerClusters).Error; err != nil {
			return nil, err
		}

		for _, schedulerCluster := range schedulerClusters {
			scheduler := models.Scheduler{}
			if err := s.db.WithContext(ctx).Preload("SchedulerCluster").First(&scheduler, models.Scheduler{
				SchedulerClusterID: schedulerCluster.ID,
				State:              models.SchedulerStateActive,
			}).Error; err != nil {
				continue
			}

			activeSchedulers = append(activeSchedulers, scheduler)
		}
	}

	if len(activeSchedulers) == 0 {
		return nil, errors.New("active schedulers not found")
	}

	return activeSchedulers, nil
}

func (s *service) findAllSchedulersInClusters(ctx context.Context, schedulerClusterIDs []uint) ([]models.Scheduler, error) {
	var activeSchedulers []models.Scheduler
	if len(schedulerClusterIDs) != 0 {
		// Find the scheduler clusters by request.
		for _, schedulerClusterID := range schedulerClusterIDs {
			schedulerCluster := models.SchedulerCluster{}
			if err := s.db.WithContext(ctx).First(&schedulerCluster, schedulerClusterID).Error; err != nil {
				return nil, fmt.Errorf("scheduler cluster id %d: %w", schedulerClusterID, err)
			}

			var schedulers []models.Scheduler
			if err := s.db.WithContext(ctx).Preload("SchedulerCluster").Find(&schedulers, models.Scheduler{
				SchedulerClusterID: schedulerCluster.ID,
				State:              models.SchedulerStateActive,
			}).Error; err != nil {
				return nil, fmt.Errorf("scheduler cluster id %d: %w", schedulerClusterID, err)
			}

			activeSchedulers = append(activeSchedulers, schedulers...)
		}
	} else {
		// Find all of the scheduler clusters that has active schedulers.
		var schedulerClusters []models.SchedulerCluster
		if err := s.db.WithContext(ctx).Find(&schedulerClusters).Error; err != nil {
			return nil, err
		}

		for _, schedulerCluster := range schedulerClusters {
			var schedulers []models.Scheduler
			if err := s.db.WithContext(ctx).Preload("SchedulerCluster").Find(&schedulers, models.Scheduler{
				SchedulerClusterID: schedulerCluster.ID,
				State:              models.SchedulerStateActive,
			}).Error; err != nil {
				continue
			}

			activeSchedulers = append(activeSchedulers, schedulers...)
		}
	}

	if len(activeSchedulers) == 0 {
		return nil, errors.New("active schedulers not found")
	}

	return activeSchedulers, nil
}

func (s *service) findAllCandidateSchedulersInClusters(ctx context.Context, schedulerClusterIDs []uint, features []string) ([]models.Scheduler, error) {
	var candidateSchedulers []models.Scheduler
	if len(schedulerClusterIDs) != 0 {
		// Find the scheduler clusters by request.
		for _, schedulerClusterID := range schedulerClusterIDs {
			schedulerCluster := models.SchedulerCluster{}
			if err := s.db.WithContext(ctx).First(&schedulerCluster, schedulerClusterID).Error; err != nil {
				return nil, fmt.Errorf("scheduler cluster id %d: %w", schedulerClusterID, err)
			}

			var schedulers []models.Scheduler
			if err := s.db.WithContext(ctx).Preload("SchedulerCluster").Find(&schedulers, models.Scheduler{
				SchedulerClusterID: schedulerCluster.ID,
				State:              models.SchedulerStateActive,
			}).Error; err != nil {
				return nil, fmt.Errorf("scheduler cluster id %d: %w", schedulerClusterID, err)
			}

			for _, scheduler := range schedulers {
				// If the features is empty, return the first scheduler directly.
				if len(features) == 0 {
					candidateSchedulers = append(candidateSchedulers, scheduler)
					break
				}

				// Scan the schedulers to find the first scheduler that supports feature.
				if slices.Contains(scheduler.Features, features...) {
					candidateSchedulers = append(candidateSchedulers, scheduler)
					break
				}
			}
		}
	} else {
		// Find all of the scheduler clusters that has active schedulers.
		var candidateSchedulerClusters []models.SchedulerCluster
		if err := s.db.WithContext(ctx).Find(&candidateSchedulerClusters).Error; err != nil {
			return nil, err
		}

		for _, candidateSchedulerCluster := range candidateSchedulerClusters {
			var schedulers []models.Scheduler
			if err := s.db.WithContext(ctx).Preload("SchedulerCluster").Find(&schedulers, models.Scheduler{
				SchedulerClusterID: candidateSchedulerCluster.ID,
				State:              models.SchedulerStateActive,
			}).Error; err != nil {
				continue
			}

			for _, scheduler := range schedulers {
				// If the features is empty, return the first scheduler directly.
				if len(features) == 0 {
					candidateSchedulers = append(candidateSchedulers, scheduler)
					break
				}

				// Scan the schedulers to find the first scheduler that supports feature.
				if slices.Contains(scheduler.Features, features...) {
					candidateSchedulers = append(candidateSchedulers, scheduler)
					break
				}
			}
		}
	}

	if len(candidateSchedulers) == 0 {
		return nil, errors.New("candidate schedulers not found")
	}

	return candidateSchedulers, nil
}

func (s *service) pollingJob(ctx context.Context, name string, id uint, groupUUID string, initBackoff float64, maxBackoff float64, maxAttempts int) {
	var (
		job models.Job
		log = logger.WithGroupAndJobID(groupUUID, fmt.Sprint(id))
	)
	if _, _, err := retry.Run(ctx, initBackoff, maxBackoff, maxAttempts, func() (any, bool, error) {
		groupJob, err := s.job.GetGroupJobState(name, groupUUID)
		if err != nil {
			err = fmt.Errorf("get group job state failed: %w", err)
			log.Error(err)
			return nil, false, err
		}

		result, err := structure.StructToMap(groupJob)
		if err != nil {
			err = fmt.Errorf("convert group job state to map failed: %w", err)
			log.Error(err)
			return nil, false, err
		}

		if err := s.db.WithContext(ctx).First(&job, id).Updates(models.Job{
			State:  groupJob.State,
			Result: result,
		}).Error; err != nil {
			err = fmt.Errorf("update job state failed: %w", err)
			log.Error(err)
			return nil, true, err
		}

		switch job.State {
		case machineryv1tasks.StateSuccess:
			// Collect CreateJobSuccessCount. metrics.
			metrics.CreateJobSuccessCount.WithLabelValues(name).Inc()

			log.Info("polling group succeeded")
			return nil, true, nil
		case machineryv1tasks.StateFailure:
			log.Error("polling group failed")
			return nil, true, nil
		default:
			msg := fmt.Sprintf("polling job state is %s", job.State)
			log.Info(msg)
			return nil, false, errors.New(msg)
		}
	}); err != nil {
		err = fmt.Errorf("polling group job failed: %w", err)
		log.Error(err)
	}

	// Polling timeout and failed.
	if job.State != machineryv1tasks.StateSuccess && job.State != machineryv1tasks.StateFailure {
		job := models.Job{}
		if err := s.db.WithContext(ctx).First(&job, id).Updates(models.Job{
			State: machineryv1tasks.StateFailure,
		}).Error; err != nil {
			err = fmt.Errorf("update job state to failure failed: %w", err)
			log.Error(err)
		}

		log.Error("polling group timeout")
	}
}

func (s *service) DestroyJob(ctx context.Context, id uint) error {
	job := models.Job{}
	if err := s.db.WithContext(ctx).First(&job, id).Error; err != nil {
		return err
	}

	if err := s.db.WithContext(ctx).Unscoped().Delete(&models.Job{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (s *service) UpdateJob(ctx context.Context, id uint, json types.UpdateJobRequest) (*models.Job, error) {
	job := models.Job{}
	if err := s.db.WithContext(ctx).Preload("SeedPeerClusters").Preload("SchedulerClusters").First(&job, id).Updates(models.Job{
		BIO:    json.BIO,
		UserID: json.UserID,
	}).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *service) GetJob(ctx context.Context, id uint) (*models.Job, error) {
	job := models.Job{}
	if err := s.db.WithContext(ctx).Preload("SeedPeerClusters").Preload("SchedulerClusters").First(&job, id).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

func (s *service) GetJobs(ctx context.Context, q types.GetJobsQuery) ([]models.Job, int64, error) {
	var count int64
	var jobs []models.Job
	if err := s.db.WithContext(ctx).Scopes(models.Paginate(q.Page, q.PerPage)).Where(&models.Job{
		Type:   q.Type,
		State:  q.State,
		UserID: q.UserID,
	}).Order("created_at DESC").Find(&jobs).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return jobs, count, nil
}

func (s *service) AddJobToSchedulerClusters(ctx context.Context, id, schedulerClusterIDs []uint) error {
	job := models.Job{}
	if err := s.db.WithContext(ctx).First(&job, id).Error; err != nil {
		return err
	}

	var schedulerClusters []*models.SchedulerCluster
	for _, schedulerClusterID := range schedulerClusterIDs {
		schedulerCluster := models.SchedulerCluster{}
		if err := s.db.WithContext(ctx).First(&schedulerCluster, schedulerClusterID).Error; err != nil {
			return err
		}
		schedulerClusters = append(schedulerClusters, &schedulerCluster)
	}

	if err := s.db.WithContext(ctx).Model(&job).Association("SchedulerClusters").Append(schedulerClusters); err != nil {
		return err
	}

	return nil
}

func (s *service) AddJobToSeedPeerClusters(ctx context.Context, id, seedPeerClusterIDs []uint) error {
	job := models.Job{}
	if err := s.db.WithContext(ctx).First(&job, id).Error; err != nil {
		return err
	}

	var seedPeerClusters []*models.SeedPeerCluster
	for _, seedPeerClusterID := range seedPeerClusterIDs {
		seedPeerCluster := models.SeedPeerCluster{}
		if err := s.db.WithContext(ctx).First(&seedPeerCluster, seedPeerClusterID).Error; err != nil {
			return err
		}
		seedPeerClusters = append(seedPeerClusters, &seedPeerCluster)
	}

	if err := s.db.WithContext(ctx).Model(&job).Association("SeedPeerClusters").Append(seedPeerClusters); err != nil {
		return err
	}

	return nil
}
