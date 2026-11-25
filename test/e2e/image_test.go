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
	"slices"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	schedulerv2 "d7y.io/api/v2/pkg/apis/scheduler/v2"

	managertypes "d7y.io/dragonfly/v2/manager/types"
	"d7y.io/dragonfly/v2/test/e2e/util"
)

var _ = Describe("Preheat and Delete Image Using Scheduler", func() {
	Context("ghcr.io/dragonflyoss/client:v1.5.12 image", func() {
		It("preheat image for linux/arm64 platform and all platforms should be ok", Label("preheat", "image"), func() {
			conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			seedClientPods := make([]*util.PodExec, 3)
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			amd64Platform := "linux/amd64"
			arm64Platform := "linux/arm64"
			schedulerClient := schedulerv2.NewSchedulerClient(conn)
			_, err = schedulerClient.PreheatImage(context.Background(), &schedulerv2.PreheatImageRequest{
				Url:      "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Platform: &arm64Platform,
				Scope:    managertypes.SingleSeedPeerScope,
			})
			Expect(err).NotTo(HaveOccurred())

			amd64TaskMetadatas := []util.TaskMetadata{
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

			arm64TaskMetadatas := []util.TaskMetadata{
				{
					ID:     "46b5a30713ed0772323dd02ebd5b50bd78c3d6c95c43ec161e222294c153eaa9",
					Sha256: "901f53de4911302661c0199f318d976e6f323d9a4a8927d37266c0e18db1df60",
				},
				{
					ID:     "51599e93d000381ccd553eb165824c6a430b56b980f172987bc6cb338a72593b",
					Sha256: "2d66c782e4847b8495b17a4b86be1159aa13d140fe4fd17cb1d536953d5114c6",
				},
				{
					ID:     "734826974c2ee8b7715c5ff585038cd7f03e13d874fa4c3c1b7a8d0870569f3b",
					Sha256: "5d33eb8f4cd8a1e776847144f5b26534d2646d9a08311cebf87448d4660b1443",
				},
				{
					ID:     "d515fa4dee4da892c373b4dc27e9c0317f603fea183d37063a326645c941fa02",
					Sha256: "419a5850574de25219e005fd1c9823c50c435fb6899614b1b244c409a1705723",
				},
				{
					ID:     "fe0261b884ce22d3556304f0d262e74c3d3d44df8ca2f6dee568afaeef5e5f5c",
					Sha256: "2ba424c08a17c5fc643e899705586f4fa1e9d4b901d84af4458e2446ae1e5fad",
				},
				{
					ID:     "85acac73e2dc1297b9ac143b47b1219c20460956e14a8f300ff0975aab71d0a1",
					Sha256: "2a2275091837f8f76481fdcb84ca8524a854de3470758e5c4e1f0e57b32e2cba",
				},
				{
					ID:     "28d7ab04fb241b216fdc82ee7f6e1bb7de28971a882a8f1be8c69d4a71038e8b",
					Sha256: "5f704c526ba80e95c21c3a57dda2eae607db44f1d619c068488219da1ef08be0",
				},
				{
					ID:     "99ba2bcdcbc6dd643349146ca4b3123310cde2be02eba3cb321c1b2e77086459",
					Sha256: "84ba1612e4f058b35f1cce29d6b08dfbda1df72ee870db3466c6c19374fc2bfb",
				},
				{
					ID:     "7e00ea536c6d2461983661680baa658861cfe21b375e8e09753ec390d0f2fa12",
					Sha256: "846332a555385acda2e7d0afdca09da745a9ab0aa202eade90bc46c246b5acbc",
				},
			}

			attestationTaskMetadatas := []util.TaskMetadata{
				{
					ID:     "f1072e6780f6296868a9ac68c9c6372d761ba3c3ed9a14253a01c2c4e7fb1d92",
					Sha256: "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
				},
				{
					ID:     "71d9c19f677d4480b2916fea8929c5c51c6416f4d2ddf05da453d07c2dc72a09",
					Sha256: "c47b50e1f145d4983ea01bff4daa932c963f1ab5240645255585ceaa709013df",
				},
				{
					ID:     "18d374c62e76244887d030b867042374393446f92c8660e2e941f73cd6fc97c6",
					Sha256: "32704f4d4e33e495935c4b440886feda3911bc8c4d7a6b644c37f89d9ffe58e5",
				},
			}

			for _, taskMetadata := range arm64TaskMetadatas {
				Eventually(func() (string, error) {
					return util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				}, 5*time.Minute, 5*time.Second).Should(Equal(taskMetadata.Sha256))
			}

			for _, taskMetadata := range slices.Concat(amd64TaskMetadatas, attestationTaskMetadatas) {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeFalse())
			}

			_, err = schedulerClient.PreheatImage(context.Background(), &schedulerv2.PreheatImageRequest{
				Url:   "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Scope: managertypes.SingleSeedPeerScope,
			})
			Expect(err).NotTo(HaveOccurred())

			for _, taskMetadata := range slices.Concat(amd64TaskMetadatas, arm64TaskMetadatas, attestationTaskMetadatas) {
				Eventually(func() (string, error) {
					return util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				}, 5*time.Minute, 5*time.Second).Should(Equal(taskMetadata.Sha256))
			}

			resp, err := schedulerClient.StatImage(context.Background(), &schedulerv2.StatImageRequest{
				Url: "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Image.Layers).To(HaveLen(len(amd64TaskMetadatas) + len(arm64TaskMetadatas) + len(attestationTaskMetadatas)))

			var urls []string
			for _, peer := range resp.Peers {
				for _, layer := range peer.CachedLayers {
					if !slices.Contains(urls, layer.Url) {
						urls = append(urls, layer.Url)
					}
				}
			}
			Expect(urls).To(HaveLen(len(resp.Image.Layers)))

			resp, err = schedulerClient.StatImage(context.Background(), &schedulerv2.StatImageRequest{
				Url:      "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Platform: &amd64Platform,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Image.Layers).To(HaveLen(len(amd64TaskMetadatas)))
		})

		It("delete image for linux/amd64 platform should be ok", Label("delete", "image"), func() {
			conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			seedClientPods := make([]*util.PodExec, 3)
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			schedulerClient := schedulerv2.NewSchedulerClient(conn)
			_, err = schedulerClient.PreheatImage(context.Background(), &schedulerv2.PreheatImageRequest{
				Url:   "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Scope: managertypes.SingleSeedPeerScope,
			})
			Expect(err).NotTo(HaveOccurred())

			amd64TaskMetadatas := []util.TaskMetadata{
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

			arm64TaskMetadatas := []util.TaskMetadata{
				{
					ID:     "46b5a30713ed0772323dd02ebd5b50bd78c3d6c95c43ec161e222294c153eaa9",
					Sha256: "901f53de4911302661c0199f318d976e6f323d9a4a8927d37266c0e18db1df60",
				},
				{
					ID:     "51599e93d000381ccd553eb165824c6a430b56b980f172987bc6cb338a72593b",
					Sha256: "2d66c782e4847b8495b17a4b86be1159aa13d140fe4fd17cb1d536953d5114c6",
				},
				{
					ID:     "734826974c2ee8b7715c5ff585038cd7f03e13d874fa4c3c1b7a8d0870569f3b",
					Sha256: "5d33eb8f4cd8a1e776847144f5b26534d2646d9a08311cebf87448d4660b1443",
				},
				{
					ID:     "d515fa4dee4da892c373b4dc27e9c0317f603fea183d37063a326645c941fa02",
					Sha256: "419a5850574de25219e005fd1c9823c50c435fb6899614b1b244c409a1705723",
				},
				{
					ID:     "fe0261b884ce22d3556304f0d262e74c3d3d44df8ca2f6dee568afaeef5e5f5c",
					Sha256: "2ba424c08a17c5fc643e899705586f4fa1e9d4b901d84af4458e2446ae1e5fad",
				},
				{
					ID:     "85acac73e2dc1297b9ac143b47b1219c20460956e14a8f300ff0975aab71d0a1",
					Sha256: "2a2275091837f8f76481fdcb84ca8524a854de3470758e5c4e1f0e57b32e2cba",
				},
				{
					ID:     "28d7ab04fb241b216fdc82ee7f6e1bb7de28971a882a8f1be8c69d4a71038e8b",
					Sha256: "5f704c526ba80e95c21c3a57dda2eae607db44f1d619c068488219da1ef08be0",
				},
				{
					ID:     "99ba2bcdcbc6dd643349146ca4b3123310cde2be02eba3cb321c1b2e77086459",
					Sha256: "84ba1612e4f058b35f1cce29d6b08dfbda1df72ee870db3466c6c19374fc2bfb",
				},
				{
					ID:     "7e00ea536c6d2461983661680baa658861cfe21b375e8e09753ec390d0f2fa12",
					Sha256: "846332a555385acda2e7d0afdca09da745a9ab0aa202eade90bc46c246b5acbc",
				},
			}

			attestationTaskMetadatas := []util.TaskMetadata{
				{
					ID:     "f1072e6780f6296868a9ac68c9c6372d761ba3c3ed9a14253a01c2c4e7fb1d92",
					Sha256: "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
				},
				{
					ID:     "71d9c19f677d4480b2916fea8929c5c51c6416f4d2ddf05da453d07c2dc72a09",
					Sha256: "c47b50e1f145d4983ea01bff4daa932c963f1ab5240645255585ceaa709013df",
				},
				{
					ID:     "18d374c62e76244887d030b867042374393446f92c8660e2e941f73cd6fc97c6",
					Sha256: "32704f4d4e33e495935c4b440886feda3911bc8c4d7a6b644c37f89d9ffe58e5",
				},
			}

			for _, taskMetadata := range slices.Concat(amd64TaskMetadatas, arm64TaskMetadatas, attestationTaskMetadatas) {
				Eventually(func() (string, error) {
					return util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				}, 5*time.Minute, 5*time.Second).Should(Equal(taskMetadata.Sha256))
			}

			var errs []error
			for _, seedClientPod := range seedClientPods {
				out, err := seedClientPod.Command("dfctl", "task", "rm", "oci://ghcr.io/dragonflyoss/client:v1.5.12", "--platform", "linux/amd64").CombinedOutput()
				fmt.Println(string(out))
				errs = append(errs, err)
			}
			Expect(errs).To(ContainElement(BeNil()))

			for _, taskMetadata := range amd64TaskMetadatas {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeFalse())
			}

			for _, taskMetadata := range slices.Concat(arm64TaskMetadatas, attestationTaskMetadatas) {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeTrue())
			}

			amd64Platform := "linux/amd64"
			resp, err := schedulerClient.StatImage(context.Background(), &schedulerv2.StatImageRequest{
				Url:      "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Platform: &amd64Platform,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Peers).To(BeEmpty())

			errs = nil
			for _, seedClientPod := range seedClientPods {
				out, err := seedClientPod.Command("dfctl", "task", "rm", "oci://ghcr.io/dragonflyoss/client:v1.5.12").CombinedOutput()
				fmt.Println(string(out))
				errs = append(errs, err)
			}
			Expect(errs).To(ContainElement(BeNil()))

			for _, taskMetadata := range slices.Concat(arm64TaskMetadatas, attestationTaskMetadatas) {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeFalse())
			}

			resp, err = schedulerClient.StatImage(context.Background(), &schedulerv2.StatImageRequest{
				Url: "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Peers).To(BeEmpty())
		})

		It("delete image should keep the layers referenced by other images", Label("delete", "image"), func() {
			conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).NotTo(HaveOccurred())
			defer conn.Close()

			seedClientPods := make([]*util.PodExec, 3)
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			schedulerClient := schedulerv2.NewSchedulerClient(conn)
			_, err = schedulerClient.PreheatImage(context.Background(), &schedulerv2.PreheatImageRequest{
				Url:   "https://ghcr.io/v2/dragonflyoss/client/manifests/sha256:cf54d0f3bda39332435c109aed0689a3a0d31f0a00e570cf09aae10dfdf5da9b",
				Scope: managertypes.SingleSeedPeerScope,
			})
			Expect(err).NotTo(HaveOccurred())

			amd64TaskMetadatas := []util.TaskMetadata{
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

			attestationTaskMetadatas := []util.TaskMetadata{
				{
					ID:     "f1072e6780f6296868a9ac68c9c6372d761ba3c3ed9a14253a01c2c4e7fb1d92",
					Sha256: "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
				},
				{
					ID:     "71d9c19f677d4480b2916fea8929c5c51c6416f4d2ddf05da453d07c2dc72a09",
					Sha256: "c47b50e1f145d4983ea01bff4daa932c963f1ab5240645255585ceaa709013df",
				},
				{
					ID:     "18d374c62e76244887d030b867042374393446f92c8660e2e941f73cd6fc97c6",
					Sha256: "32704f4d4e33e495935c4b440886feda3911bc8c4d7a6b644c37f89d9ffe58e5",
				},
			}

			for _, taskMetadata := range amd64TaskMetadatas {
				Eventually(func() (string, error) {
					return util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				}, 5*time.Minute, 5*time.Second).Should(Equal(taskMetadata.Sha256))
			}

			_, err = schedulerClient.PreheatImage(context.Background(), &schedulerv2.PreheatImageRequest{
				Url:   "https://ghcr.io/v2/dragonflyoss/client/manifests/v1.5.12",
				Scope: managertypes.SingleSeedPeerScope,
			})
			Expect(err).NotTo(HaveOccurred())

			// Wait for the attestation layers, which are preheated after the linux/amd64 layers,
			// because PreheatImage returns before the seed peers record the image.
			for _, taskMetadata := range attestationTaskMetadatas {
				Eventually(func() (string, error) {
					return util.CalculateSha256ByTaskID(seedClientPods, taskMetadata.ID)
				}, 5*time.Minute, 5*time.Second).Should(Equal(taskMetadata.Sha256))
			}

			var errs []error
			for _, seedClientPod := range seedClientPods {
				out, err := seedClientPod.Command("dfctl", "task", "rm", "oci://ghcr.io/dragonflyoss/client:v1.5.12").CombinedOutput()
				fmt.Println(string(out))
				errs = append(errs, err)
			}
			Expect(errs).To(ContainElement(BeNil()))

			for _, taskMetadata := range amd64TaskMetadatas {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeTrue())
			}

			for _, taskMetadata := range attestationTaskMetadatas {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeFalse())
			}

			errs = nil
			for _, seedClientPod := range seedClientPods {
				out, err := seedClientPod.Command("dfctl", "task", "rm", "oci://ghcr.io/dragonflyoss/client@sha256:cf54d0f3bda39332435c109aed0689a3a0d31f0a00e570cf09aae10dfdf5da9b").CombinedOutput()
				fmt.Println(string(out))
				errs = append(errs, err)
			}
			Expect(errs).To(ContainElement(BeNil()))

			for _, taskMetadata := range amd64TaskMetadatas {
				Expect(util.CheckFilesExist(seedClientPods, taskMetadata.ID)).To(BeFalse())
			}
		})
	})
})
