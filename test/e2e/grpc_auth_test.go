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

package e2e

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	managerv1 "d7y.io/api/v2/pkg/apis/manager/v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"d7y.io/dragonfly/v2/test/e2e/util"
)

/*
The gRPC token authentication configuration is located in
`test/testdata/charts/config-jwt-auth.yaml`. Manager and scheduler share the
same secret and issuer; the scheduler attaches a token to every call to the
manager through the TokenProvider (see pkg/rpc/auth).

This suite verifies that, with authentication enabled:

  - The scheduler's dynconfig fetch and announcer heartbeat towards the
    manager keep working (the manager console still sees the schedulers).
  - Calls without a token are rejected with codes.Unauthenticated.
  - The gRPC health checking service stays exempt so that probes are
    unaffected.
*/

// managerGRPCTarget is the manager gRPC service endpoint of the e2e chart.
const managerGRPCTarget = "dragonfly-manager.dragonfly-system.svc:65003"

var _ = Describe("gRPC Token Authentication", Label("grpc-auth"), func() {
	Context("scheduler talking to manager with shared secret", func() {
		It("scheduler should announce itself to manager", Label("grpc-auth", "dynconfig"), func() {
			// The scheduler announcer registers itself in the manager
			// database through authenticated gRPC calls (dynconfig fetch
			// and announce). Poll the manager REST API until it reports
			// schedulers of cluster 1.
			Eventually(func() error {
				pod, err := util.ManagerExec(0)
				if err != nil {
					return err
				}

				out, err := pod.CurlCommand("GET", map[string]string{
					"Authorization": "Bearer " + util.GetPersonalAccessToken(),
				}, nil, "http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/schedulers?scheduler_cluster_id=1").CombinedOutput()
				if err != nil {
					return fmt.Errorf("curl manager error: %s, output: %s", err, string(out))
				}

				if !strings.Contains(string(out), "\"id\"") {
					return fmt.Errorf("no scheduler announced in manager, output: %s", string(out))
				}

				return nil
			}).WithTimeout(3 * time.Minute).WithPolling(10 * time.Second).Should(Succeed())
		})

		It("scheduler pods should stay ready without restarts", Label("grpc-auth"), func() {
			// Authentication must not break liveness/readiness probes: the
			// gRPC health service is always exempt from token checks, so a
			// failing scheduler (crash looping on auth errors) shows up as
			// restarts.
			Expect(util.Servers).To(HaveKey(util.SchedulerServerName))

			for i := range util.Servers[util.SchedulerServerName].Replicas {
				out, err := util.KubeCtlCommand("-n", util.DragonflyNamespace,
					"get", "pod", "-l", fmt.Sprintf("component=%s", util.SchedulerServerName),
					"-o", fmt.Sprintf("jsonpath='{.items[%d].status.containerStatuses[0].restartCount}'", i)).CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), string(out))
				Expect(strings.Trim(string(out), "'")).To(Equal("0"))
			}
		})
	})

	Context("unauthenticated gRPC calls", func() {
		It("should be rejected by the manager", Label("grpc-auth"), func() {
			// Dial the manager gRPC port from the test process. The e2e
			// runner is outside the cluster but the port is reachable
			// through the service type used by the chart.
			conn, err := grpc.NewClient(managerGRPCTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			client := managerv1.NewManagerClient(conn)
			_, err = client.GetScheduler(context.Background(), &managerv1.GetSchedulerRequest{Hostname: "e2e-auth-test"})
			Expect(err).To(HaveOccurred())
			Expect(status.Code(err)).To(Equal(codes.Unauthenticated))
		})

		It("should leave the health service reachable", Label("grpc-auth"), func() {
			// The health checking service is always exempt: no token, but
			// still a valid gRPC health response.
			conn, err := grpc.NewClient(managerGRPCTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			resp, err := healthpb.NewHealthClient(conn).Check(context.Background(), &healthpb.HealthCheckRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetStatus()).To(Equal(healthpb.HealthCheckResponse_SERVING))
		})
	})
})
