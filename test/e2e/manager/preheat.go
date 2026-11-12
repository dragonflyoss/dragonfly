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

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint

	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/models"
	"d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/pkg/structure"
	"d7y.io/dragonfly/v2/test/e2e/util"
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
					URLs: []string{testFile.GetDownloadURL()},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + util.GetPersonalAccessToken()}, req,
				"http://127.0.0.1:8080/oapi/v1/jobs").CombinedOutput()
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
			for i := range 3 {
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
					URLs: []string{testFile.GetDownloadURL()},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + util.GetPersonalAccessToken()}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/oapi/v1/jobs").CombinedOutput()
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
			for i := range 3 {
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
					URLs: []string{testFile.GetDownloadURL()},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + util.GetPersonalAccessToken()}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/oapi/v1/jobs").CombinedOutput()
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
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
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

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + util.GetPersonalAccessToken()}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/oapi/v1/jobs").CombinedOutput()
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
					ID:     "a711f05d33845e2e9deffcfcc5adf082d7c6e97e3e3a881d193d9aae38f092a8",
					Sha256: "a711f05d33845e2e9deffcfcc5adf082d7c6e97e3e3a881d193d9aae38f092a8",
				},
				{
					ID:     "f643e116a03d9604c344edb345d7592c48cc00f2a4848aaf773411f4fb30d2f5",
					Sha256: "f643e116a03d9604c344edb345d7592c48cc00f2a4848aaf773411f4fb30d2f5",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				seedClientPods := make([]*util.PodExec, 3)
				for i := range 3 {
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

	Context("ghcr.io/dragonflyoss/client:v1.5.11 image", func() {
		It("preheat image for linux/amd64 platform should be ok", Label("preheat", "image"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type:     "image",
					URL:      "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.11",
					Platform: "linux/amd64",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + util.GetPersonalAccessToken()}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/oapi/v1/jobs").CombinedOutput()
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
					ID:     "2ada26cb7c739b559bd588eb535e921a9378ab00881f57f374a13ac56c6cb917",
					Sha256: "2ada26cb7c739b559bd588eb535e921a9378ab00881f57f374a13ac56c6cb917",
				},
				{
					ID:     "3c5439d0ae62993af4cc4e21b8c777297e81bc40e200d4633c85a3e0f977a3bc",
					Sha256: "3c5439d0ae62993af4cc4e21b8c777297e81bc40e200d4633c85a3e0f977a3bc",
				},
				{
					ID:     "d2154306f880c2808019d75c3e98335bb1ec3d1504954a571366871f2e2bb70f",
					Sha256: "d2154306f880c2808019d75c3e98335bb1ec3d1504954a571366871f2e2bb70f",
				},
				{
					ID:     "8bb494b4454f9cd978fcc7157a17b063bc53d39a592677e7a65893ca0a19ace3",
					Sha256: "8bb494b4454f9cd978fcc7157a17b063bc53d39a592677e7a65893ca0a19ace3",
				},
				{
					ID:     "4cf2c707ac17a61f90e03c69f8acf5ca8058857992a26a210f80aaaeac3c2c96",
					Sha256: "4cf2c707ac17a61f90e03c69f8acf5ca8058857992a26a210f80aaaeac3c2c96",
				},
				{
					ID:     "3d3277d8190820abbb1d03f9b9953a73f6343acc95ee1d8b07ce2c47e5c520ad",
					Sha256: "3d3277d8190820abbb1d03f9b9953a73f6343acc95ee1d8b07ce2c47e5c520ad",
				},
				{
					ID:     "7a139ea27144cbf092ace5aff1ae4d985960fdcb5c24f1b6901aba003c82527f",
					Sha256: "7a139ea27144cbf092ace5aff1ae4d985960fdcb5c24f1b6901aba003c82527f",
				},
				{
					ID:     "3d1a614f216826d6da7f6aaae88504323edc2f5e0167706fd7aa0feb363202dc",
					Sha256: "3d1a614f216826d6da7f6aaae88504323edc2f5e0167706fd7aa0feb363202dc",
				},
				{
					ID:     "619599c4b948cfc6dc2573c78ceed92634701a8043797f492a24adbe47dcc9bb",
					Sha256: "619599c4b948cfc6dc2573c78ceed92634701a8043797f492a24adbe47dcc9bb",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				seedClientPods := make([]*util.PodExec, 3)
				for i := range 3 {
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
					URL:      "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.11",
					Platform: "linux/arm64",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			out, err := managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json", "Authorization": "Bearer " + util.GetPersonalAccessToken()}, req,
				"http://dragonfly-manager.dragonfly-system.svc:8080/oapi/v1/jobs").CombinedOutput()
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
					ID:     "1e6cd1c82932fae5aeb6f2573d6637ad9951212d17c0d8a39f3b99d7868cf3b7",
					Sha256: "1e6cd1c82932fae5aeb6f2573d6637ad9951212d17c0d8a39f3b99d7868cf3b7",
				},
				{
					ID:     "f791e47cf52abf7b364878f2141547e3781072e437f599f6632849d88fca8ca4",
					Sha256: "f791e47cf52abf7b364878f2141547e3781072e437f599f6632849d88fca8ca4",
				},
				{
					ID:     "11b1f920f55e03df30ae0d398f2a8f4f5178e50126fff50b9299e13b96f7140d",
					Sha256: "11b1f920f55e03df30ae0d398f2a8f4f5178e50126fff50b9299e13b96f7140d",
				},
				{
					ID:     "2f035ef55838059a85d0d80c8e5370b91b578515272be8f2df7f65cc050e29b1",
					Sha256: "2f035ef55838059a85d0d80c8e5370b91b578515272be8f2df7f65cc050e29b1",
				},
				{
					ID:     "baf53db0cd0097d0e142824dad138c1143e20a6fcdad6cfa371c93c2de3dea18",
					Sha256: "baf53db0cd0097d0e142824dad138c1143e20a6fcdad6cfa371c93c2de3dea18",
				},
				{
					ID:     "205ef3ca16a253d4537b42957a54b21504ed6ce940668a0ee275c43ccb5ffe72",
					Sha256: "205ef3ca16a253d4537b42957a54b21504ed6ce940668a0ee275c43ccb5ffe72",
				},
				{
					ID:     "a25d21ddcde045c95cc342b72c1b17b13acb41d619fdc880e41b924b671e98a1",
					Sha256: "a25d21ddcde045c95cc342b72c1b17b13acb41d619fdc880e41b924b671e98a1",
				},
				{
					ID:     "dfd2ce883983b79cf790d562172c24d37210e7345466cff49bd0f0bf2deaf789",
					Sha256: "dfd2ce883983b79cf790d562172c24d37210e7345466cff49bd0f0bf2deaf789",
				},
				{
					ID:     "c94f8f7656ae6024f078471006c270dc4bababbe77d9f39b598663896bb8a178",
					Sha256: "c94f8f7656ae6024f078471006c270dc4bababbe77d9f39b598663896bb8a178",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				seedClientPods := make([]*util.PodExec, 3)
				for i := range 3 {
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
