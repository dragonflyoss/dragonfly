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

import (
	"os"

	"d7y.io/dragonfly/v2/pkg/idgen"
)

// File represents a file.
type File struct {
	// info is the local file info.
	info os.FileInfo
	// downloadURL is the download URL of the file from remote file server.
	downloadURL string
}

// GetInfo returns the file info.
func (f *File) GetInfo() os.FileInfo {
	return f.info
}

// GetSha256 returns the sha256 of the file content.
func (f *File) GetSha256() string {
	// the file name is the sha256 of the file content
	return f.info.Name()
}

// GetDownloadURL returns the download URL of the file from remote file server.
func (f *File) GetDownloadURL() string {
	return f.downloadURL
}

// GetTaskID returns the task id of the file.
func (f *File) GetTaskID(opts ...TaskIDOption) string {
	taskIDOpts := &TaskIDOptions{
		url: f.downloadURL,
	}

	for _, opt := range opts {
		opt(taskIDOpts)
	}

	return idgen.TaskIDV2(taskIDOpts.url, taskIDOpts.tag, taskIDOpts.application, taskIDOpts.filteredQueryParams)
}
