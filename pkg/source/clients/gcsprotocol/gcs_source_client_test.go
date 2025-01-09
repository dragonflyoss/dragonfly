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

package gcsprotocol

import (
	"io"
	"net/url"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/assert"

	"d7y.io/dragonfly/v2/pkg/source"
)

func TestGetContentLength(t *testing.T) {
	testCases := []struct {
		name    string
		request *source.Request
		object  fakestorage.Object
		err     error
		length  int
	}{
		{
			name: "normal file",
			request: &source.Request{
				URL: &url.URL{
					Scheme: "gcs",
					Opaque: "",
					User:   nil,
					Host:   "fake-bucket",
					Path:   "fake-dir/fake-object",
				},
				Header: nil,
			},
			object: fakestorage.Object{
				ObjectAttrs: fakestorage.ObjectAttrs{
					BucketName: "fake-bucket",
					Name:       "fake-dir/fake-object",
				},
				Content: []byte("test"),
			},
			err:    nil,
			length: 4,
		},
		{
			name: "normal file with range",
			request: &source.Request{
				URL: &url.URL{
					Scheme: "gcs",
					Opaque: "",
					User:   nil,
					Host:   "fake-bucket",
					Path:   "fake-dir/fake-object",
				},
				Header: map[string][]string{
					headers.Range: {"bytes=0-1"},
				},
			},
			object: fakestorage.Object{
				ObjectAttrs: fakestorage.ObjectAttrs{
					BucketName: "fake-bucket",
					Name:       "fake-dir/fake-object",
				},
				Content: []byte("test"),
			},
			err:    nil,
			length: 2,
		},
		{
			name: "file not found",
			request: &source.Request{
				URL: &url.URL{
					Scheme: "gcs",
					Opaque: "",
					User:   nil,
					Host:   "fake-bucket",
					Path:   "fake-dir/fake-object",
				},
				Header: nil,
			},
			object: fakestorage.Object{
				ObjectAttrs: fakestorage.ObjectAttrs{
					BucketName: "fake-bucket",
					Name:       "fake-dir/fake-object-x",
				},
				Content: []byte("test"),
			},
			err:    storage.ErrObjectNotExist,
			length: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := fakestorage.Options{
				InitialObjects:      []fakestorage.Object{tc.object},
				StorageRoot:         "",
				Seed:                "",
				Scheme:              "",
				Host:                "localhost",
				Port:                0,
				NoListener:          false,
				ExternalURL:         "",
				PublicHost:          "",
				AllowedCORSHeaders:  nil,
				Writer:              nil,
				EventOptions:        fakestorage.EventManagerOptions{},
				BucketsLocation:     "",
				CertificateLocation: "",
				PrivateKeyLocation:  "",
			}

			server, err := fakestorage.NewServerWithOptions(opts)
			assert.Nil(t, err)
			client := &gcsSourceClient{}
			client.httpClient = server.HTTPClient()
			defer server.Stop()

			length, err := client.GetContentLength(tc.request)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.length, int(length))
		})
	}
}

func TestDownload(t *testing.T) {
	testCases := []struct {
		name    string
		request *source.Request
		object  fakestorage.Object
		err     error
		data    []byte
	}{
		{
			name: "normal file",
			request: &source.Request{
				URL: &url.URL{
					Scheme: "gcs",
					Opaque: "",
					User:   nil,
					Host:   "fake-bucket",
					Path:   "fake-dir/fake-object",
				},
				Header: nil,
			},
			object: fakestorage.Object{
				ObjectAttrs: fakestorage.ObjectAttrs{
					BucketName: "fake-bucket",
					Name:       "fake-dir/fake-object",
				},
				Content: []byte("test"),
			},
			err:  nil,
			data: []byte("test"),
		},
		{
			name: "normal file with range",
			request: &source.Request{
				URL: &url.URL{
					Scheme: "gcs",
					Opaque: "",
					User:   nil,
					Host:   "fake-bucket",
					Path:   "fake-dir/fake-object",
				},
				Header: map[string][]string{
					headers.Range: {"bytes=0-1"},
				},
			},
			object: fakestorage.Object{
				ObjectAttrs: fakestorage.ObjectAttrs{
					BucketName: "fake-bucket",
					Name:       "fake-dir/fake-object",
				},
				Content: []byte("test"),
			},
			err:  nil,
			data: []byte("te"),
		},
		{
			name: "file not found",
			request: &source.Request{
				URL: &url.URL{
					Scheme: "gcs",
					Opaque: "",
					User:   nil,
					Host:   "fake-bucket",
					Path:   "fake-dir/fake-object",
				},
				Header: nil,
			},
			object: fakestorage.Object{
				ObjectAttrs: fakestorage.ObjectAttrs{
					BucketName: "fake-bucket",
					Name:       "fake-dir/fake-object-x",
				},
				Content: []byte("test"),
			},
			err: storage.ErrObjectNotExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := fakestorage.Options{
				InitialObjects:      []fakestorage.Object{tc.object},
				StorageRoot:         "",
				Seed:                "",
				Scheme:              "",
				Host:                "localhost",
				Port:                0,
				NoListener:          false,
				ExternalURL:         "",
				PublicHost:          "",
				AllowedCORSHeaders:  nil,
				Writer:              nil,
				EventOptions:        fakestorage.EventManagerOptions{},
				BucketsLocation:     "",
				CertificateLocation: "",
				PrivateKeyLocation:  "",
			}

			server, err := fakestorage.NewServerWithOptions(opts)
			assert.Nil(t, err)
			client := &gcsSourceClient{}
			client.httpClient = server.HTTPClient()
			defer server.Stop()

			response, err := client.Download(tc.request)
			assert.Equal(t, tc.err, err)
			if err == nil {
				data, err := io.ReadAll(response.Body)
				assert.Nil(t, err)
				assert.Equal(t, tc.data, data)
			}
		})
	}
}
