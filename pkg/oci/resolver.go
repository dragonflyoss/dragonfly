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
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/containerd/platforms"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/ocischema"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	registryclient "github.com/docker/distribution/registry/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Manifest is the manifest resolved from the image.
type Manifest struct {
	// URL is the manifest url referenced by the digest of the manifest.
	URL string

	// Platform is the platform that the manifest is resolved for in the format "os/arch".
	// It is empty when all the platforms are resolved and the image is not a multi-platform
	// image index.
	Platform string

	// BlobURLs is the blob urls (config and layers) of the manifest.
	BlobURLs []string
}

// ResolveManifests fetches and resolves container image manifests from a registry by the platform matcher.
// It constructs an HTTP request to retrieve the manifest, handles authentication via headers, and processes the response
// to return the manifests matched by the matcher with the manifest urls referenced by their digests and their blob urls,
// excluding the manifest list (image index) itself.
func (c *AuthClient) ResolveManifests(ctx context.Context, ref *Reference, header http.Header, matcher platforms.Matcher) ([]*Manifest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ref.ManifestURL(), nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	// Handle response.
	if resp.StatusCode == http.StatusNotModified {
		return nil, distribution.ErrManifestNotModified
	} else if !registryclient.SuccessStatus(resp.StatusCode) {
		return nil, registryclient.HandleErrorResponse(resp)
	}

	ctHeader := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	manifest, desc, err := distribution.UnmarshalManifest(ctHeader, body)
	if err != nil {
		return nil, err
	}

	switch v := manifest.(type) {
	case *schema1.SignedManifest, *schema2.DeserializedManifest, *ocischema.DeserializedManifest:
		digestRef := *ref
		digestRef.Reference = desc.Digest.String()

		// A blob referenced more than once by the manifest is resolved once.
		var blobURLs []string
		for _, blob := range v.References() {
			if blobURL := ref.BlobURL(blob.Digest.String()); !slices.Contains(blobURLs, blobURL) {
				blobURLs = append(blobURLs, blobURL)
			}
		}

		return []*Manifest{{URL: digestRef.ManifestURL(), BlobURLs: blobURLs}}, nil
	case *manifestlist.DeserializedManifestList:
		var manifests []*Manifest
		for _, desc := range v.Manifests {
			// Skip the manifest not matched by the platform matcher.
			if !matcher.Match(specs.Platform{
				Architecture: desc.Platform.Architecture,
				OS:           desc.Platform.OS,
				OSVersion:    desc.Platform.OSVersion,
				OSFeatures:   desc.Platform.OSFeatures,
				Variant:      desc.Platform.Variant,
			}) {
				continue
			}

			digestRef := *ref
			digestRef.Reference = desc.Digest.String()
			children, err := c.ResolveManifests(ctx, &digestRef, header.Clone(), matcher)
			if err != nil {
				return nil, err
			}

			// Set the platform of the image index entry to the manifests resolved from it.
			if desc.Platform.OS != "" && desc.Platform.Architecture != "" {
				platform := platforms.Normalize(specs.Platform{OS: desc.Platform.OS, Architecture: desc.Platform.Architecture})
				for _, child := range children {
					child.Platform = fmt.Sprintf("%s/%s", platform.OS, platform.Architecture)
				}
			}

			manifests = append(manifests, children...)
		}

		return manifests, nil
	}

	return nil, errors.New("unknown manifest type")
}
