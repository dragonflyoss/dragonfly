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

// Weights for persistent cache parent evaluation
type scoreWeights struct {
	FinishedPiece           float64
	ParentHostUploadSuccess float64
	FreeUpload              float64
	HostType                float64
	IDCAffinity             float64
	LocationAffinity        float64
}

// Default weights
var defaultScoreWeights = scoreWeights{
	FinishedPiece:           0.2,
	ParentHostUploadSuccess: 0.2,
	FreeUpload:              0.15,
	HostType:                0.15,
	IDCAffinity:             0.15,
	LocationAffinity:        0.15,
}

// PersistentCache Weights
type persistentCacheScoreWeights struct {
	FinishedPiece    float64
	IDCAffinity      float64
	LocationAffinity float64
}

// Default PersistentCache weights
var defaultPersistentCacheScoreWeights = persistentCacheScoreWeights{
	FinishedPiece:    0.2,
	IDCAffinity:      0.4,
	LocationAffinity: 0.4,
}

// EvaluateParents sorts parents by evaluating multiple feature scores.
func (e *evaluatorBase) EvaluateParents(parents []*standard.Peer, child *standard.Peer, totalPieceCount uint32) []*standard.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluateParents(parents[i], child, totalPieceCount) > e.evaluateParents(parents[j], child, totalPieceCount)
		},
	)

	return parents
}

// evaluateParents calculates the total score for a parent node, step by step for readability.
func (e *evaluatorBase) evaluateParents(parent *standard.Peer, child *standard.Peer, totalPieceCount uint32) float64 {
	w := defaultScoreWeights

	pieceScore := e.calculatePieceScore(parent.FinishedPieces.Count(), child.FinishedPieces.Count(), totalPieceCount)
	uploadSuccessScore := e.calculateParentHostUploadSuccessScore(parent.Host.UploadCount.Load(), parent.Host.UploadFailedCount.Load())
	freeUploadScore := e.calculateFreeUploadScore(parent.Host)
	hostTypeScore := e.calculateHostTypeScore(parent)
	idcAffinityScore := e.calculateIDCAffinityScore(parent.Host.Network.IDC, child.Host.Network.IDC)
	locationAffinityScore := e.calculateMultiElementAffinityScore(parent.Host.Network.Location, child.Host.Network.Location)

	totalScore := w.FinishedPiece*pieceScore +
		w.ParentHostUploadSuccess*uploadSuccessScore +
		w.FreeUpload*freeUploadScore +
		w.HostType*hostTypeScore +
		w.IDCAffinity*idcAffinityScore +
		w.LocationAffinity*locationAffinityScore

	return totalScore
}

// EvaluatePersistentCacheParents sorts persistent cache parents by evaluating multiple feature scores.
func (e *evaluatorBase) EvaluatePersistentCacheParents(parents []*persistentcache.Peer, child *persistentcache.Peer, totalPieceCount uint32) []*persistentcache.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluatePersistentCacheParents(parents[i], child, totalPieceCount) > e.evaluatePersistentCacheParents(parents[j], child, totalPieceCount)
		},
	)

	return parents
}

// evaluatePersistentCacheParents calculates the total score for a persistent cache parent, step by step for readability.
func (e *evaluatorBase) evaluatePersistentCacheParents(parent *persistentcache.Peer, child *persistentcache.Peer, totalPieceCount uint32) float64 {
	w := defaultPersistentCacheScoreWeights

	pieceScore := e.calculatePieceScore(parent.FinishedPieces.Count(), child.FinishedPieces.Count(), totalPieceCount)
	idcAffinityScore := e.calculateIDCAffinityScore(parent.Host.Network.IDC, child.Host.Network.IDC)
	locationAffinityScore := e.calculateMultiElementAffinityScore(parent.Host.Network.Location, child.Host.Network.Location)

	totalScore := w.FinishedPiece*pieceScore +
		w.IDCAffinity*idcAffinityScore +
		w.LocationAffinity*locationAffinityScore

	return totalScore
}

// calculatePieceScore 0.0~unlimited larger and better.
func (e *evaluatorBase) calculatePieceScore(parentFinishedPieceCount uint, childFinishedPieceCount uint, totalPieceCount uint32) float64 {
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

// calculateMultiElementAffinityScore returns a score between 0.0 and 1.0, higher is better.
func (e *evaluatorBase) calculateMultiElementAffinityScore(dst, src string) float64 {
	if dst == "" || src == "" {
		return minScore
	}

	if strings.EqualFold(dst, src) {
		return maxScore
	}

	// Calculate the number of multi-level elements (separated by |) that match.
	score := 0
	dstElements := strings.Split(dst, types.AffinitySeparator)
	srcElements := strings.Split(src, types.AffinitySeparator)
	elementLen := min(len(dstElements), len(srcElements))
	// Maximum element length is 5.
	elementLen = min(elementLen, maxElementLen)

	for i := 0; i < elementLen; i++ {
		if !strings.EqualFold(dstElements[i], srcElements[i]) {
			break
		}
		score++
	}

	return float64(score) / float64(maxElementLen)
}
