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
	"testing"

	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/pkg/idgen"
)

func TestPreheatTaskID(t *testing.T) {
	pieceLength := uint64(1024)

	tests := []struct {
		name                string
		url                 string
		pieceLength         *uint64
		tag                 string
		application         string
		filteredQueryParams []string
		expectIsBlobURL     bool
		expectTaskID        func(t *testing.T, taskID string)
	}{
		{
			name:                "blob URL uses blob digest-based task ID",
			url:                 "https://registry.example.com/v2/library/nginx/blobs/sha256:c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: nil,
			expectIsBlobURL:     true,
			expectTaskID: func(t *testing.T, taskID string) {
				assert.Equal(t, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4", taskID)
			},
		},
		{
			name:                "blob URL with query params uses blob digest-based task ID",
			url:                 "https://registry.example.com/v2/library/nginx/blobs/sha256:c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4?foo=bar",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: nil,
			expectIsBlobURL:     true,
			expectTaskID: func(t *testing.T, taskID string) {
				assert.Equal(t, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4", taskID)
			},
		},
		{
			name:                "blob URL with nested repository uses blob digest-based task ID",
			url:                 "https://registry.example.com/v2/org/team/nginx/blobs/sha256:c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: nil,
			expectIsBlobURL:     true,
			expectTaskID: func(t *testing.T, taskID string) {
				assert.Equal(t, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4", taskID)
			},
		},
		{
			name:                "same blob digest from different registries produces same task ID",
			url:                 "https://another-registry.example.com/v2/library/nginx/blobs/sha256:c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: nil,
			expectIsBlobURL:     true,
			expectTaskID: func(t *testing.T, taskID string) {
				assert.Equal(t, "c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4", taskID)
			},
		},
		{
			name:                "non-blob URL uses URL-based task ID",
			url:                 "https://example.com/file.tar.gz",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: nil,
			expectIsBlobURL:     false,
			expectTaskID: func(t *testing.T, taskID string) {
				expected := idgen.TaskIDV2ByURLBased("https://example.com/file.tar.gz", &pieceLength, "foo", "bar", nil, "")
				assert.Equal(t, expected, taskID)
			},
		},
		{
			name:                "non-blob URL with filtered query params uses URL-based task ID",
			url:                 "https://example.com/file.tar.gz?foo=foo&bar=bar",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: []string{"foo", "bar"},
			expectIsBlobURL:     false,
			expectTaskID: func(t *testing.T, taskID string) {
				expected := idgen.TaskIDV2ByURLBased("https://example.com/file.tar.gz?foo=foo&bar=bar", &pieceLength, "foo", "bar", []string{"foo", "bar"}, "")
				assert.Equal(t, expected, taskID)
			},
		},
		{
			name:                "manifest URL is not a blob URL",
			url:                 "https://registry.example.com/v2/library/nginx/manifests/latest",
			pieceLength:         &pieceLength,
			tag:                 "foo",
			application:         "bar",
			filteredQueryParams: nil,
			expectIsBlobURL:     false,
			expectTaskID: func(t *testing.T, taskID string) {
				expected := idgen.TaskIDV2ByURLBased("https://registry.example.com/v2/library/nginx/manifests/latest", &pieceLength, "foo", "bar", nil, "")
				assert.Equal(t, expected, taskID)
			},
		},
		{
			name:                "empty URL is not a blob URL",
			url:                 "",
			pieceLength:         &pieceLength,
			tag:                 "",
			application:         "",
			filteredQueryParams: nil,
			expectIsBlobURL:     false,
			expectTaskID: func(t *testing.T, taskID string) {
				expected := idgen.TaskIDV2ByURLBased("", &pieceLength, "", "", nil, "")
				assert.Equal(t, expected, taskID)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			taskID, isBlobURL := preheatTaskID(tc.url, tc.pieceLength, tc.tag, tc.application, tc.filteredQueryParams)
			assert.Equal(t, tc.expectIsBlobURL, isBlobURL)
			tc.expectTaskID(t, taskID)
		})
	}
}
