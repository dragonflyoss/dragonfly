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

var _ = Describe("Download With RateLimit", func() {
	Context("/bin/kubectl file", func() {
		It("should be ok", Label("download", "rateLimit"), func() {
			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			start := time.Now()
			out, err := clientPod.Command("sh", "-c", fmt.Sprintf("curl -x 127.0.0.1:4001 -H 'X-Dragonfly-Tag: download-rate-limit' %s --output %s", util.GetFileURL("/bin/kubectl"), util.GetOutputPath("kubectl-proxy"))).CombinedOutput()
			elapsed := time.Since(start)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			// The download rate limit is 1MiB/s, and the file size is about 44MiB, so the download time should be greater than 44s.
			Expect(elapsed).Should(BeNumerically(">", 44*time.Second))
			fmt.Println(string(out))

			fileMetadata := util.FileMetadata{
				ID:     "6b94aa169d11ba67c1852f14c4df39e5ebf84d1c0d7a6f6ad0ba2f3547883f4f",
				Sha256: "327b4022d0bfd1d5e9c0701d4a3f989a536f7e6e865e102dcd77c7e7adb31f9a",
			}

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, fileMetadata.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileMetadata.Sha256).To(Equal(sha256sum))

			sha256sum, err = util.CalculateSha256ByOutput([]*util.PodExec{clientPod}, util.GetOutputPath("kubectl-proxy"))
			Expect(err).NotTo(HaveOccurred())
			Expect(fileMetadata.Sha256).To(Equal(sha256sum))
		})
	})
})
