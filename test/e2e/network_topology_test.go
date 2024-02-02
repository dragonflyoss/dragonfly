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

package e2e

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint

	"d7y.io/dragonfly/v2/test/e2e/e2eutil"
)

var _ = Describe("Evaluator with networkTopology", func() {
	Context("networkTopology", func() {
		It("check networkTopology in redis", Label("networkTopology"), func() {
			mode := os.Getenv("DRAGONFLY_COMPATIBILITY_E2E_TEST_MODE")
			if mode == schedulerCompatibilityTestMode {
				fmt.Println("networkTopology is disable, skip")
				return
			}
			Expect(waitForProbedInNetworkTopology()).Should(BeTrue())
		})
	})
})

// getRedisExec get redis pod.
func getRedisExec() *e2eutil.PodExec {
	out, err := e2eutil.KubeCtlCommand("-n", dragonflyNamespace, "get", "pod", "-l", "app.kubernetes.io/name=redis",
		"-o", "jsonpath='{range .items[0]}{.metadata.name}{end}'").CombinedOutput()
	podName := strings.Trim(string(out), "'")
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(podName)
	Expect(strings.HasPrefix(podName, "dragonfly-redis")).Should(BeTrue())
	return e2eutil.NewPodExec(dragonflyNamespace, podName, "redis")
}

func waitForProbedInNetworkTopology() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	redisPod := getRedisExec()

	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			out, err := redisPod.Command("redis-cli", "-a", "dragonfly", "-n", "3", "dbsize").CombinedOutput()
			Expect(err).NotTo(HaveOccurred())
			keyNumber, err := strconv.Atoi(strings.Split(string(out), "\n")[1])
			if keyNumber != 0 && err == nil {
				return true
			}
		}
	}
}
