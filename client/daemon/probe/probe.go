/*
 *     Copyright 2022 The Dragonfly Authors
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

//go:generate mockgen -destination mocks/probe_mock.go -source probe.go -package mocks

package probe

import (
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "d7y.io/api/pkg/apis/common/v1"
	schedulerv1 "d7y.io/api/pkg/apis/scheduler/v1"

	"d7y.io/dragonfly/v2/client/config"
	logger "d7y.io/dragonfly/v2/internal/dflog"
	pkgbalancer "d7y.io/dragonfly/v2/pkg/balancer"
	"d7y.io/dragonfly/v2/pkg/net/ping"
	schedulerclient "d7y.io/dragonfly/v2/pkg/rpc/scheduler/client"
)

type Probe interface {
	// Serve starts probe server.
	Serve()

	// Stop stops probe server.
	Stop()
}

type probe struct {
	config             *config.DaemonOption
	hostID             string
	daemonPort         int32
	daemonDownloadPort int32
	schedulerClient    schedulerclient.V1
	done               chan struct{}
}

// NewProbe returns a new Probe interface.
func NewProbe(cfg *config.DaemonOption, hostID string, daemonPort int32, daemonDownloadPort int32,
	schedulerClient schedulerclient.V1) (Probe, error) {
	return &probe{
		config:             cfg,
		hostID:             hostID,
		daemonPort:         daemonPort,
		daemonDownloadPort: daemonDownloadPort,
		schedulerClient:    schedulerClient,
		done:               make(chan struct{}),
	}, nil
}

// Serve starts probe server.
func (p *probe) Serve() {
	logger.Info("collect probes and upload probes to scheduler")
	p.uploadProbesToScheduler()
}

// Stop stops probe server.
func (p *probe) Stop() {
	close(p.done)
}

// uploadProbesToScheduler collects probes and uploads probes to scheduler.
func (p *probe) uploadProbesToScheduler() {
	host := &v1.Host{
		Id:           p.hostID,
		Ip:           p.config.Host.AdvertiseIP.String(),
		Hostname:     p.config.Host.Hostname,
		Port:         p.daemonPort,
		DownloadPort: p.daemonDownloadPort,
		Location:     p.config.Host.Location,
		Idc:          p.config.Host.IDC,
	}

	ctx := context.WithValue(context.Background(), pkgbalancer.ContextKey, host.Id)
	tick := time.NewTicker(p.config.NetworkTopology.Probe.Interval)
	for {
		select {
		case <-tick.C:
			stream, err := p.schedulerClient.SyncProbes(ctx, &schedulerv1.SyncProbesRequest{
				Host: host,
				Request: &schedulerv1.SyncProbesRequest_ProbeStartedRequest{
					ProbeStartedRequest: &schedulerv1.ProbeStartedRequest{},
				},
			})
			if err != nil {
				continue
			}

			var syncProbesResponse *schedulerv1.SyncProbesResponse
			for {
				response, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}

					logger.Errorf("receive error: %s", err.Error())
					continue
				}

				syncProbesResponse = response
			}

			probes, failedProbes := p.collectProbes(syncProbesResponse.Hosts)
			eg := errgroup.Group{}
			eg.Go(func() error {
				if len(probes) > 0 {
					if err := stream.Send(&schedulerv1.SyncProbesRequest{
						Host: host,
						Request: &schedulerv1.SyncProbesRequest_ProbeFinishedRequest{
							ProbeFinishedRequest: &schedulerv1.ProbeFinishedRequest{
								Probes: probes,
							},
						},
					}); err != nil {
						return fmt.Errorf("synchronize finished probe: %w", err)
					}
				}

				return nil
			})

			eg.Go(func() error {
				if len(failedProbes) > 0 {
					if err := stream.Send(&schedulerv1.SyncProbesRequest{
						Host: host,
						Request: &schedulerv1.SyncProbesRequest_ProbeFailedRequest{
							ProbeFailedRequest: &schedulerv1.ProbeFailedRequest{
								Probes: failedProbes,
							},
						},
					}); err != nil {
						return fmt.Errorf("synchronize failed probe: %w", err)
					}
				}

				return nil
			})

			if err := eg.Wait(); err != nil {
				logger.Error(err)
			}
		case <-p.done:
			return
		}
	}
}

// collectProbes probes hosts and collect probes and failedProbes.
func (p *probe) collectProbes(desthosts []*v1.Host) ([]*schedulerv1.Probe, []*schedulerv1.FailedProbe) {
	probes := make([]*schedulerv1.Probe, 0)
	failedProbes := make([]*schedulerv1.FailedProbe, 0)
	for _, desthost := range desthosts {
		statistics, err := ping.Ping(desthost.Ip)
		if err != nil {
			failedProbes = append(failedProbes, &schedulerv1.FailedProbe{
				Host: &v1.Host{
					Id:           desthost.Id,
					Ip:           desthost.Ip,
					Hostname:     desthost.Hostname,
					Port:         desthost.Port,
					DownloadPort: desthost.DownloadPort,
					Location:     desthost.Location,
					Idc:          desthost.Idc,
				},
				Description: err.Error(),
			})

			continue
		}

		probes = append(probes, &schedulerv1.Probe{
			Host: &v1.Host{
				Id:           desthost.Id,
				Ip:           desthost.Ip,
				Hostname:     desthost.Hostname,
				Port:         desthost.Port,
				DownloadPort: desthost.DownloadPort,
				Location:     desthost.Location,
				Idc:          desthost.Idc,
			},
			Rtt:       durationpb.New(statistics.AvgRtt),
			CreatedAt: timestamppb.New(time.Now()),
		})
	}

	return probes, failedProbes
}
