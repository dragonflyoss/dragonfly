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
	"math/big"

	"github.com/montanaflynn/stats"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/scheduler/resource/persistent"
	"d7y.io/dragonfly/v2/scheduler/resource/persistentcache"
	"d7y.io/dragonfly/v2/scheduler/resource/standard"
)

const (
	// DefaultAlgorithm is a rule-based scheduling algorithm.
	DefaultAlgorithm = "default"

	// PluginAlgorithm is a scheduling algorithm based on plugin extension.
	PluginAlgorithm = "plugin"
)

const (
	// Maximum score.
	maxScore float64 = 1

	// Minimum score.
	minScore = 0
)

const (
	// Maximum number of elements.
	maxElementLen = 5

	// If the number of samples is greater than or equal to 2000,
	// it is close to the normal distribution.
	normalDistributionLen = 2000

	// When costs len is greater than or equal to 2,
	// the last cost can be compared and calculated.
	minAvailableCostLen = 2
)

// Evaluator is an interface that evaluates the parents.
type Evaluator interface {
	// EvaluateParents sorts and returns a list of parent peers ordered by their suitability as download sources.
	// Parents are ranked from best to worst based on a comprehensive multi-dimensional evaluation.
	EvaluateParents(parents []*standard.Peer, child *standard.Peer) []*standard.Peer

	// IsBadParent determines whether a peer is unsuitable to be selected as a parent
	// for downloading pieces.
	IsBadParent(peer *standard.Peer) bool

	// EvaluatePersistentParents sorts and returns a list of parent peers ordered by their suitability as download sources.
	// Parents are ranked from best to worst based on a comprehensive multi-dimensional evaluation.
	EvaluatePersistentParents(parents []*persistent.Peer, child *persistent.Peer) []*persistent.Peer

	// IsBadPersistentParent determines whether a peer is unsuitable to be selected as a persistent parent
	// for downloading pieces.
	IsBadPersistentParent(peer *persistent.Peer) bool

	// EvaluatePersistentCacheParents sorts and returns a list of parent peers ordered by their suitability as download sources.
	// Parents are ranked from best to worst based on a comprehensive multi-dimensional evaluation.
	EvaluatePersistentCacheParents(parents []*persistentcache.Peer, child *persistentcache.Peer) []*persistentcache.Peer

	// IsBadPersistentCacheParent determines whether a peer is unsuitable to be selected as a persistent cache parent
	// for downloading pieces.
	IsBadPersistentCacheParent(peer *persistentcache.Peer) bool
}

// evaluator is an implementation of Evaluator.
type evaluator struct{}

// New returns a new Evaluator.
func New(algorithm string, pluginDir string) Evaluator {
	switch algorithm {
	case PluginAlgorithm:
		if plugin, err := LoadPlugin(pluginDir); err == nil {
			return plugin
		}
	case DefaultAlgorithm:
		return newEvaluatorDefault()
	}

	return newEvaluatorDefault()
}

// IsBadParent determines whether a peer is unsuitable to be selected as a parent
// for downloading pieces. It evaluates peers based on their current state and historical
// download performance metrics.
//
// A peer is considered a bad parent if:
//  1. It is in an unsuitable state (Failed, Leave, Pending, or any Received* state).
//  2. Its recent download costs indicate poor performance based on statistical analysis.
//
// For performance evaluation, the function uses two strategies:
//   - If sample size is small (< normalDistributionLen): Uses a simple threshold check
//     where the last cost must not exceed 20 times the mean of previous costs.
//   - If sample size is sufficient (>= normalDistributionLen): Applies the three-sigma rule
//     where costs beyond mean + 3*standard_deviation are considered bad.
func (e *evaluator) IsBadParent(peer *standard.Peer) bool {
	if peer.FSM.Is(standard.PeerStateFailed) || peer.FSM.Is(standard.PeerStateLeave) || peer.FSM.Is(standard.PeerStatePending) ||
		peer.FSM.Is(standard.PeerStateReceivedTiny) || peer.FSM.Is(standard.PeerStateReceivedSmall) ||
		peer.FSM.Is(standard.PeerStateReceivedNormal) || peer.FSM.Is(standard.PeerStateReceivedEmpty) {
		peer.Log.Debugf("peer is bad node because peer status is %s", peer.FSM.Current())
		return true
	}

	// Determine whether to bad node based on piece download costs.
	costs := stats.LoadRawData(peer.PieceCosts())
	len := len(costs)
	// Peer has not finished downloading enough piece.
	if len < minAvailableCostLen {
		logger.Debugf("peer %s has not finished downloading enough piece, it can't be bad node", peer.ID)
		return false
	}

	lastCost := costs[len-1]
	mean, _ := stats.Mean(costs[:len-1]) // nolint: errcheck

	// Download costs does not meet the normal distribution,
	// if the last cost is twenty times more than mean, it is bad node.
	if len < normalDistributionLen {
		isBadParent := big.NewFloat(lastCost).Cmp(big.NewFloat(mean*20)) > 0
		logger.Debugf("peer %s mean is %.2f and it is bad node: %t", peer.ID, mean, isBadParent)
		return isBadParent
	}

	// Download costs satisfies the normal distribution,
	// last cost falling outside of three-sigma effect need to be adjusted parent,
	// refer to https://en.wikipedia.org/wiki/68%E2%80%9395%E2%80%9399.7_rule.
	stdev, _ := stats.StandardDeviation(costs[:len-1]) // nolint: errcheck
	isBadParent := big.NewFloat(lastCost).Cmp(big.NewFloat(mean+3*stdev)) > 0
	logger.Debugf("peer %s meet the normal distribution, costs mean is %.2f and standard deviation is %.2f, peer is bad node: %t",
		peer.ID, mean, stdev, isBadParent)
	return isBadParent
}

// IsBadPersistentParent determines whether a persistent peer is unsuitable to be selected as a parent
// for replication. It evaluates the peer based on its current state.
//
// A persistent peer is considered a bad parent if it is in an unsuitable state:
// - Pending: waiting to start.
// - Uploading: currently uploading data.
// - ReceivedEmpty: received an empty response.
// - ReceivedNormal: received normal data but not yet completed.
// - Failed: has failed.
func (e *evaluator) IsBadPersistentParent(peer *persistent.Peer) bool {
	if peer.FSM.Is(persistent.PeerStatePending) || peer.FSM.Is(persistent.PeerStateUploading) || peer.FSM.Is(persistent.PeerStateReceivedEmpty) ||
		peer.FSM.Is(persistent.PeerStateReceivedNormal) || peer.FSM.Is(persistent.PeerStateFailed) {
		peer.Log.Debugf("persistent peer is bad node because peer status is %s", peer.FSM.Current())
		return true
	}

	return false
}

// IsBadPersistentCacheParent determines whether a persistent cache peer is unsuitable to be selected as a parent
// for replication. It evaluates the peer based on its current state.
//
// A persistent cache peer is considered a bad parent if it is in an unsuitable state:
// - Pending: waiting to start.
// - Uploading: currently uploading data.
// - ReceivedEmpty: received an empty response.
// - ReceivedNormal: received normal data but not yet completed.
// - Failed: has failed.
func (e *evaluator) IsBadPersistentCacheParent(peer *persistentcache.Peer) bool {
	if peer.FSM.Is(persistentcache.PeerStatePending) || peer.FSM.Is(persistentcache.PeerStateUploading) || peer.FSM.Is(persistentcache.PeerStateReceivedEmpty) ||
		peer.FSM.Is(persistentcache.PeerStateReceivedNormal) || peer.FSM.Is(persistentcache.PeerStateFailed) {
		peer.Log.Debugf("persistent cache peer is bad node because peer status is %s", peer.FSM.Current())
		return true
	}

	return false
}
