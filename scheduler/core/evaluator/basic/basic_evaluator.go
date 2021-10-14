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

package basic

import (
	"math"
	"strings"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/pkg/util/mathutils"
	"d7y.io/dragonfly/v2/pkg/util/stringutils"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/core/evaluator"
	"d7y.io/dragonfly/v2/scheduler/supervisor"
)

type baseEvaluator struct {
	cfg *config.SchedulerConfig
}

func NewEvaluator(cfg *config.SchedulerConfig) evaluator.Evaluator {
	eval := &baseEvaluator{cfg: cfg}
	logger.Debugf("create basic evaluator successfully")
	return eval
}

func (eval *baseEvaluator) NeedAdjustParent(peer *supervisor.Peer) bool {
	if peer.Host.IsCDN {
		return false
	}

	parent, ok := peer.GetParent()
	if !ok && !peer.IsDone() {
		logger.Debugf("peer %s need adjust parent because it has not parent and status is %s", peer.ID, peer.GetStatus())
		return true
	}

	// TODO Check whether the parent node is in the blacklist
	if ok && eval.IsBadNode(parent) {
		logger.Debugf("peer %s need adjust parent because it current parent is bad", peer.ID)
		return true
	}

	if ok && parent.IsLeave() {
		logger.Debugf("peer %s need adjust parent because it current parent is status is leave", peer.ID)
		return true
	}

	costs := peer.GetPieceCosts()
	if len(costs) < 4 {
		return false
	}

	avgCost, lastCost := getAvgAndLastCost(costs, 4)
	// TODO adjust policy
	if (avgCost * 20) < lastCost {
		logger.Debugf("peer %s need adjust parent because it latest download cost is too time consuming", peer.ID)
		return true
	}

	return false
}

func (eval *baseEvaluator) IsBadNode(peer *supervisor.Peer) bool {
	if peer.IsBad() {
		logger.Debugf("peer %s is bad because it's status is %s", peer.ID, peer.GetStatus())
		return true
	}

	costs := peer.GetPieceCosts()
	if len(costs) < 4 {
		return false
	}

	avgCost, lastCost := getAvgAndLastCost(costs, 4)
	if avgCost*40 < lastCost && !peer.Host.IsCDN {
		logger.Debugf("peer %s is bad because recent pieces have taken too long to download avg[%d] last[%d]", peer.ID, avgCost, lastCost)
		return true
	}

	return false
}

// Evaluate The bigger, the better
func (eval *baseEvaluator) Evaluate(parent *supervisor.Peer, child *supervisor.Peer) float64 {
	profits := getProfits(parent, child)

	load := getHostLoad(parent.Host)

	dist := getDistance(parent, child)

	return profits * load * dist
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
	diff := dst.TotalPieceCount.Load() - src.TotalPieceCount.Load()
	depth := dst.GetTreeDepth()

	return float64(int(diff+1)*src.GetTreeNodeCount()) / float64(depth*depth)
}

// getHostLoad 0.0~1.0 larger and better
func getHostLoad(host *supervisor.Host) float64 {
	return 1.0 - host.GetUploadLoadPercent()
}

// getDistance 0.0~1.0 larger and better
func getDistance(dst *supervisor.Peer, src *supervisor.Peer) float64 {
	topScore := 80.0
	score := 40.0
	if dst.Host == src.Host {
		score = topScore
	} else {
		if stringutils.IsNotBlank(src.Host.SecurityDomain) && stringutils.IsNotBlank(dst.Host.SecurityDomain) && src.Host.SecurityDomain != dst.Host.
			SecurityDomain {
			score = 0.0
		} else {
			if stringutils.IsNotBlank(src.Host.IDC) && dst.Host.IDC == src.Host.IDC {
				maxScore := topScore * 0.95 // distinguish from the same host
				if src.Host.NetTopology != "" && src.Host.NetTopology == dst.Host.NetTopology {
					score = maxScore
				} else {
					srcNetTopologies := strings.Split(src.Host.NetTopology, "|")
					dstNetTopologies := strings.Split(dst.Host.NetTopology, "|")
					minLen := mathutils.MinInt(len(srcNetTopologies), len(dstNetTopologies))
					fragmentLength := (maxScore - score) / (math.Pow(2, float64(minLen)) - 1)
					for i := 0; i < minLen; i++ {
						if srcNetTopologies[i] != dstNetTopologies[i] {
							break
						}
						score = score + fragmentLength*math.Pow(2, float64(i))
					}
				}
			} else if stringutils.IsNotBlank(src.Host.Location) && stringutils.IsNotBlank(dst.Host.Location) {
				maxScore := topScore * 0.9 // distinguish from the same host and same idc
				if src.Host.Location == dst.Host.Location {
					score = maxScore
				} else {
					srcLocations := strings.Split(src.Host.Location, "|")
					dstLocations := strings.Split(dst.Host.Location, "|")
					minLen := mathutils.MinInt(len(srcLocations), len(dstLocations))
					fragmentLength := (maxScore - score) / (math.Pow(2, float64(minLen)) - 1)
					for i := 0; i < minLen; i++ {
						if srcLocations[i] != dstLocations[i] {
							break
						}
						score = score + fragmentLength*math.Pow(2, float64(i))
					}
				}
			}
		}
	}
	return score / topScore
}
