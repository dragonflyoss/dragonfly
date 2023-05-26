/*
 *     Copyright 2023 The Dragonfly Authors
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

//go:generate mockgen -destination mocks/probes_mock.go -source probes.go -package mocks

package networktopology

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"

	pkgredis "d7y.io/dragonfly/v2/pkg/redis"
	"d7y.io/dragonfly/v2/scheduler/config"
	"d7y.io/dragonfly/v2/scheduler/resource"
)

const (
	// DefaultMovingAverageWeight is the weight of the moving average.
	DefaultMovingAverageWeight = 0.1
)

// Probe is the probe metadata.
type Probe struct {
	// Host metadata.
	Host *resource.Host `json:"host"`

	// RTT is the round-trip time sent via this pinger.
	RTT time.Duration `json:"rtt"`

	// CreatedAt is the time to create probe.
	CreatedAt time.Time `json:"createdAt"`
}

// NewProbe creates a new probe instance.
func NewProbe(host *resource.Host, rtt time.Duration, createdAt time.Time) *Probe {
	return &Probe{
		Host:      host,
		RTT:       rtt,
		CreatedAt: createdAt,
	}
}

// Probes is the interface to store probes.
type Probes interface {
	// Peek returns the oldest probe without removing it.
	Peek() (*Probe, error)

	// Enqueue enqueues probe into the queue.
	Enqueue(*Probe) error

	// Dequeue removes and returns the oldest probe.
	Dequeue() (*Probe, error)

	// Length gets the length of probes.
	Length() (int64, error)

	// CreatedAt is the creation time of probes.
	CreatedAt() (time.Time, error)

	// UpdatedAt is the updated time to store probe.
	UpdatedAt() (time.Time, error)

	// AverageRTT is the average round-trip time of probes.
	AverageRTT() (time.Duration, error)
}

// probes is the implementation of Probes.
type probes struct {
	// config is the probe config.
	config config.ProbeConfig

	// rdb is redis universal client interface.
	rdb redis.UniversalClient

	// srcHostID is the source host id.
	srcHostID string

	// destHostID is the destination host id.
	destHostID string
}

// NewProbes creates a probes interface.
func NewProbes(cfg config.ProbeConfig, rdb redis.UniversalClient, srcHostID string, destHostID string) Probes {
	return &probes{
		config:     cfg,
		rdb:        rdb,
		srcHostID:  srcHostID,
		destHostID: destHostID,
	}
}

// Peek returns the oldest probe without removing it.
func (p *probes) Peek() (*Probe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	rawProbe, err := p.rdb.LIndex(ctx, pkgredis.MakeProbesKeyInScheduler(p.srcHostID, p.destHostID), 0).Bytes()
	if err != nil {
		return nil, err
	}

	probe := &Probe{}
	if err = json.Unmarshal(rawProbe, probe); err != nil {
		return nil, err
	}

	return probe, err
}

// Enqueue enqueues probe into the queue.
func (p *probes) Enqueue(probe *Probe) error {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	length, err := p.Length()
	if err != nil {
		return err
	}

	if length == int64(p.config.QueueLength) {
		if _, err := p.Dequeue(); err != nil {
			return err
		}
	}

	data, err := json.Marshal(probe)
	if err != nil {
		return err
	}

	if err := p.rdb.RPush(ctx, pkgredis.MakeProbesKeyInScheduler(p.srcHostID, p.destHostID), data).Err(); err != nil {
		return err
	}

	if length == 0 {
		if _, err := p.rdb.Pipelined(context.Background(), func(rdb redis.Pipeliner) error {
			rdb.HSet(context.Background(), pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "averageRTT", probe.RTT.Nanoseconds())
			rdb.HSet(context.Background(), pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "createdAt", probe.CreatedAt.Format(time.RFC3339Nano))
			rdb.HSet(context.Background(), pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "updatedAt", probe.CreatedAt.Format(time.RFC3339Nano))
			return nil
		}); err != nil {
			return err
		}

		return nil
	}

	values, err := p.rdb.LRange(context.Background(), pkgredis.MakeProbesKeyInScheduler(p.srcHostID, p.destHostID), 0, -1).Result()
	if err != nil {
		return err
	}

	var averageRTT time.Duration
	for index, value := range values {
		probe := &Probe{}
		if err = json.Unmarshal([]byte(value), probe); err != nil {
			return err
		}

		if index == 0 {
			averageRTT = probe.RTT
			continue
		}

		averageRTT = time.Duration(float64(averageRTT)*DefaultMovingAverageWeight +
			float64(probe.RTT)*(1-DefaultMovingAverageWeight))
	}

	if _, err := p.rdb.Pipelined(context.Background(), func(rdb redis.Pipeliner) error {
		rdb.HSet(context.Background(), pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "averageRTT", averageRTT.Nanoseconds())
		rdb.HSet(context.Background(), pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "updatedAt", probe.CreatedAt.Format(time.RFC3339Nano))
		return nil
	}); err != nil {
		return err
	}

	return nil
}

// Dequeue removes and returns the oldest probe.
func (p *probes) Dequeue() (*Probe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	rawProbe, err := p.rdb.LPop(ctx, pkgredis.MakeProbesKeyInScheduler(p.srcHostID, p.destHostID)).Bytes()
	if err != nil {
		return nil, err
	}

	probe := &Probe{}
	if err = json.Unmarshal(rawProbe, probe); err != nil {
		return nil, err
	}

	return probe, nil
}

// Length gets the length of probes.
func (p *probes) Length() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return p.rdb.LLen(ctx, pkgredis.MakeProbesKeyInScheduler(p.srcHostID, p.destHostID)).Result()
}

// CreatedAt is the creation time of probes.
func (p *probes) CreatedAt() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return p.rdb.HGet(ctx, pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "createdAt").Time()
}

// UpdatedAt is the updated time to store probe.
func (p *probes) UpdatedAt() (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	return p.rdb.HGet(ctx, pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "updatedAt").Time()
}

// AverageRTT is the average round-trip time of probes.
func (p *probes) AverageRTT() (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	nano, err := p.rdb.HGet(ctx, pkgredis.MakeNetworkTopologyKeyInScheduler(p.srcHostID, p.destHostID), "averageRTT").Uint64()
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(nano), nil
}
