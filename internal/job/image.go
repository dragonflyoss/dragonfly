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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/manifestlist"
	"github.com/docker/distribution/manifest/ocischema"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/manifest/schema2"
	registryclient "github.com/docker/distribution/registry/client"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	typesregistry "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/registry"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"go.opentelemetry.io/otel/trace"

	"d7y.io/dragonfly/v2/manager/config"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
)

// defaultRegistryTimeout is the default timeout for registry requests.
const defaultRegistryTimeout = 1 * time.Minute

// defaultHTTPTransport is the default http transport.
var defaultHTTPTransport = &http.Transport{
	MaxIdleConns:        400,
	MaxIdleConnsPerHost: 20,
	MaxConnsPerHost:     50,
	IdleConnTimeout:     120 * time.Second,
}

// accessURLRegexp is the regular expression for parsing access url.
var accessURLRegexp, _ = regexp.Compile("^(.*)://(.*)/v2/(.*)/manifests/(.*)")

// preheatImage is image information for preheat.
type preheatImage struct {
	protocol string
	domain   string
	name     string
	tag      string
}

func (p *preheatImage) manifestURL() string {
	return fmt.Sprintf("%s://%s/v2/%s/manifests/%s", p.protocol, p.domain, p.name, p.tag)
}

func (p *preheatImage) blobsURL(digest string) string {
	return fmt.Sprintf("%s://%s/v2/%s/blobs/%s", p.protocol, p.domain, p.name, digest)
}

// parseManifestURL parses a container image manifest URL into its components.
// It extracts the protocol, domain, image name, and tag from the provided URL
// using a regular expression.
func parseManifestURL(url string) (*preheatImage, error) {
	r := accessURLRegexp.FindStringSubmatch(url)
	if len(r) != 5 {
		return nil, errors.New("parse access url failed")
	}

	return &preheatImage{
		protocol: r[1],
		domain:   r[2],
		name:     r[3],
		tag:      r[4],
	}, nil
}

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
	Platform string

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
}

// Image implements the interface for handling container images.
type Image interface {
	// CreatePreheatRequestsByManifestURL generates a list of preheat requests for a container image
	// by fetching and parsing its manifest from a registry. It handles authentication, platform-specific
	// manifest filtering, and layer extraction for preheating.
	CreatePreheatRequestsByManifestURL(ctx context.Context, req *ManifestRequest) ([]*PreheatRequest, error)
}

// image is the implementation of the Image interface.
type image struct{}

// NewImage creates a new instance of the Image interface.
func NewImage() Image {
	return &image{}
}

// CreatePreheatRequestsByManifestURL generates a list of preheat requests for a container image
// by fetching and parsing its manifest from a registry. It handles authentication, platform-specific
// manifest filtering, and layer extraction for preheating.
func (i *image) CreatePreheatRequestsByManifestURL(ctx context.Context, req *ManifestRequest) ([]*PreheatRequest, error) {
	ctx, span := tracer.Start(ctx, config.SpanGetLayers, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	// Parse image manifest url.
	image, err := parseManifestURL(req.URL)
	if err != nil {
		return nil, err
	}

	// Background:
	// Harbor uses the V1 preheat request and will carry the auth info in the headers.
	options := []imageAuthClientOption{}
	header := nethttp.MapToHeader(req.Headers)
	if token := header.Get("Authorization"); len(token) > 0 {
		options = append(options, withIssuedToken(token))
		header.Set("Authorization", token)
	}

	httpClient := &http.Client{
		Timeout: defaultRegistryTimeout,
		Transport: &http.Transport{
			DialContext:         nethttp.NewSafeDialer().DialContext,
			TLSClientConfig:     &tls.Config{RootCAs: req.RootCAs, InsecureSkipVerify: req.InsecureSkipVerify},
			MaxIdleConns:        defaultHTTPTransport.MaxIdleConns,
			MaxIdleConnsPerHost: defaultHTTPTransport.MaxIdleConnsPerHost,
			MaxConnsPerHost:     defaultHTTPTransport.MaxConnsPerHost,
			IdleConnTimeout:     defaultHTTPTransport.IdleConnTimeout,
		},
	}

	// Init docker auth client.
	client, err := newImageAuthClient(image, httpClient, &typesregistry.AuthConfig{Username: req.Username, Password: req.Password}, options...)
	if err != nil {
		return nil, err
	}

	// Get platform.
	platform := platforms.DefaultSpec()
	if req.Platform != "" {
		platform, err = platforms.Parse(req.Platform)
		if err != nil {
			return nil, err
		}
	}

	// Resolve manifests for the image.
	manifests, err := resolveManifests(ctx, client, image, header.Clone(), platform)
	if err != nil {
		return nil, err
	}

	// No matching manifest for platform in the manifest list entries
	if len(manifests) == 0 {
		return nil, fmt.Errorf("no matching manifest for platform %s", platforms.Format(platform))
	}

	// Set authorization header
	header.Set("Authorization", client.GetAuthToken())

	// Build preheat requests for container image layers from the provided manifests.
	preheatRequest, err := buildPreheatRequestFromManifests(manifests, req, header.Clone(), image)
	if err != nil {
		return nil, err
	}

	return preheatRequest, nil
}

// resolveManifests fetches and resolves container image manifests from a registry for a specified platform.
// It constructs an HTTP request to retrieve the manifest, handles authentication via headers, and processes the response
// to return manifests matching the given platform. Supports single manifests and manifest lists.
func resolveManifests(ctx context.Context, client *imageAuthClient, image *preheatImage, header http.Header, platform specs.Platform) ([]distribution.Manifest, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, image.manifestURL(), nil)
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

	resp, err := client.Do(req)
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

	// Unmarshal manifest.
	manifest, _, err := distribution.UnmarshalManifest(ctHeader, body)
	if err != nil {
		return nil, err
	}

	switch v := manifest.(type) {
	case *schema1.SignedManifest, *schema2.DeserializedManifest, *ocischema.DeserializedManifest:
		return []distribution.Manifest{v}, nil
	case *manifestlist.DeserializedManifestList:
		var result []distribution.Manifest
		for _, v := range filterManifests(v.Manifests, platform) {
			image.tag = v.Digest.String()
			manifests, err := resolveManifests(ctx, client, image, header.Clone(), platform)
			if err != nil {
				return nil, err
			}

			result = append(result, manifests...)
		}

		return result, nil
	}

	return nil, errors.New("unknown manifest type")
}

// filterManifests filters a list of manifest descriptors to return only those matching the specified platform.
// It compares the architecture and operating system of each manifest descriptor against the target platform.
func filterManifests(manifests []manifestlist.ManifestDescriptor, platform specs.Platform) []manifestlist.ManifestDescriptor {
	var matches []manifestlist.ManifestDescriptor
	for _, desc := range manifests {
		if desc.Platform.Architecture == platform.Architecture && desc.Platform.OS == platform.OS {
			matches = append(matches, desc)
		}
	}

	return matches
}

// buildPreheatRequestFromManifests constructs preheat requests for container image layers from
// the provided manifests. It extracts layer URLs from the manifests and builds a PreheatRequest
// using the specified arguments, headers, and TLS settings.
func buildPreheatRequestFromManifests(manifests []distribution.Manifest, req *ManifestRequest, header http.Header, image *preheatImage) ([]*PreheatRequest, error) {
	var certificateChain [][]byte
	if req.RootCAs != nil {
		certificateChain = req.RootCAs.Subjects()
	}

	var layerURLs []string
	for _, m := range manifests {
		for _, v := range m.References() {
			header.Set("Accept", v.MediaType)
			layerURLs = append(layerURLs, image.blobsURL(v.Digest.String()))
		}
	}

	layers := &PreheatRequest{
		URLs:                layerURLs,
		PieceLength:         req.PieceLength,
		Tag:                 req.Tag,
		Application:         req.Application,
		FilteredQueryParams: req.FilteredQueryParams,
		Headers:             nethttp.HeaderToMap(header),
		Scope:               req.Scope,
		IPs:                 req.IPs,
		Percentage:          req.Percentage,
		Count:               req.Count,
		ConcurrentTaskCount: req.ConcurrentTaskCount,
		ConcurrentPeerCount: req.ConcurrentPeerCount,
		CertificateChain:    certificateChain,
		InsecureSkipVerify:  req.InsecureSkipVerify,
		Timeout:             req.Timeout,
	}

	return []*PreheatRequest{layers}, nil
}

// imageAuthClientOption is an option for imageAuthClient.
type imageAuthClientOption func(*imageAuthClient)

// withIssuedToken sets the issuedToken for imageAuthClient.
func withIssuedToken(token string) imageAuthClientOption {
	return func(c *imageAuthClient) {
		c.issuedToken = token
	}
}

// imageAuthClient is a client for image authentication.
type imageAuthClient struct {
	// issuedToken is the issued token specified in header from user request,
	// there is no need to go through v2 authentication to get a new token
	// if the token is not empty, just use this token to access v2 API directly.
	issuedToken string

	// httpClient is the http client.
	httpClient *http.Client

	// authConfig is the auth config.
	authConfig *typesregistry.AuthConfig

	// interceptorTokenHandler is the token interceptor.
	interceptorTokenHandler *interceptorTokenHandler
}

// newImageAuthClient creates a new imageAuthClient.
func newImageAuthClient(image *preheatImage, httpClient *http.Client, authConfig *typesregistry.AuthConfig, opts ...imageAuthClientOption) (*imageAuthClient, error) {
	d := &imageAuthClient{
		httpClient:              httpClient,
		authConfig:              authConfig,
		interceptorTokenHandler: newInterceptorTokenHandler(),
	}

	for _, opt := range opts {
		opt(d)
	}

	if len(d.issuedToken) > 0 {
		return d, nil
	}

	// New a challenge manager for the supported authentication types.
	challengeManager, err := registry.PingV2Registry(&url.URL{Scheme: image.protocol, Host: image.domain}, d.httpClient.Transport)
	if err != nil {
		return nil, err
	}

	// New a credential store which always returns the same credential values.
	creds := registry.NewStaticCredentialStore(d.authConfig)

	// Transport with authentication.
	d.httpClient.Transport = transport.NewTransport(
		d.httpClient.Transport,
		auth.NewAuthorizer(
			challengeManager,
			auth.NewTokenHandlerWithOptions(auth.TokenHandlerOptions{
				Transport:   d.httpClient.Transport,
				Credentials: creds,
				Scopes: []auth.Scope{auth.RepositoryScope{
					Repository: image.name,
					Actions:    []string{"pull"},
				}},
				ClientID: registry.AuthClientID,
			}),
			d.interceptorTokenHandler,
			auth.NewBasicHandler(creds),
		),
	)

	return d, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (d *imageAuthClient) Do(req *http.Request) (*http.Response, error) {
	return d.httpClient.Do(req)
}

// GetAuthToken returns the bearer token.
func (d *imageAuthClient) GetAuthToken() string {
	if len(d.issuedToken) > 0 {
		return d.issuedToken
	}

	return d.interceptorTokenHandler.GetAuthToken()
}

// interceptorTokenHandler is a token interceptor intercept bearer token from auth handler.
type interceptorTokenHandler struct {
	auth.AuthenticationHandler
	token string
}

// NewInterceptorTokenHandler returns a new InterceptorTokenHandler.
func newInterceptorTokenHandler() *interceptorTokenHandler {
	return &interceptorTokenHandler{}
}

// Scheme returns the authentication scheme.
func (h *interceptorTokenHandler) Scheme() string {
	return "bearer"
}

// AuthorizeRequest sets the Authorization header on the request.
func (h *interceptorTokenHandler) AuthorizeRequest(req *http.Request, params map[string]string) error {
	h.token = req.Header.Get("Authorization")
	return nil
}

// GetAuthToken returns the bearer token.
func (h *interceptorTokenHandler) GetAuthToken() string {
	return h.token
}
