/*
 *     Copyright 2025 The Dragonfly Authors
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
	// defaultHostTypeWeight is the weight of host type.
	defaultIDCAffinityWeight = 0.2

	// defaultLocationAffinityWeight is the weight of location affinity.
	defaultLocationAffinityWeight = 0.1

	// defaultHostTypeWeight is the weight of host type.
	defaultHostTypeWeight = 0.1
)

// evaluatorDefault is an implementation of Evaluator.
type evaluatorDefault struct {
	evaluator
}

// newEvaluatorDefault returns a new EvaluatorDefault.
func newEvaluatorDefault() Evaluator {
	return &evaluatorDefault{}
}

// EvaluateParents sort parents by evaluating multiple feature scores.
func (e *evaluatorDefault) EvaluateParents(parents []*standard.Peer, child *standard.Peer, totalPieceCount uint32) []*standard.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluateParents(parents[i], child, totalPieceCount) > e.evaluateParents(parents[j], child, totalPieceCount)
		},
	)

	return parents
}

// evaluateParents sort parents by evaluating multiple feature scores.
func (e *evaluatorDefault) evaluateParents(parent *standard.Peer, child *standard.Peer, totalPieceCount uint32) float64 {
	parentLocation := parent.Host.Network.Location
	parentIDC := parent.Host.Network.IDC
	childLocation := child.Host.Network.Location
	childIDC := child.Host.Network.IDC

	return defaultHostTypeWeight*e.calculateHostTypeScore(parent) +
		defaultIDCAffinityWeight*e.calculateIDCAffinityScore(parentIDC, childIDC) +
		defaultLocationAffinityWeight*e.calculateMultiElementAffinityScore(parentLocation, childLocation)
}

// EvaluatePersistentCacheParents sort persistent cache parents by evaluating multiple feature scores.
func (e *evaluatorDefault) EvaluatePersistentCacheParents(parents []*persistentcache.Peer, child *persistentcache.Peer, totalPieceCount uint32) []*persistentcache.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluatePersistentCacheParents(parents[i], child, totalPieceCount) > e.evaluatePersistentCacheParents(parents[j], child, totalPieceCount)
		},
	)

	return parents
}

// evaluatePersistentCacheParents sort persistent cache parents by evaluating multiple feature scores.
func (e *evaluatorDefault) evaluatePersistentCacheParents(parent *persistentcache.Peer, child *persistentcache.Peer, totalPieceCount uint32) float64 {
	parentLocation := parent.Host.Network.Location
	parentIDC := parent.Host.Network.IDC
	childLocation := child.Host.Network.Location
	childIDC := child.Host.Network.IDC

	return defaultIDCAffinityWeight*e.calculateIDCAffinityScore(parentIDC, childIDC) +
		defaultLocationAffinityWeight*e.calculateMultiElementAffinityScore(parentLocation, childLocation)
}

// calculateHostTypeScore 0.0~1.0 larger and better.
func (e *evaluatorDefault) calculateHostTypeScore(peer *standard.Peer) float64 {
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
func (e *evaluatorDefault) calculateIDCAffinityScore(dst, src string) float64 {
	if dst == "" || src == "" {
		return minScore
	}

	if strings.EqualFold(dst, src) {
		return maxScore
	}

	return minScore
}

// calculateMultiElementAffinityScore 0.0~1.0 larger and better.
func (e *evaluatorDefault) calculateMultiElementAffinityScore(dst, src string) float64 {
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
	elementLen = min(len(dstElements), len(srcElements))

	// Maximum element length is 5.
	elementLen = min(elementLen, maxElementLen)

	for i := range elementLen {
		if !strings.EqualFold(dstElements[i], srcElements[i]) {
			break
		}

		score++
	}

	return float64(score) / float64(maxElementLen)
}
