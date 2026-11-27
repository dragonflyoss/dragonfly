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
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint

	"d7y.io/dragonfly/v2/test/e2e/util"
)

var _ = Describe("Schedulers Restarting", func() {
	Context("rollout restart", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize1MiB)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("download should be ok after the schedulers restart", Label("scheduler", "restart"), func() {
			// Restart every scheduler, the statefulset recreates the pods with the
			// same names and new ips, so the clients must pick them up again.
			out, err := util.KubeCtlCommand("-n", util.DragonflyNamespace, "rollout", "restart", "statefulset", util.SchedulerStatefulSetName).CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			out, err = util.KubeCtlCommand("-n", util.DragonflyNamespace, "rollout", "status", "statefulset", util.SchedulerStatefulSetName, "--timeout=10m").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			// Wait for the clients to refresh the schedulers and announce themselves
			// to the restarted schedulers.
			Eventually(func() error {
				out, err := clientPod.Command("sh", "-c", fmt.Sprintf("dfget %s --disable-back-to-source --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
				fmt.Println(string(out), err)
				return err
			}, 5*time.Minute, 15*time.Second).Should(Succeed())

			// The download above proves one scheduler tracks the seed peers again,
			// wait for the rest: the seed peers discover a restarted scheduler within
			// the dynconfig refresh interval and announce within the announce interval,
			// and the schedulers of the compatibility tests rebuild the hashring of
			// seed peers every 30 seconds after that.
			time.Sleep(90 * time.Second)

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			sha256sum, err = util.CalculateSha256ByOutput([]*util.PodExec{clientPod}, testFile.GetOutputPath())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
		})
	})

	It("should reuse seed caches after all schedulers restart", Label("scheduler", "restart", "seed-cache"), SpecTimeout(20*time.Minute), func(specCtx SpecContext) {
		// Use a unique, complete file larger than the client's small-request
		// scheduling bypass. Only seed peers see this task before the restart.
		testFile, err := util.GetFileServer().GenerateFile(32 * util.FileSize1MiB)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			Expect(util.GetFileServer().DeleteFile(testFile.GetInfo())).To(Succeed())
		})

		clientPod, err := util.ClientExec()
		Expect(err).NotTo(HaveOccurred())

		seedPods := make([]*util.PodExec, util.Servers[util.SeedClientServerName].Replicas)
		for i := range seedPods {
			seedPods[i], err = util.SeedClientExec(i)
			Expect(err).NotTo(HaveOccurred())
		}

		peerPods := append([]*util.PodExec{clientPod}, seedPods...)
		peerNames := []string{"client"}
		for i := range seedPods {
			peerNames = append(peerNames, fmt.Sprintf("seed-%d", i))
		}
		for _, pod := range peerPods {
			DeferCleanup(func() {
				cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				out, err := pod.CommandContext(cleanupCtx, "rm", "-f", testFile.GetOutputPath()).CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), string(out))
			})
		}

		By("fully caching the same task on every seed peer")
		for i, seedPod := range seedPods {
			out, err := seedPod.CommandContext(specCtx, "dfget", testFile.GetDownloadURL(), "--output", testFile.GetOutputPath()).CombinedOutput()
			fmt.Printf("preheat seed-%d: %s\n", i, out)
			Expect(err).NotTo(HaveOccurred(), string(out))

			// The helper returns the first cache it finds, so check seeds one at
			// a time instead of accidentally accepting only one populated seed.
			digest, err := util.CalculateSha256ByTaskIDContext(specCtx, []*util.PodExec{seedPod}, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(digest).To(Equal(testFile.GetSha256()), "seed-%d cache", i)
		}

		captureTraffic := func() []util.DownloadTrafficSnapshot {
			snapshots := make([]util.DownloadTrafficSnapshot, len(peerPods))
			for i, pod := range peerPods {
				snapshots[i], err = util.CaptureDownloadTraffic(pod)
				Expect(err).NotTo(HaveOccurred(), "capture %s download traffic", peerNames[i])
			}
			return snapshots
		}
		beforeRestart := captureTraffic()
		oldSchedulers, err := util.CaptureSchedulerPodUIDs()
		Expect(err).NotTo(HaveOccurred())

		By("replacing every scheduler while retaining all peer processes and caches")
		out, err := util.KubeCtlCommandContext(specCtx, "-n", util.DragonflyNamespace, "rollout", "restart", "statefulset", util.SchedulerStatefulSetName).CombinedOutput()
		fmt.Println(string(out))
		Expect(err).NotTo(HaveOccurred(), string(out))
		out, err = util.KubeCtlCommandContext(specCtx, "-n", util.DragonflyNamespace, "rollout", "status", "statefulset", util.SchedulerStatefulSetName, "--timeout=10m").CombinedOutput()
		fmt.Println(string(out))
		Expect(err).NotTo(HaveOccurred(), string(out))

		ctx, cancel := context.WithTimeout(specCtx, 5*time.Minute)
		defer cancel()
		Expect(util.WaitForSchedulerReplacement(ctx, oldSchedulers)).To(Succeed())
		Expect(util.WaitForSchedulerHosts(ctx, peerPods...)).To(Succeed())

		// Readiness checks must not download this task: retrying dfget can
		// populate the ordinary client's cache and conceal a failed recovery.
		taskID := testFile.GetTaskID()
		out, err = clientPod.CommandContext(specCtx, "sh", "-c", `test ! -e "/var/lib/dragonfly/content/tasks/$1/$2"`, "sh", taskID[:3], taskID).CombinedOutput()
		Expect(err).NotTo(HaveOccurred(), "ordinary client must have no cache for %s: %s", taskID, out)

		beforeDownload := captureTraffic()
		for i := range peerPods {
			// A peer restart could destroy the cache or reset the counters and
			// invalidate the premise of the test, even if the download succeeds.
			_, err := util.DownloadTrafficDelta(beforeRestart[i], beforeDownload[i])
			Expect(err).NotTo(HaveOccurred(), "%s must retain its process across scheduler restarts", peerNames[i])
		}

		By("downloading the target exactly once from a cold ordinary client")
		out, err = clientPod.CommandContext(specCtx, "dfget", testFile.GetDownloadURL(), "--disable-back-to-source", "--output", testFile.GetOutputPath()).CombinedOutput()
		fmt.Println(string(out))
		Expect(err).NotTo(HaveOccurred(), string(out))
		afterDownload := captureTraffic()

		By("proving P2P reuse without another source download on any peer")
		for i := range peerPods {
			delta, err := util.DownloadTrafficDelta(beforeDownload[i], afterDownload[i])
			Expect(err).NotTo(HaveOccurred(), "%s download traffic", peerNames[i])
			fmt.Printf("%s download traffic delta: %+v\n", peerNames[i], delta)
			Expect(delta.BackToSource).To(BeZero(), "%s must not redownload source pieces", peerNames[i])
			if i == 0 {
				Expect(delta.RemotePeer).To(Equal(uint64(testFile.GetInfo().Size())), "ordinary client must receive the complete file from peers")
			}
		}

		digest, err := util.CalculateSha256ByTaskIDContext(specCtx, []*util.PodExec{clientPod}, taskID)
		Expect(err).NotTo(HaveOccurred())
		Expect(digest).To(Equal(testFile.GetSha256()))
		digest, err = util.CalculateSha256ByOutputContext(specCtx, []*util.PodExec{clientPod}, testFile.GetOutputPath())
		Expect(err).NotTo(HaveOccurred())
		Expect(digest).To(Equal(testFile.GetSha256()))
	})
})
