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

package cdn

import (
	"d7y.io/dragonfly/v2/cdnsystem/types"
	logger "d7y.io/dragonfly/v2/pkg/dflog"
	"d7y.io/dragonfly/v2/pkg/structure/maputils"
	"d7y.io/dragonfly/v2/pkg/util/rangeutils"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

const RangeHeaderName = "Range"

// download downloads the file from the original address and
// sets the "Range" header to the unDownloaded file range.
//
// If the returned error is nil, the Response will contain a non-nil
// Body which the caller is expected to close.
func (cm *Manager) download(task *types.SeedTask, detectResult *cacheResult) (io.ReadCloser, map[string]string, error) {
	headers := maputils.DeepCopyMap(nil, task.Headers)
	if detectResult.breakPoint > 0 {
		breakRange, err := rangeutils.GetBreakRange(detectResult.breakPoint, task.SourceFileLength)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to calculate the breakRange")
		}
		// check if Range in header? if Range already in Header, priority use this range
		if !hasRange(headers) {
			headers[RangeHeaderName] = fmt.Sprintf("bytes=%s", breakRange)
		}
	}
	logger.WithTaskID(task.TaskId).Infof("start download url %s at range:%d-%d: %s with header: %+v", task.Url, detectResult.breakPoint,
		task.SourceFileLength, task.Headers)
	return cm.resourceClient.Download(task.Url, headers)
}

func hasRange(headers map[string]string) bool {
	if headers == nil {
		return false
	}
	_, ok := headers["Range"]
	return ok
}
