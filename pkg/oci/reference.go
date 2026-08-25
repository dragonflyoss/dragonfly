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

// Package oci resolves container image manifests from a registry with v2
// authentication, used for preheating images.
package oci

import (
	"errors"
	"fmt"
	"regexp"
)

// manifestURLRegexp is the regular expression for parsing manifest urls.
var manifestURLRegexp = regexp.MustCompile("^(.*)://(.*)/v2/(.*)/manifests/(.*)")

// Reference is the set of registry coordinates of an image, following the
// naming of the OCI distribution spec.
type Reference struct {
	// Scheme is the registry scheme (http or https).
	Scheme string

	// Registry is the registry host.
	Registry string

	// Repository is the repository name.
	Repository string

	// Reference is the image tag or digest.
	Reference string
}

// ManifestURL returns the manifest url of the reference.
func (r *Reference) ManifestURL() string {
	return fmt.Sprintf("%s://%s/v2/%s/manifests/%s", r.Scheme, r.Registry, r.Repository, r.Reference)
}

// BlobURL returns the blob url of the reference for the given digest.
func (r *Reference) BlobURL(digest string) string {
	return fmt.Sprintf("%s://%s/v2/%s/blobs/%s", r.Scheme, r.Registry, r.Repository, digest)
}

// ParseManifestURL parses a container image manifest URL into a reference.
// It extracts the scheme, registry, repository, and tag or digest from the
// provided URL.
func ParseManifestURL(url string) (*Reference, error) {
	matches := manifestURLRegexp.FindStringSubmatch(url)
	if len(matches) != 5 {
		return nil, errors.New("parse access url failed")
	}

	return &Reference{
		Scheme:     matches[1],
		Registry:   matches[2],
		Repository: matches[3],
		Reference:  matches[4],
	}, nil
}
