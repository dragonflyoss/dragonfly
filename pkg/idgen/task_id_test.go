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
	"testing"

	"github.com/stretchr/testify/assert"

	commonv1 "d7y.io/api/v2/pkg/apis/common/v1"
)

func TestTaskIDV1(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		meta   *commonv1.UrlMeta
		expect func(t *testing.T, d any)
	}{
		{
			name: "generate taskID with url",
			url:  "https://example.com",
			meta: nil,
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "100680ad546ce6a577f42f52df33b4cfdca756859e664b8d7de329b150d09ce9")
			},
		},
		{
			name: "generate taskID with meta",
			url:  "https://example.com",
			meta: &commonv1.UrlMeta{
				Range:  "foo",
				Digest: "bar",
				Tag:    "",
			},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "aeee0e0a2a0c75130582641353c539aaf9011a0088b31347f7588e70e449a3e0")
			},
		},
		{
			name: "generate taskID with filter",
			url:  "https://example.com?foo=foo&bar=bar",
			meta: &commonv1.UrlMeta{
				Tag:    "foo",
				Filter: "foo&bar",
			},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "2773851c628744fb7933003195db436ce397c1722920696c4274ff804d86920b")
			},
		},
		{
			name: "generate taskID with tag",
			url:  "https://example.com",
			meta: &commonv1.UrlMeta{
				Tag: "foo",
			},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "2773851c628744fb7933003195db436ce397c1722920696c4274ff804d86920b")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, TaskIDV1(tc.url, tc.meta))
		})
	}
}

func TestTaskIDV2ByURLBased(t *testing.T) {
	pieceLength := uint64(1024)

	tests := []struct {
		name        string
		url         string
		pieceLength *uint64
		tag         string
		application string
		filters     []string
		revision    string
		expect      func(t *testing.T, d any)
	}{
		{
			name:        "generate taskID",
			url:         "https://example.com",
			pieceLength: &pieceLength,
			tag:         "foo",
			application: "bar",
			filters:     []string{},
			revision:    "v1.0",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "5844f27a257287e9b734256bb25603d8005422ced8c0377f15063ec11963b25f")
			},
		},
		{
			name:        "generate taskID with tag and application",
			url:         "https://example.com",
			tag:         "foo",
			application: "bar",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "06408fbf247ddaca478f8cb9565fe5591c28efd0994b8fea80a6a87d3203c5ca")
			},
		},
		{
			name: "generate taskID with tag",
			url:  "https://example.com",
			tag:  "foo",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "3c3f230ef9f191dd2821510346a7bc138e4894bee9aee184ba250a3040701d2a")
			},
		},
		{
			name:        "generate taskID with application",
			url:         "https://example.com",
			application: "bar",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "c9f9261b7305c24371244f9f149f5d4589ed601348fdf22d7f6f4b10658fdba2")
			},
		},
		{
			name:        "generate taskID with pieceLength",
			url:         "https://example.com",
			pieceLength: &pieceLength,
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "9f7c9aafbc6f30f8f41a96ca77eeae80c5b60964b3034b0ee43ccf7b2f9e52b8")
			},
		},
		{
			name:    "generate taskID with filters",
			url:     "https://example.com?foo=foo&bar=bar",
			filters: []string{"foo", "bar"},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "457b4328cde278e422c9e243f7bfd1e97f511fec43a80f535cf6b0ef6b086776")
			},
		},
		{
			name:     "generate taskID with revision",
			url:      "https://example.com",
			revision: "v1.0",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "b171331534b80e0bf91da38ebbfcdbf4d177898f4b9beac44f14733e3f004d4e")
			},
		},
		{
			name:        "generate taskID with sorted query params",
			url:         "https://example.com/file.txt?z=9&b=2&a=1",
			tag:         "foo",
			application: "bar",
			filters:     []string{"z"},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "8b3f6e9b9b8fe20903bced565cfd1d0aaef354a4c17573f0c2c1979210443f9d")
			},
		},
		{
			name:    "generate taskID with same key query params keeping order",
			url:     "https://example.com/file.txt?b=2&a=1&b=1",
			filters: []string{"c"},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "7c8801d0596be5e8f9449d5c4af23866c72fe5205119c0e5912981f3b16a37aa")
			},
		},
		{
			name:        "generate taskID with escaped query params",
			url:         "https://example.com/file.txt?k=a b&m=x*y&n=c~d",
			pieceLength: &pieceLength,
			filters:     []string{"none"},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "6196a6846023f6d3c1e4d30f6c86f3d4186e4c664a33e5692b0e04e49b26a9af")
			},
		},
		{
			name:    "generate taskID with all query params filtered",
			url:     "https://example.com/file.txt?a=1&b=2",
			tag:     "foo",
			filters: []string{"a", "b"},
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "c8f4b41117329d54af920010394f6f607bac707e933ab2f18d372e3dd4c7fcb3")
			},
		},
		{
			name: "generate taskID with raw url when no filters",
			url:  "https://example.com/file.txt?b=2&a=1",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "980ee327518ccc5a7c30703e1a2232e8ba9047b39431f940636c85b6146f8b9a")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, TaskIDV2ByURLBased(tc.url, tc.pieceLength, tc.tag, tc.application, tc.filters, tc.revision))
		})
	}
}

func TestTaskIDV2ByContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		expect  func(t *testing.T, d any)
	}{
		{
			name:    "generate taskID",
			content: "This is a test file",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "e2d0fe1585a63ec6009c8016ff8dda8b17719a637405a4e23c0ff81339148249")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, TaskIDV2ByContent(tc.content))
		})
	}
}

func TestPersistentCacheTaskIDbyContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		expect  func(t *testing.T, d any)
	}{
		{
			name:    "generate persistentCacheTaskID",
			content: "This is a test file",
			expect: func(t *testing.T, d any) {
				assert := assert.New(t)
				assert.Equal(d, "e2d0fe1585a63ec6009c8016ff8dda8b17719a637405a4e23c0ff81339148249")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, PersistentCacheTaskIDByContent(tc.content))
		})
	}
}

func TestIsBlobURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		expect func(t *testing.T, ok bool)
	}{
		{
			name: "http blob url",
			url:  "http://registry.example.com/v2/library/ubuntu/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name: "blob url with nested repository",
			url:  "https://registry.io/v2/org/team/project/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name: "blob url with port and query params",
			url:  "http://localhost:5000/v2/myrepo/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e?ns=docker.io",
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name: "url without blobs path",
			url:  "https://registry.example.com/v2/library/ubuntu/manifests/latest",
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
		{
			name: "plain url",
			url:  "https://example.com/file.txt",
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, IsBlobURL(tc.url))
		})
	}
}

func TestTaskIDV2ByBlobDigest(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		expect func(t *testing.T, d string, err error)
	}{
		{
			name: "generate taskID by sha256 blob digest",
			url:  "http://registry.example.com/v2/library/ubuntu/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, "b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e")
			},
		},
		{
			name: "generate taskID by sha512 blob digest",
			url:  "https://registry.example.com/v2/myorg/myrepo/blobs/sha512:94381a28e8c039fedfa78de025158a068226c3ccd041b22c2c8e73fc993584e9b167d9ae32bc8b372c66701c808ab134e0768c8f16b9a3e61eec1ccf8faa9db8",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, "94381a28e8c039fedfa78de025158a068226c3ccd041b22c2c8e73fc993584e9b167d9ae32bc8b372c66701c808ab134e0768c8f16b9a3e61eec1ccf8faa9db8")
			},
		},
		{
			name: "generate taskID by blob digest with query params",
			url:  "http://localhost:5000/v2/myrepo/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e?ns=docker.io",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, "b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e")
			},
		},
		{
			name: "not a blob url",
			url:  "https://example.com/file.txt",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "invalid digest length",
			url:  "http://registry.example.com/v2/library/ubuntu/blobs/sha256:abc",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "unsupported digest algorithm",
			url:  "http://registry.example.com/v2/library/ubuntu/blobs/md5:8a04994a666b4e4b20a2fd9e5a44f44c",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := TaskIDV2ByBlobDigest(tc.url)
			tc.expect(t, d, err)
		})
	}
}

func TestTaskIDV2(t *testing.T) {
	tests := []struct {
		name                        string
		url                         string
		content                     string
		enableTaskIDBasedBlobDigest bool
		expect                      func(t *testing.T, d string, err error)
	}{
		{
			name:                        "generate taskID by content",
			url:                         "http://registry.example.com/v2/library/ubuntu/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			content:                     "This is a content",
			enableTaskIDBasedBlobDigest: true,
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, TaskIDV2ByContent("This is a content"))
			},
		},
		{
			name:                        "generate taskID by blob digest",
			url:                         "http://registry.example.com/v2/library/ubuntu/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			enableTaskIDBasedBlobDigest: true,
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, "b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e")
			},
		},
		{
			name: "generate taskID by url based when blob digest is disabled",
			url:  "http://registry.example.com/v2/library/ubuntu/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e",
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, TaskIDV2ByURLBased("http://registry.example.com/v2/library/ubuntu/blobs/sha256:b2c366cce7e68013d5441c6326d5a3e1b12aeb5ed58564d0fd3fa089bc29cb6e", nil, "foo", "bar", nil, ""))
			},
		},
		{
			name:                        "generate taskID by url based for non blob url",
			url:                         "https://example.com/file.txt",
			enableTaskIDBasedBlobDigest: true,
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(d, TaskIDV2ByURLBased("https://example.com/file.txt", nil, "foo", "bar", nil, ""))
			},
		},
		{
			name:                        "unsupported digest algorithm",
			url:                         "http://registry.example.com/v2/library/ubuntu/blobs/md5:8a04994a666b4e4b20a2fd9e5a44f44c",
			enableTaskIDBasedBlobDigest: true,
			expect: func(t *testing.T, d string, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, err := TaskIDV2(tc.url, nil, "foo", "bar", nil, tc.content, tc.enableTaskIDBasedBlobDigest)
			tc.expect(t, d, err)
		})
	}
}
