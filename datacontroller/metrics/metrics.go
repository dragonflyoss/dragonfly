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

// Package metrics registers the metrics of the data controller with the
// controller-runtime registry, next to the built-in reconciler metrics.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	ctrlmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"

	"d7y.io/dragonfly/v2/pkg/types"
)

var (
	// JobsCreated counts the manager jobs created, by job type and result.
	JobsCreated = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: types.MetricsNamespace,
		Subsystem: types.DataControllerName,
		Name:      "jobs_created_total",
		Help:      "Counter of manager jobs created by the data controller.",
	}, []string{"type", "result"})

	// OperationsCompleted counts the finished operations, by operation type and outcome.
	OperationsCompleted = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: types.MetricsNamespace,
		Subsystem: types.DataControllerName,
		Name:      "operations_completed_total",
		Help:      "Counter of dataset operations completed by the data controller.",
	}, []string{"type", "outcome"})
)

func init() {
	ctrlmetrics.Registry.MustRegister(JobsCreated, OperationsCompleted)
}
