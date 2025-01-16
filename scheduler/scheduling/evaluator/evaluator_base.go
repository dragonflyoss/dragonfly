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
	"sort"
	"strings"

	"d7y.io/dragonfly/v2/pkg/math"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
)

const (
	// Finished piece weight.
	finishedPieceWeight float64 = 0.2

	// Parent's host upload success weight.
	parentHostUploadSuccessWeight = 0.2

	// Free upload weight.
	freeUploadWeight = 0.15

	// Host type weight.
	hostTypeWeight = 0.15

	// IDC affinity weight.
	idcAffinityWeight = 0.15

	// Location affinity weight.
	locationAffinityWeight = 0.15
)

// evaluatorBase is an implementation of Evaluator.
type evaluatorBase struct {
	evaluator
}

// NewEvaluatorBase returns a new EvaluatorBase.
func newEvaluatorBase() Evaluator {
	return &evaluatorBase{}
}

// EvaluateParents sort parents by evaluating multiple feature scores.
func (e *evaluatorBase) EvaluateParents(parents []*standard.Peer, child *standard.Peer, totalPieceCount int32) []*standard.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluateParents(parents[i], child, totalPieceCount) > e.evaluateParents(parents[j], child, totalPieceCount)
		},
	)

	return parents
}

// evaluateParents sort parents by evaluating multiple feature scores.
func (e *evaluatorBase) evaluateParents(parent *standard.Peer, child *standard.Peer, totalPieceCount int32) float64 {
	parentLocation := parent.Host.Network.Location
	parentIDC := parent.Host.Network.IDC
	childLocation := child.Host.Network.Location
	childIDC := child.Host.Network.IDC

	return finishedPieceWeight*e.calculatePieceScore(parent.FinishedPieces.Count(), child.FinishedPieces.Count(), totalPieceCount) +
		parentHostUploadSuccessWeight*e.calculateParentHostUploadSuccessScore(parent.Host.UploadCount.Load(), parent.Host.UploadFailedCount.Load()) +
		freeUploadWeight*e.calculateFreeUploadScore(parent.Host) +
		hostTypeWeight*e.calculateHostTypeScore(parent) +
		idcAffinityWeight*e.calculateIDCAffinityScore(parentIDC, childIDC) +
		locationAffinityWeight*e.calculateMultiElementAffinityScore(parentLocation, childLocation)
}

// EvaluateParents sort parents by evaluating multiple feature scores.
func (e *evaluatorBase) EvaluatePersistentCacheParents(parents []*persistentcache.Peer, child *persistentcache.Peer, totalPieceCount int32) []*persistentcache.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluatePersistentCacheParents(parents[i], child, totalPieceCount) > e.evaluatePersistentCacheParents(parents[j], child, totalPieceCount)
		},
	)

	return parents
}

// evaluatePersistentCacheParents sort persistent cache parents by evaluating multiple feature scores.
func (e *evaluatorBase) evaluatePersistentCacheParents(parent *persistentcache.Peer, child *persistentcache.Peer, totalPieceCount int32) float64 {
	parentLocation := parent.Host.Network.Location
	parentIDC := parent.Host.Network.IDC
	childLocation := child.Host.Network.Location
	childIDC := child.Host.Network.IDC

	return finishedPieceWeight*e.calculatePieceScore(parent.FinishedPieces.Count(), child.FinishedPieces.Count(), totalPieceCount) +
		idcAffinityWeight*e.calculateIDCAffinityScore(parentIDC, childIDC) +
		locationAffinityWeight*e.calculateMultiElementAffinityScore(parentLocation, childLocation)
}

// calculatePieceScore 0.0~unlimited larger and better.
func (e *evaluatorBase) calculatePieceScore(parentFinishedPieceCount uint, childFinishedPieceCount uint, totalPieceCount int32) float64 {
	// If the total piece is determined, normalize the number of
	// pieces downloaded by the parent node.
	if totalPieceCount > 0 {
		return float64(parentFinishedPieceCount) / float64(totalPieceCount)
	}

	// Use the difference between the parent node and the child node to
	// download the piece to roughly represent the piece score.
	return float64(parentFinishedPieceCount) - float64(childFinishedPieceCount)
}

// calculateParentHostUploadSuccessScore 0.0~unlimited larger and better.
func (e *evaluatorBase) calculateParentHostUploadSuccessScore(parentUploadCount int64, parentUploadFailedCount int64) float64 {
	if parentUploadCount < parentUploadFailedCount {
		return minScore
	}

	// Host has not been scheduled, then it is scheduled first.
	if parentUploadCount == 0 && parentUploadFailedCount == 0 {
		return maxScore
	}

	return float64(parentUploadCount-parentUploadFailedCount) / float64(parentUploadCount)
}

// calculateFreeUploadScore 0.0~1.0 larger and better.
func (e *evaluatorBase) calculateFreeUploadScore(host *standard.Host) float64 {
	ConcurrentUploadLimit := host.ConcurrentUploadLimit.Load()
	freeUploadCount := host.FreeUploadCount()
	if ConcurrentUploadLimit > 0 && freeUploadCount > 0 {
		return float64(freeUploadCount) / float64(ConcurrentUploadLimit)
	}

	return minScore
}

// calculateHostTypeScore 0.0~1.0 larger and better.
func (e *evaluatorBase) calculateHostTypeScore(peer *standard.Peer) float64 {
	// When the task is downloaded for the first time,
	// peer will be scheduled to seed peer first,
	// otherwise it will be scheduled to dfdaemon first.
	if peer.Host.Type != types.HostTypeNormal {
		if peer.FSM.Is(standard.PeerStateReceivedNormal) ||
			peer.FSM.Is(standard.PeerStateRunning) {
			return maxScore
		}

		return minScore
	}

	return maxScore * 0.5
}

// calculateIDCAffinityScore 0.0~1.0 larger and better.
func (e *evaluatorBase) calculateIDCAffinityScore(dst, src string) float64 {
	if dst == "" || src == "" {
		return minScore
	}

	if strings.EqualFold(dst, src) {
		return maxScore
	}

	return minScore
}

// calculateMultiElementAffinityScore 0.0~1.0 larger and better.
func (e *evaluatorBase) calculateMultiElementAffinityScore(dst, src string) float64 {
	if dst == "" || src == "" {
		return minScore
	}

	if strings.EqualFold(dst, src) {
		return maxScore
	}

	// Calculate the number of multi-element matches divided by "|".
	var score, elementLen int
	dstElements := strings.Split(dst, types.AffinitySeparator)
	srcElements := strings.Split(src, types.AffinitySeparator)
	elementLen = math.Min(len(dstElements), len(srcElements))

	// Maximum element length is 5.
	if elementLen > maxElementLen {
		elementLen = maxElementLen
	}

	for i := 0; i < elementLen; i++ {
		if !strings.EqualFold(dstElements[i], srcElements[i]) {
			break
		}

		score++
	}

	return float64(score) / float64(maxElementLen)
}
