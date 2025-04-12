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

package manager

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint

	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/structure"
	"d7y.io/dragonfly/v2/test/e2e/v2/util"
)

var _ = Describe("Preheat with Manager", func() {
	Context("1MiB file", func() {
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

		It("preheat files should be ok", Label("preheat", "file"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type: "file",
					URL:  testFile.GetDownloadURL(),
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://127.0.0.1:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			seedClientPods := make([]*util.PodExec, 3)
			for i := 0; i < 3; i++ {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
		})
	})

	Context("10MiB file", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize10MiB)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("preheat files should be ok", Label("preheat", "file"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type: "file",
					URL:  testFile.GetDownloadURL(),
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			seedClientPods := make([]*util.PodExec, 3)
			for i := 0; i < 3; i++ {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
		})
	})

	Context("100MiB file", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize100MiB)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("preheat files should be ok", Label("preheat", "file"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type: "file",
					URL:  testFile.GetDownloadURL(),
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			seedClientPods := make([]*util.PodExec, 3)
			for i := 0; i < 3; i++ {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
		})
	})

	Context("10MiB file in cache", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize10MiB)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("preheat files in cache should be ok", Label("preheat", "file", "cache"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type:        "file",
					URL:         testFile.GetDownloadURL(),
					Scope:       "single_seed_peer",
					LoadToCache: true,
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			var preheatedSeedClient *util.PodExec
			var preheatedSeedClientIndex int
			var otherSeedClientIndices []int

			for i := 0; i < 3; i++ {
				seedClient, err := util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())

				taskIDCmd := fmt.Sprintf("grep -a '%s' /var/log/dragonfly/dfdaemon/dfdaemon.log", testFile.GetTaskID())
				successCmd := fmt.Sprintf("%s | grep -a 'download task succeeded'", taskIDCmd)

				out, err = seedClient.Command("bash", "-c", successCmd).CombinedOutput()
				if err == nil && len(out) > 0 {
					preheatedSeedClient = seedClient
					preheatedSeedClientIndex = i
					fmt.Printf("Found preheated seed client: %d\n", i)
					break
				}
			}

			Expect(preheatedSeedClient).NotTo(BeNil())

			for i := 0; i < 3; i++ {
				if i != preheatedSeedClientIndex {
					otherSeedClientIndices = append(otherSeedClientIndices, i)
				}
			}

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{preheatedSeedClient}, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			req, err = structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type:  "file",
					URL:   testFile.GetDownloadURL(),
					Scope: "all_seed_peers",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err = managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job = &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done = waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			seedClientPods := make([]*util.PodExec, 2)
			for i, idx := range otherSeedClientIndices {
				seedClientPods[i], err = util.SeedClientExec(idx)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err = util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			out, err = preheatedSeedClient.Command("cat", "/var/log/dragonfly/dfdaemon/dfdaemon.log").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			logs := string(out)

			Expect(logs).To(ContainSubstring(fmt.Sprintf("put task to cache: %s", testFile.GetTaskID())))

			pieceRegex := regexp.MustCompile(`pieces: \[([0-9, ]+)\]`)
			pieceMatches := pieceRegex.FindStringSubmatch(logs)
			Expect(pieceMatches).NotTo(BeNil())
			pieceNumbers := strings.Split(strings.ReplaceAll(pieceMatches[1], " ", ""), ",")

			for _, number := range pieceNumbers {
				pieceID := fmt.Sprintf("%s-%s", testFile.GetTaskID(), number)
				putPieceCacheLog := fmt.Sprintf("put piece to cache: %s", pieceID)
				Expect(logs).To(ContainSubstring(putPieceCacheLog))

				getPieceCacheLog := fmt.Sprintf("get piece from cache: %s", pieceID)
				Expect(logs).To(ContainSubstring(getPieceCacheLog))
			}
		})
	})

	Context("ghcr.io/dragonflyoss/busybox:v1.35.0 image", func() {
		It("preheat image should be ok", Label("preheat", "image"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type: "image",
					URL:  "https://ghcr.io/v2/dragonflyoss/busybox/manifests/1.35.0",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			Expect(err).NotTo(HaveOccurred())
			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "50f91088cd60a49bf565396c1f7ba1f8b9d32c9f1541d2e06b085c780030805a",
					Sha256: "a711f05d33845e2e9deffcfcc5adf082d7c6e97e3e3a881d193d9aae38f092a8",
				},
				{
					ID:     "a048c8c66a652030fba31514735ff51288c54f72f878e363948303ded065e199",
					Sha256: "f643e116a03d9604c344edb345d7592c48cc00f2a4848aaf773411f4fb30d2f5",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				seedClientPods := make([]*util.PodExec, 3)
				for i := 0; i < 3; i++ {
					seedClientPods[i], err = util.SeedClientExec(i)
					fmt.Println(err)
					Expect(err).NotTo(HaveOccurred())
				}

				sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(taskMetadata.Sha256).To(Equal(sha256sum))
			}
		})
	})

	Context("ghcr.io/dragonflyoss/scheduler:v2.1.0 image", func() {
		It("preheat image for linux/amd64 platform should be ok", Label("preheat", "image"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type:     "image",
					URL:      "https://ghcr.io/v2/dragonflyoss/scheduler/manifests/v2.1.0",
					Platform: "linux/amd64",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			Expect(err).NotTo(HaveOccurred())
			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "1a2cf94649557d3e978f355bcc8eaffadd4c8d41e334fc93adc3ea49488d1200",
					Sha256: "f1f1039835051ecc04909f939530e86a20f02d2ce5ad7a81c0fa3616f7303944",
				},
				{
					ID:     "ba71de48fb521ec52a4dff60106af10512eb21a40b8fd1756467af19592e6a74",
					Sha256: "871ab018db94b4ae7b137764837bc4504393a60656ba187189e985cd809064f7",
				},
				{
					ID:     "3cf396cec1d067374714655cf98f5e505626d1b6cdec094b1469b3c98569e3d1",
					Sha256: "f1a1d290795d904815786e41d39a41dc1af5de68a9e9020baba8bd83b32d8f95",
				},
				{
					ID:     "b04072342af19c13fcda146a71727edd03442d9851ca185d297f53442ba12a22",
					Sha256: "f1ffc4b5459e82dc8e7ddd1d1a2ec469e85a1f076090c22851a1f2ce6f71e1a6",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				seedClientPods := make([]*util.PodExec, 3)
				for i := 0; i < 3; i++ {
					seedClientPods[i], err = util.SeedClientExec(i)
					fmt.Println(err)
					Expect(err).NotTo(HaveOccurred())
				}

				sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(taskMetadata.Sha256).To(Equal(sha256sum))
			}
		})

		It("preheat image for linux/arm64 platform should be ok", Label("preheat", "image"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type:     "image",
					URL:      "https://ghcr.io/v2/dragonflyoss/scheduler/manifests/v2.2.0",
					Platform: "linux/arm64",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job := &models.Job{}
			err = json.Unmarshal(out, job)
			Expect(err).NotTo(HaveOccurred())
			done := waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "d83d909b2ec5660323a4225149fb87886f9897cf6098ffcb108a77ac252facaf",
					Sha256: "9986a736f7d3d24bb01b0a560fa0f19c4b57e56c646e1f998941529d28710e6b",
				},
				{
					ID:     "0eb92ed9ea3561063efceaf414e2a5431410d5ba76523bb2e0a4793fdf9a38b1",
					Sha256: "f7307687fd72fb79eadd7f38f8cb9675b76480e32365a5d282a06f788944e9f2",
				},
				{
					ID:     "45e6cddeb6d117e61312e3913779b88e3ac8c37f8df2f0814bf2d8cda6f784d7",
					Sha256: "fc5951fb196d09e569f4592b50e3a71ad01d11da229b8a500fea278eba0170c5",
				},
				{
					ID:     "e1ab22f7ff59940064e2ffde854fe3e0fcd23267268ff382696875c49844c985",
					Sha256: "c7c72808bf776cd122bdaf4630a4a35ea319603d6a3b6cbffddd4c7fd6d2d269",
				},
				{
					ID:     "56ee9d6d3bcac9e6543073e175ba2a0ee5d37c27d6e1a7daabeea44b7fbdd4a4",
					Sha256: "edbf1aa1d62d9c17605c1ee2d9dff43489bc0f8ae056367734386c35bfae226a",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				seedClientPods := make([]*util.PodExec, 3)
				for i := 0; i < 3; i++ {
					seedClientPods[i], err = util.SeedClientExec(i)
					fmt.Println(err)
					Expect(err).NotTo(HaveOccurred())
				}

				sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(taskMetadata.Sha256).To(Equal(sha256sum))
			}
		})
	})
})
