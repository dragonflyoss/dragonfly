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

package dfnet

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNetAddr_String(t *testing.T) {
	tests := []struct {
		name   string
		addr   NetAddr
		expect func(t *testing.T, s string)
	}{
		{
			name: "tcp",
			addr: NetAddr{Type: TCP, Addr: "127.0.0.1:8002"},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("dns:///127.0.0.1:8002", s)
			},
		},
		{
			name: "unix",
			addr: NetAddr{Type: UNIX, Addr: "/var/run/dfdaemon.sock"},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("unix:///var/run/dfdaemon.sock", s)
			},
		},
		{
			name: "vsock",
			addr: NetAddr{Type: VSOCK, Addr: "2:8002"},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("vsock://2:8002", s)
			},
		},
		{
			name: "empty type falls back to dns",
			addr: NetAddr{Addr: "127.0.0.1:8002"},
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("dns:///127.0.0.1:8002", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, tc.addr.String())
		})
	}
}

func TestNetAddr_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		expect func(t *testing.T, n NetAddr, err error)
	}{
		{
			name: "bare string becomes tcp",
			data: `"127.0.0.1:8002"`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: TCP, Addr: "127.0.0.1:8002"}, n)
			},
		},
		{
			name: "mapping with unix type",
			data: `{"type":"unix","addr":"/var/run/dfdaemon.sock"}`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: UNIX, Addr: "/var/run/dfdaemon.sock"}, n)
			},
		},
		{
			name: "mapping with vsock type",
			data: `{"type":"vsock","addr":"2:8002"}`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: VSOCK, Addr: "2:8002"}, n)
			},
		},
		{
			name: "mapping without type",
			data: `{"addr":"127.0.0.1:8002"}`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Addr: "127.0.0.1:8002"}, n)
			},
		},
		{
			name: "number",
			data: `8002`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "sequence",
			data: `["127.0.0.1:8002"]`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "boolean",
			data: `true`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "mapping with non string type",
			data: `{"type":1,"addr":"127.0.0.1:8002"}`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "malformed json",
			data: `{"type":`,
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var n NetAddr
			err := json.Unmarshal([]byte(tc.data), &n)
			tc.expect(t, n, err)
		})
	}
}

func TestNetAddr_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		expect func(t *testing.T, n NetAddr, err error)
	}{
		{
			name: "bare scalar becomes tcp",
			data: "127.0.0.1:8002",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: TCP, Addr: "127.0.0.1:8002"}, n)
			},
		},
		{
			name: "numeric scalar becomes tcp",
			data: "8002",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: TCP, Addr: "8002"}, n)
			},
		},
		{
			name: "mapping with unix type",
			data: "type: unix\naddr: /var/run/dfdaemon.sock",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: UNIX, Addr: "/var/run/dfdaemon.sock"}, n)
			},
		},
		{
			name: "mapping with vsock type",
			data: "type: vsock\naddr: 2:8002",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Type: VSOCK, Addr: "2:8002"}, n)
			},
		},
		{
			name: "mapping without type",
			data: "addr: 127.0.0.1:8002",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.NoError(err)
				assert.Equal(NetAddr{Addr: "127.0.0.1:8002"}, n)
			},
		},
		{
			name: "sequence",
			data: "- 127.0.0.1:8002",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "mapping with sequence type",
			data: "type: [unix]\naddr: /var/run/dfdaemon.sock",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
		{
			name: "malformed yaml",
			data: "type: unix\n  addr: [",
			expect: func(t *testing.T, n NetAddr, err error) {
				assert := assert.New(t)
				assert.Error(err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var n NetAddr
			err := yaml.Unmarshal([]byte(tc.data), &n)
			tc.expect(t, n, err)
		})
	}
}
