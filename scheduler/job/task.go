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

package job

import (
	pkgdigest "d7y.io/dragonfly/v2/pkg/digest"
	"d7y.io/dragonfly/v2/pkg/idgen"
)

// preheatTaskID computes the task ID for a preheat URL.
// If the URL is an OCI blob URL, it uses the blob digest-based task ID
// and returns true for isBlobURL; otherwise it falls back to URL-based task ID.
func preheatTaskID(url string, pieceLength *uint64, tag, application string, filteredQueryParams []string) (taskID string, isBlobURL bool) {
	isBlobURL = pkgdigest.IsBlobURL(url)
	if isBlobURL {
		taskID = idgen.TaskIDByBlobDigest(url)
	} else {
		taskID = idgen.TaskIDV2ByURLBased(url, pieceLength, tag, application, filteredQueryParams, "")
	}
	return
}
