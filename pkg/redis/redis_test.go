/*
 *     Copyright 2024 The Dragonfly Authors
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

package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsEnabled(t *testing.T) {
	tests := []struct {
		name   string
		addrs  []string
		expect func(t *testing.T, ok bool)
	}{
		{
			name:  "check redis is enabled",
			addrs: []string{"172.0.0.1"},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.True(ok)
			},
		},
		{
			name:  "addrs is empty",
			addrs: []string{},
			expect: func(t *testing.T, ok bool) {
				assert := assert.New(t)
				assert.False(ok)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, IsEnabled(tc.addrs))
		})
	}
}

func Test_MakeNamespaceKeyInManager(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		expect    func(t *testing.T, s string)
	}{
		{
			name:      "make namespace key in manager",
			namespace: "namespace",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:namespace", s)
			},
		},
		{
			name:      "namespace is empty",
			namespace: "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeNamespaceKeyInManager(tc.namespace))
		})
	}
}

func Test_MakeKeyInManager(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		id        string
		expect    func(t *testing.T, s string)
	}{
		{
			name:      "make key in manager",
			namespace: "namespace",
			id:        "foo",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:namespace:foo", s)
			},
		},
		{
			name:      "namespace is empty",
			namespace: "",
			id:        "foo",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager::foo", s)
			},
		},
		{
			name:      "key is empty",
			namespace: "namespace",
			id:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:namespace:", s)
			},
		},
		{
			name:      "namespace and key are empty",
			namespace: "",
			id:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager::", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeKeyInManager(tc.namespace, tc.id))
		})
	}
}

func Test_MakeSeedPeerKeyInManager(t *testing.T) {
	tests := []struct {
		name      string
		clusterID uint
		hostname  string
		ip        string
		expect    func(t *testing.T, s string)
	}{
		{
			name:      "make seed peer key in manager",
			clusterID: 1,
			hostname:  "foo",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:seed-peers:1-foo-127.0.0.1", s)
			},
		},
		{
			name:      "hostname is empty",
			clusterID: 1,
			hostname:  "",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:seed-peers:1--127.0.0.1", s)
			},
		},
		{
			name:      "ip is empty",
			clusterID: 1,
			hostname:  "bar",
			ip:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:seed-peers:1-bar-", s)
			},
		},
		{
			name:      "hostname and ip are empty",
			clusterID: 1,
			hostname:  "",
			ip:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:seed-peers:1--", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeSeedPeerKeyInManager(tc.clusterID, tc.hostname, tc.ip))
		})
	}
}

func Test_MakeSchedulerKeyInManager(t *testing.T) {
	tests := []struct {
		name      string
		clusterID uint
		hostname  string
		ip        string
		expect    func(t *testing.T, s string)
	}{
		{
			name:      "make scheduler key in manager",
			clusterID: 1,
			hostname:  "bar",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:schedulers:1-bar-127.0.0.1", s)
			},
		},
		{
			name:      "hostname is empty",
			clusterID: 1,
			hostname:  "",
			ip:        "127.0.0.1",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:schedulers:1--127.0.0.1", s)
			},
		},
		{
			name:      "ip is empty",
			clusterID: 1,
			hostname:  "bar",
			ip:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:schedulers:1-bar-", s)
			},
		},
		{
			name:      "hostname and ip are empty",
			clusterID: 1,
			hostname:  "",
			ip:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:schedulers:1--", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeSchedulerKeyInManager(tc.clusterID, tc.hostname, tc.ip))
		})
	}
}

func Test_MakeSeedPeersKeyForPeerInManager(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		ip       string
		expect   func(t *testing.T, s string)
	}{
		{
			name:     "make seed peer key for peer in manager",
			hostname: "bar",
			ip:       "127.0.0.1",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:bar-127.0.0.1:seed-peers", s)
			},
		},
		{
			name:     "hostname is empty",
			hostname: "",
			ip:       "127.0.0.1",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:-127.0.0.1:seed-peers", s)
			},
		},
		{
			name:     "ip is empty",
			hostname: "bar",
			ip:       "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:bar-:seed-peers", s)
			},
		},
		{
			name:     "hostname and ip are empty",
			hostname: "",
			ip:       "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:-:seed-peers", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeSeedPeersKeyForPeerInManager(tc.hostname, tc.ip))
		})
	}
}

func Test_MakeSchedulersKeyForPeerInManager(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		ip       string
		version  string
		expect   func(t *testing.T, s string)
	}{
		{
			name:     "make scheduler key for peer in manager",
			hostname: "bar",
			ip:       "127.0.0.1",
			version:  "0.1.0",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:bar-127.0.0.1-0.1.0:schedulers", s)
			},
		},
		{
			name:     "hostname is empty",
			hostname: "",
			ip:       "127.0.0.1",
			version:  "0.1.0",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:-127.0.0.1-0.1.0:schedulers", s)
			},
		},
		{
			name:     "ip is empty",
			hostname: "bar",
			ip:       "",
			version:  "0.1.0",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:bar--0.1.0:schedulers", s)
			},
		},
		{
			name:     "version is empty",
			hostname: "bar",
			ip:       "127.0.0.1",
			version:  "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:bar-127.0.0.1-:schedulers", s)
			},
		},
		{
			name:     "hostname, ip and version are empty",
			hostname: "",
			ip:       "",
			version:  "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("manager:peers:--:schedulers", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeSchedulersKeyForPeerInManager(tc.hostname, tc.ip, tc.version))
		})
	}
}

func Test_MakeApplicationsKeyInManager(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("manager:applications", MakeApplicationsKeyInManager())
}

func Test_MakeNamespaceKeyInScheduler(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		expect    func(t *testing.T, s string)
	}{
		{
			name:      "make namespace key in scheduler",
			namespace: "baz",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("scheduler:baz", s)
			},
		},
		{
			name:      "namespace is empty",
			namespace: "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("scheduler:", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeNamespaceKeyInScheduler(tc.namespace))
		})
	}
}

func Test_MakeKeyInScheduler(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		id        string
		expect    func(t *testing.T, s string)
	}{
		{
			name:      "make key in scheduler",
			namespace: "bas",
			id:        "id",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("scheduler:bas:id", s)
			},
		},
		{
			name:      "namespace is empty",
			namespace: "",
			id:        "id",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("scheduler::id", s)
			},
		},
		{
			name:      "id is empty",
			namespace: "bas",
			id:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("scheduler:bas:", s)
			},
		},
		{
			name:      "namespace and id are empty",
			namespace: "",
			id:        "",
			expect: func(t *testing.T, s string) {
				assert := assert.New(t)
				assert.Equal("scheduler::", s)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.expect(t, MakeKeyInScheduler(tc.namespace, tc.id))
		})
	}
}
