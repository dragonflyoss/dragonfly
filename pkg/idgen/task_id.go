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

package idgen

import (
	"slices"
	"strconv"
	"strings"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	pkgdigest "d7y.io/dragonfly/v2/pkg/digest"
	neturl "d7y.io/dragonfly/v2/pkg/net/url"
	pkgstrings "d7y.io/dragonfly/v2/pkg/strings"
)

// DefaultFilteredQueryParams is the default filtered query params to generate the task id.
var DefaultFilteredQueryParams []string = slices.Concat(
	S3FilteredQueryParams,
	GcsFilteredQueryParams,
	OssFilteredQueryParams,
	ObsFilteredQueryParams,
	CosFilteredQueryParams,
	ContainerdFilteredQueryParams,
)

var (
	// S3FilteredQueryParams is the default filtered query params with s3 protocol to generate the task id.
	S3FilteredQueryParams = []string{
		"X-Amz-Algorithm",
		"X-Amz-Credential",
		"X-Amz-Date",
		"X-Amz-Expires",
		"X-Amz-SignedHeaders",
		"X-Amz-Signature",
		"X-Amz-Security-Token",
		"X-Amz-User-Agent",
	}

	// GcsFilteredQueryParams is the filtered query params with gcs protocol to generate the task id.
	GcsFilteredQueryParams = []string{
		"X-Goog-Algorithm",
		"X-Goog-Credential",
		"X-Goog-Date",
		"X-Goog-Expires",
		"X-Goog-SignedHeaders",
		"X-Goog-Signature",
	}

	// OssFilteredQueryParams is the default filtered query params with oss protocol to generate the task id.
	OssFilteredQueryParams = []string{
		"OSSAccessKeyId",
		"Expires",
		"Signature",
		"SecurityToken",
	}

	// ObsFilteredQueryParams is the default filtered query params with obs protocol to generate the task id.
	ObsFilteredQueryParams = []string{
		"AccessKeyId",
		"Signature",
		"Expires",
		"X-Obs-Date",
		"X-Obs-Security-Token",
	}

	// CosFilteredQueryParams is the default filtered query params with cos protocol to generate the task id.
	CosFilteredQueryParams = []string{
		"q-sign-algorithm",
		"q-ak",
		"q-sign-time",
		"q-key-time",
		"q-header-list",
		"q-url-param-list",
		"q-signature",
		"x-cos-security-token",
	}

	// ContainerdFilteredQueryParams is the default filtered query params with containerd to generate the task id.
	ContainerdFilteredQueryParams = []string{
		"ns",
	}
)

const (
	// FilteredQueryParamsSeparator is the separator of filtered query params.
	FilteredQueryParamsSeparator = "&"
)

// TaskIDV1 generates v1 version of task id.
// filter is separated by & character.
func TaskIDV1(url string, meta *commonv1.UrlMeta) string {
	return taskIDV1(url, meta, false)
}

// ParentTaskIDV1 generates v1 version of parent task id, but without range.
// this func is used to check the parent tasks for ranged requests
func ParentTaskIDV1(url string, meta *commonv1.UrlMeta) string {
	return taskIDV1(url, meta, true)
}

// taskIDV1 generates v1 version of task id.
// filter is separated by & character.
func taskIDV1(url string, meta *commonv1.UrlMeta, ignoreRange bool) string {
	if meta == nil {
		return pkgdigest.SHA256FromStrings(url)
	}

	filteredQueryParams := ParseFilteredQueryParams(meta.Filter)

	var (
		u   string
		err error
	)
	u, err = neturl.FilterQueryParams(url, filteredQueryParams)
	if err != nil {
		u = ""
	}

	data := []string{u}
	if meta.Digest != "" {
		data = append(data, meta.Digest)
	}

	if !ignoreRange && meta.Range != "" {
		data = append(data, meta.Range)
	}

	if meta.Tag != "" {
		data = append(data, meta.Tag)
	}

	if meta.Application != "" {
		data = append(data, meta.Application)
	}

	return pkgdigest.SHA256FromStrings(data...)
}

// ParseFilteredQueryParams parses filtered query params.
func ParseFilteredQueryParams(rawFilteredQueryParams string) []string {
	if pkgstrings.IsBlank(rawFilteredQueryParams) {
		return nil
	}

	return strings.Split(rawFilteredQueryParams, FilteredQueryParamsSeparator)
}

// FormatFilteredQueryParams formats a slice of strings into a filtered query params string.
func FormatFilteredQueryParams(params []string) string {
	return strings.Join(params, FilteredQueryParamsSeparator)
}

// TaskIDV2ByURLBased generates v2 version of task id by url based.
func TaskIDV2ByURLBased(url string, pieceLength *uint64, tag, application string, filteredQueryParams []string) string {
	url, err := neturl.FilterQueryParams(url, filteredQueryParams)
	if err != nil {
		url = ""
	}

	params := []string{url, tag, application}
	if pieceLength != nil {
		params = append(params, strconv.FormatUint(*pieceLength, 10))
	}

	params = append(params, commonv2.TaskType_STANDARD.String())
	return pkgdigest.SHA256FromStrings(params...)
}

// TaskIDV2ByContent generates v2 version of task id by content.
func TaskIDV2ByContent(content string) string {
	return pkgdigest.SHA256FromStrings(content)
}

// PersistentCacheTaskIDByContent generates persistent cache task id by content.
func PersistentCacheTaskIDByContent(content string) string {
	return pkgdigest.CRC32FromStrings(content)
}
