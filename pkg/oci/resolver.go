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
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/containerd/platforms"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/ocischema"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	registryclient "github.com/docker/distribution/registry/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// ResolveManifests fetches and resolves container image manifests from a registry for a specified platform.
// It constructs an HTTP request to retrieve the manifest, handles authentication via headers, and processes the response
// to return manifests matching the given platform, along with the manifest urls referenced by the digests of the
// matched platform manifests, excluding the manifest list (image index) itself.
func (c *AuthClient) ResolveManifests(ctx context.Context, ref *Reference, header http.Header, platform specs.Platform) ([]distribution.Manifest, []string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ref.ManifestURL(), nil)
	if err != nil {
		return nil, nil, err
	}

	// Set header from the user request.
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Set accept header with media types.
	for _, mediaType := range distribution.ManifestMediaTypes() {
		req.Header.Add("Accept", mediaType)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Handle response.
	if resp.StatusCode == http.StatusNotModified {
		return nil, nil, distribution.ErrManifestNotModified
	} else if !registryclient.SuccessStatus(resp.StatusCode) {
		return nil, nil, registryclient.HandleErrorResponse(resp)
	}

	ctHeader := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	manifest, desc, err := distribution.UnmarshalManifest(ctHeader, body)
	if err != nil {
		return nil, nil, err
	}

	switch v := manifest.(type) {
	case *schema1.SignedManifest, *schema2.DeserializedManifest, *ocischema.DeserializedManifest:
		digestRef := *ref
		digestRef.Reference = desc.Digest.String()
		return []distribution.Manifest{v}, []string{digestRef.ManifestURL()}, nil
	case *manifestlist.DeserializedManifestList:
		var result []distribution.Manifest
		var manifestURLs []string
		for _, desc := range filterManifests(v.Manifests, platform) {
			ref.Reference = desc.Digest.String()
			manifests, childManifestURLs, err := c.ResolveManifests(ctx, ref, header.Clone(), platform)
			if err != nil {
				return nil, nil, err
			}

			result = append(result, manifests...)
			manifestURLs = append(manifestURLs, childManifestURLs...)
		}

		return result, manifestURLs, nil
	}

	return nil, nil, errors.New("unknown manifest type")
}

// filterManifests filters a list of manifest descriptors to return only those
// matching the specified platform, using the containerd platform matcher.
func filterManifests(manifests []manifestlist.ManifestDescriptor, platform specs.Platform) []manifestlist.ManifestDescriptor {
	matcher := platforms.Only(platform)
	var matches []manifestlist.ManifestDescriptor
	for _, desc := range manifests {
		if matcher.Match(specs.Platform{
			Architecture: desc.Platform.Architecture,
			OS:           desc.Platform.OS,
			OSVersion:    desc.Platform.OSVersion,
			OSFeatures:   desc.Platform.OSFeatures,
			Variant:      desc.Platform.Variant,
		}) {
			matches = append(matches, desc)
		}
	}

	return matches
}
