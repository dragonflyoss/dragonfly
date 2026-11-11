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
					ID:     "a5040eb77de7f771cb3ce3ecb2ebb61af124d3341e0b5c6854b7e220eb0dc680",
					Sha256: "a711f05d33845e2e9deffcfcc5adf082d7c6e97e3e3a881d193d9aae38f092a8",
				},
				{
					ID:     "f9f24ea0c08c3637d2d5770fbf80a201f69482226de7cf4490bb1c540ac51b37",
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
					ID:     "f435132946228f1e6492959df897a0e566a65e1b0c38add695ab96815b92ab11",
					Sha256: "2ada26cb7c739b559bd588eb535e921a9378ab00881f57f374a13ac56c6cb917",
				},
				{
					ID:     "36392811b51b808ac4b057c7915ac18116c726c5dd4b488b650841fe85e9e1a8",
					Sha256: "3c5439d0ae62993af4cc4e21b8c777297e81bc40e200d4633c85a3e0f977a3bc",
				},
				{
					ID:     "53aa11bf52bb7b224dc843d3465400f3718b5658aec7ca855daebe4ab45ff9c8",
					Sha256: "d2154306f880c2808019d75c3e98335bb1ec3d1504954a571366871f2e2bb70f",
				},
				{
					ID:     "d3297c34a454c520f2bab867f7af0d070f5f5f74dcaa7bcc8c731ed82f935e70",
					Sha256: "8bb494b4454f9cd978fcc7157a17b063bc53d39a592677e7a65893ca0a19ace3",
				},
				{
					ID:     "84a8a4250c11092f3ce657987d586faa1c28b614e944ea894e443487d1bee5e4",
					Sha256: "4cf2c707ac17a61f90e03c69f8acf5ca8058857992a26a210f80aaaeac3c2c96",
				},
				{
					ID:     "eda414c80da060fd075a5816393db73899042e61018aa105b2512e185441f80d",
					Sha256: "3d3277d8190820abbb1d03f9b9953a73f6343acc95ee1d8b07ce2c47e5c520ad",
				},
				{
					ID:     "0d33568fc34073132c3edcae34f8ebdfb08e00a54d9312204b2be368bb93a075",
					Sha256: "7a139ea27144cbf092ace5aff1ae4d985960fdcb5c24f1b6901aba003c82527f",
				},
				{
					ID:     "7154e17db0a98a4bc13d67c350df6b0e158f556d3e405df3af7c01c42fdaf037",
					Sha256: "3d1a614f216826d6da7f6aaae88504323edc2f5e0167706fd7aa0feb363202dc",
				},
				{
					ID:     "d932f4beee2a70e2d8bc50546dd68891fdfad4f4ccbd1062ed73dc7172a5e652",
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
					ID:     "dea387e7a58987cdc6f430e0c7f4849d51ec67df15e7c86b8280cdaeac4f5064",
					Sha256: "1e6cd1c82932fae5aeb6f2573d6637ad9951212d17c0d8a39f3b99d7868cf3b7",
				},
				{
					ID:     "425b0c891572f70fedcff62f15684906cb97cfa15a494b2dd3617c976647a3c5",
					Sha256: "f791e47cf52abf7b364878f2141547e3781072e437f599f6632849d88fca8ca4",
				},
				{
					ID:     "7ca25ce03f2adc739d1b1942b5f393f30388afbdfae5fb7b73c7041d02e7a190",
					Sha256: "11b1f920f55e03df30ae0d398f2a8f4f5178e50126fff50b9299e13b96f7140d",
				},
				{
					ID:     "862be127e77273e82c29dd3b3032fc12b5e0a467b3603469c7ac5378fe60ee58",
					Sha256: "2f035ef55838059a85d0d80c8e5370b91b578515272be8f2df7f65cc050e29b1",
				},
				{
					ID:     "a8864cabb27331693fac1c99895f81a9e2b8f052d73deb53de91a9d082b488d5",
					Sha256: "baf53db0cd0097d0e142824dad138c1143e20a6fcdad6cfa371c93c2de3dea18",
				},
				{
					ID:     "6cb5b485da002fea45a28df78dbbaafcf05fab9d62a920d93a3baa7147107278",
					Sha256: "205ef3ca16a253d4537b42957a54b21504ed6ce940668a0ee275c43ccb5ffe72",
				},
				{
					ID:     "5c1b63850953da3d6be0d55f32eb31e10288ce193e328e54cb6dd580ac36f979",
					Sha256: "a25d21ddcde045c95cc342b72c1b17b13acb41d619fdc880e41b924b671e98a1",
				},
				{
					ID:     "112ad58251b94d91e3ba1884e41c167a727ceac1b8e209358a6c2507bd1cf8aa",
					Sha256: "dfd2ce883983b79cf790d562172c24d37210e7345466cff49bd0f0bf2deaf789",
				},
				{
					ID:     "976839ff5ded8efcce560e21fe39b25133974c301c620604619f428daf0fd438",
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
