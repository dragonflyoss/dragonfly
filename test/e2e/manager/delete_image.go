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

var _ = Describe("DeleteImage with Manager", func() {
	Context("ghcr.io/dragonflyoss/client:v1.5.12 image", func() {
		It("delete image should keep the layers referenced by other images", Label("preheat", "delete", "image"), func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			for _, url := range []string{
				"https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				"https://ghcr.io/v2/dragonflyoss/client/manifests/sha256:cf54d0f3bda39332435c109aed0689a3a0d31f0a00e570cf09aae10dfdf5da9b",
			} {
				req, err := structure.StructToMap(types.CreatePreheatJobRequest{
					Type: internaljob.PreheatJob,
					Args: types.PreheatArgs{
						Type:     "image",
						URL:      url,
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
			}

			seedClientPods := make([]*util.PodExec, 3)
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "8a90ec4f3df5df15c4e16536b65ce8c8fb119c93ca88d979e669802801c63ff0",
					Sha256: "9eb2e5ef987d268ebe2d9f1617fd3398f5a170aec8546574ba72d7e2362f89c1",
				},
				{
					ID:     "c62a3077a84e6240511462a96c7a663984af0dda4636d8c8e1a1e83ede969737",
					Sha256: "dd051e1b8c92ff101fa3071ab97ea4924885942f15644e01c3dc095764e863e3",
				},
				{
					ID:     "8baa97095c3007702a3f67778ae145d839e4fc619d78f5c70f03012788f7b576",
					Sha256: "524382ca7728c7aedabc0faecbe7dadb6f4f0a1dc1bbdbe9d14e00307cae36fa",
				},
				{
					ID:     "0ff3cdd4a102091086fe05777c4e4eda9b068f8e7995224f8be967e5871c7445",
					Sha256: "64bbbf0249fd6001f88d224e6c752712d4d4e156e2ec3a388791051d8070675c",
				},
				{
					ID:     "826844ba48fc65834d77526d1fcd474289903da839c98af65f23a7220bd91116",
					Sha256: "a4d359a97e3b1f92542a0fa7fdc5e6c1566bb208276d8879460ba3b1a2783f63",
				},
				{
					ID:     "a162ed06ba97d9fca583707c0c7aba7ae16bc30da0a539a4c45fecb8139018cb",
					Sha256: "f7d8f0305de41ab5f3a54fe2c23c8da74e16fccbdec15933da32d4343715e65f",
				},
				{
					ID:     "cade76879f28877e3d9ddb92bee759b37021ebc2971cafa4ebb72edaf5f9e574",
					Sha256: "b866e5602329fd674dd5e831b94910ad6f881b232fd2d273da8538668c6b86f9",
				},
				{
					ID:     "e0c2af65866697edd06e26ce3503eeb234a4399246cab826269257fd7e7485ed",
					Sha256: "4b4cbff4fd9b9d8b0beb4d61b91e6497b3d9bf16742e23ed1d46b164074a6740",
				},
				{
					ID:     "77be657e62346dcd482a3803589922476c827bb8e18b84000d0a755c61e3d156",
					Sha256: "f1b3ecfe00babdadef5ab7b857343e016a684a33f714defdc3e41a99274367a9",
				},
			}

			for _, taskMetadata := range taskMetadatas {
				sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(taskMetadata.Sha256).To(Equal(sha256sum))
			}

			for _, seedClientPod := range seedClientPods {
				out, err := seedClientPod.Command("dfctl", "task", "rm", "oci://ghcr.io/dragonflyoss/client:v1.5.12").CombinedOutput()
				fmt.Println(string(out))
				Expect(err).NotTo(HaveOccurred())
			}

			for _, taskMetadata := range taskMetadatas {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeTrue())
			}

			for _, seedClientPod := range seedClientPods {
				out, err := seedClientPod.Command("dfctl", "task", "rm", "oci://ghcr.io/dragonflyoss/client@sha256:cf54d0f3bda39332435c109aed0689a3a0d31f0a00e570cf09aae10dfdf5da9b").CombinedOutput()
				fmt.Println(string(out))
				Expect(err).NotTo(HaveOccurred())
			}

			for _, taskMetadata := range taskMetadatas {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeFalse())
			}
		})
	})
})
