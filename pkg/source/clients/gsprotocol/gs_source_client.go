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

package gsprotocol

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/go-http-utils/headers"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	nethttp "d7y.io/dragonfly/v2/pkg/net/http"
	"d7y.io/dragonfly/v2/pkg/source"
)

const GSScheme = "gs"

const (
	// GS credentials json
	gsCredentialsJSONBase64 = "gsCredentialsJSONBase64"
)

var _ source.ResourceClient = (*gsSourceClient)(nil)

func init() {
	source.RegisterBuilder(GSScheme, source.NewPlainResourceClientBuilder(Builder))
}

func Builder(optionYaml []byte) (source.ResourceClient, source.RequestAdapter, []source.Hook, error) {
	client := &gsSourceClient{}
	return client, client.adaptor, nil, nil
}

// gsSourceClient is an implementation of the interface of source.ResourceClient.
type gsSourceClient struct {
	httpClient *http.Client
}

func (s *gsSourceClient) adaptor(request *source.Request) *source.Request {
	clonedRequest := request.Clone(request.Context())
	if request.Header.Get(source.Range) != "" {
		clonedRequest.Header.Set(headers.Range, fmt.Sprintf("bytes=%s", request.Header.Get(source.Range)))
		clonedRequest.Header.Del(source.Range)
	}
	return clonedRequest
}

func (s *gsSourceClient) newGSClient(request *source.Request) (*storage.Client, error) {
	gsCredentialsJSON, err := base64.StdEncoding.DecodeString(request.Header.Get(gsCredentialsJSONBase64))
	if err != nil {
		return nil, err
	}
	opts := []option.ClientOption{
		option.WithScopes(storage.ScopeReadOnly),
		option.WithCredentialsJSON(gsCredentialsJSON),
	}
	if s.httpClient != nil {
		opts = append(opts, option.WithHTTPClient(s.httpClient))
	}
	return storage.NewClient(request.Context(), opts...)
}

// GetContentLength get length of resource content
// return source.UnknownSourceFileLen if response status is not StatusOK and StatusPartialContent
func (s *gsSourceClient) GetContentLength(request *source.Request) (int64, error) {
	client, err := s.newGSClient(request)
	if err != nil {
		return -1, err
	}

	objAttrs, err := client.Bucket(request.URL.Host).Object(request.URL.Path).Attrs(request.Context())
	if err != nil {
		return -1, err
	}

	rangeHeader := request.Header.Get(headers.Range)
	if rangeHeader != "" {
		rgs, err := nethttp.ParseRange(rangeHeader, objAttrs.Size)
		if err != nil {
			return -1, err
		}
		if len(rgs) != 1 {
			return -1, err
		}
		return rgs[0].Length, nil
	}

	if err != nil {
		return -1, err
	}
	return objAttrs.Size, nil
}

// IsSupportRange checks if resource supports breakpoint continuation
// return false if response status is not StatusPartialContent
func (s *gsSourceClient) IsSupportRange(request *source.Request) (bool, error) {
	return true, nil
}

// IsExpired checks if a resource received or stored is the same.
// return false and non-nil err to prevent the source from exploding if
// fails to get the result, it is considered that the source has not expired
func (s *gsSourceClient) IsExpired(request *source.Request, info *source.ExpireInfo) (bool, error) {
	return false, fmt.Errorf("not implemented") // TODO: Implement
}

// Download downloads from source
func (s *gsSourceClient) Download(request *source.Request) (*source.Response, error) {
	var (
		err    error
		client *storage.Client
		reader *storage.Reader
	)

	client, err = s.newGSClient(request)
	if err != nil {
		return nil, err
	}

	rangeHeader := request.Header.Get(headers.Range)
	bucket := request.URL.Host
	object := strings.TrimPrefix(request.URL.Path, "/")

	if rangeHeader == "" {
		reader, err = client.Bucket(bucket).Object(object).NewReader(request.Context())
	} else {
		var objAttrs *storage.ObjectAttrs
		objAttrs, err = client.Bucket(bucket).Object(object).Attrs(request.Context())
		if err != nil {
			return nil, err
		}

		var rgs []nethttp.Range
		rgs, err = nethttp.ParseRange(rangeHeader, objAttrs.Size)
		if err != nil {
			return nil, err
		}
		if len(rgs) != 1 {
			return nil, err
		}

		reader, err = client.Bucket(bucket).Object(object).
			NewRangeReader(request.Context(), rgs[0].Start, rgs[0].Length)
	}

	if err != nil {
		// TODO parse error details
		return nil, err
	}

	hdr := source.Header{}
	if reader.Attrs.CacheControl != "" {
		hdr[headers.CacheControl] = []string{reader.Attrs.CacheControl}
	}

	return &source.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		Header:        hdr,
		Body:          reader,
		ContentLength: reader.Attrs.Size,
		Validate: func() error {
			return nil
		},
	}, nil
}

// GetLastModified gets last modified timestamp milliseconds of resource
func (s *gsSourceClient) GetLastModified(request *source.Request) (int64, error) {
	client, err := s.newGSClient(request)
	if err != nil {
		return -1, err
	}

	bucket := request.URL.Host
	object := strings.TrimPrefix(request.URL.Path, "/")
	objectAttrs, err := client.Bucket(bucket).Object(object).Attrs(request.Context())
	if err != nil {
		return -1, err
	}
	return objectAttrs.Updated.UnixMilli(), nil
}

func (s *gsSourceClient) List(request *source.Request) (urls []source.URLEntry, err error) {
	client, err := s.newGSClient(request)
	if err != nil {
		return nil, fmt.Errorf("get gs client: %w", err)
	}

	// list all files
	prefix := buildListPrefix(request.URL.Path)
	it := client.Bucket(request.URL.Host).Objects(
		request.Context(), &storage.Query{
			Prefix: prefix,
		})

	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}

		url := *request.URL
		url.Path = attrs.Name
		urls = append(urls, source.URLEntry{URL: &url, Name: url.Path, IsDir: false})
	}

	return urls, nil
}

func addTrailingSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

func buildListPrefix(s string) string {
	s = addTrailingSlash(s)

	// s3 objects id should not start with '/'
	return strings.TrimLeft(s, "/")
}
