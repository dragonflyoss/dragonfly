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

package util

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

const downloadTrafficMetric = "dragonfly_client_download_traffic"

// DownloadTraffic counts successfully downloaded piece bytes by their source.
// It does not include metadata probes or failed source-download attempts.
type DownloadTraffic struct {
	RemotePeer   uint64
	BackToSource uint64
}

// DownloadTrafficSnapshot associates counters with the container that served them.
type DownloadTrafficSnapshot struct {
	Traffic      DownloadTraffic
	PodUID       string
	RestartCount int32
}

// CaptureDownloadTraffic reads counters without resetting them. The PodExec must
// identify a container explicitly so its restart count matches the metrics source.
func CaptureDownloadTraffic(pod *PodExec) (DownloadTrafficSnapshot, error) {
	if pod == nil || pod.namespace == "" || pod.name == "" || pod.container == "" {
		return DownloadTrafficSnapshot{}, fmt.Errorf("download traffic requires a pod namespace, name and container")
	}

	// Bound the host kubectl processes too: curl's timeout only applies after
	// kubectl has connected to the pod and started the remote command.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	uid, restartCount, err := downloadTrafficPodState(ctx, pod)
	if err != nil {
		return DownloadTrafficSnapshot{}, err
	}

	out, err := pod.CommandContext(ctx, "/usr/bin/curl", "--fail", "--silent", "--show-error",
		"--connect-timeout", "5", "--max-time", "15", "http://127.0.0.1:4002/metrics").Output()
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return DownloadTrafficSnapshot{}, fmt.Errorf("fetch download traffic from pod %s/%s: %w", pod.namespace, pod.name, err)
	}

	traffic, err := parseDownloadTraffic(string(out))
	if err != nil {
		return DownloadTrafficSnapshot{}, fmt.Errorf("parse download traffic from pod %s/%s: %w", pod.namespace, pod.name, err)
	}

	// A replacement or restart during the scrape must not produce a valid snapshot.
	afterUID, afterRestartCount, err := downloadTrafficPodState(ctx, pod)
	if err != nil {
		return DownloadTrafficSnapshot{}, err
	}
	if uid != afterUID || restartCount != afterRestartCount {
		return DownloadTrafficSnapshot{}, fmt.Errorf("pod %s/%s or container %s restarted while capturing download traffic", pod.namespace, pod.name, pod.container)
	}

	return DownloadTrafficSnapshot{Traffic: traffic, PodUID: uid, RestartCount: restartCount}, nil
}

// DownloadTrafficDelta rejects snapshots whose counters cannot be compared safely.
func DownloadTrafficDelta(before, after DownloadTrafficSnapshot) (DownloadTraffic, error) {
	if before.PodUID == "" || after.PodUID == "" {
		return DownloadTraffic{}, fmt.Errorf("download traffic snapshot is missing the pod UID")
	}
	if before.PodUID != after.PodUID {
		return DownloadTraffic{}, fmt.Errorf("pod changed while measuring download traffic: %s -> %s", before.PodUID, after.PodUID)
	}
	if before.RestartCount < 0 || after.RestartCount < 0 || before.RestartCount != after.RestartCount {
		return DownloadTraffic{}, fmt.Errorf("container restart count changed or is invalid: %d -> %d", before.RestartCount, after.RestartCount)
	}
	if after.Traffic.RemotePeer < before.Traffic.RemotePeer || after.Traffic.BackToSource < before.Traffic.BackToSource {
		return DownloadTraffic{}, fmt.Errorf("download traffic counter decreased: %+v -> %+v", before.Traffic, after.Traffic)
	}

	return DownloadTraffic{
		RemotePeer:   after.Traffic.RemotePeer - before.Traffic.RemotePeer,
		BackToSource: after.Traffic.BackToSource - before.Traffic.BackToSource,
	}, nil
}

func downloadTrafficPodState(ctx context.Context, pod *PodExec) (string, int32, error) {
	out, err := KubeCtlCommandContext(ctx, "-n", pod.namespace, "get", "pod", pod.name, "-o", "json", "--request-timeout=15s").Output()
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return "", 0, fmt.Errorf("read download traffic pod %s/%s: %w", pod.namespace, pod.name, err)
	}

	var state struct {
		Metadata struct {
			UID string `json:"uid"`
		} `json:"metadata"`
		Status struct {
			ContainerStatuses []struct {
				Name         string `json:"name"`
				RestartCount *int32 `json:"restartCount"`
			} `json:"containerStatuses"`
		} `json:"status"`
	}
	if err := json.Unmarshal(out, &state); err != nil {
		return "", 0, fmt.Errorf("decode download traffic pod %s/%s: %w", pod.namespace, pod.name, err)
	}
	if state.Metadata.UID == "" {
		return "", 0, fmt.Errorf("pod %s/%s is missing its UID", pod.namespace, pod.name)
	}
	for _, container := range state.Status.ContainerStatuses {
		if container.Name != pod.container {
			continue
		}
		if container.RestartCount == nil || *container.RestartCount < 0 {
			return "", 0, fmt.Errorf("pod %s/%s container %s has no valid restart count", pod.namespace, pod.name, pod.container)
		}
		return state.Metadata.UID, *container.RestartCount, nil
	}

	return "", 0, fmt.Errorf("pod %s/%s has no status for container %s", pod.namespace, pod.name, pod.container)
}

func parseDownloadTraffic(input string) (DownloadTraffic, error) {
	parser := expfmt.NewTextParser(model.UTF8Validation)
	families, err := parser.TextToMetricFamilies(strings.NewReader(input))
	if err != nil {
		return DownloadTraffic{}, err
	}
	family, found := families[downloadTrafficMetric]
	if !found {
		// CounterVec series do not exist until their label values are first used.
		return DownloadTraffic{}, nil
	}
	if family.GetType() != dto.MetricType_COUNTER {
		return DownloadTraffic{}, fmt.Errorf("%s is not a counter", downloadTrafficMetric)
	}

	var traffic DownloadTraffic
	for _, metric := range family.GetMetric() {
		if metric.GetCounter() == nil {
			return DownloadTraffic{}, fmt.Errorf("%s sample is missing its counter", downloadTrafficMetric)
		}
		value := metric.GetCounter().GetValue()
		if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || math.Trunc(value) != value || value >= math.Exp2(64) {
			return DownloadTraffic{}, fmt.Errorf("%s has invalid byte count %v", downloadTrafficMetric, value)
		}

		var trafficType string
		for _, label := range metric.GetLabel() {
			if label.GetName() == "type" {
				trafficType = label.GetValue()
				break
			}
		}
		var total *uint64
		switch trafficType {
		case "REMOTE_PEER":
			total = &traffic.RemotePeer
		case "BACK_TO_SOURCE":
			total = &traffic.BackToSource
		case "":
			return DownloadTraffic{}, fmt.Errorf("%s sample is missing its type label", downloadTrafficMetric)
		default:
			continue
		}
		bytes := uint64(value)
		if bytes > ^uint64(0)-*total {
			return DownloadTraffic{}, fmt.Errorf("%s %s byte count overflows uint64", downloadTrafficMetric, trafficType)
		}
		*total += bytes
	}

	return traffic, nil
}
