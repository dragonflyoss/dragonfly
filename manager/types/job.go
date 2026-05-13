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

package types

import "time"

const (
	// SingleSeedPeerScope represents the scope that only single seed peer will be preheated.
	SingleSeedPeerScope = "single_seed_peer"

	// AllSeedPeersScope represents the scope that all seed peers will be preheated.
	AllSeedPeersScope = "all_seed_peers"

	// AllPeersScope represents the scope that all peers will be preheated.
	AllPeersScope = "all_peers"
)

const (
	// DefaultPreheatConcurrentPeerCount is the default concurrent peer count for preheating all peers.
	DefaultPreheatConcurrentPeerCount = 500

	// DefaultPreheatConcurrentTaskCount is the default concurrent task count for preheating all peers.
	DefaultPreheatConcurrentTaskCount = 8

	// DefaultPreheatConcurrentLayerCount is the default concurrent layer count for getting image distribution.
	DefaultPreheatConcurrentLayerCount = 8

	// DefaultGetTaskConcurrentPeerCount is the default concurrent peer count for getting task.
	DefaultGetTaskConcurrentPeerCount = 500

	// DefaultJobTimeout is the default timeout for executing job.
	DefaultJobTimeout = 60 * time.Minute
)

type CreateJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// Args is the arguments of the job.
	Args map[string]any `json:"args" binding:"omitempty"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`

	// SeedPeerClusterIDs is the seed peer cluster ids of the job.
	SeedPeerClusterIDs []uint `json:"seed_peer_cluster_ids" binding:"omitempty"`

	// SchedulerClusterIDs is the scheduler cluster ids of the job.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`
}

type UpdateJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`
}

type JobParams struct {
	// Type is the type of the job.
	ID uint `uri:"id" binding:"required"`
}

type GetJobsQuery struct {
	// Type is the type of the job.
	Type string `form:"type" binding:"omitempty"`

	// State is the state of the job.
	State string `form:"state" binding:"omitempty,oneof=PENDING RECEIVED STARTED RETRY SUCCESS FAILURE"`

	// UserID is the user id of the job.
	UserID uint `form:"user_id" binding:"omitempty"`

	// Page is the page number of the job list.
	Page int `form:"page" binding:"omitempty,gte=1"`

	// PerPage is the item count per page of the job list.
	PerPage int `form:"per_page" binding:"omitempty,gte=1,lte=10000000"`
}

type CreatePreheatJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// Args is the arguments of the preheating job.
	Args PreheatArgs `json:"args" binding:"omitempty"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`

	// SchedulerClusterIDs is the scheduler cluster ids of the job.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`
}

type PreheatArgs struct {
	// Type is the preheating type, support image and file.
	Type string `json:"type" binding:"required,oneof=image file"`

	// URL is the image manifest url for preheating.
	URL string `json:"url" binding:"omitempty"`

	// URLs is the file urls for preheating, only support file type.
	URLs []string `json:"urls" binding:"omitempty"`

	// PieceLength is the piece length(bytes) for downloading file. The value needs to
	// be greater than 4MiB (4,194,304 bytes) and less than 64MiB (67,108,864 bytes),
	// for example: 4194304(4mib), 8388608(8mib). If the piece length is not specified,
	// the piece length will be calculated according to the file size.
	PieceLength *uint64 `json:"piece_length" binding:"omitempty,gte=4194304,lte=67108864"`

	// Tag is the tag for preheating.
	Tag string `json:"tag" binding:"omitempty"`

	// Application is the application string for preheating.
	Application string `json:"application" binding:"omitempty"`

	// FilteredQueryParams is the filtered query params for preheating.
	FilteredQueryParams string `json:"filtered_query_params" binding:"omitempty"`

	// Headers is the http headers for authentication.
	Headers map[string]string `json:"headers" binding:"omitempty"`

	// Username is the username for authentication.
	Username string `json:"username" binding:"omitempty"`

	// Password is the password for authentication.
	Password string `json:"password" binding:"omitempty"`

	// The image type preheating task can specify the image architecture type. eg: linux/amd64.
	Platform string `json:"platform" binding:"omitempty"`

	// Scope is the scope for preheating, default is single_seed_peer.
	Scope string `json:"scope" binding:"omitempty"`

	// IPs is a list of specific peer IPs for preheating.
	// This field has the highest priority: if provided, both 'Count' and 'Percentage' will be ignored.
	// Applies to 'all_peers' and 'all_seed_peers' scopes.
	IPs []string `json:"ips" binding:"omitempty,gte=1,lte=100"`

	// Percentage is the percentage of available peers to preheat.
	// This field has the lowest priority and is only used if both 'IPs' and 'Count' are not provided.
	// It must be a value between 1 and 100 (inclusive) if provided.
	// Applies to 'all_peers' and 'all_seed_peers' scopes.
	Percentage *uint32 `json:"percentage" binding:"omitempty,gte=1,lte=100"`

	// Count is the desired number of peers to preheat.
	// This field is used only when 'IPs' is not specified. It has priority over 'Percentage'.
	// It must be a value between 1 and 200 (inclusive) if provided.
	// Applies to 'all_peers' and 'all_seed_peers' scopes.
	Count *uint32 `json:"count" binding:"omitempty,gte=1,lte=200"`

	// ConcurrentTaskCount specifies the maximum number of tasks (e.g., image layers) to preheat concurrently.
	// For example, if preheating 100 layers with ConcurrentTaskCount set to 10, up to 10 layers are processed simultaneously.
	// If ConcurrentPeerCount is 10 for 1000 peers, each layer is preheated by 10 peers concurrently.
	// Default is 8, maximum is 100.
	ConcurrentTaskCount int64 `json:"concurrent_task_count" binding:"omitempty,gte=1,lte=100"`

	// ConcurrentPeerCount specifies the maximum number of peers to preheat concurrently for a single task (e.g., an image layer).
	// For example, if preheating a layer with ConcurrentPeerCount set to 10, up to 10 peers process that layer simultaneously.
	// Default is 500, maximum is 1000.
	ConcurrentPeerCount int64 `json:"concurrent_peer_count" binding:"omitempty,gte=1,lte=1000"`

	// Timeout is the timeout for preheating, default is 60 minutes.
	Timeout time.Duration `json:"timeout" binding:"omitempty"`

	// EnableTaskIDBasedBlobDigest indicates whether to use the blob digest for task ID calculation
	// when downloading from OCI registries. When enabled for OCI blob URLs (e.g., /v2/<name>/blobs/sha256:<digest>),
	// the task ID is derived from the blob digest rather than the full URL. This enables deduplication across
	// registries - the same blob from different registries shares one task ID, eliminating redundant downloads
	// and storage. Default is false for backward compatibility.
	EnableTaskIDBasedBlobDigest *bool `json:"enable_task_id_based_blob_digest" binding:"omitempty"`
}

type CreateSyncPeersJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`

	// SchedulerClusterIDs is the scheduler cluster ids of the job.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`
}

type CreateGetTaskJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// Args is the arguments of the job.
	Args GetTaskArgs `json:"args" binding:"omitempty"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`

	// SchedulerClusterIDs is the scheduler cluster ids of the job.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`
}

type GetTaskArgs struct {
	// TaskID is the task id for getting.
	TaskID string `json:"task_id" binding:"omitempty"`

	// URL is the download url of the task.
	URL string `json:"url" binding:"omitempty"`

	// PieceLength is the piece length(bytes) for downloading file. The value needs to
	// be greater than 4MiB (4,194,304 bytes) and less than 64MiB (67,108,864 bytes),
	// for example: 4194304(4mib), 8388608(8mib). If the piece length is not specified,
	// the piece length will be calculated according to the file size.
	PieceLength *uint64 `json:"piece_length" binding:"omitempty,gte=4194304,lte=67108864"`

	// Tag is the tag of the task.
	Tag string `json:"tag" binding:"omitempty"`

	// Application is the application of the task.
	Application string `json:"application" binding:"omitempty"`

	// FilteredQueryParams is the filtered query params of the task.
	FilteredQueryParams string `json:"filtered_query_params" binding:"omitempty"`

	// ContentForCalculatingTaskID is the content used to calculate the task id.
	// If ContentForCalculatingTaskID is set, use its value to calculate the task ID.
	// Otherwise, calculate the task ID based on url, piece_length, tag, application, and filtered_query_params.
	ContentForCalculatingTaskID *string `json:"content_for_calculating_task_id" binding:"omitempty"`

	// ConcurrentPeerCount specifies the maximum number of peers stat concurrently for a single task (e.g., an image layer).
	// For example, if stat a layer with ConcurrentPeerCount set to 10, up to 10 peers process that layer simultaneously.
	// Default is 500, maximum is 1000.
	ConcurrentPeerCount int64 `json:"concurrent_peer_count" binding:"omitempty,gte=1,lte=1000"`

	// Timeout is the timeout for getting task, default is 60 minutes.
	Timeout time.Duration `json:"timeout" binding:"omitempty"`
}

type CreateGetImageDistributionJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// Args is the arguments of the job.
	Args GetImageDistributionArgs `json:"args" binding:"omitempty"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`

	// SchedulerClusterIDs is the scheduler cluster ids of the job.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`
}

type GetImageDistributionArgs struct {
	// URL is the image manifest url of the task.
	URL string `json:"url" binding:"required"`

	// PieceLength is the piece length(bytes) for downloading image blobs. The value needs to
	// be greater than 4MiB (4,194,304 bytes) and less than 64MiB (67,108,864 bytes),
	// for example: 4194304(4mib), 8388608(8mib). If the piece length is not specified,
	// the piece length will be calculated according to the file size.
	PieceLength *uint64 `json:"piece_length" binding:"omitempty,gte=4194304,lte=67108864"`

	// Tag is the tag of the task.
	Tag string `json:"tag" binding:"omitempty"`

	// Application is the application of the task.
	Application string `json:"application" binding:"omitempty"`

	// FilteredQueryParams is the filtered query params of the task.
	FilteredQueryParams string `json:"filtered_query_params" binding:"omitempty"`

	// Headers is the http headers for authentication.
	Headers map[string]string `json:"headers" binding:"omitempty"`

	// Username is the username for authentication.
	Username string `json:"username" binding:"omitempty"`

	// Password is the password for authentication.
	Password string `json:"password" binding:"omitempty"`

	// The image type preheating task can specify the image architecture type. eg: linux/amd64.
	Platform string `json:"platform" binding:"omitempty"`

	// ConcurrentLayerCount specifies the maximum number of layers to get concurrently.
	ConcurrentLayerCount int64 `json:"concurrent_layer_count" binding:"omitempty,gte=1,lte=100"`

	// ConcurrentPeerCount specifies the maximum number of peers stat concurrently for a single task (e.g., an image layer).
	// For example, if stat a layer with ConcurrentPeerCount set to 10, up to 10 peers process that layer simultaneously.
	// Default is 500, maximum is 1000.
	ConcurrentPeerCount int64 `json:"concurrent_peer_count" binding:"omitempty,gte=1,lte=1000"`

	// Timeout is the timeout for getting image distribution, default is 60 minutes.
	Timeout time.Duration `json:"timeout" binding:"omitempty"`
}

// CreateGetImageDistributionJobResponse is the response for creating a get image job.
type CreateGetImageDistributionJobResponse struct {
	// Image is the image information.
	Image Image `json:"image"`

	// Peers is the peers that have downloaded the image.
	Peers []Peer `json:"peers"`
}

// Peer represents a peer in the get image job.
type Peer struct {
	// IP is the IP address of the peer.
	IP string `json:"ip"`

	// Hostname is the hostname of the peer.
	Hostname string `json:"hostname"`

	// CachedLayers is the list of layers that the peer has downloaded.
	CachedLayers []Layer `json:"cached_layers"`

	// SchedulerClusterID is the scheduler cluster id of the peer.
	SchedulerClusterID uint `json:"scheduler_cluster_id"`
}

// Image represents the image information.
type Image struct {
	// Layers is the list of layers of the image.
	Layers []Layer `json:"layers"`
}

// Layer represents a layer of the image.
type Layer struct {
	// URL is the URL of the layer.
	URL string `json:"url"`
}

type CreateDeleteTaskJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// Args is the arguments of the job.
	Args DeleteTaskArgs `json:"args" binding:"omitempty"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`

	// SchedulerClusterIDs is the scheduler cluster ids of the job.
	SchedulerClusterIDs []uint `json:"scheduler_cluster_ids" binding:"omitempty"`
}

type DeleteTaskArgs struct {
	// TaskID is the task id for deleting.
	TaskID string `json:"task_id" binding:"omitempty"`

	// URL is the download url of the task.
	URL string `json:"url" binding:"omitempty"`

	// PieceLength is the piece length(bytes) for downloading file. The value needs to
	// be greater than 4MiB (4,194,304 bytes) and less than 64MiB (67,108,864 bytes),
	// for example: 4194304(4mib), 8388608(8mib). If the piece length is not specified,
	// the piece length will be calculated according to the file size.
	PieceLength *uint64 `json:"piece_length" binding:"omitempty,gte=4194304,lte=67108864"`

	// Tag is the tag of the task.
	Tag string `json:"tag" binding:"omitempty"`

	// Application is the application of the task.
	Application string `json:"application" binding:"omitempty"`

	// FilteredQueryParams is the filtered query params of the task.
	FilteredQueryParams string `json:"filtered_query_params" binding:"omitempty"`

	// Timeout is the timeout for deleting, default is 60 minutes.
	Timeout time.Duration `json:"timeout" binding:"omitempty"`

	// ContentForCalculatingTaskID is the content used to calculate the task id.
	// If ContentForCalculatingTaskID is set, use its value to calculate the task ID.
	// Otherwise, calculate the task ID based on url, piece_length, tag, application, and filtered_query_params.
	ContentForCalculatingTaskID *string `json:"content_for_calculating_task_id" binding:"omitempty"`
}

type CreateGCJobRequest struct {
	// BIO is the description of the job.
	BIO string `json:"bio" binding:"omitempty"`

	// Type is the type of the job.
	Type string `json:"type" binding:"required"`

	// Args is the arguments of the gc.
	Args GCArgs `json:"args" binding:"required"`

	// UserID is the user id of the job.
	UserID uint `json:"user_id" binding:"omitempty"`
}

type GCArgs struct {
	// Type is the type of the job.
	Type string `json:"type" binding:"required,oneof=audit job"`
}
