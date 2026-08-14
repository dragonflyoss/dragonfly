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
	"fmt"
	"net/http"
	"time"

	"github.com/containerd/platforms"
	"github.com/distribution/reference"
)

// defaultRegistryTimeout is the default timeout for registry requests.
const defaultRegistryTimeout = 1 * time.Minute

// defaultHost maps the docker.io domain to its actual registry host,
// following containerd's remotes/docker resolver.
func defaultHost(domain string) string {
	if domain == "docker.io" {
		return "registry-1.docker.io"
	}

	return domain
}

// ParseImage parses an image reference (e.g., "docker.io/library/nginx:latest")
// into a registry reference. The reference is normalized with docker
// conventions, so "nginx:latest" resolves to "docker.io/library/nginx:latest".
func ParseImage(image string) (*Reference, error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, fmt.Errorf("invalid image reference: %w", err)
	}
	named = reference.TagNameOnly(named)

	var tag string
	switch ref := named.(type) {
	case reference.Digested:
		tag = ref.Digest().String()
	case reference.Tagged:
		tag = ref.Tag()
	default:
		return nil, fmt.Errorf("invalid image reference: %s", image)
	}

	return &Reference{
		Scheme:     "https",
		Registry:   defaultHost(reference.Domain(named)),
		Repository: reference.Path(named),
		Reference:  tag,
	}, nil
}

// resolveImageOptions holds the configurable settings for ResolveImage.
type resolveImageOptions struct {
	username string
	password string
	platform string
}

// ResolveImageOption configures ResolveImage.
type ResolveImageOption func(o *resolveImageOptions)

// WithAuth sets the username and password for registry authentication,
// anonymous access is used when they are empty.
func WithAuth(username, password string) ResolveImageOption {
	return func(o *resolveImageOptions) {
		o.username = username
		o.password = password
	}
}

// WithPlatform sets the target platform in the format "os/arch"
// (e.g., "linux/amd64"), the current platform is used when it is empty.
func WithPlatform(platform string) ResolveImageOption {
	return func(o *resolveImageOptions) {
		o.platform = platform
	}
}

// ResolveImage resolves the image manifest (including multi-platform image
// indexes) from the registry and returns the blob urls (config and layers)
// along with the authorization token for downloading them.
func ResolveImage(ctx context.Context, image string, opts ...ResolveImageOption) (blobURLs []string, token string, err error) {
	o := &resolveImageOptions{}
	for _, opt := range opts {
		opt(o)
	}

	ref, err := ParseImage(image)
	if err != nil {
		return nil, "", err
	}

	platform := platforms.DefaultSpec()
	if o.platform != "" {
		platform, err = platforms.Parse(o.platform)
		if err != nil {
			return nil, "", fmt.Errorf("invalid platform format %q, expected \"os/arch\" (e.g., \"linux/amd64\"): %w", o.platform, err)
		}
	}

	httpClient := &http.Client{
		Timeout:   defaultRegistryTimeout,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	client, err := NewAuthClient(ref, httpClient, o.username, o.password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to authenticate with registry: %w", err)
	}

	manifests, err := client.ResolveManifests(ctx, ref, make(http.Header), platform)
	if err != nil {
		return nil, "", fmt.Errorf("failed to pull image manifest: %w", err)
	}

	if len(manifests) == 0 {
		return nil, "", fmt.Errorf("no matching manifest for platform %s", platforms.Format(platform))
	}

	for _, manifest := range manifests {
		for _, desc := range manifest.References() {
			blobURLs = append(blobURLs, ref.BlobURL(desc.Digest.String()))
		}
	}

	return blobURLs, client.AuthToken(), nil
}
