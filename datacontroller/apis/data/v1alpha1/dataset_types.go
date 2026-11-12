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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatasetFinalizer is the finalizer added to datasets so that the data they
// distributed can be purged from the P2P network before the object disappears.
const DatasetFinalizer = "data.d7y.io/dataset-lifecycle"

// DataSourceType is the kind of data a source refers to.
// +kubebuilder:validation:Enum=file;image
type DataSourceType string

const (
	// DataSourceTypeFile is a source made of one or more plain URLs (HTTP(S), object storage, HDFS).
	DataSourceTypeFile DataSourceType = "file"

	// DataSourceTypeImage is a source that refers to an OCI image manifest URL; every layer
	// referenced by the manifest is managed.
	DataSourceTypeImage DataSourceType = "image"
)

// DistributionScope is the set of peers the data is pushed to.
// +kubebuilder:validation:Enum=single_seed_peer;all_seed_peers;all_peers
type DistributionScope string

const (
	// DistributionScopeSingleSeedPeer distributes the data to one seed peer per scheduler cluster.
	DistributionScopeSingleSeedPeer DistributionScope = "single_seed_peer"

	// DistributionScopeAllSeedPeers distributes the data to all seed peers.
	DistributionScopeAllSeedPeers DistributionScope = "all_seed_peers"

	// DistributionScopeAllPeers distributes the data to all peers.
	DistributionScopeAllPeers DistributionScope = "all_peers"
)

// ExpireAction is what happens to the distributed data when its TTL elapses.
// +kubebuilder:validation:Enum=Purge;Retain
type ExpireAction string

const (
	// ExpireActionPurge removes the data from the P2P network when it expires.
	ExpireActionPurge ExpireAction = "Purge"

	// ExpireActionRetain keeps the data on the peers and only marks the dataset as expired.
	ExpireActionRetain ExpireAction = "Retain"
)

// DeletionPolicy is what happens to the distributed data when the Dataset object is deleted.
// +kubebuilder:validation:Enum=Purge;Retain
type DeletionPolicy string

const (
	// DeletionPolicyPurge removes the data from the P2P network before the Dataset is deleted.
	DeletionPolicyPurge DeletionPolicy = "Purge"

	// DeletionPolicyRetain leaves the data on the peers when the Dataset is deleted.
	DeletionPolicyRetain DeletionPolicy = "Retain"
)

// DataSource describes one piece of data managed as part of a Dataset.
type DataSource struct {
	// Name identifies the source within the dataset and is used to correlate status entries.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name"`

	// Type is the kind of source, "file" or "image".
	// +kubebuilder:default=file
	Type DataSourceType `json:"type,omitempty"`

	// URL is the image manifest URL for image sources or a single file URL for file sources.
	// +optional
	URL string `json:"url,omitempty"`

	// URLs is the list of file URLs for file sources. It is ignored for image sources.
	// +optional
	URLs []string `json:"urls,omitempty"`

	// Tag distinguishes tasks with the same URL from each other.
	// +optional
	Tag string `json:"tag,omitempty"`

	// Application is the application name the tasks are attributed to.
	// +optional
	Application string `json:"application,omitempty"`

	// Platform selects the image platform to distribute, for example linux/amd64. Image sources only.
	// +optional
	Platform string `json:"platform,omitempty"`

	// PieceLength is the piece length in bytes used to download the data. Must be between 4MiB and 64MiB.
	// +kubebuilder:validation:Minimum=4194304
	// +kubebuilder:validation:Maximum=67108864
	// +optional
	PieceLength *uint64 `json:"pieceLength,omitempty"`

	// FilteredQueryParams is the comma separated list of query parameters ignored when computing task IDs.
	// +optional
	FilteredQueryParams string `json:"filteredQueryParams,omitempty"`

	// Headers are HTTP headers sent when fetching the data.
	// +optional
	Headers map[string]string `json:"headers,omitempty"`

	// HeadersSecretRef references a Secret in the same namespace whose data entries are all sent as HTTP headers.
	// Headers from the Secret take precedence over Headers.
	// +optional
	HeadersSecretRef *corev1.LocalObjectReference `json:"headersSecretRef,omitempty"`

	// AuthSecretRef references a Secret in the same namespace holding "username" and "password" entries
	// used for basic authentication against the origin.
	// +optional
	AuthSecretRef *corev1.LocalObjectReference `json:"authSecretRef,omitempty"`

	// ObjectStorage configures access to object storage backends such as s3, gcs or oss.
	// +optional
	ObjectStorage *ObjectStorageSpec `json:"objectStorage,omitempty"`

	// HDFS configures access to HDFS backends.
	// +optional
	HDFS *HDFSSpec `json:"hdfs,omitempty"`

	// EnableTaskIDBasedBlobDigest derives task IDs of OCI blob URLs from the blob digest so the same
	// blob pulled from different registries shares one task.
	// +optional
	EnableTaskIDBasedBlobDigest bool `json:"enableTaskIDBasedBlobDigest,omitempty"`
}

// ObjectStorageSpec configures access to an object storage backend.
type ObjectStorageSpec struct {
	// Region of the bucket.
	// +optional
	Region string `json:"region,omitempty"`

	// Endpoint of the object storage service.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// PredefinedACL applied to the objects.
	// +optional
	PredefinedACL string `json:"predefinedACL,omitempty"`

	// InsecureSkipVerify disables TLS verification against the object storage service.
	// +optional
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`

	// CredentialsSecretRef references a Secret in the same namespace holding the "accessKeyID",
	// "accessKeySecret" and optional "sessionToken", "securityToken" and "credentialPath" entries.
	// +optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
}

// HDFSSpec configures access to an HDFS backend.
type HDFSSpec struct {
	// DelegationTokenSecretRef references a Secret in the same namespace holding a "delegationToken" entry.
	// +optional
	DelegationTokenSecretRef *corev1.LocalObjectReference `json:"delegationTokenSecretRef,omitempty"`
}

// DistributionPolicy controls how data is pushed into the P2P network.
type DistributionPolicy struct {
	// Scope is the set of peers the data is distributed to.
	// +kubebuilder:default=single_seed_peer
	// +optional
	Scope DistributionScope `json:"scope,omitempty"`

	// IPs restricts distribution to these peer IPs. It takes precedence over Count and Percentage.
	// +kubebuilder:validation:MaxItems=100
	// +optional
	IPs []string `json:"ips,omitempty"`

	// Count is the number of peers to distribute to. It takes precedence over Percentage.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=200
	// +optional
	Count *uint32 `json:"count,omitempty"`

	// Percentage is the percentage of available peers to distribute to.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +optional
	Percentage *uint32 `json:"percentage,omitempty"`

	// ConcurrentTaskCount is the maximum number of tasks (for example image layers) distributed concurrently.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	// +optional
	ConcurrentTaskCount *int64 `json:"concurrentTaskCount,omitempty"`

	// ConcurrentPeerCount is the maximum number of peers a single task is distributed to concurrently.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +optional
	ConcurrentPeerCount *int64 `json:"concurrentPeerCount,omitempty"`

	// SchedulerClusterIDs restricts the operation to these scheduler clusters. Empty means all clusters.
	// +optional
	SchedulerClusterIDs []uint `json:"schedulerClusterIDs,omitempty"`

	// Timeout bounds the duration of every job created for the dataset. Defaults to the manager default.
	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`

	// RetryLimit is the number of times a failed operation is retried before the source is marked failed.
	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=0
	// +optional
	RetryLimit *int32 `json:"retryLimit,omitempty"`
}

// LifecycleSpec describes what happens to the data over time. It is used inline in a
// Dataset and as the body of a DataLifecyclePolicy. Fields set on a Dataset override
// fields set by a matching DataLifecyclePolicy.
type LifecycleSpec struct {
	// TTL is how long the data is kept after its last successful distribution. When it elapses the
	// ExpireAction is applied and the dataset is marked expired.
	// +optional
	TTL *metav1.Duration `json:"ttl,omitempty"`

	// ExpireAction is applied when the TTL elapses.
	// +optional
	ExpireAction *ExpireAction `json:"expireAction,omitempty"`

	// RefreshInterval re-distributes the data periodically after the last successful distribution,
	// keeping it warm on the peers. Mutually exclusive with RefreshSchedule.
	// +optional
	RefreshInterval *metav1.Duration `json:"refreshInterval,omitempty"`

	// RefreshSchedule re-distributes the data on a cron schedule. Mutually exclusive with RefreshInterval.
	// +optional
	RefreshSchedule *string `json:"refreshSchedule,omitempty"`

	// VerifyInterval periodically checks which peers still hold the data and reports it in the status.
	// +optional
	VerifyInterval *metav1.Duration `json:"verifyInterval,omitempty"`

	// DeletionPolicy is applied when the Dataset object is deleted.
	// +optional
	DeletionPolicy *DeletionPolicy `json:"deletionPolicy,omitempty"`
}

// ManagerRef points at the Dragonfly manager that operates the P2P network.
type ManagerRef struct {
	// Endpoint is the base URL of the manager REST API, for example http://dragonfly-manager.dragonfly-system.svc:8080.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// TokenSecretRef references a Secret key in the same namespace holding a manager personal access
	// token with the "job" scope. The key defaults to "token".
	// +optional
	TokenSecretRef *corev1.SecretKeySelector `json:"tokenSecretRef,omitempty"`

	// InsecureSkipTLSVerify disables TLS verification against the manager.
	// +optional
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
}

// DatasetSpec is the desired state of a Dataset.
type DatasetSpec struct {
	// Sources is the data that makes up the dataset.
	// +kubebuilder:validation:MinItems=1
	// +listType=map
	// +listMapKey=name
	Sources []DataSource `json:"sources"`

	// Distribution controls how the data is pushed into the P2P network.
	// +optional
	Distribution DistributionPolicy `json:"distribution,omitempty"`

	// Lifecycle controls what happens to the data over time. Fields set here override
	// the matching DataLifecyclePolicy.
	// +optional
	Lifecycle *LifecycleSpec `json:"lifecycle,omitempty"`

	// ManagerRef overrides the manager configured on the controller.
	// +optional
	ManagerRef *ManagerRef `json:"managerRef,omitempty"`

	// Suspend stops the controller from starting new operations for this dataset. In-flight
	// operations keep being tracked.
	// +optional
	Suspend bool `json:"suspend,omitempty"`
}

// DatasetPhase is a coarse summary of the dataset state.
type DatasetPhase string

const (
	// DatasetPhasePending means the dataset has not been distributed yet.
	DatasetPhasePending DatasetPhase = "Pending"

	// DatasetPhaseDistributing means at least one source is being distributed.
	DatasetPhaseDistributing DatasetPhase = "Distributing"

	// DatasetPhaseReady means every source has been distributed successfully.
	DatasetPhaseReady DatasetPhase = "Ready"

	// DatasetPhaseDegraded means some sources are distributed and others failed, or verification
	// found data missing from the peers.
	DatasetPhaseDegraded DatasetPhase = "Degraded"

	// DatasetPhaseFailed means every source failed to distribute.
	DatasetPhaseFailed DatasetPhase = "Failed"

	// DatasetPhaseExpired means the TTL elapsed.
	DatasetPhaseExpired DatasetPhase = "Expired"

	// DatasetPhasePurging means data is being removed from the peers.
	DatasetPhasePurging DatasetPhase = "Purging"

	// DatasetPhaseSuspended means the dataset is suspended.
	DatasetPhaseSuspended DatasetPhase = "Suspended"
)

// DataSourceState is the state of a single source.
type DataSourceState string

const (
	// DataSourceStatePending means the source has not been distributed yet.
	DataSourceStatePending DataSourceState = "Pending"

	// DataSourceStateDistributing means a distribution operation is in flight.
	DataSourceStateDistributing DataSourceState = "Distributing"

	// DataSourceStateReady means the last distribution succeeded.
	DataSourceStateReady DataSourceState = "Ready"

	// DataSourceStateDegraded means verification found the data missing from some peers.
	DataSourceStateDegraded DataSourceState = "Degraded"

	// DataSourceStateFailed means distribution failed and retries are exhausted.
	DataSourceStateFailed DataSourceState = "Failed"

	// DataSourceStateExpired means the TTL elapsed.
	DataSourceStateExpired DataSourceState = "Expired"

	// DataSourceStatePurging means a purge operation is in flight.
	DataSourceStatePurging DataSourceState = "Purging"

	// DataSourceStatePurged means the data was removed from the peers.
	DataSourceStatePurged DataSourceState = "Purged"
)

// OperationType is the kind of operation run against the P2P network.
type OperationType string

const (
	// OperationTypeDistribute preheats the data onto the peers.
	OperationTypeDistribute OperationType = "Distribute"

	// OperationTypeVerify checks which peers hold the data.
	OperationTypeVerify OperationType = "Verify"

	// OperationTypePurge removes the data from the peers.
	OperationTypePurge OperationType = "Purge"
)

// JobReference points at a job in the Dragonfly manager.
type JobReference struct {
	// ID is the job ID in the manager.
	ID uint `json:"id"`

	// Type is the manager job type, for example preheat, get_task or delete_task.
	Type string `json:"type"`

	// State is the last observed job state.
	// +optional
	State string `json:"state,omitempty"`

	// URL is the task URL the job targets, when it targets a single task.
	// +optional
	URL string `json:"url,omitempty"`
}

// OperationStatus is an operation in flight for a source.
type OperationStatus struct {
	// Type of the operation.
	Type OperationType `json:"type"`

	// Jobs created in the manager for the operation.
	// +optional
	Jobs []JobReference `json:"jobs,omitempty"`

	// StartedAt is when the operation was started.
	StartedAt metav1.Time `json:"startedAt"`
}

// TaskStatus is the observed state of one task (a file URL or an image layer) of a source.
type TaskStatus struct {
	// URL of the task.
	URL string `json:"url"`

	// ID is the task ID in the P2P network.
	// +optional
	ID string `json:"id,omitempty"`

	// PeerCount is the number of peers holding the task at the last verification.
	// +optional
	PeerCount *int32 `json:"peerCount,omitempty"`
}

// DataSourceStatus is the observed state of a source.
type DataSourceStatus struct {
	// Name of the source.
	Name string `json:"name"`

	// State of the source.
	State DataSourceState `json:"state"`

	// SpecHash is the hash of the source and distribution spec at the last distribution.
	// +optional
	SpecHash string `json:"specHash,omitempty"`

	// Operation in flight, if any.
	// +optional
	Operation *OperationStatus `json:"operation,omitempty"`

	// Tasks distributed for the source.
	// +optional
	Tasks []TaskStatus `json:"tasks,omitempty"`

	// LastDistributionTime is when the source was last distributed successfully.
	// +optional
	LastDistributionTime *metav1.Time `json:"lastDistributionTime,omitempty"`

	// LastVerificationTime is when the source was last verified.
	// +optional
	LastVerificationTime *metav1.Time `json:"lastVerificationTime,omitempty"`

	// Retries is the number of consecutive failed attempts of the current operation type.
	// +optional
	Retries int32 `json:"retries,omitempty"`

	// NextRetryTime is when the failed operation is retried.
	// +optional
	NextRetryTime *metav1.Time `json:"nextRetryTime,omitempty"`

	// Reason is a machine readable reason for the state.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable description of the state.
	// +optional
	Message string `json:"message,omitempty"`
}

// Condition types reported on a Dataset.
const (
	// ConditionReady is true when every source is distributed and nothing is in flight.
	ConditionReady = "Ready"

	// ConditionDistributed is true when every source was distributed successfully at least once.
	ConditionDistributed = "Distributed"

	// ConditionVerified is true when the last verification found every task on at least one peer.
	ConditionVerified = "Verified"

	// ConditionExpired is true when the TTL elapsed.
	ConditionExpired = "Expired"
)

// DatasetStatus is the observed state of a Dataset.
type DatasetStatus struct {
	// ObservedGeneration is the generation last reconciled.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Phase is a coarse summary of the state.
	// +optional
	Phase DatasetPhase `json:"phase,omitempty"`

	// Conditions describe the state in detail.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Sources is the per-source state.
	// +listType=map
	// +listMapKey=name
	// +optional
	Sources []DataSourceStatus `json:"sources,omitempty"`

	// AppliedPolicy is the name of the DataLifecyclePolicy in effect, if any.
	// +optional
	AppliedPolicy string `json:"appliedPolicy,omitempty"`

	// LastDistributionTime is when every source was last distributed successfully.
	// +optional
	LastDistributionTime *metav1.Time `json:"lastDistributionTime,omitempty"`

	// NextRefreshTime is when the data will be re-distributed.
	// +optional
	NextRefreshTime *metav1.Time `json:"nextRefreshTime,omitempty"`

	// ExpirationTime is when the TTL elapses.
	// +optional
	ExpirationTime *metav1.Time `json:"expirationTime,omitempty"`

	// LastVerificationTime is when the data was last verified.
	// +optional
	LastVerificationTime *metav1.Time `json:"lastVerificationTime,omitempty"`

	// NextVerificationTime is when the data will be verified next.
	// +optional
	NextVerificationTime *metav1.Time `json:"nextVerificationTime,omitempty"`
}

// Dataset declares data that must be present in the Dragonfly P2P network and how its
// lifecycle is managed.
//
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ds;dset,categories=dragonfly
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Policy",type=string,JSONPath=`.status.appliedPolicy`
// +kubebuilder:printcolumn:name="Distributed",type=date,JSONPath=`.status.lastDistributionTime`
// +kubebuilder:printcolumn:name="Expires",type=date,JSONPath=`.status.expirationTime`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
type Dataset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatasetSpec   `json:"spec,omitempty"`
	Status DatasetStatus `json:"status,omitempty"`
}

// DatasetList is a list of Dataset.
//
// +kubebuilder:object:root=true
type DatasetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dataset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dataset{}, &DatasetList{})
}
