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
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint

	"d7y.io/dragonfly/v2/test/e2e/util"
)

var _ = Describe("Containerd with CRI support", func() {
	Context("ghcr.io/dragonflyoss/dragonfly/manager:v2.3.0 image", func() {
		It("pull should be ok", Label("containerd", "pull"), func() {
			out, err := util.CriCtlCommand("pull", "ghcr.io/dragonflyoss/dragonfly/manager:v2.3.0").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "1e280c1c7902615133dd7b1b201e0eff56b34928fbcc5e6eb89373af59e39bc6",
					Sha256: "662513304ab057e389327f7aa180e2731a19d7e6e816e26bea219909923a24e0",
				},
				{
					ID:     "80c5252267e388eb28cc6d09be7fdd2c9b0a9adde03d8564ff4be7c1b65c467a",
					Sha256: "99f960a3d990bc51293c661d50aa89af1467e4dbc17525812e3c85cc98a3ed46",
				},
				{
					ID:     "4f426384b6db5fec0fcd36a07770c68a49a0c35d508601c35b3507eca70d648e",
					Sha256: "900b06713d54060058d46f6371960ca55f2f0a765aa12ea3642c477870603e4b",
				},
				{
					ID:     "5719330d5a5afe9beea08a163aa64a8f89c33cc988497b08317740554b22d518",
					Sha256: "1fa75e2072193ff41dd869524a0498204a8cadeca337c747f6c5568b153cb752",
				},
				{
					ID:     "a289f8be37877fbd94ebb26441bef598efef5d0eedffd49ab2471d408c961336",
					Sha256: "0a9a5dfd008f05ebc27e4790db0709a29e527690c21bcbcd01481eaeb6bb49dc",
				},
				{
					ID:     "3c6ccc32bfc32bb4b3a53d5b427f2ea65fbb0d49c69a5dffd8e3fc8f18f8d83b",
					Sha256: "86664370769e9b5b7ab333bbd63120eab25bb6bd9d921cc85d5a472fec06b177",
				},
				{
					ID:     "897c55b33262cf79ac76ea6f2f34fac6df06c07d56e67d19e8f670f4317a22d5",
					Sha256: "3237f56cf6b106ed7e877bc4d8bb657d0fbbab12803abcab47be91003602572e",
				},
			}

			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			for _, taskMetadata := range taskMetadatas {
				sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, taskMetadata.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(taskMetadata.Sha256).To(Equal(sha256sum))
			}

			time.Sleep(1 * time.Second)
			seedClientPods := make([]*util.PodExec, 3)
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}
			for _, taskMetadata := range taskMetadatas {
				sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(taskMetadata.Sha256).To(Equal(sha256sum))
			}
		})

		It("rmi should be ok", Label("containerd", "rmi"), func() {
			out, err := util.CriCtlCommand("rmi", "ghcr.io/dragonflyoss/dragonfly/manager:v2.3.0").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
