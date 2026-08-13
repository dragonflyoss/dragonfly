/*
 *     Copyright 2026 The Dragonfly Authors
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

package oci

import (
	"testing"

	"github.com/docker/distribution/manifest/manifestlist"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/stretchr/testify/assert"
)

func TestFilterManifests(t *testing.T) {
	manifests := []manifestlist.ManifestDescriptor{
		{Platform: manifestlist.PlatformSpec{OS: "linux", Architecture: "amd64"}},
		{Platform: manifestlist.PlatformSpec{OS: "linux", Architecture: "arm64", Variant: "v8"}},
		{Platform: manifestlist.PlatformSpec{OS: "windows", Architecture: "amd64"}},
	}

	tests := []struct {
		name     string
		platform specs.Platform
		expect   func(t *testing.T, matches []manifestlist.ManifestDescriptor)
	}{
		{
			name:     "match linux/amd64",
			platform: specs.Platform{OS: "linux", Architecture: "amd64"},
			expect: func(t *testing.T, matches []manifestlist.ManifestDescriptor) {
				assert := assert.New(t)
				assert.Len(matches, 1)
				assert.Equal("amd64", matches[0].Platform.Architecture)
				assert.Equal("linux", matches[0].Platform.OS)
			},
		},
		{
			name:     "match linux/arm64 with variant",
			platform: specs.Platform{OS: "linux", Architecture: "arm64"},
			expect: func(t *testing.T, matches []manifestlist.ManifestDescriptor) {
				assert := assert.New(t)
				assert.Len(matches, 1)
				assert.Equal("arm64", matches[0].Platform.Architecture)
			},
		},
		{
			name:     "no match",
			platform: specs.Platform{OS: "linux", Architecture: "riscv64"},
			expect: func(t *testing.T, matches []manifestlist.ManifestDescriptor) {
				assert := assert.New(t)
				assert.Empty(matches)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, filterManifests(manifests, tc.platform))
		})
	}
}
