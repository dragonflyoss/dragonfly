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

package util

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	"d7y.io/dragonfly/v2/pkg/idgen"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
)

type FileSize = uint64

const (
	// FileSize1MiB represents the size of 1MiB.
	FileSize1MiB FileSize = 1024 * 1024

	// FileSize10MiB represents the size of 10MiB.
	FileSize10MiB FileSize = 10 * FileSize1MiB

	// FileSize100MiB represents the size of 100MiB.
	FileSize100MiB FileSize = 100 * FileSize1MiB

	// FileSize500MiB represents the size of 500MiB.
	FileSize500MiB FileSize = 500 * FileSize1MiB

	// FileSize1GiB represents the size of 1GiB.
	FileSize1GiB FileSize = 1024 * FileSize1MiB
)

// File represents a file.
type File struct {
	// info is the local file info.
	info os.FileInfo

	// localPath is the local path of the file.
	localPath string

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

// GetRangeSha256 returns the sha256 of the range part of the file content.
func (f *File) GetRangeSha256(r string, fileSize int64) string {
	// parse the range header
	parsedRange := nethttp.MustParseRange(fmt.Sprintf("bytes=%s", r), fileSize)
	start, end := parsedRange.Start, parsedRange.Start+parsedRange.Length-1

	file, err := os.Open(f.localPath)
	if err != nil {
		fmt.Printf("open file %s error: %v\n", f.localPath, err)
		return ""
	}
	defer file.Close()

	if _, err = file.Seek(start, 0); err != nil {
		fmt.Printf("seek file %s to %d error: %v\n", f.localPath, start, err)
		return ""
	}

	limitedReader := io.LimitReader(file, end-start+1)
	bufferedReader := bufio.NewReader(limitedReader)

	hash := sha256.New()
	if _, err = io.Copy(hash, bufferedReader); err != nil {
		fmt.Printf("copy file %s content error: %v\n", f.localPath, err)
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// GetDownloadURL returns the download URL of the file from remote file server.
func (f *File) GetDownloadURL() string {
	return f.downloadURL
}

// GetTaskID returns the task id of the file.
func (f *File) GetTaskID(opts ...TaskIDOption) string {
	taskIDOptions := &taskID{
		url: f.downloadURL,
	}

	for _, opt := range opts {
		opt(taskIDOptions)
	}

	return idgen.TaskIDV2(taskIDOptions.url, taskIDOptions.pieceLength, taskIDOptions.tag, taskIDOptions.application, taskIDOptions.filteredQueryParams)
}

// GetOutputPath returns the output path of the file.
func (f *File) GetOutputPath() string {
	return path.Join("/tmp", f.info.Name())
}

// parseRangeHeader parses the Range header value and returns start and end positions.
func parseRangeHeader(rangeHeader string, fileSize int64) (start, end int64) {
	if rangeHeader == "" {
		return 0, fileSize - 1
	}

	// split the range header by "-"
	parts := strings.Split(rangeHeader, "-")

	// handle different range formats
	switch {
	case parts[0] == "": // -N: last N bytes
		end = fileSize - 1
		bytes, _ := strconv.ParseInt(parts[1], 10, 64)
		start = fileSize - bytes
		if start < 0 {
			start = 0
		}

	case parts[1] == "": // N-: from N to end
		start, _ = strconv.ParseInt(parts[0], 10, 64)
		end = fileSize - 1

	default: // N-M: from N to M
		start, _ = strconv.ParseInt(parts[0], 10, 64)
		end, _ = strconv.ParseInt(parts[1], 10, 64)
		if end >= fileSize {
			end = fileSize - 1
		}
	}

	return start, end
}
