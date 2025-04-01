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

//go:generate mockgen -destination mocks/scheduling_mock.go -source scheduling.go -package mocks

package scheduling

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"
	schedulerv1 "d7y.io/api/v2/pkg/apis/scheduler/v1"
	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	"d7y.io/dragonfly/v2/pkg/container/set"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
	"d7y.io/dragonfly/v2/scheduler/scheduling/evaluator"
)

type Scheduling interface {
	// ScheduleCandidateParents schedules candidate parents to the normal peer to download the task.
	// Used only in v2 version of the grpc.
	ScheduleCandidateParents(context.Context, *standard.Peer, set.SafeSet[string]) error

	// ScheduleParentAndCandidateParents schedules a parent and candidate parents to the normal peer to download the task.
	// Used only in v1 version of the grpc.
	ScheduleParentAndCandidateParents(context.Context, *standard.Peer, set.SafeSet[string])

	// FindCandidateParents finds candidate parents for the peer to download the task.
	// Used only in v2 version of the grpc.
	FindCandidateParents(context.Context, *standard.Peer, set.SafeSet[string]) ([]*standard.Peer, bool)

	// FindParentAndCandidateParents finds a parent and candidate parents for the peer to download the task.
	// Used only in v1 version of the grpc.
	FindParentAndCandidateParents(context.Context, *standard.Peer, set.SafeSet[string]) ([]*standard.Peer, bool)

	// FindSuccessParent finds success parent for the peer to download the task.
	FindSuccessParent(context.Context, *standard.Peer, set.SafeSet[string]) (*standard.Peer, bool)

	// FindReplicatePersistentCacheHosts finds replicate persistent cache hosts for the peer to replicate the task. It will compare the current
	// persistent replica count with the persistent replica count and try to find enough parents. Then function will return the cached replicate parents,
	// the replicate hosts without cache and found flag.
	FindReplicatePersistentCacheHosts(context.Context, *persistentcache.Task, set.SafeSet[string]) ([]*persistentcache.Peer, []*persistentcache.Host, bool)

	// FindCandidatePersistentCacheParents finds candidate persistent cache parents for the peer to download the task.
	FindCandidatePersistentCacheParents(context.Context, *persistentcache.Peer, set.SafeSet[string]) ([]*persistentcache.Peer, bool)
}

type scheduling struct {
	// Evaluator interface.
	evaluator evaluator.Evaluator

	// Scheduler configuration.
	config *config.SchedulerConfig

	// Persistent cache resource.
	persistentCacheResource persistentcache.Resource

	// Scheduler dynamic configuration.
	dynconfig config.DynconfigInterface
}

func New(cfg *config.SchedulerConfig, persistentCacheResource persistentcache.Resource, dynconfig config.DynconfigInterface, pluginDir string) Scheduling {
	return &scheduling{
		evaluator:               evaluator.New(cfg.Algorithm, pluginDir),
		config:                  cfg,
		persistentCacheResource: persistentCacheResource,
		dynconfig:               dynconfig,
	}
}

// ScheduleCandidateParents schedules candidate parents to the normal peer.
// Used only in v2 version of the grpc.
func (s *scheduling) ScheduleCandidateParents(ctx context.Context, peer *standard.Peer, blocklist set.SafeSet[string]) error {
	var n int
	for {
		select {
		case <-ctx.Done():
			peer.Log.Infof("context was done")
			return ctx.Err()
		default:
		}

		// Scheduling will send NeedBackToSourceResponse to peer.
		//
		// Condition 1: Peer's NeedBackToSource is true.
		// Condition 2: Scheduling exceeds the RetryBackToSourceLimit.
		if peer.Task.CanBackToSource() {
			// Check condition 1:
			// Peer's NeedBackToSource is true.
			if peer.NeedBackToSource.Load() {
				stream, loaded := peer.LoadAnnouncePeerStream()
				if !loaded {
					peer.Log.Error("load stream failed")
					return status.Error(codes.FailedPrecondition, "load stream failed")
				}

				// Send NeedBackToSourceResponse to peer.
				peer.Log.Infof("send NeedBackToSourceResponse, because of peer's NeedBackToSource is %t", peer.NeedBackToSource.Load())
				description := fmt.Sprintf("peer's NeedBackToSource is %t", peer.NeedBackToSource.Load())
				if err := stream.Send(&schedulerv2.AnnouncePeerResponse{
					Response: &schedulerv2.AnnouncePeerResponse_NeedBackToSourceResponse{
						NeedBackToSourceResponse: &schedulerv2.NeedBackToSourceResponse{
							Description: &description,
						},
					},
				}); err != nil {
					peer.Log.Error(err)
					return status.Error(codes.FailedPrecondition, err.Error())
				}

				return nil
			}

			// Check condition 2:
			// The number of retry scheduling is greater than RetryBackToSourceLimit
			if n >= s.config.RetryBackToSourceLimit {
				stream, loaded := peer.LoadAnnouncePeerStream()
				if !loaded {
					peer.Log.Error("load stream failed")
					return status.Error(codes.FailedPrecondition, "load stream failed")
				}

				// Send NeedBackToSourceResponse to peer.
				peer.Log.Infof("send NeedBackToSourceResponse, because of scheduling exceeded RetryBackToSourceLimit %d", s.config.RetryBackToSourceLimit)
				description := "scheduling exceeded RetryBackToSourceLimit"
				if err := stream.Send(&schedulerv2.AnnouncePeerResponse{
					Response: &schedulerv2.AnnouncePeerResponse_NeedBackToSourceResponse{
						NeedBackToSourceResponse: &schedulerv2.NeedBackToSourceResponse{
							Description: &description,
						},
					},
				}); err != nil {
					peer.Log.Error(err)
					return status.Error(codes.FailedPrecondition, err.Error())
				}

				return nil
			}
		}

		// Scheduling will return schedule failed.
		//
		// Condition 1: Scheduling exceeds the RetryLimit.
		if n >= s.config.RetryLimit {
			peer.Log.Errorf("scheduling failed, because of scheduling exceeded RetryLimit %d", s.config.RetryLimit)
			return status.Error(codes.FailedPrecondition, "scheduling exceeded RetryLimit")
		}

		// Scheduling will send NormalTaskResponse to peer.
		//
		// Condition 1: Scheduling can find candidate parents.
		if err := peer.Task.DeletePeerInEdges(peer.ID); err != nil {
			peer.Log.Error(err)
			return status.Error(codes.Internal, err.Error())
		}

		// Find candidate parents.
		candidateParents, found := s.FindCandidateParents(ctx, peer, blocklist)
		if !found {
			n++
			peer.Log.Infof("scheduling failed in %d times, because of candidate parents not found", n)

			// Sleep to avoid hot looping.
			time.Sleep(s.config.RetryInterval)
			continue
		}

		// Load AnnouncePeerStream from peer.
		stream, loaded := peer.LoadAnnouncePeerStream()
		if !loaded {
			if err := peer.Task.DeletePeerInEdges(peer.ID); err != nil {
				msg := fmt.Sprintf("peer deletes inedges failed: %s", err.Error())
				peer.Log.Error(msg)
				return status.Error(codes.Internal, msg)
			}

			peer.Log.Error("load stream failed")
			return status.Error(codes.FailedPrecondition, "load stream failed")
		}

		// Send NormalTaskResponse to peer.
		peer.Log.Info("send NormalTaskResponse")
		if err := stream.Send(&schedulerv2.AnnouncePeerResponse{
			Response: constructSuccessNormalTaskResponse(candidateParents),
		}); err != nil {
			peer.Log.Error(err)
			return status.Error(codes.FailedPrecondition, err.Error())
		}

		// Add edge from parent to peer.
		for _, candidateParent := range candidateParents {
			if err := peer.Task.AddPeerEdge(candidateParent, peer); err != nil {
				peer.Log.Warnf("peer adds edge failed: %s", err.Error())
				continue
			}
		}

		peer.Log.Infof("scheduling success in %d times", n+1)
		return nil
	}
}

// ScheduleParentAndCandidateParents schedules a parent and candidate parents to a peer.
// Used only in v1 version of the grpc.
func (s *scheduling) ScheduleParentAndCandidateParents(ctx context.Context, peer *standard.Peer, blocklist set.SafeSet[string]) {
	var n int
	for {
		select {
		case <-ctx.Done():
			peer.Log.Infof("context was done")
			return
		default:
		}

		// Scheduling will send Code_SchedNeedBackSource to peer.
		//
		// Condition 1: Peer's NeedBackToSource is true.
		// Condition 2: Scheduling exceeds the RetryBackToSourceLimit.
		if peer.Task.CanBackToSource() {
			// Check condition 1:
			// Peer's NeedBackToSource is true.
			if peer.NeedBackToSource.Load() {
				stream, loaded := peer.LoadReportPieceResultStream()
				if !loaded {
					peer.Log.Error("load stream failed")
					return
				}

				// Send Code_SchedNeedBackSource to peer.
				if err := stream.Send(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource}); err != nil {
					peer.Log.Error(err)
					return
				}
				peer.Log.Infof("send Code_SchedNeedBackSource to peer, because of peer's NeedBackToSource is %t", peer.NeedBackToSource.Load())

				if err := peer.FSM.Event(ctx, standard.PeerEventDownloadBackToSource); err != nil {
					peer.Log.Errorf("peer fsm event failed: %s", err.Error())
					return
				}

				// If the task state is TaskStateFailed,
				// peer back-to-source and reset task state to TaskStateRunning.
				if peer.Task.FSM.Is(standard.TaskStateFailed) {
					if err := peer.Task.FSM.Event(ctx, standard.TaskEventDownload); err != nil {
						peer.Task.Log.Errorf("task fsm event failed: %s", err.Error())
						return
					}
				}

				return
			}

			// Check condition 2:
			// The number of retry scheduling is greater than RetryBackToSourceLimit
			if n >= s.config.RetryBackToSourceLimit {
				stream, loaded := peer.LoadReportPieceResultStream()
				if !loaded {
					peer.Log.Error("load stream failed")
					return
				}

				// Send Code_SchedNeedBackSource peer.
				if err := stream.Send(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedNeedBackSource}); err != nil {
					peer.Log.Error(err)
					return
				}
				peer.Log.Infof("send Code_SchedNeedBackSource to peer, because of scheduling exceeded RetryBackToSourceLimit %d", s.config.RetryBackToSourceLimit)

				if err := peer.FSM.Event(ctx, standard.PeerEventDownloadBackToSource); err != nil {
					peer.Log.Errorf("peer fsm event failed: %s", err.Error())
					return
				}

				// If the task state is TaskStateFailed,
				// peer back-to-source and reset task state to TaskStateRunning.
				if peer.Task.FSM.Is(standard.TaskStateFailed) {
					if err := peer.Task.FSM.Event(ctx, standard.TaskEventDownload); err != nil {
						peer.Task.Log.Errorf("task fsm event failed: %s", err.Error())
						return
					}
				}

				return
			}
		}

		// Scheduling will send Code_SchedTaskStatusError to peer.
		//
		// Condition 1: Scheduling exceeds the RetryLimit.
		if n >= s.config.RetryLimit {
			stream, loaded := peer.LoadReportPieceResultStream()
			if !loaded {
				peer.Log.Error("load stream failed")
				return
			}

			// Send Code_SchedTaskStatusError to peer.
			if err := stream.Send(&schedulerv1.PeerPacket{Code: commonv1.Code_SchedTaskStatusError}); err != nil {
				peer.Log.Error(err)
				return
			}

			peer.Log.Errorf("send SchedulePeerFailed to peer, because of scheduling exceeded RetryLimit %d", s.config.RetryLimit)
			return
		}

		// Scheduling will send PeerPacket to peer.
		//
		// Condition 1: Scheduling can find candidate parents.
		if err := peer.Task.DeletePeerInEdges(peer.ID); err != nil {
			n++
			peer.Log.Errorf("scheduling failed in %d times, because of %s", n, err.Error())

			// Sleep to avoid hot looping.
			time.Sleep(s.config.RetryInterval)
			continue
		}

		// Find candidate parents.
		candidateParents, found := s.FindCandidateParents(ctx, peer, blocklist)
		if !found {
			n++
			peer.Log.Infof("scheduling failed in %d times, because of candidate parents not found", n)

			// Sleep to avoid hot looping.
			time.Sleep(s.config.RetryInterval)
			continue
		}

		// Load ReportPieceResultStream from peer.
		stream, loaded := peer.LoadReportPieceResultStream()
		if !loaded {
			n++
			peer.Log.Errorf("scheduling failed in %d times, because of loading peer stream failed", n)

			if err := peer.Task.DeletePeerInEdges(peer.ID); err != nil {
				peer.Log.Errorf("peer deletes inedges failed: %s", err.Error())
				return
			}

			return
		}

		// Send PeerPacket to peer.
		peer.Log.Info("send PeerPacket to peer")
		if err := stream.Send(constructSuccessPeerPacket(peer, candidateParents[0], candidateParents[1:])); err != nil {
			n++
			peer.Log.Errorf("scheduling failed in %d times, because of %s", n, err.Error())

			if err := peer.Task.DeletePeerInEdges(peer.ID); err != nil {
				peer.Log.Errorf("peer deletes inedges failed: %s", err.Error())
				return
			}

			return
		}

		// Add edge from parent to peer.
		for _, candidateParent := range candidateParents {
			if err := peer.Task.AddPeerEdge(candidateParent, peer); err != nil {
				peer.Log.Debugf("peer adds edge failed: %s", err.Error())
				continue
			}
		}

		peer.Log.Infof("scheduling success in %d times", n+1)
		return
	}
}

// FindCandidateParents finds candidate parents for the peer.
func (s *scheduling) FindCandidateParents(ctx context.Context, peer *standard.Peer, blocklist set.SafeSet[string]) ([]*standard.Peer, bool) {
	// Only PeerStateReceivedNormal and PeerStateRunning peers need to be rescheduled,
	// and other states including the PeerStateBackToSource indicate that
	// they have been scheduled.
	if !(peer.FSM.Is(standard.PeerStateReceivedNormal) || peer.FSM.Is(standard.PeerStateRunning)) {
		peer.Log.Infof("peer state is %s, can not schedule parent", peer.FSM.Current())
		return []*standard.Peer{}, false
	}

	// Find the candidate parent that can be scheduled.
	candidateParents := s.filterCandidateParents(peer, blocklist)
	if len(candidateParents) == 0 {
		peer.Log.Info("can not find candidate parents")
		return []*standard.Peer{}, false
	}

	// Sort candidate parents by evaluation score.
	taskTotalPieceCount := peer.Task.TotalPieceCount.Load()
	candidateParents = s.evaluator.EvaluateParents(candidateParents, peer, uint32(taskTotalPieceCount))

	// Get the parents with candidateParentLimit.
	candidateParentLimit := config.DefaultSchedulerCandidateParentLimit
	if config, err := s.dynconfig.GetSchedulerClusterConfig(); err == nil {
		if config.CandidateParentLimit > 0 {
			candidateParentLimit = int(config.CandidateParentLimit)
		}
	}

	if len(candidateParents) > candidateParentLimit {
		candidateParents = candidateParents[:candidateParentLimit]
	}

	var parentIDs []string
	for _, candidateParent := range candidateParents {
		parentIDs = append(parentIDs, candidateParent.ID)
	}

	peer.Log.Infof("scheduling candidate parents is %#v", parentIDs)
	return candidateParents, true
}

// FindParentAndCandidateParents finds a parent and candidate parents for the peer.
func (s *scheduling) FindParentAndCandidateParents(ctx context.Context, peer *standard.Peer, blocklist set.SafeSet[string]) ([]*standard.Peer, bool) {
	// Only PeerStateRunning peers need to be rescheduled,
	// and other states including the PeerStateBackToSource indicate that
	// they have been scheduled.
	if !peer.FSM.Is(standard.PeerStateRunning) {
		peer.Log.Infof("peer state is %s, can not schedule parent", peer.FSM.Current())
		return []*standard.Peer{}, false
	}

	// Find the candidate parent that can be scheduled.
	candidateParents := s.filterCandidateParents(peer, blocklist)
	if len(candidateParents) == 0 {
		peer.Log.Info("can not find candidate parents")
		return []*standard.Peer{}, false
	}

	// Sort candidate parents by evaluation score.
	taskTotalPieceCount := peer.Task.TotalPieceCount.Load()
	candidateParents = s.evaluator.EvaluateParents(candidateParents, peer, uint32(taskTotalPieceCount))

	// Get the parents with candidateParentLimit.
	candidateParentLimit := config.DefaultSchedulerCandidateParentLimit
	if config, err := s.dynconfig.GetSchedulerClusterConfig(); err == nil {
		if config.CandidateParentLimit > 0 {
			candidateParentLimit = int(config.CandidateParentLimit)
		}
	}

	if len(candidateParents) > candidateParentLimit {
		candidateParents = candidateParents[:candidateParentLimit]
	}

	var parentIDs []string
	for _, candidateParent := range candidateParents {
		parentIDs = append(parentIDs, candidateParent.ID)
	}

	peer.Log.Infof("scheduling candidate parents is %#v", parentIDs)
	return candidateParents, true
}

// FindSuccessParent finds success parent for the peer.
func (s *scheduling) FindSuccessParent(ctx context.Context, peer *standard.Peer, blocklist set.SafeSet[string]) (*standard.Peer, bool) {
	// Only PeerStateRunning peers need to be rescheduled,
	// and other states including the PeerStateBackToSource indicate that
	// they have been scheduled.
	if !peer.FSM.Is(standard.PeerStateRunning) {
		peer.Log.Infof("peer state is %s, can not schedule parent", peer.FSM.Current())
		return nil, false
	}

	// Find the candidate parent that can be scheduled.
	candidateParents := s.filterCandidateParents(peer, blocklist)
	if len(candidateParents) == 0 {
		peer.Log.Info("can not find candidate parents")
		return nil, false
	}

	var successParents []*standard.Peer
	for _, candidateParent := range candidateParents {
		if candidateParent.FSM.Is(standard.PeerStateSucceeded) {
			successParents = append(successParents, candidateParent)
		}
	}

	// Sort candidate parents by evaluation score.
	taskTotalPieceCount := peer.Task.TotalPieceCount.Load()
	successParents = s.evaluator.EvaluateParents(successParents, peer, uint32(taskTotalPieceCount))

	peer.Log.Infof("scheduling success parent is %s", successParents[0].ID)
	return successParents[0], true
}

// filterCandidateParents filters the candidate parents that can be scheduled.
func (s *scheduling) filterCandidateParents(peer *standard.Peer, blocklist set.SafeSet[string]) []*standard.Peer {
	filterParentLimit := config.DefaultSchedulerFilterParentLimit
	if config, err := s.dynconfig.GetSchedulerClusterConfig(); err == nil {
		if config.FilterParentLimit > 0 {
			filterParentLimit = int(config.FilterParentLimit)
		}
	}

	var (
		candidateParents   []*standard.Peer
		candidateParentIDs []string
	)
	for _, candidateParent := range peer.Task.LoadRandomPeers(uint(filterParentLimit)) {
		// Candidate parent is in blocklist.
		if blocklist.Contains(candidateParent.ID) {
			peer.Log.Debugf("parent %s host %s is not selected because it is in blocklist", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Candidate parent is disable shared.
		if candidateParent.Host.DisableShared {
			peer.Log.Debugf("parent %s host %s is not selected because it is disable shared", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Candidate parent host is not allowed to be the same as the peer host,
		// because dfdaemon cannot handle the situation
		// where two tasks are downloading and downloading each other.
		if peer.Host.ID == candidateParent.Host.ID {
			peer.Log.Debugf("parent %s host %s is the same as peer host", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Candidate parent can not find in dag.
		inDegree, err := peer.Task.PeerInDegree(candidateParent.ID)
		if err != nil {
			peer.Log.Debugf("can not find parent %s host %s vertex in dag", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Parent can be parent of the peer:
		// Condition 1: Parent has parent.
		// Condition 2: Parent has been back-to-source.
		// Condition 3: Parent has been succeeded.
		// Condition 4: Parent is seed peer.
		if candidateParent.Host.Type == types.HostTypeNormal && inDegree == 0 && !candidateParent.FSM.Is(standard.PeerStateBackToSource) &&
			!candidateParent.FSM.Is(standard.PeerStateSucceeded) {
			peer.Log.Debugf("parent %s host %s is not selected, because its download state is %d %d %s",
				candidateParent.ID, candidateParent.Host.ID, inDegree, int(candidateParent.Host.Type), candidateParent.FSM.Current())
			continue
		}

		// Candidate parent is bad parent.
		if s.evaluator.IsBadParent(candidateParent) {
			peer.Log.Debugf("parent %s host %s is not selected because it is bad node", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Candidate parent's free upload is empty.
		if candidateParent.Host.FreeUploadCount() <= 0 {
			peer.Log.Debugf("parent %s host %s is not selected because its free upload is empty, upload limit is %d, upload count is %d",
				candidateParent.ID, candidateParent.Host.ID, candidateParent.Host.ConcurrentUploadLimit.Load(), candidateParent.Host.ConcurrentUploadCount.Load())
			continue
		}

		// Candidate parent can add edge with peer.
		if !peer.Task.CanAddPeerEdge(candidateParent.ID, peer.ID) {
			peer.Log.Debugf("can not add edge with parent %s host %s", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		candidateParents = append(candidateParents, candidateParent)
		candidateParentIDs = append(candidateParentIDs, candidateParent.ID)
	}

	peer.Log.Infof("filter candidate parents is %#v", candidateParentIDs)
	return candidateParents
}

// FindReplicatePersistentCacheHosts finds replicate persistent cache hosts for the peer to replicate the task. It will compare the current
// persistent replica count with the persistent replica count and try to find enough parents. Then function will return the cached replicate parents,
// the replicate hosts without cache and found flag.
func (s *scheduling) FindReplicatePersistentCacheHosts(ctx context.Context, task *persistentcache.Task, blocklist set.SafeSet[string]) ([]*persistentcache.Peer, []*persistentcache.Host, bool) {
	currentPersistentReplicaCount, err := s.persistentCacheResource.TaskManager().LoadCurrentPersistentReplicaCount(ctx, task.ID)
	if err != nil {
		task.Log.Errorf("load current persistent replica count failed %s", err)
		return nil, nil, false
	}

	needPersistentReplicaCount := int(task.PersistentReplicaCount - currentPersistentReplicaCount)
	if needPersistentReplicaCount <= 0 {
		task.Log.Infof("persistent cache task %s has enough persistent replica count %d", task.ID, task.PersistentReplicaCount)
		return nil, nil, false
	}

	var (
		replicateHosts           []*persistentcache.Host
		replicateHostIDs         []string
		cachedReplicateParents   []*persistentcache.Peer
		cachedReplicateParentIDs []string
	)
	cachedParents := s.filterCachedReplicatePersistentCacheParents(ctx, task, blocklist)
	cachedParentsCount := len(cachedParents)

	// If the number of cached parents is greater than or equal to the number of persistent replica count,
	// return the cached parents directly and no need to find the replicate hosts without cache.
	if cachedParentsCount >= needPersistentReplicaCount {
		for _, cachedParent := range cachedParents[:needPersistentReplicaCount] {
			cachedReplicateParents = append(cachedReplicateParents, cachedParent)
			cachedReplicateParentIDs = append(cachedReplicateParentIDs, cachedParent.ID)
		}

		task.Log.Infof("find cached parents is %#v", cachedReplicateParentIDs)
		return cachedReplicateParents, nil, true
	}

	// If cached parents are not enough, append the replicate cached parents and find the replicate hosts without cache.
	if cachedParentsCount > 0 {
		for _, cachedParent := range cachedParents {
			cachedReplicateParents = append(cachedReplicateParents, cachedParent)
			cachedReplicateParentIDs = append(cachedReplicateParentIDs, cachedParent.ID)
			blocklist.Add(cachedParent.Host.ID)
		}
	}

	// Load all current persistent peers and add them to the blocklist to avoid scheduling the same host.
	currentPersistentPeers, err := s.persistentCacheResource.PeerManager().LoadPersistentAllByTaskID(ctx, task.ID)
	if err != nil {
		task.Log.Errorf("load all persistent cache peers failed: %s", err.Error())
		return nil, nil, false
	}

	for _, currentPersistentPeer := range currentPersistentPeers {
		blocklist.Add(currentPersistentPeer.Host.ID)
	}

	// Find the replicate hosts without cache. Calculate the number of persistent replicas needed without considering the cache.
	// Formula: Needed persistent replica count without cache = Total persistent replica count - Current persistent replica count - Cached parents count.
	needPersistentReplicaCount -= cachedParentsCount
	hosts := s.filterReplicatePersistentCacheHosts(ctx, task, needPersistentReplicaCount, blocklist)
	for _, host := range hosts {
		replicateHosts = append(replicateHosts, host)
		replicateHostIDs = append(replicateHostIDs, host.ID)
	}

	if len(cachedReplicateParents) == 0 && len(replicateHosts) == 0 {
		task.Log.Info("can not find replicate hosts")
		return nil, nil, false
	}

	task.Log.Infof("find cached parents is %#v and hosts is %#v", cachedReplicateParentIDs, replicateHostIDs)
	return cachedReplicateParents, replicateHosts, true
}

// FindCandidatePersistentCacheParents finds candidate persistent cache parents for the peer to download the task.
func (s *scheduling) FindCandidatePersistentCacheParents(ctx context.Context, peer *persistentcache.Peer, blocklist set.SafeSet[string]) ([]*persistentcache.Peer, bool) {
	// Find the candidate parent that can be scheduled.
	candidateParents := s.filterCandidatePersistentCacheParents(ctx, peer, blocklist)
	if len(candidateParents) == 0 {
		peer.Log.Info("can not find candidate persistent cache parents")
		return candidateParents, false
	}

	// Sort candidate parents by evaluation score.
	candidateParents = s.evaluator.EvaluatePersistentCacheParents(candidateParents, peer, peer.Task.TotalPieceCount)

	// Get the parents with candidateParentLimit.
	candidateParentLimit := config.DefaultSchedulerCandidateParentLimit
	if config, err := s.dynconfig.GetSchedulerClusterConfig(); err == nil {
		if config.CandidateParentLimit > 0 {
			candidateParentLimit = int(config.CandidateParentLimit)
		}
	}

	if len(candidateParents) > candidateParentLimit {
		candidateParents = candidateParents[:candidateParentLimit]
	}

	var parentIDs []string
	for _, candidateParent := range candidateParents {
		parentIDs = append(parentIDs, candidateParent.ID)
	}

	peer.Log.Infof("scheduling candidate persistent cache parents is %#v", parentIDs)
	return candidateParents, true
}

// filterCandidatePersistentCacheParents filters the candidate persistent cache parents that can be scheduled.
func (s *scheduling) filterCandidatePersistentCacheParents(ctx context.Context, peer *persistentcache.Peer, blocklist set.SafeSet[string]) []*persistentcache.Peer {
	parents, err := s.persistentCacheResource.PeerManager().LoadAllByTaskID(ctx, peer.Task.ID)
	if err != nil {
		peer.Log.Errorf("load all persistent cache parents failed: %s", err.Error())
		return nil
	}

	var (
		candidateParents   []*persistentcache.Peer
		candidateParentIDs []string
	)
	for _, candidateParent := range parents {
		// Candidate persistent cache parent is in blocklist.
		if blocklist.Contains(candidateParent.ID) {
			peer.Log.Debugf("persistent cache parent %s host %s is not selected because it is in blocklist", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Candidate persistent cache parent host is not allowed to be the same as the peer host,
		if peer.Host.ID == candidateParent.Host.ID {
			peer.Log.Debugf("persistent cache parent %s host %s is the same as peer host", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		// Candidate persistent cache parent is bad parent.
		if s.evaluator.IsBadPersistentCacheParent(candidateParent) {
			peer.Log.Debugf("persistent cache parent %s host %s is not selected because it is bad node", candidateParent.ID, candidateParent.Host.ID)
			continue
		}

		candidateParents = append(candidateParents, candidateParent)
		candidateParentIDs = append(candidateParentIDs, candidateParent.ID)
	}

	peer.Log.Infof("filter candidate persistent cache parents is %#v", candidateParentIDs)
	return candidateParents
}

// filterCachedReplicatePersistentCacheHosts filters the cached replicate persistent cache parents that can be scheduled.
func (s *scheduling) filterCachedReplicatePersistentCacheParents(ctx context.Context, task *persistentcache.Task, blocklist set.SafeSet[string]) []*persistentcache.Peer {
	parents, err := s.persistentCacheResource.PeerManager().LoadAllByTaskID(ctx, task.ID)
	if err != nil {
		task.Log.Errorf("load all persistent cache parents failed: %s", err.Error())
		return nil
	}

	var (
		replicateParents   []*persistentcache.Peer
		replicateParentIDs []string
	)
	for _, replicateParent := range parents {
		// Candidate persistent cache parent is in blocklist.
		if blocklist.Contains(replicateParent.ID) {
			task.Log.Debugf("persistent cache parent %s host %s is not selected because it is in blocklist", replicateParent.ID, replicateParent.Host.ID)
			continue
		}

		// If the parent is persistent, it cannot be selected.
		if replicateParent.Persistent {
			task.Log.Debugf("persistent cache parent %s host %s is not selected because it is persistent", replicateParent.ID, replicateParent.Host.ID)
			continue
		}

		// If the parent is not succeeded, it cannot be selected.
		if !replicateParent.FSM.Is(persistentcache.PeerStateSucceeded) {
			task.Log.Debugf("persistent cache parent %s host %s is not selected because its download state is %s", replicateParent.ID, replicateParent.Host.ID, replicateParent.FSM.Current())
			continue
		}

		// If the host is disable shared, it cannot be selected.
		if replicateParent.Host.DisableShared {
			task.Log.Debugf("persistent cache parent %s host %s is not selected because it is disable shared", replicateParent.ID, replicateParent.Host.ID)
			continue
		}

		replicateParents = append(replicateParents, replicateParent)
		replicateParentIDs = append(replicateParentIDs, replicateParent.ID)
	}

	task.Log.Infof("filter cached parents is %#v", replicateParentIDs)
	return replicateParents
}

// filterReplicatePersistentCacheHosts filters the replicate persistent cache hosts that can be scheduled.
func (s *scheduling) filterReplicatePersistentCacheHosts(ctx context.Context, task *persistentcache.Task, count int, blocklist set.SafeSet[string]) []*persistentcache.Host {
	hosts, err := s.persistentCacheResource.HostManager().LoadRandom(ctx, count, blocklist)
	if err != nil {
		task.Log.Errorf("load all persistent cache hosts failed: %s", err.Error())
		return nil
	}

	var (
		replicateHosts   []*persistentcache.Host
		replicateHostIDs []string
	)
	for _, host := range hosts {
		// If the host is disable shared, it cannot be selected.
		if host.DisableShared {
			task.Log.Debugf("persistent cache host %s is not selected because it is disable shared", host.ID)
			continue
		}

		// If the available disk space is not enough, it cannot be selected.
		if host.Disk.Free < task.ContentLength {
			task.Log.Debugf("persistent cache host %s is not selected because its free disk space is not enough, free disk is %d, content length is %d",
				host.ID, host.Disk.Free, task.ContentLength)
			continue
		}

		replicateHosts = append(replicateHosts, host)
		replicateHostIDs = append(replicateHostIDs, host.ID)
	}

	task.Log.Infof("filter hosts is %#v", replicateHostIDs)
	return replicateHosts
}

// constructSuccessNormalTaskResponse constructs scheduling successful response of the normal task.
// Used only in v2 version of the grpc.
func constructSuccessNormalTaskResponse(candidateParents []*standard.Peer) *schedulerv2.AnnouncePeerResponse_NormalTaskResponse {
	var parents []*commonv2.Peer
	for _, candidateParent := range candidateParents {
		parent := &commonv2.Peer{
			Id:               candidateParent.ID,
			Priority:         candidateParent.Priority,
			Cost:             durationpb.New(candidateParent.Cost.Load()),
			State:            candidateParent.FSM.Current(),
			NeedBackToSource: candidateParent.NeedBackToSource.Load(),
			CreatedAt:        timestamppb.New(candidateParent.CreatedAt.Load()),
			UpdatedAt:        timestamppb.New(candidateParent.UpdatedAt.Load()),
		}

		// Set range to parent.
		if candidateParent.Range != nil {
			parent.Range = &commonv2.Range{
				Start:  uint64(candidateParent.Range.Start),
				Length: uint64(candidateParent.Range.Length),
			}
		}

		// Set task to parent.
		parent.Task = &commonv2.Task{
			Id:                  candidateParent.Task.ID,
			Type:                candidateParent.Task.Type,
			Url:                 candidateParent.Task.URL,
			Tag:                 &candidateParent.Task.Tag,
			Application:         &candidateParent.Task.Application,
			FilteredQueryParams: candidateParent.Task.FilteredQueryParams,
			RequestHeader:       candidateParent.Task.Header,
			ContentLength:       uint64(candidateParent.Task.ContentLength.Load()),
			PieceCount:          uint32(candidateParent.Task.TotalPieceCount.Load()),
			SizeScope:           candidateParent.Task.SizeScope(),
			State:               candidateParent.Task.FSM.Current(),
			PeerCount:           uint32(candidateParent.Task.PeerCount()),
			CreatedAt:           timestamppb.New(candidateParent.Task.CreatedAt.Load()),
			UpdatedAt:           timestamppb.New(candidateParent.Task.UpdatedAt.Load()),
		}

		// Set digest to parent task.
		if candidateParent.Task.Digest != nil {
			dgst := candidateParent.Task.Digest.String()
			parent.Task.Digest = &dgst
		}

		// Set host to parent.
		parent.Host = &commonv2.Host{
			Id:              candidateParent.Host.ID,
			Type:            uint32(candidateParent.Host.Type),
			Hostname:        candidateParent.Host.Hostname,
			Ip:              candidateParent.Host.IP,
			Port:            candidateParent.Host.Port,
			DownloadPort:    candidateParent.Host.DownloadPort,
			Os:              candidateParent.Host.OS,
			Platform:        candidateParent.Host.Platform,
			PlatformFamily:  candidateParent.Host.PlatformFamily,
			PlatformVersion: candidateParent.Host.PlatformVersion,
			KernelVersion:   candidateParent.Host.KernelVersion,
			Cpu: &commonv2.CPU{
				LogicalCount:   candidateParent.Host.CPU.LogicalCount,
				PhysicalCount:  candidateParent.Host.CPU.PhysicalCount,
				Percent:        candidateParent.Host.CPU.Percent,
				ProcessPercent: candidateParent.Host.CPU.ProcessPercent,
				Times: &commonv2.CPUTimes{
					User:      candidateParent.Host.CPU.Times.User,
					System:    candidateParent.Host.CPU.Times.System,
					Idle:      candidateParent.Host.CPU.Times.Idle,
					Nice:      candidateParent.Host.CPU.Times.Nice,
					Iowait:    candidateParent.Host.CPU.Times.Iowait,
					Irq:       candidateParent.Host.CPU.Times.Irq,
					Softirq:   candidateParent.Host.CPU.Times.Softirq,
					Steal:     candidateParent.Host.CPU.Times.Steal,
					Guest:     candidateParent.Host.CPU.Times.Guest,
					GuestNice: candidateParent.Host.CPU.Times.GuestNice,
				},
			},
			Memory: &commonv2.Memory{
				Total:              candidateParent.Host.Memory.Total,
				Available:          candidateParent.Host.Memory.Available,
				Used:               candidateParent.Host.Memory.Used,
				UsedPercent:        candidateParent.Host.Memory.UsedPercent,
				ProcessUsedPercent: candidateParent.Host.Memory.ProcessUsedPercent,
				Free:               candidateParent.Host.Memory.Free,
			},
			Network: &commonv2.Network{
				TcpConnectionCount:       candidateParent.Host.Network.TCPConnectionCount,
				UploadTcpConnectionCount: candidateParent.Host.Network.UploadTCPConnectionCount,
				Location:                 &candidateParent.Host.Network.Location,
				Idc:                      &candidateParent.Host.Network.IDC,
				DownloadRate:             candidateParent.Host.Network.DownloadRate,
				DownloadRateLimit:        candidateParent.Host.Network.DownloadRateLimit,
				UploadRate:               candidateParent.Host.Network.UploadRate,
				UploadRateLimit:          candidateParent.Host.Network.UploadRateLimit,
			},
			Disk: &commonv2.Disk{
				Total:             candidateParent.Host.Disk.Total,
				Free:              candidateParent.Host.Disk.Free,
				Used:              candidateParent.Host.Disk.Used,
				UsedPercent:       candidateParent.Host.Disk.UsedPercent,
				InodesTotal:       candidateParent.Host.Disk.InodesTotal,
				InodesUsed:        candidateParent.Host.Disk.InodesUsed,
				InodesFree:        candidateParent.Host.Disk.InodesFree,
				InodesUsedPercent: candidateParent.Host.Disk.InodesUsedPercent,
				WriteBandwidth:    candidateParent.Host.Disk.WriteBandwidth,
				ReadBandwidth:     candidateParent.Host.Disk.ReadBandwidth,
			},
			Build: &commonv2.Build{
				GitVersion: candidateParent.Host.Build.GitVersion,
				GitCommit:  &candidateParent.Host.Build.GitCommit,
				GoVersion:  &candidateParent.Host.Build.GoVersion,
				Platform:   &candidateParent.Host.Build.Platform,
			},
		}

		parents = append(parents, parent)
	}

	return &schedulerv2.AnnouncePeerResponse_NormalTaskResponse{
		NormalTaskResponse: &schedulerv2.NormalTaskResponse{
			CandidateParents: parents,
		},
	}
}

// constructSuccessPeerPacket constructs peer successful packet.
// Used only in v1 version of the grpc.
func constructSuccessPeerPacket(peer *standard.Peer, parent *standard.Peer, candidateParents []*standard.Peer) *schedulerv1.PeerPacket {
	var parents []*schedulerv1.PeerPacket_DestPeer
	for _, candidateParent := range candidateParents {
		parents = append(parents, &schedulerv1.PeerPacket_DestPeer{
			Ip:      candidateParent.Host.IP,
			RpcPort: candidateParent.Host.Port,
			PeerId:  candidateParent.ID,
		})
	}

	return &schedulerv1.PeerPacket{
		TaskId: peer.Task.ID,
		SrcPid: peer.ID,
		MainPeer: &schedulerv1.PeerPacket_DestPeer{
			Ip:      parent.Host.IP,
			RpcPort: parent.Host.Port,
			PeerId:  parent.ID,
		},
		CandidatePeers: parents,
		Code:           commonv1.Code_Success,
	}
}
