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

package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	"d7y.io/dragonfly/v2/pkg/types"
)

const (
	schedulerWaitTimeout    = 5 * time.Minute
	schedulerPollInterval   = time.Second
	schedulerCommandTimeout = 15 * time.Second
	schedulerRPCTimeout     = 5 * time.Second
)

// schedulerPod contains only the Kubernetes fields needed by the restart test.
type schedulerPod struct {
	Metadata struct {
		Name              string  `json:"name"`
		UID               string  `json:"uid"`
		DeletionTimestamp *string `json:"deletionTimestamp"`
	} `json:"metadata"`
	Status struct {
		Phase  string `json:"phase"`
		PodIP  string `json:"podIP"`
		PodIPs []struct {
			IP string `json:"ip"`
		} `json:"podIPs"`
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"conditions"`
	} `json:"status"`
}

func (p schedulerPod) ready() bool {
	if p.Metadata.DeletionTimestamp != nil || p.Status.Phase != "Running" {
		return false
	}
	for _, condition := range p.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}

func schedulerKubectlJSON(ctx context.Context, target any, args ...string) error {
	ctx, cancel := context.WithTimeout(ctx, schedulerCommandTimeout)
	defer cancel()
	args = append([]string{"--request-timeout=10s", "-n", DragonflyNamespace}, args...)
	out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	if err := json.Unmarshal(out, target); err != nil {
		return fmt.Errorf("decode kubectl %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func listSchedulerPods(ctx context.Context) ([]schedulerPod, error) {
	var list struct {
		Items []schedulerPod `json:"items"`
	}
	if err := schedulerKubectlJSON(ctx, &list, "get", "pods", "-l", "component=scheduler", "-o", "json"); err != nil {
		return nil, err
	}
	if len(list.Items) != Servers[SchedulerServerName].Replicas {
		return nil, fmt.Errorf("found %d scheduler pods, want %d", len(list.Items), Servers[SchedulerServerName].Replicas)
	}
	return list.Items, nil
}

// CaptureSchedulerPodUIDs records the pods that must be replaced by a restart.
// It only reads Kubernetes state and never initiates a restart.
func CaptureSchedulerPodUIDs() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), schedulerCommandTimeout)
	defer cancel()
	pods, err := listSchedulerPods(ctx)
	if err != nil {
		return nil, err
	}
	uids := make(map[string]string, len(pods))
	for _, pod := range pods {
		if pod.Metadata.Name == "" || pod.Metadata.UID == "" || !pod.ready() {
			return nil, fmt.Errorf("scheduler pod %q is not Ready or has no UID", pod.Metadata.Name)
		}
		uids[pod.Metadata.Name] = pod.Metadata.UID
	}
	return uids, nil
}

func pollScheduler(ctx context.Context, description string, check func() error) error {
	ticker := time.NewTicker(schedulerPollInterval)
	defer ticker.Stop()
	for {
		err := check()
		if err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%s: %w; last observation: %v", description, ctx.Err(), err)
		case <-ticker.C:
		}
	}
}

// WaitForSchedulerReplacement waits until every captured StatefulSet pod has a
// different UID and all replacement pods are Ready. It does not restart pods.
func WaitForSchedulerReplacement(ctx context.Context, before map[string]string) error {
	if len(before) == 0 {
		return fmt.Errorf("no scheduler pod UIDs were captured")
	}
	ctx, cancel := context.WithTimeout(ctx, schedulerWaitTimeout)
	defer cancel()
	return pollScheduler(ctx, "wait for scheduler replacement", func() error {
		pods, err := listSchedulerPods(ctx)
		if err != nil {
			return err
		}
		if len(pods) != len(before) {
			return fmt.Errorf("found %d scheduler pods, captured %d", len(pods), len(before))
		}
		for _, pod := range pods {
			oldUID, ok := before[pod.Metadata.Name]
			if !ok {
				return fmt.Errorf("unexpected scheduler pod %s", pod.Metadata.Name)
			}
			if pod.Metadata.UID == "" || pod.Metadata.UID == oldUID {
				return fmt.Errorf("scheduler pod %s has not been replaced (UID %s)", pod.Metadata.Name, pod.Metadata.UID)
			}
			if !pod.ready() {
				return fmt.Errorf("replacement scheduler pod %s is not Ready", pod.Metadata.Name)
			}
		}
		return nil
	})
}

type schedulerExpectedHost struct {
	name string
	ips  map[string]bool
	typ  types.HostType
}

func expectedSchedulerHosts(ctx context.Context, peers []*PodExec) ([]schedulerExpectedHost, error) {
	if len(peers) == 0 {
		return nil, fmt.Errorf("no client or seed-client pods supplied")
	}
	hosts := make([]schedulerExpectedHost, 0, len(peers))
	for _, peer := range peers {
		if peer == nil || peer.namespace != DragonflyNamespace || peer.name == "" {
			return nil, fmt.Errorf("expected a client pod in %s", DragonflyNamespace)
		}
		host := schedulerExpectedHost{name: peer.name, ips: make(map[string]bool)}
		switch peer.container {
		case ClientServerName:
			host.typ = types.HostTypeNormal
		case SeedClientServerName:
			host.typ = types.HostTypeSuperSeed
		default:
			return nil, fmt.Errorf("pod %s has unsupported container %q", peer.name, peer.container)
		}
		var pod schedulerPod
		if err := schedulerKubectlJSON(ctx, &pod, "get", "pod", peer.name, "-o", "json"); err != nil {
			return nil, err
		}
		if !pod.ready() {
			return nil, fmt.Errorf("client pod %s is not Ready", peer.name)
		}
		// Accept either pod address in the dual-stack/IPv6 E2E environments.
		if ip := net.ParseIP(pod.Status.PodIP); ip != nil {
			host.ips[ip.String()] = true
		}
		for _, address := range pod.Status.PodIPs {
			if ip := net.ParseIP(address.IP); ip != nil {
				host.ips[ip.String()] = true
			}
		}
		if len(host.ips) == 0 {
			return nil, fmt.Errorf("client pod %s has no pod IP", peer.name)
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

type schedulerForward struct {
	address string
	cancel  context.CancelFunc
	done    chan struct{}
	err     error
	mu      sync.Mutex
	output  string
	pending string
	ready   chan string
}

// Write consumes kubectl output while keeping only a bounded diagnostic tail.
// Using an exec writer lets Wait also wait for the output to be drained.
func (f *schedulerForward) Write(data []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.output += string(data)
	if len(f.output) > 4096 {
		f.output = f.output[len(f.output)-4096:]
	}
	f.pending += string(data)
	for {
		end := strings.IndexByte(f.pending, '\n')
		if end < 0 {
			break
		}
		line := f.pending[:end]
		f.pending = f.pending[end+1:]
		fields := strings.Fields(line)
		if len(fields) >= 3 && strings.HasPrefix(line, "Forwarding from 127.0.0.1:") {
			select {
			case f.ready <- fields[2]:
			default:
			}
		}
	}
	if len(f.pending) > 4096 {
		f.pending = f.pending[len(f.pending)-4096:]
	}
	return len(data), nil
}

func (f *schedulerForward) diagnostic() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return strings.TrimSpace(f.output)
}

func (f *schedulerForward) close() {
	f.cancel()
	<-f.done
}

func startSchedulerForward(ctx context.Context, pod string) (*schedulerForward, error) {
	ctx, cancel := context.WithCancel(ctx)
	f := &schedulerForward{cancel: cancel, done: make(chan struct{}), ready: make(chan string, 1)}
	// Ask kubectl to allocate an unused loopback port, avoiding fixed-port races.
	cmd := exec.CommandContext(ctx, "kubectl", "-n", DragonflyNamespace, "port-forward", "--address=127.0.0.1", "pod/"+pod, ":8002")
	cmd.WaitDelay = schedulerRPCTimeout
	cmd.Stdout = f
	cmd.Stderr = f
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("start port-forward for scheduler %s: %w", pod, err)
	}
	go func() {
		f.err = cmd.Wait()
		close(f.done)
	}()
	timer := time.NewTimer(schedulerCommandTimeout)
	defer timer.Stop()
	select {
	case f.address = <-f.ready:
		return f, nil
	case <-f.done:
		cancel()
		return nil, fmt.Errorf("port-forward for scheduler %s exited: %v: %s", pod, f.err, f.diagnostic())
	case <-timer.C:
		f.close()
		return nil, fmt.Errorf("port-forward for scheduler %s did not become ready: %s", pod, f.diagnostic())
	case <-ctx.Done():
		f.close()
		return nil, fmt.Errorf("port-forward for scheduler %s: %w: %s", pod, ctx.Err(), f.diagnostic())
	}
}

// WaitForSchedulerHosts independently queries every Ready scheduler until all
// supplied seed/client pods have re-announced their addresses and host types.
// It never downloads a task. ListHosts exposes registration, not the internal
// seed hashring; older schedulers may refresh that ring asynchronously.
func WaitForSchedulerHosts(ctx context.Context, peers ...*PodExec) error {
	ctx, cancel := context.WithTimeout(ctx, schedulerWaitTimeout)
	defer cancel()
	var pods []schedulerPod
	if err := pollScheduler(ctx, "wait for Ready scheduler pods", func() error {
		var err error
		pods, err = listSchedulerPods(ctx)
		if err != nil {
			return err
		}
		for _, pod := range pods {
			if !pod.ready() {
				return fmt.Errorf("scheduler pod %s is not Ready", pod.Metadata.Name)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	hosts, err := expectedSchedulerHosts(ctx, peers)
	if err != nil {
		return err
	}
	type endpoint struct {
		name    string
		forward *schedulerForward
		client  schedulerv2.SchedulerClient
	}
	endpoints := make([]endpoint, 0, len(pods))
	for _, pod := range pods {
		forward, err := startSchedulerForward(ctx, pod.Metadata.Name)
		if err != nil {
			return err
		}
		defer forward.close()
		conn, err := grpc.NewClient(forward.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return fmt.Errorf("connect scheduler %s: %w", pod.Metadata.Name, err)
		}
		defer func() { _ = conn.Close() }()
		endpoints = append(endpoints, endpoint{pod.Metadata.Name, forward, schedulerv2.NewSchedulerClient(conn)})
	}
	return pollScheduler(ctx, "wait for client hosts in every scheduler", func() error {
		for _, endpoint := range endpoints {
			select {
			case <-endpoint.forward.done:
				return fmt.Errorf("scheduler %s port-forward exited: %v: %s", endpoint.name, endpoint.forward.err, endpoint.forward.diagnostic())
			default:
			}
			rpcCtx, rpcCancel := context.WithTimeout(ctx, schedulerRPCTimeout)
			response, err := endpoint.client.ListHosts(rpcCtx, &schedulerv2.ListHostsRequest{})
			rpcCancel()
			if err != nil {
				return fmt.Errorf("scheduler %s ListHosts: %w", endpoint.name, err)
			}
			for _, expected := range hosts {
				found := false
				for _, actual := range response.GetHosts() {
					ip := net.ParseIP(actual.GetIp())
					if ip != nil && expected.ips[ip.String()] && types.HostType(actual.GetType()) == expected.typ {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("scheduler %s has not registered %s (%s, pod IPs %v)", endpoint.name, expected.name, expected.typ.Name(), expected.ips)
				}
			}
		}
		return nil
	})
}
