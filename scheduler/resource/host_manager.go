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

//go:generate mockgen -destination host_manager_mock.go -source host_manager.go -package resource

package resource

import (
	"sync"

	pkggc "d7y.io/dragonfly/v2/pkg/gc"
	"d7y.io/dragonfly/v2/pkg/types"
	"d7y.io/dragonfly/v2/scheduler/config"
)

const (
	// GC host id.
	GCHostID = "host"
)

// HostManager is the interface used for host manager.
type HostManager interface {
	// Load returns host for a key.
	Load(string) (*Host, bool)

	// Store sets host.
	Store(*Host)

	// LoadOrStore returns host the key if present.
	// Otherwise, it stores and returns the given host.
	// The loaded result is true if the host was loaded, false if stored.
	LoadOrStore(*Host) (*Host, bool)

	// Delete deletes host for a key.
	Delete(string)

	// FilterProbedHosts returns the specified number of hosts randomly,
	// or returns all hosts when the number in the map is insufficient.
	FilterProbedHosts(srcHostID string, count int) []*Host

	// Try to reclaim host.
	RunGC() error
}

// hostManager contains content for host manager.
type hostManager struct {
	// Host sync map.
	*sync.Map
}

// New host manager interface.
func newHostManager(cfg *config.GCConfig, gc pkggc.GC) (HostManager, error) {
	h := &hostManager{
		Map: &sync.Map{},
	}

	if err := gc.Add(pkggc.Task{
		ID:       GCHostID,
		Interval: cfg.HostGCInterval,
		Timeout:  cfg.HostGCInterval,
		Runner:   h,
	}); err != nil {
		return nil, err
	}

	return h, nil
}

// Load returns host for a key.
func (h *hostManager) Load(key string) (*Host, bool) {
	rawHost, loaded := h.Map.Load(key)
	if !loaded {
		return nil, false
	}

	return rawHost.(*Host), loaded
}

// Store sets host.
func (h *hostManager) Store(host *Host) {
	h.Map.Store(host.ID, host)
}

// LoadOrStore returns host the key if present.
// Otherwise, it stores and returns the given host.
// The loaded result is true if the host was loaded, false if stored.
func (h *hostManager) LoadOrStore(host *Host) (*Host, bool) {
	rawHost, loaded := h.Map.LoadOrStore(host.ID, host)
	return rawHost.(*Host), loaded
}

// Delete deletes host for a key.
func (h *hostManager) Delete(key string) {
	h.Map.Delete(key)
}

// FilterProbedHosts returns the specified number of hosts randomly,
// or returns all hosts when the number in the map is insufficient.
func (h *hostManager) FilterProbedHosts(hostID string, count int) []*Host {
	var hosts []*Host
	h.Map.Range(func(key, value any) bool {
		if len(hosts) >= count {
			return false
		}

		if key != hostID {
			hosts = append(hosts, value.(*Host))
		}

		return true
	})

	return hosts
}

// Try to reclaim host.
func (h *hostManager) RunGC() error {
	h.Map.Range(func(_, value any) bool {
		host, ok := value.(*Host)
		if !ok {
			host.Log.Error("invalid host")
			return true
		}

		if host.PeerCount.Load() == 0 &&
			host.ConcurrentUploadCount.Load() == 0 &&
			host.Type == types.HostTypeNormal {
			host.Log.Info("host has been reclaimed")
			h.Delete(host.ID)
		}

		return true
	})

	return nil
}
