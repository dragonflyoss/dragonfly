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

//go:generate mockgen -destination mocks/preheat_mock.go -source preheat.go -package mocks

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
	machineryv1tasks "github.com/dragonflyoss/machinery/v1/tasks"
	"github.com/google/uuid"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"go.opentelemetry.io/otel/trace"

	logger "d7y.io/dragonfly/v2/internal/dflog"
	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/config"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
)

// preheatImage is an image for preheat.
type PreheatType string

// String returns the string representation of PreheatType.
func (p PreheatType) String() string {
	return string(p)
}

const (
	// PreheatImageType is image type of preheat job.
	PreheatImageType PreheatType = "image"

	// PreheatFileType is file type of preheat job.
	PreheatFileType PreheatType = "file"
)

// defaultHTTPTransport is the default http transport.
var defaultHTTPTransport = &http.Transport{
	MaxIdleConns:        400,
	MaxIdleConnsPerHost: 20,
	MaxConnsPerHost:     50,
	IdleConnTimeout:     120 * time.Second,
}

// Preheat is an interface for preheat job.
type Preheat interface {
	// CreatePreheat creates a preheat job.
	CreatePreheat(context.Context, []models.Scheduler, types.PreheatArgs) (*internaljob.GroupJobState, error)
}

// preheat is an implementation of Preheat.
type preheat struct {
	job                *internaljob.Job
	rootCAs            *x509.CertPool
	insecureSkipVerify bool
	registryTimeout    time.Duration
}

// newPreheat creates a new Preheat.
func newPreheat(job *internaljob.Job, registryTimeout time.Duration, rootCAs *x509.CertPool, insecureSkipVerify bool) Preheat {
	return &preheat{
		job:                job,
		rootCAs:            rootCAs,
		insecureSkipVerify: insecureSkipVerify,
		registryTimeout:    registryTimeout,
	}
}

// CreatePreheat creates a preheat job.
func (p *preheat) CreatePreheat(ctx context.Context, schedulers []models.Scheduler, json types.PreheatArgs) (*internaljob.GroupJobState, error) {
	var span trace.Span
	ctx, span = tracer.Start(ctx, config.SpanPreheat, trace.WithSpanKind(trace.SpanKindProducer))
	span.SetAttributes(config.AttributePreheatType.String(json.Type))
	span.SetAttributes(config.AttributePreheatURL.String(json.URL))
	defer span.End()

	// Generate download files.
	var files []internaljob.PreheatRequest
	var err error
	switch PreheatType(json.Type) {
	case PreheatImageType:
		// Image preheat only supports to preheat single image.
		if json.URL == "" {
			return nil, errors.New("invalid params: url is required")
		}

		files, err = GetImageLayers(ctx, json, p.registryTimeout, p.rootCAs, p.insecureSkipVerify)
		if err != nil {
			return nil, err
		}
	case PreheatFileType:
		// File preheat supports to preheat multiple files and single file.
		if json.URL == "" && len(json.URLs) == 0 {
			return nil, errors.New("invalid params: url or urls is required")
		}

		var certificateChain [][]byte
		if p.rootCAs != nil {
			certificateChain = p.rootCAs.Subjects()
		}

		files = append(files, internaljob.PreheatRequest{
			URLs:                append(json.URLs, json.URL),
			PieceLength:         json.PieceLength,
			Tag:                 json.Tag,
			Application:         json.Application,
			FilteredQueryParams: json.FilteredQueryParams,
			Headers:             json.Headers,
			Scope:               json.Scope,
			Percentage:          json.Percentage,
			Count:               json.Count,
			ConcurrentCount:     json.ConcurrentCount,
			CertificateChain:    certificateChain,
			InsecureSkipVerify:  p.insecureSkipVerify,
			Timeout:             json.Timeout,
			LoadToCache:         json.LoadToCache,
		})

	default:
		return nil, errors.New("unknown preheat type")
	}

	// Initialize queues.
	queues, err := getSchedulerQueues(schedulers)
	if err != nil {
		return nil, err
	}

	return p.createGroupJob(ctx, files, queues)
}

// createGroupJob creates a group job.
func (p *preheat) createGroupJob(ctx context.Context, files []internaljob.PreheatRequest, queues []internaljob.Queue) (*internaljob.GroupJobState, error) {
	groupUUID := fmt.Sprintf("group_%s", uuid.New().String())
	var signatures []*machineryv1tasks.Signature
	for _, queue := range queues {
		for _, file := range files {
			file.GroupUUID = groupUUID
			taskUUID := fmt.Sprintf("task_%s", uuid.New().String())
			file.TaskUUID = taskUUID

			args, err := internaljob.MarshalRequest(file)
			if err != nil {
				logger.Errorf("[preheat]: preheat marshal request: %v, error: %v", file, err)
				continue
			}

			signatures = append(signatures, &machineryv1tasks.Signature{
				UUID:       taskUUID,
				Name:       internaljob.PreheatJob,
				RoutingKey: queue.String(),
				Args:       args,
			})
		}
	}

	group, err := machineryv1tasks.NewGroup(signatures...)
	if err != nil {
		return nil, err
	}

	var tasks []machineryv1tasks.Signature
	for _, signature := range signatures {
		tasks = append(tasks, *signature)
	}

	group.GroupUUID = groupUUID
	logger.Infof("[preheat]: create preheat group %s in queues %v, tasks: %#v", group.GroupUUID, queues, tasks)
	if _, err := p.job.Server.SendGroupWithContext(ctx, group, 50); err != nil {
		logger.Errorf("[preheat]: create preheat group %s failed", group.GroupUUID, err)
		return nil, err
	}

	return &internaljob.GroupJobState{
		GroupUUID: group.GroupUUID,
		State:     machineryv1tasks.StatePending,
		CreatedAt: time.Now(),
	}, nil
}

func GetImageLayers(ctx context.Context, args types.PreheatArgs, registryTimeout time.Duration, rootCAs *x509.CertPool, insecureSkipVerify bool) ([]internaljob.PreheatRequest, error) {
	ctx, span := tracer.Start(ctx, config.SpanGetLayers, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	// Parse image manifest url.
	image, err := parseManifestURL(args.URL)
	if err != nil {
		return nil, err
	}

	// Background:
	// Harbor uses the V1 preheat request and will carry the auth info in the headers.
	options := []imageAuthClientOption{}
	header := nethttp.MapToHeader(args.Headers)
	if token := header.Get("Authorization"); len(token) > 0 {
		options = append(options, withIssuedToken(token))
		header.Set("Authorization", token)
	}

	httpClient := &http.Client{
		Timeout: registryTimeout,
		Transport: &http.Transport{
			DialContext:         nethttp.NewSafeDialer().DialContext,
			TLSClientConfig:     &tls.Config{RootCAs: rootCAs, InsecureSkipVerify: insecureSkipVerify},
			MaxIdleConns:        defaultHTTPTransport.MaxIdleConns,
			MaxIdleConnsPerHost: defaultHTTPTransport.MaxIdleConnsPerHost,
			MaxConnsPerHost:     defaultHTTPTransport.MaxConnsPerHost,
			IdleConnTimeout:     defaultHTTPTransport.IdleConnTimeout,
		},
	}

	// Init docker auth client.
	client, err := newImageAuthClient(image, httpClient, &typesregistry.AuthConfig{Username: args.Username, Password: args.Password}, options...)
	if err != nil {
		return nil, err
	}

	// Get platform.
	platform := platforms.DefaultSpec()
	if args.Platform != "" {
		platform, err = platforms.Parse(args.Platform)
		if err != nil {
			return nil, err
		}
	}

	// Get manifests.
	manifests, err := getManifests(ctx, client, image, header.Clone(), platform)
	if err != nil {
		return nil, err
	}

	// no matching manifest for platform in the manifest list entries
	if len(manifests) == 0 {
		return nil, fmt.Errorf("no matching manifest for platform %s", platforms.Format(platform))
	}

	// set authorization header
	header.Set("Authorization", client.GetAuthToken())

	// parse image layers to preheat
	layers, err := parseLayers(manifests, args, header.Clone(), image, rootCAs, insecureSkipVerify)
	if err != nil {
		return nil, err
	}

	return layers, nil
}

// getManifests gets manifests of image.
func getManifests(ctx context.Context, client *imageAuthClient, image *preheatImage, header http.Header, platform specs.Platform) ([]distribution.Manifest, error) {
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
			manifests, err := getManifests(ctx, client, image, header.Clone(), platform)
			if err != nil {
				return nil, err
			}

			result = append(result, manifests...)
		}

		return result, nil
	}

	return nil, errors.New("unknown manifest type")
}

// filterManifests filters manifests with platform.
func filterManifests(manifests []manifestlist.ManifestDescriptor, platform specs.Platform) []manifestlist.ManifestDescriptor {
	var matches []manifestlist.ManifestDescriptor
	for _, desc := range manifests {
		if desc.Platform.Architecture == platform.Architecture && desc.Platform.OS == platform.OS {
			matches = append(matches, desc)
		}
	}

	return matches
}

// parseLayers parses layers of image.
func parseLayers(manifests []distribution.Manifest, args types.PreheatArgs, header http.Header, image *preheatImage, rootCAs *x509.CertPool, insecureSkipVerify bool) ([]internaljob.PreheatRequest, error) {
	var certificateChain [][]byte
	if rootCAs != nil {
		certificateChain = rootCAs.Subjects()
	}

	var layerURLs []string
	for _, m := range manifests {
		for _, v := range m.References() {
			header.Set("Accept", v.MediaType)
			layerURLs = append(layerURLs, image.blobsURL(v.Digest.String()))
		}
	}

	layers := internaljob.PreheatRequest{
		URLs:                layerURLs,
		PieceLength:         args.PieceLength,
		Tag:                 args.Tag,
		Application:         args.Application,
		FilteredQueryParams: args.FilteredQueryParams,
		Headers:             nethttp.HeaderToMap(header),
		Scope:               args.Scope,
		Percentage:          args.Percentage,
		Count:               args.Count,
		ConcurrentCount:     args.ConcurrentCount,
		CertificateChain:    certificateChain,
		InsecureSkipVerify:  insecureSkipVerify,
		Timeout:             args.Timeout,
	}

	return []internaljob.PreheatRequest{layers}, nil
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

	// Return earlier if issued token is not empty.
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

// interceptorTokenHandler is a token interceptor
// intercept bearer token from auth handler.
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
