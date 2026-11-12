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

//go:generate mockgen -destination mocks/job_mock.go -source image.go -package mocks

package job

import (
	"context"
	"crypto/x509"
	"slices"
	"time"

	"go.opentelemetry.io/otel/trace"

	commonv2 "d7y.io/api/v2/pkg/apis/common/v2"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	"d7y.io/dragonfly/v2/manager/config"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
	pkgoci "d7y.io/dragonfly/v2/pkg/oci"
)

type ManifestRequest struct {
	// URL is the image manifest url for preheating.
	URL string

	// PieceLength is the piece length(bytes) for downloading file. The value needs to
	// be greater than 4MiB (4,194,304 bytes) and less than 64MiB (67,108,864 bytes),
	// for example: 4194304(4mib), 8388608(8mib). If the piece length is not specified,
	// the piece length will be calculated according to the file size.
	PieceLength *uint64

	// Tag is the tag for preheating.
	Tag string

	// Application is the application string for preheating.
	Application string

	// FilteredQueryParams is the filtered query params for preheating.
	FilteredQueryParams string

	// Headers is the http headers for authentication.
	Headers map[string]string

	// Username is the username for authentication.
	Username string

	// Password is the password for authentication.
	Password string

	// The image type preheating task can specify the image architecture type. eg: linux/amd64.
	// If it is nil, all the platforms of the image are preheated. If it is empty, the default
	// platform is preheated.
	Platform *string

	// Scope is the scope for preheating, default is single_seed_peer.
	Scope string

	// IPs is a list of specific peer IPs for preheating.
	// This field has the highest priority: if provided, both 'Count' and 'Percentage' will be ignored.
	// Applies to 'all_peers' and 'all_seed_peers' scopes.
	IPs []string

	// Percentage is the percentage of available peers to preheat.
	// This field has the lowest priority and is only used if both 'IPs' and 'Count' are not provided.
	// It must be a value between 1 and 100 (inclusive) if provided.
	// Applies to 'all_peers' and 'all_seed_peers' scopes.
	Percentage *uint32

	// Count is the desired number of peers to preheat.
	// This field is used only when 'IPs' is not specified. It has priority over 'Percentage'.
	// It must be a value between 1 and 200 (inclusive) if provided.
	// Applies to 'all_peers' and 'all_seed_peers' scopes.
	Count *uint32

	// ConcurrentTaskCount specifies the maximum number of tasks (e.g., image layers) to preheat concurrently.
	// For example, if preheating 100 layers with ConcurrentTaskCount set to 10, up to 10 layers are processed simultaneously.
	// If ConcurrentPeerCount is 10 for 1000 peers, each layer is preheated by 10 peers concurrently.
	// Default is 8, maximum is 100.
	ConcurrentTaskCount int64

	// ConcurrentPeerCount specifies the maximum number of peers to preheat concurrently for a single task (e.g., an image layer).
	// For example, if preheating a layer with ConcurrentPeerCount set to 10, up to 10 peers process that layer simultaneously.
	// Default is 500, maximum is 1000.
	ConcurrentPeerCount int64

	// Timeout is the timeout for preheating, default is 30 minutes.
	Timeout time.Duration

	// RootCAs is the root CAs for TLS verification.
	RootCAs *x509.CertPool

	// InsecureSkipVerify indicates whether to skip TLS verification.
	InsecureSkipVerify bool

	// EnableTaskIDBasedBlobDigest indicates whether to use the blob digest for task ID calculation
	// when the blob url is an OCI blob url (e.g., /v2/<name>/blobs/sha256:<digest>).
	EnableTaskIDBasedBlobDigest bool
}

// Image implements the interface for handling container images.
type Image interface {
	// CreatePreheatRequestsByManifestURL generates a list of preheat requests for a container image
	// by fetching and parsing its manifest from a registry. It handles authentication, platform-specific
	// manifest filtering, and layer extraction for preheating. Each preheat request preheats the layers
	// of one platform of the image.
	CreatePreheatRequestsByManifestURL(ctx context.Context, req *ManifestRequest) ([]*PreheatRequest, error)

	// GetAuthorizationByManifestURL returns the authorization of the registry for a container image,
	// the bearer token or the basic credentials. It is empty if the registry needs no authorization.
	GetAuthorizationByManifestURL(ctx context.Context, req *ManifestRequest) (string, error)
}

// image is the implementation of the Image interface.
type image struct{}

// NewImage creates a new instance of the Image interface.
func NewImage() Image {
	return &image{}
}

// CreatePreheatRequestsByManifestURL generates a list of preheat requests for a container image
// by fetching and parsing its manifest from a registry. It handles authentication, platform-specific
// manifest filtering, and layer extraction for preheating. Each preheat request preheats the layers
// of one platform of the image.
func (i *image) CreatePreheatRequestsByManifestURL(ctx context.Context, req *ManifestRequest) ([]*PreheatRequest, error) {
	ctx, span := tracer.Start(ctx, config.SpanGetLayers, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	// Parse image manifest url.
	ref, err := pkgoci.ParseManifestURL(req.URL)
	if err != nil {
		return nil, err
	}

	// Resolve the manifests and the authorization token. Harbor uses the V1
	// preheat request and carries the auth info in the headers, which is used
	// as the issued token by the resolver. All the platforms of the image are
	// resolved if the platform is not specified.
	header := nethttp.MapToHeader(req.Headers)
	opts := []pkgoci.ResolveOption{
		pkgoci.WithAuth(req.Username, req.Password),
		pkgoci.WithHeader(header.Clone()),
	}

	if req.Platform != nil {
		opts = append(opts, pkgoci.WithPlatform(*req.Platform))
	}

	manifests, token, err := pkgoci.Resolve(ctx, ref, opts...)
	if err != nil {
		return nil, err
	}

	// Set authorization header, the registry may need no authorization.
	if token != "" {
		header.Set("Authorization", token)
	}

	var certificateChain [][]byte
	if req.RootCAs != nil {
		certificateChain = req.RootCAs.Subjects() //nolint:staticcheck
	}

	// Create one preheat request for each platform, so the peers record the layers by the
	// image and the platform to delete the image.
	var preheatRequests []*PreheatRequest
	for _, manifest := range manifests {
		if index := slices.IndexFunc(preheatRequests, func(preheatRequest *PreheatRequest) bool {
			return preheatRequest.Image.GetPlatform() == manifest.Platform
		}); index >= 0 {
			for _, url := range manifest.BlobURLs {
				if !slices.Contains(preheatRequests[index].URLs, url) {
					preheatRequests[index].URLs = append(preheatRequests[index].URLs, url)
				}
			}

			continue
		}

		preheatRequest := &PreheatRequest{
			URLs:                        manifest.BlobURLs,
			PieceLength:                 req.PieceLength,
			Tag:                         req.Tag,
			Application:                 req.Application,
			FilteredQueryParams:         req.FilteredQueryParams,
			Headers:                     nethttp.HeaderToMap(header),
			Scope:                       req.Scope,
			IPs:                         req.IPs,
			Percentage:                  req.Percentage,
			Count:                       req.Count,
			ConcurrentTaskCount:         req.ConcurrentTaskCount,
			ConcurrentPeerCount:         req.ConcurrentPeerCount,
			CertificateChain:            certificateChain,
			InsecureSkipVerify:          req.InsecureSkipVerify,
			Timeout:                     req.Timeout,
			EnableTaskIDBasedBlobDigest: req.EnableTaskIDBasedBlobDigest,
			Image:                       &commonv2.Image{Reference: ref.String()},
		}

		if manifest.Platform != "" {
			preheatRequest.Image.Platform = &manifest.Platform
		}

		preheatRequests = append(preheatRequests, preheatRequest)
	}

	if err := ref.Validate(); err != nil {
		logger.Warnf("preheat without the image, the image reference %s is invalid: %s", ref.String(), err.Error())
		for _, preheatRequest := range preheatRequests {
			preheatRequest.Image = nil
		}
	}

	return preheatRequests, nil
}

// GetAuthorizationByManifestURL returns the authorization of the registry for a container image,
// the bearer token or the basic credentials. It is empty if the registry needs no authorization.
func (i *image) GetAuthorizationByManifestURL(ctx context.Context, req *ManifestRequest) (string, error) {
	// Parse image manifest url.
	ref, err := pkgoci.ParseManifestURL(req.URL)
	if err != nil {
		return "", err
	}

	// Harbor uses the V1 preheat request and carries the auth info in the headers,
	// which is used as the issued token.
	var opts []pkgoci.Option
	if issuedToken := nethttp.MapToHeader(req.Headers).Get("Authorization"); issuedToken != "" {
		opts = append(opts, pkgoci.WithIssuedToken(issuedToken))
	}

	client, err := pkgoci.NewAuthClient(ref, pkgoci.DefaultHTTPClient(), req.Username, req.Password, opts...)
	if err != nil {
		return "", err
	}

	return client.Authorization()
}
