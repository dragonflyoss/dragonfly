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

	"d7y.io/dragonfly/v2/test/e2e/v2/util"
)

var _ = Describe("Containerd with CRI support", func() {
	Context("ghcr.io/dragonflyoss/manager:v2.1.0 image", func() {
		It("pull should be ok", Label("containerd", "pull"), func() {
			out, err := util.CriCtlCommand("pull", "ghcr.io/dragonflyoss/manager:v2.1.0").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "f6da796ed137cd53c9759b634aa80e2b93b5c66f7dc23db016d739d63496a430",
					Sha256: "ca51217de9012bffe54390f1a91365af22a06279a3f2b3e57d4d2dc99b989588",
				},
				{
					ID:     "9176125ce8d7b67276d077ab15af1a5809080f02f16efce2a3d2c4fce600bedd",
					Sha256: "0d816dfc0753b877a04e3df93557bd3597fc7d0e308726655b14401c22a3b92a",
				},
				{
					ID:     "f5d1cbede9a750c5f3014d8aefef5f57a12ee4ee8d48e6201d19d42a69226b2c",
					Sha256: "b5941d5a445040d3a792e5be361ca42989d97fc30ff53031f3004ccea8e44520",
				},
				{
					ID:     "17d49b1b992fde658d1abfce6bae5877f48de760aba17cf573bd0e40e6a15001",
					Sha256: "c1d6d1b2d5a367259e6e51a7f4d1ccd66a28cc9940d6599d8a8ea9544dd4b4a8",
				},
				{
					ID:     "d68cb057073fb4eff2989e9110e7127d6812415db223a6658e8e71f3686b8acb",
					Sha256: "2a1bc4e0f20bb5ed9a2197ecffde7eace4a9b9179048614205d025df73ba97c7",
				},
				{
					ID:     "5ad1cd4c767450ad216b95bc438d4c228f0945cd247fb43f737ae9beac6fba4b",
					Sha256: "078ea4eebc352a499d7bb6ff65fab1325226e524acac89a9db922ad91cab88f1",
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
			out, err := util.CriCtlCommand("rmi", "ghcr.io/dragonflyoss/manager:v2.1.0").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("ghcr.io/dragonflyoss/scheduler:v2.0.0 image", func() {
		It("pull should be ok", Label("containerd", "pull"), func() {
			out, err := util.CriCtlCommand("pull", "ghcr.io/dragonflyoss/scheduler:v2.0.0").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "6f55276e68f872fba78843019fd12a8461a14ccf4ac375057680123d3e4d8cbe",
					Sha256: "0f4277a6444fbaf4eb5a7f39103e281dd57969953c7425edc7c8d4aa419347eb",
				},
				{
					ID:     "83b91ece0ebe89274fd99aed64fd4257bccde245e8da0fbbcdf7c4da8d5533de",
					Sha256: "e55b67c1d5660c34dcb0d8e6923d0a50695a4f0d94f858353069bae17d0bfdea",
				},
				{
					ID:     "6b6a238a08ddaed6594284cff238cdaf0cd0364570023e4731eefca3aa27a002",
					Sha256: "8572bc8fb8a32061648dd183b2c0451c82be1bd053a4ea8fae991436b92faebb",
				},
				{
					ID:     "8e373a2dbd9a733585db5478a0d15bab88ec3d4c60f0ce064ef4d7453c97a05f",
					Sha256: "88bfc12bad0cc91b2d47de4c7a755f6547b750256cc4c8b284e07aae13e4e041",
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
			out, err := util.CriCtlCommand("rmi", "ghcr.io/dragonflyoss/scheduler:v2.0.0").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("ghcr.io/dragonflyoss/client:v0.1.30 image", func() {
		It("pull should be ok", Label("containerd", "pull"), func() {
			out, err := util.CriCtlCommand("pull", "ghcr.io/dragonflyoss/client:v0.1.30").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "eae30e4f11c1f9e46305de15d560f161eca4339358684150bbc6b3636242fa08",
					Sha256: "c8071d0de0f5bb17fde217dafdc9d2813ce9db77e60f6233bcd32f1c8888b121",
				},
				{
					ID:     "fc4538dc3b9f25e601a59302e6edcf7cb9eaf97c4ecda5fd3f600de736e3e5dd",
					Sha256: "e964513726885fa2f977425fc889eabbe25c9fa47e7a4b0ec5e2baef96290f47",
				},
				{
					ID:     "53cc7eb8cfe8f396d7dcaa25163bdd65720e4ac5d39a70d8a595abde26b86f22",
					Sha256: "0e304933d7eae4674e05b3bc409f236c65077e2b7055119bbd66ff613fe5e1ad",
				},
				{
					ID:     "8961dce0ceecdf7eeace18af8c9e5404b279b5cd541f9220c42e049c738d00b9",
					Sha256: "53b01ef3d5d676a8514ded6b469932e33d84738e5e00932ca124382a8567c44b",
				},
				{
					ID:     "829047ff8b3ef37c68764e73ad50dfbf3e657419fa6bf1b373c46ad5dfc48768",
					Sha256: "c9d959fc168ad8bdc9a021066eb9c1dd4de8e860c03619a88d8ba0ff5479d9ea",
				},
				{
					ID:     "7f784cfc82e45628e6cd35fae45b61a0e9ce66b772542f03d650534a34a8297c",
					Sha256: "b6acfae843b58bf14369ebbeafa96af5352cde9a89f8255ca51f92b233a6e405",
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
			out, err := util.CriCtlCommand("rmi", "ghcr.io/dragonflyoss/client:v0.1.30").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("ghcr.io/dragonflyoss/dfinit:v0.1.30 image", func() {
		It("pull should be ok", Label("containerd", "pull"), func() {
			out, err := util.CriCtlCommand("pull", "ghcr.io/dragonflyoss/dfinit:v0.1.30").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())

			taskMetadatas := []util.TaskMetadata{
				{
					ID:     "ec17104ac49cf85fced3652f854e09b0cf0948fb232b136e3f57be3f9aeb77b2",
					Sha256: "c58d97dd21c3b3121f262a1fbb5a278f77ab85dba7a02b819e710f34683cf746",
				},
				{
					ID:     "2453086dcc76db72d09a710be96625d53a41c78747fd56c6ea1f2f340691b731",
					Sha256: "2ff0ae26fa61a2b0f88f470a8e50f7623ea48b224eb072a5878a20d663d5307d",
				},
				{
					ID:     "4f89c2f02c36714b19cde131a654babc558c261abb6ee2c94ff54b5fcd4c8daf",
					Sha256: "b1826117441e607acd3b98c93cdb16759c2cc2240852055b8a2b5860f3204f1e",
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
			out, err := util.CriCtlCommand("rmi", "ghcr.io/dragonflyoss/dfinit:v0.1.30").CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
