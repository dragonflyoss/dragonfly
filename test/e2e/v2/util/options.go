/*
 *     Copyright 2024 The Dragonfly Authors
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

// FileServerOptions represents the options of the file server.
type FileServerOptions struct {
	// endpoint is the endpoint of the file server.
	endpoint string
	// localDir is the local temporary directory of the file server.
	localDir string
}

// FileServerOption is the type of the options of the file server.
type FileServerOption func(*FileServerOptions)

// WithFileServerEndpoint sets the endpoint of the file server.
func WithFileServerEndpoint(endpoint string) FileServerOption {
	return func(o *FileServerOptions) {
		o.endpoint = endpoint
	}
}

// WithFileServerLocalDir sets the local temporary directory of the file server.
func WithFileServerLocalDir(localDir string) FileServerOption {
	return func(o *FileServerOptions) {
		o.localDir = localDir
	}
}

// TaskIDOptions represents the options of the task id.
type TaskIDOptions struct {
	// url is the url of the download task.
	url string
	// tag is the tag of the download task.
	tag string
	// appliccation is the application of the download task.
	application string
	// filteredQueryParams is the filtered query params of the download task.
	filteredQueryParams []string
}

// TaskIDOption is the type of the options of the task id.
type TaskIDOption func(*TaskIDOptions)

// WithTaskIDURL sets the url of the download task.
func WithTaskIDURL(url string) TaskIDOption {
	return func(o *TaskIDOptions) {
		o.url = url
	}
}

// WithTaskIDTag sets the tag of the download task.
func WithTaskIDTag(tag string) TaskIDOption {
	return func(o *TaskIDOptions) {
		o.tag = tag
	}
}

// WithTaskIDApplication sets the application of the download task.
func WithTaskIDApplication(application string) TaskIDOption {
	return func(o *TaskIDOptions) {
		o.application = application
	}
}

// WithTaskIDFilteredQueryParams sets the filtered query params of the download task.
func WithTaskIDFilteredQueryParams(filteredQueryParams []string) TaskIDOption {
	return func(o *TaskIDOptions) {
		o.filteredQueryParams = filteredQueryParams
	}
}
