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

package dynamic

import (
	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/core/evaluator"
	"d7y.io/dragonfly/v2/scheduler/supervisor"
)

type dynamicEvaluator struct {
	cfg *config.SchedulerConfig
}

var args map[string]float64

func init() {
	args = make(map[string]float64)
	args["load"] = 0.5
	args["dist"] = 0.5
}

func NewEvaluator(cfg *config.SchedulerConfig) evaluator.Evaluator {
	eval := &dynamicEvaluator{cfg: cfg}
	logger.Debugf("create dynamic evaluator successfully")
	return eval
}

func (eval *dynamicEvaluator) NeedAdjustParent(peer *supervisor.Peer) bool {
	if peer.Host.CDN {
		return false
	}

	if peer.GetParent() == nil && !peer.IsDone() {
		logger.Debugf("peer %s need adjust parent because it has not parent and status is %s", peer.PeerID, peer.GetStatus())
		return true
	}
	// TODO Check whether the parent node is in the blacklist
	if peer.GetParent() != nil && eval.IsBadNode(peer.GetParent()) {
		logger.Debugf("peer %s need adjust parent because it current parent is bad", peer.PeerID)
		return true
	}

	if peer.GetParent() != nil && peer.GetParent().IsLeave() {
		logger.Debugf("peer %s need adjust parent because it current parent is status is leave", peer.PeerID)
		return true
	}

	costHistory := peer.GetCostHistory()
	if len(costHistory) < 4 {
		return false
	}

	avgCost, lastCost := getAvgAndLastCost(costHistory, 4)
	// TODO adjust policy
	if (avgCost * 20) < lastCost {
		logger.Debugf("peer %s need adjust parent because it latest download cost is too time consuming", peer.PeerID)
		return true
	}
	return false
}

func (eval *dynamicEvaluator) IsBadNode(peer *supervisor.Peer) bool {
	if peer.IsBad() {
		logger.Debugf("peer %s is bad because status is %s", peer.PeerID, peer.GetStatus())
		return true
	}
	costHistory := peer.GetCostHistory()
	if len(costHistory) < 4 {
		return false
	}

	avgCost, lastCost := getAvgAndLastCost(costHistory, 4)

	if avgCost*40 < lastCost && !peer.Host.CDN {
		logger.Debugf("peer %s is bad because recent pieces have taken too long to download avg[%d] last[%d]", peer.PeerID, avgCost, lastCost)
		return true
	}

	return false
}

// Evaluate The bigger, the better
func (eval *dynamicEvaluator) Evaluate(dst *supervisor.Peer, src *supervisor.Peer) float64 {
	profits := getProfits(dst, src)

	load := getHostLoad(dst.Host)

	dist := getDistance(dst, src)

	score := profits + load*args["load"] + dist*args["dist"]
	return score
}

func getAvgAndLastCost(list []int, splitPos int) (avgCost, lastCost int) {
	length := len(list)
	totalCost := 0
	for i, cost := range list {
		totalCost += cost
		if length-i <= splitPos {
			lastCost += cost
		}
	}
	avgCost = totalCost / length
	lastCost = lastCost / splitPos
	return
}

// getProfits 0.0~unlimited larger and better
func getProfits(dst *supervisor.Peer, src *supervisor.Peer) float64 {
	diff := supervisor.GetDiffPieceNum(dst, src)
	if diff > 0 {
		return float64(dst.Task.PieceTotal + 1 - diff)
	} else {
		return float64(diff)
	}
}

// getHostLoad 0.0~1.0 larger and better
func getHostLoad(host *supervisor.PeerHost) float64 {
	return 1.0 - host.GetUploadLoadPercent()
}

// getDistance 0.0~1.0 larger and better
func getDistance(dst *supervisor.Peer, src *supervisor.Peer) float64 {
	hostDist := 40.0
	if dst.Host == src.Host {
		hostDist = 0.0
	} else {
		if src.Host.NetTopology != "" && dst.Host.NetTopology == src.Host.NetTopology {
			hostDist = 10.0
		} else if src.Host.IDC != "" && dst.Host.IDC == src.Host.IDC {
			hostDist = 20.0
		} else if dst.Host.SecurityDomain != src.Host.SecurityDomain {
			hostDist = 80.0
		}
	}
	return 1.0 - hostDist/80.0
}
