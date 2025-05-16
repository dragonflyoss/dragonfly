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

package storage

import (
	"errors"
	"os"
	"path/filepath"

	"d7y.io/dragonfly/v2/pkg/idgen"
)

// statTask is the options for stat task.
type statTask struct {
	// path is the base path of the storage.
	path string
	// url is the url of the task.
	url string
	// contentLength is the content length of the task.(optional)
	contentLength *uint64
	// pieceLength is the piece length of the task.(optional)
	pieceLength *uint64
	// tag is the tag of the task.(optional)
	tag string
	// application is the application of the task.(optional)
	application string
	// filteredQueryParams is the filtered query params of the task.(optional)
	filteredQueryParams []string
}

type StatTaskOption func(*statTask)

func WithStatTaskContentLength(contentLength *uint64) StatTaskOption {
	return func(task *statTask) {
		task.contentLength = contentLength
	}
}

func WithStatTaskPieceLength(pieceLength *uint64) StatTaskOption {
	return func(task *statTask) {
		task.pieceLength = pieceLength
	}
}

func WithStatTaskTag(tag string) StatTaskOption {
	return func(task *statTask) {
		task.tag = tag
	}
}

func WithStatTaskApplication(application string) StatTaskOption {
	return func(task *statTask) {
		task.application = application
	}
}

func WithStatTaskFilteredQueryParams(filteredQueryParams []string) StatTaskOption {
	return func(task *statTask) {
		task.filteredQueryParams = filteredQueryParams
	}
}

// StatTask stats the task by the given parameters.
func StatTask(path, url string, opts ...StatTaskOption) (os.FileInfo, error) {
	st := &statTask{
		path:                path,
		url:                 url,
		filteredQueryParams: idgen.DefaultFilteredQueryParams,
	}

	for _, opt := range opts {
		opt(st)
	}

	// Validate and mutate the options for stat task.
	if st.path == "" || st.url == "" {
		return nil, errors.New("path and url are required")
	}

	if st.contentLength == nil && st.pieceLength == nil {
		return nil, errors.New("either contentLength or pieceLength must be specified")
	}

	if st.pieceLength == nil {
		// Calculate pieceLength from contentLength if not specified.
		pieceLength := CalculatePieceLength(*st.contentLength)
		st.pieceLength = &pieceLength
	}

	taskID := idgen.TaskIDV2ByURLBased(st.url, st.pieceLength, st.tag, st.application, st.filteredQueryParams)

	// Construct the file path.
	filePath := filepath.Join(path, "content/tasks", taskID[0:3], taskID)
	return os.Stat(filePath)
}
