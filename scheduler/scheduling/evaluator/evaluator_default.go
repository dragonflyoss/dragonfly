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
	// defaultLoadQualityWeight is the weight of load quality.
	defaultLoadQualityWeight = 0.6

	// defaultHostTypeWeight is the weight of host type.
	defaultIDCAffinityWeight = 0.2

	// defaultLocationAffinityWeight is the weight of location affinity.
	defaultLocationAffinityWeight = 0.1

	// defaultHostTypeWeight is the weight of host type.
	defaultHostTypeWeight = 0.1
)

const (
	// defaultPeakBandwidthUsageWeight is the weight of peak bandwidth usage.
	defaultPeakBandwidthUsageWeight = 0.5

	// defaultBandwidthDurationWeight is the weight of bandwidth duration.
	defaultBandwidthDurationWeight = 0.3

	// defaultConcurrencyWeight is the weight of concurrency.
	defaultConcurrencyWeight = 0.2
)

const (
	// defaultIDCAffinityWeightForPersistentCacheTask is the weight of host type for persistent cache task.
	defaultIDCAffinityWeightForPersistentCacheTask = 0.7

	// defaultLocationAffinityWeightForPersistentCacheTask is the weight of location affinity for persistent cache task.
	defaultLocationAffinityWeightForPersistentCacheTask = 0.3
)

// evaluatorDefault is the default evaluator implementation.
type evaluatorDefault struct {
	evaluator
}

// newEvaluatorDefault returns a new EvaluatorDefault.
func newEvaluatorDefault() Evaluator {
	return &evaluatorDefault{}
}

// EvaluateParents sorts and returns a list of parent peers ordered by their suitability as download sources.
// Parents are ranked from best to worst based on a comprehensive multi-dimensional evaluation.
//
// This function evaluates each parent peer using multiple metrics including load quality (bandwidth
// usage, sustained load, and concurrency), IDC affinity, location affinity, and host
// type. The parents are then sorted in descending order by their total scores, with the highest-scoring (most suitable)
// parents appearing first in the returned slice.
func (e *evaluatorDefault) EvaluateParents(parents []*standard.Peer, child *standard.Peer) []*standard.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluateParents(parents[i], child) > e.evaluateParents(parents[j], child)
		},
	)

	return parents
}

// evaluateParents evaluates and scores a parent peer for selection as a download source for a child peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates a better parent candidate.
//
// This function combines four key metrics to comprehensively evaluate parent peer quality:
// 1. Load Quality Score: Evaluates the parent's current bandwidth, sustained load, and concurrency.
// 2. IDC Affinity Score: Measures network proximity based on Internet Data Center(IDC) affinity.
// 3. Location Affinity Score: Measures geographic proximity based on location affinity.
// 4. Host Type Score: Evaluates whether the parent is a seed peer or normal peer(preferred).
//
// Formula: TotalScore = (IDCAffinityScore * 0.2) + (LocationAffinityScore * 0.1) + (LoadQualityScore * 0.6) + (HostTypeScore * 0.1)
func (e *evaluatorDefault) evaluateParents(parent *standard.Peer, child *standard.Peer) float64 {
	loadQualityScore := e.calculateLoadQualityScore(parent)
	idcAffinityScore := e.calculateIDCAffinityScore(parent.Host.Network.IDC, child.Host.Network.IDC)
	locationAffinityScore := e.calculateLocationAffinityScore(parent.Host.Network.Location, child.Host.Network.Location)
	hostTypeScore := e.calculateHostTypeScore(parent)
	parent.Log.Debugf("[evaluator] evaluate parent: loadQualityScore=%.4f, idcAffinityScore=%.4f, locationAffinityScore=%.4f, hostTypeScore=%.4f",
		loadQualityScore, idcAffinityScore, locationAffinityScore, hostTypeScore)

	return defaultLoadQualityWeight*loadQualityScore + defaultIDCAffinityWeight*idcAffinityScore +
		defaultLocationAffinityWeight*locationAffinityScore + defaultHostTypeWeight*hostTypeScore
}

// calculateLoadQualityScore calculates the overall load quality score for a parent peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates better availability and capacity.
//
// This function combines three key metrics to evaluate the parent peer's load quality:
// 1. Peak Bandwidth Usage: Measures instantaneous bandwidth occupancy.
// 2. Bandwidth Duration: Reflects the sustained impact of downloads over a duration window.
// 3. Concurrent Efficiency: Evaluates concurrent upload capacity and potential overhead.
//
// The score is calculated as a weighted sum of these three component scores, allowing for
// a comprehensive assessment that balances immediate capacity, sustained load, and concurrency.
//
// Formula: LoadQualityScore = (PeakBandwidthUsage * 0.5) + (BandwidthDurationRatio * 0.3) + (ConcurrentEfficiency * 0.2)
func (e *evaluatorDefault) calculateLoadQualityScore(parent *standard.Peer) float64 {
	peakBandwidthUsageScore := e.calculatePeakBandwidthUsageScore(parent)
	bandwidthDurationScore := e.calculateBandwidthDurationScore(parent)
	concurrencyScore := e.calculateConcurrencyScore(parent)
	parent.Log.Debugf("[evaluator] calculate load quality: peakBandwidthUsageScore=%.4f, bandwidthDurationScore=%.4f, concurrencyScore=%.4f",
		peakBandwidthUsageScore, bandwidthDurationScore, concurrencyScore)

	return defaultPeakBandwidthUsageWeight*peakBandwidthUsageScore + defaultBandwidthDurationWeight*bandwidthDurationScore + defaultConcurrencyWeight*concurrencyScore
}

// calculatePeakBandwidthUsageScore calculates the peak bandwidth usage score for a parent peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates better availability.
//
// The score represents the remaining bandwidth capacity of the parent peer. A score closer to 1.0
// indicates more available bandwidth, making the peer more suitable for serving additional requests.
//
// The function returns minScore (0.0) if:
// - The maximum transmit bandwidth is 0 (unconfigured or invalid).
// - The current transmit bandwidth has reached or exceeded the maximum capacity.
//
// The bandwidth usage ratio is weighted higher (0.5) in the P2P download speed evaluation
// because instantaneous bandwidth occupancy is highly sensitive to parent upload bandwidth peaks.
// Bandwidth Delay Product (BDP) is not used because Dragonfly is designed for intranet
// environments where RTT is typically less than 1ms.
//
// Formula: PeakBandwidthUsageScore = Max(0, 1 - (TxBandwidth / MaxTxBandwidth))
func (e *evaluatorDefault) calculatePeakBandwidthUsageScore(parent *standard.Peer) float64 {
	parent.Log.Debugf("[evaluator] calculate peak bandwidth usage: maxTxBandwidth=%d, txBandwidth=%d",
		parent.Host.Network.MaxTxBandwidth, parent.Host.TxBandwidth.Load())

	maxTxBandwidth := parent.Host.Network.MaxTxBandwidth
	if maxTxBandwidth == 0 {
		return minScore
	}

	txBandwidth := parent.Host.TxBandwidth.Load()
	if txBandwidth >= maxTxBandwidth {
		return minScore
	}

	bandwidthUsageRatio := float64(txBandwidth) / float64(maxTxBandwidth)
	if bandwidthUsageRatio >= 1.0 {
		return minScore
	}

	return 1.0 - bandwidthUsageRatio
}

// calculateBandwidthDurationScore calculates the bandwidth duration score for a parent peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates better availability.
//
// The score represents the remaining bandwidth capacity of the parent peer over a fixed duration window.
// A score closer to 1.0 indicates that the peer has more available bandwidth capacity to serve
// additional requests without saturation.
//
// The function returns minScore (0.0) if:
// - The maximum transmit bandwidth is 0 (unconfigured or invalid).
// - The bandwidth duration ratio has reached or exceeded 1.0 (fully saturated).
//
// The function returns maxScore (1.0) if:
// - No content has been uploaded yet (uploadContentLength is 0).
//
// The bandwidth duration ratio is calculated using a 60-second duration window, which is an empirical
// value chosen to smooth out short-term fluctuations. The ratio reflects the sustained impact of
// upload tasks on network bandwidth over this period. Task size affects the overall download duration,
// reflecting the sustained impact of different-sized tasks on the network.
//
// Formula: BandwidthDurationScore = Max(0, 1 - (UploadContentLength / MaxTxBandwidth / DurationWindow))
func (e *evaluatorDefault) calculateBandwidthDurationScore(parent *standard.Peer) float64 {
	parent.Log.Debugf("[evaluator] calculate bandwidth duration: maxTxBandwidth=%d, uploadContentLength=%d",
		parent.Host.Network.MaxTxBandwidth, parent.Host.UploadContentLength.Load())

	maxTxBandwidth := parent.Host.Network.MaxTxBandwidth
	if maxTxBandwidth == 0 {
		return minScore
	}

	uploadContentLength := parent.Host.UploadContentLength.Load()
	if uploadContentLength == 0 {
		return maxScore
	}

	// Duration window over which upload bandwidth utilization is evaluated, uint is seconds.
	const durationWindow = 60
	durationRatio := float64(parent.Host.UploadContentLength.Load()*8) / float64(maxTxBandwidth) / durationWindow
	if durationRatio >= 1.0 {
		return minScore
	}

	return 1.0 - durationRatio
}

// calculateConcurrencyScore calculates the concurrent efficiency score for a parent peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates better availability.
//
// The score represents the concurrent upload capacity of the parent peer. It evaluates whether
// the peer can handle the current number of concurrent piece uploads without exceeding its
// bandwidth capacity. When the number of concurrent pieces reaches the bandwidth limit, further
// increases can cause penalties, reflecting excessive concurrency overhead.
//
// The function returns minScore (0.0) if:
// - The maximum transmit bandwidth is 0 (unconfigured or invalid).
// - The pieceCountNeeded is less than or equal to 1.0 (insufficient bandwidth capacity).
//
// The function returns maxScore (1.0) if:
// - No pieces are currently being uploaded concurrently (concurrentUploadPieceCount is 0).
// - The concurrency ratio is less than or equal to 1.0 (within optimal capacity).
//
// The calculation assumes:
// - Piece length sets to 16 MiB is an empirical value for efficient I/O throughput.
// - PieceCountNeeded represents the number of pieces that can fully saturate the bandwidth.
//
// Formula:
// - PieceCountNeeded = MaxParentTransmitBandwidth / PieceLength
// - If ConcurrentUploadPieceCount / PieceCountNeeded <= 1: ConcurrentEfficiency = 1 (maxScore)
// - If ConcurrentUploadPieceCount / PieceCountNeeded > 1: ConcurrentEfficiency = PieceCountNeeded / ConcurrentUploadPieceCount
func (e *evaluatorDefault) calculateConcurrencyScore(parent *standard.Peer) float64 {
	parent.Log.Debugf("[evaluator] calculate concurrency: maxTxBandwidth=%d, concurrentUploadPieceCount=%d",
		parent.Host.Network.MaxTxBandwidth, parent.Host.ConcurrentUploadPieceCount.Load())

	maxTxBandwidth := parent.Host.Network.MaxTxBandwidth
	if maxTxBandwidth == 0 {
		return minScore
	}

	const (
		pieceLength = 16 * 1024 * 1024
	)
	pieceCountNeeded := float64(maxTxBandwidth) / pieceLength / 8
	if pieceCountNeeded <= 1.0 {
		return minScore
	}

	concurrentUploadPieceCount := parent.Host.ConcurrentUploadPieceCount.Load()
	if concurrentUploadPieceCount == 0 {
		return maxScore
	}

	concurrencyRatio := float64(concurrentUploadPieceCount) / pieceCountNeeded
	if concurrencyRatio <= 1.0 {
		return maxScore
	}

	return 1 / concurrencyRatio
}

// EvaluatePersistentCacheParents sorts and returns a list of parent peers ordered by their suitability as download sources.
// Parents are ranked from best to worst based on a comprehensive multi-dimensional evaluation.
//
// This function evaluates each parent peer using multiple metrics including IDC affinity and location affinity.
// The parents are then sorted in descending order by their total scores, with the highest-scoring (most suitable)
// parents appearing first in the returned slice.
func (e *evaluatorDefault) EvaluatePersistentCacheParents(parents []*persistentcache.Peer, child *persistentcache.Peer) []*persistentcache.Peer {
	sort.Slice(
		parents,
		func(i, j int) bool {
			return e.evaluatePersistentCacheParents(parents[i], child) > e.evaluatePersistentCacheParents(parents[j], child)
		},
	)

	return parents
}

// evaluatePersistentCacheParents evaluates and scores a parent peer for selection as a download source for a child peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates a better parent candidate.
//
// This function combines two key metrics to comprehensively evaluate parent peer quality:
// 1. IDC Affinity Score: Measures network proximity based on Internet Data Center(IDC) affinity.
// 2. Location Affinity Score: Measures geographic proximity based on location affinity.
//
// Formula: TotalScore = (IDCAffinityScore * 0.7) + (LocationAffinityScore * 0.3)
func (e *evaluatorDefault) evaluatePersistentCacheParents(parent *persistentcache.Peer, child *persistentcache.Peer) float64 {
	idcAffinityScore := e.calculateIDCAffinityScore(parent.Host.Network.IDC, child.Host.Network.IDC)
	locationAffinityScore := e.calculateLocationAffinityScore(parent.Host.Network.Location, child.Host.Network.Location)
	parent.Log.Debugf("[evaluator] evaluate persistent cache parent: idcAffinityScore=%.4f, locationAffinityScore=%.4f",
		idcAffinityScore, locationAffinityScore)

	return defaultIDCAffinityWeightForPersistentCacheTask*idcAffinityScore +
		defaultLocationAffinityWeightForPersistentCacheTask*locationAffinityScore
}

// calculateIDCAffinityScore calculates the IDC affinity score between two peers.
// The score ranges from 0.0 to 1.0, where a higher value indicates better network affinity.
//
// IDC affinity is particularly important in enterprise and cloud environments where intra-IDC
// transfers are significantly faster and more cost-effective than cross-IDC transfers.
func (e *evaluatorDefault) calculateIDCAffinityScore(dst, src string) float64 {
	if dst == "" || src == "" {
		return minScore
	}

	if strings.EqualFold(dst, src) {
		return maxScore
	}

	return minScore
}

// calculateLocationAffinityScore calculates the geographic location affinity score between two peers.
// The score ranges from 0.0 to 1.0, where a higher value indicates better geographic proximity.
//
// This function evaluates hierarchical location matching using a multi-element location string
// format separated by "|" (e.g., "country|province|city|zone|cluster"). The score is calculated
// based on how many location elements match from left to right, rewarding closer geographic proximity.
//
// Location hierarchy example:
// "China|Beijing|Haidian|Zone-A|Cluster-1" vs "China|Beijing|Chaoyang|Zone-B|Cluster-2"
// Matches: Country (China) ✓, Province (Beijing) ✓, City (Haidian vs Chaoyang) ✗
// Score: 2 / 5 = 0.4
func (e *evaluatorDefault) calculateLocationAffinityScore(dst, src string) float64 {
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

// calculateHostTypeScore calculates the host type score for a parent peer.
// The score ranges from 0.0 to 1.0, where a higher value indicates a more preferred peer type.
//
// For first-time task downloads, the scheduler prioritizes seed peers to bootstrap the P2P network
// with reliable sources. For subsequent downloads, the scheduler prefers normal peers to maximize
// resource utilization and reduce load on seed peers.
func (e *evaluatorDefault) calculateHostTypeScore(peer *standard.Peer) float64 {
	if peer.Host.Type != types.HostTypeNormal {
		if peer.FSM.Is(standard.PeerStateReceivedNormal) ||
			peer.FSM.Is(standard.PeerStateRunning) {
			return maxScore
		}

		return minScore
	}

	return maxScore * 0.5
}
