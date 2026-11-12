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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DataLifecyclePolicySpec is the desired state of a DataLifecyclePolicy.
type DataLifecyclePolicySpec struct {
	// DatasetSelector selects the Datasets in the same namespace the policy applies to.
	// An empty selector matches every Dataset in the namespace.
	// +optional
	DatasetSelector *metav1.LabelSelector `json:"datasetSelector,omitempty"`

	// Priority breaks ties when several policies match a Dataset; the highest priority wins,
	// then the lexically smallest name.
	// +optional
	Priority int32 `json:"priority,omitempty"`

	// Lifecycle is the lifecycle applied to the matching Datasets.
	Lifecycle LifecycleSpec `json:"lifecycle"`
}

// DataLifecyclePolicyStatus is the observed state of a DataLifecyclePolicy.
type DataLifecyclePolicyStatus struct {
	// ObservedGeneration is the generation last reconciled.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// MatchedDatasets is the number of Datasets currently governed by the policy.
	// +optional
	MatchedDatasets int32 `json:"matchedDatasets,omitempty"`

	// Conditions describe the state in detail.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// DataLifecyclePolicy applies a lifecycle to a set of Datasets selected by labels. Fields
// set directly on a Dataset take precedence over the policy.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=dlp,categories=dragonfly
// +kubebuilder:printcolumn:name="Priority",type=integer,JSONPath=`.spec.priority`
// +kubebuilder:printcolumn:name="TTL",type=string,JSONPath=`.spec.lifecycle.ttl`
// +kubebuilder:printcolumn:name="Refresh",type=string,JSONPath=`.spec.lifecycle.refreshInterval`
// +kubebuilder:printcolumn:name="Matched",type=integer,JSONPath=`.status.matchedDatasets`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type DataLifecyclePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataLifecyclePolicySpec   `json:"spec,omitempty"`
	Status DataLifecyclePolicyStatus `json:"status,omitempty"`
}

// DataLifecyclePolicyList is a list of DataLifecyclePolicy.
//
// +kubebuilder:object:root=true
type DataLifecyclePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DataLifecyclePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DataLifecyclePolicy{}, &DataLifecyclePolicyList{})
}
