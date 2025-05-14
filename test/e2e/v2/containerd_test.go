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
					ID:     "1f047bc5cd298700f1190641992556c3958ad3bf389105dedd537e7913b3f8dc",
					Sha256: "ca51217de9012bffe54390f1a91365af22a06279a3f2b3e57d4d2dc99b989588",
				},
				{
					ID:     "77e5c810d83f5be37314b512ab4acd78341e923c50d4da47c9e1199abfd8da9f",
					Sha256: "0d816dfc0753b877a04e3df93557bd3597fc7d0e308726655b14401c22a3b92a",
				},
				{
					ID:     "562184a74d7dc841ce714921d394a95d540534bd8518b20ea0f2bac9ab040694",
					Sha256: "b5941d5a445040d3a792e5be361ca42989d97fc30ff53031f3004ccea8e44520",
				},
				{
					ID:     "ca428437ad023f7f6e1174d95b1d352d3c834846a0f1de5ae8dce5cf3d4f5aa5",
					Sha256: "2a1bc4e0f20bb5ed9a2197ecffde7eace4a9b9179048614205d025df73ba97c7",
				},
				{
					ID:     "53a25c54374bec29c3b0f313f421ae5d139b27ed3b51bfe1fbfeaa86a08937b5",
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
					ID:     "fd68f82ea453803d0e9fb2976f5f034e0ab90d6aa5856b1adc40039212d92aa8",
					Sha256: "0f4277a6444fbaf4eb5a7f39103e281dd57969953c7425edc7c8d4aa419347eb",
				},
				{
					ID:     "785f97d679062050a12d0489104f13a3e8fe68922e17b95aaed6249d0b983ab9",
					Sha256: "e55b67c1d5660c34dcb0d8e6923d0a50695a4f0d94f858353069bae17d0bfdea",
				},
				{
					ID:     "04f978cb5d8062dfccf0e70774c05c8f2c95ee385b3349fe6efbb6d2b15c6d02",
					Sha256: "8572bc8fb8a32061648dd183b2c0451c82be1bd053a4ea8fae991436b92faebb",
				},
				{
					ID:     "95c7239527dd27df223708d2fa4a9ca2d764d44341ebe4dab600d346502058be",
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
					ID:     "4977d204abab1b21a677daf2f705934caa7f5f21c364a7625ad230c933b41e72",
					Sha256: "c8071d0de0f5bb17fde217dafdc9d2813ce9db77e60f6233bcd32f1c8888b121",
				},
				{
					ID:     "cd9e68218c4f69d8b71cec3e9d9342b9665088994dd2049e733a430aaf564d8b",
					Sha256: "e964513726885fa2f977425fc889eabbe25c9fa47e7a4b0ec5e2baef96290f47",
				},
				{
					ID:     "94ad52f16ab43c874ab6a0a6480aefe34970139e9c444472fbb01d76b0ec5bda",
					Sha256: "0e304933d7eae4674e05b3bc409f236c65077e2b7055119bbd66ff613fe5e1ad",
				},
				{
					ID:     "a294e4746f268f08345c522c03e86e3c787b22ff27c2de8cabffd52d14080226",
					Sha256: "53b01ef3d5d676a8514ded6b469932e33d84738e5e00932ca124382a8567c44b",
				},
				{
					ID:     "0036c704a45de385e36c7752e4a3496e9f3f50c37ac5cc4ce737d57aceb3eb12",
					Sha256: "c9d959fc168ad8bdc9a021066eb9c1dd4de8e860c03619a88d8ba0ff5479d9ea",
				},
				{
					ID:     "efcea01bf4082b7425e027a436e6a284d6e46860e053e673678758a906f28660",
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
					ID:     "0721e454cdc3cf4a6285b8f44638387af69475cc76dcf81b9cd4ae618268b921",
					Sha256: "c58d97dd21c3b3121f262a1fbb5a278f77ab85dba7a02b819e710f34683cf746",
				},
				{
					ID:     "fcae534b9d18697437c096c9442a1800e427a7d8aed80de335a93bef42d881de",
					Sha256: "2ff0ae26fa61a2b0f88f470a8e50f7623ea48b224eb072a5878a20d663d5307d",
				},
				{
					ID:     "1ce41bca3caca3d74f262b6e2fbccca6629714f6a31c26fe5c3fe941566b7777",
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
