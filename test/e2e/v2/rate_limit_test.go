/*
 *     Copyright 2025 The Dragonfly Authors
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
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint

	"d7y.io/dragonfly/v2/test/e2e/v2/util"
)

/*
The rate limit configuration is located in the `test/testdata/charts/config-v2-rate-limit.yaml`.
By default, the rate limit is set to `1MiB`. However, due to the nature of the download process,
where pieces are downloaded concurrently, the rate-limiting mechanism cannot achieve
100% precision. As a result, the actual download speed will may exceed the theoretical
rate limit.

To ensure the stability of the following end-to-end (e2e) tests, the theoretical download
time is adjusted by multiplying it with a factor of 0.5. This adjustment accounts for
the discrepancy between the theoretical and actual speeds.

For example:
- If the rate limit is set to 1 MiB/s and the file size is 100 MiB, the theoretical
  download time is `100 MiB / 1 MiB/s = 100 seconds`.

  Here are the expected download times(JettiFactor: 0.5):

  Minimum expected download time: 100 seconds * (1-0.5) = 50 seconds
  Maximum expected download time: 100 seconds * (1+0.5) = 150 seconds


This adjustment ensures that the e2e tests remain consistent and reliable while accounting
for the inherent imprecision of the rate-limiting mechanism.
*/

var _ = Describe("Download Using Dfget With Rate Limit", func() {
	Context("50MiB file", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize10MiB * 5)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should download successfully in over 25 seconds", func() {
			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			startAt := time.Now()
			out, err := clientPod.Command("sh", "-c", fmt.Sprintf("dfget %s --disable-back-to-source --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
			elapsed := time.Since(startAt)
			fmt.Println(string(out), err)
			Expect(err).NotTo(HaveOccurred())
			Expect(elapsed).Should(BeNumerically(">", 25*time.Second))

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			sha256sum, err = util.CalculateSha256ByOutput([]*util.PodExec{clientPod}, testFile.GetOutputPath())
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

		It("should download successfully in over 50 seconds", func() {
			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			startAt := time.Now()
			out, err := clientPod.Command("sh", "-c", fmt.Sprintf("dfget %s --disable-back-to-source --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
			elapsed := time.Since(startAt)
			fmt.Println(string(out), err)
			Expect(err).NotTo(HaveOccurred())
			Expect(elapsed).Should(BeNumerically(">", 50*time.Second))

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			sha256sum, err = util.CalculateSha256ByOutput([]*util.PodExec{clientPod}, testFile.GetOutputPath())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
		})
	})
})

var _ = Describe("Download Using Proxy With Rate Limit", func() {
	Context("50MiB file", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize10MiB * 5)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should download successfully in over 25 seconds", Label("proxy", "download"), func() {
			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			startAt := time.Now()
			out, err := clientPod.Command("sh", "-c", fmt.Sprintf("curl -x 127.0.0.1:4001 -H 'X-Dragonfly-Tag: proxy' %s --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
			elapsed := time.Since(startAt)
			fmt.Println(string(out), err)
			Expect(err).NotTo(HaveOccurred())
			Expect(elapsed).Should(BeNumerically(">", 25*time.Second))

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, testFile.GetTaskID(util.WithTaskIDTag("proxy")))
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			sha256sum, err = util.CalculateSha256ByOutput([]*util.PodExec{clientPod}, testFile.GetOutputPath())
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

		It("should download successfully in over 50 seconds", Label("proxy", "download"), func() {
			clientPod, err := util.ClientExec()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			startAt := time.Now()
			out, err := clientPod.Command("sh", "-c", fmt.Sprintf("curl -x 127.0.0.1:4001 -H 'X-Dragonfly-Tag: proxy' %s --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
			elapsed := time.Since(startAt)
			fmt.Println(string(out), err)
			Expect(err).NotTo(HaveOccurred())
			Expect(elapsed).Should(BeNumerically(">", 50*time.Second))

			sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{clientPod}, testFile.GetTaskID(util.WithTaskIDTag("proxy")))
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			sha256sum, err = util.CalculateSha256ByOutput([]*util.PodExec{clientPod}, testFile.GetOutputPath())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))
		})
	})
})

var _ = Describe("Download Using Prefetch Proxy With Rate Limit", func() {
	Context("50MiB file and set range header bytes=100-200", func() {
		var (
			testFile *util.File
			err      error
		)

		BeforeEach(func() {
			testFile, err = util.GetFileServer().GenerateFile(util.FileSize10MiB * 5)
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile).NotTo(BeNil())
		})

		AfterEach(func() {
			err = util.GetFileServer().DeleteFile(testFile.GetInfo())
			Expect(err).NotTo(HaveOccurred())
		})

		It("download should be ok", Label("prefetch-proxy", "download", "range: bytes=100-200"), func() {
			seedClientPod, err := util.SeedClientExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			out, err := seedClientPod.Command("sh", "-c", fmt.Sprintf("curl -x 127.0.0.1:4001 -r 100-200 -H 'X-Dragonfly-Tag: prefetch-proxy-bytes-100-200' %s --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			sha256sum, err := util.CalculateSha256ByOutput([]*util.PodExec{seedClientPod}, testFile.GetOutputPath())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetRangeSha256("100-200", testFile.GetInfo().Size())).To(Equal(sha256sum))

			// Prefetch should not be completed within 25 seconds.
			Consistently(func() error {
				sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{seedClientPod}, testFile.GetTaskID(util.WithTaskIDTag("prefetch-proxy-bytes-100-200")))
				if err != nil {
					return err
				}
				// Prefetch should not be completed, so the sha256sum should not be equal.
				if testFile.GetSha256() == sha256sum {
					return errors.New("prefetch should not be completed, but it seems done as the sha256sum is equal")
				}

				return nil
			}, 25*time.Second, 5*time.Second).ShouldNot(HaveOccurred())

			// Prefetch should be eventually completed within 75 seconds(so wait for another 50 seconds).
			Eventually(func() string {
				sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{seedClientPod}, testFile.GetTaskID(util.WithTaskIDTag("prefetch-proxy-bytes-100-200")))
				if err != nil {
					return err.Error()
				}
				// Eventually, the sha256sum should be equal.
				return sha256sum
			}, 50*time.Second, 5*time.Second).Should(Equal(testFile.GetSha256()))
		})
	})

	Context("100MiB file and set range header bytes=100-200", func() {
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

		It("download should be ok", Label("prefetch-proxy", "download", "range: bytes=100-200"), func() {
			seedClientPod, err := util.SeedClientExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			out, err := seedClientPod.Command("sh", "-c", fmt.Sprintf("curl -x 127.0.0.1:4001 -r 100-200 -H 'X-Dragonfly-Tag: prefetch-proxy-bytes-100-200' %s --output %s", testFile.GetDownloadURL(), testFile.GetOutputPath())).CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			sha256sum, err := util.CalculateSha256ByOutput([]*util.PodExec{seedClientPod}, testFile.GetOutputPath())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetRangeSha256("100-200", testFile.GetInfo().Size())).To(Equal(sha256sum))

			// Prefetch should not be completed within 80 seconds.
			Consistently(func() error {
				sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{seedClientPod}, testFile.GetTaskID(util.WithTaskIDTag("prefetch-proxy-bytes-100-200")))
				if err != nil {
					return err
				}
				// Prefetch should not be completed, so the sha256sum should not be equal.
				if testFile.GetSha256() == sha256sum {
					return errors.New("prefetch should not be completed, but it seems done as the sha256sum is equal")
				}

				return nil
			}, 50*time.Second, 10*time.Second).ShouldNot(HaveOccurred())

			// Prefetch should be eventually completed within 150 seconds(so wait for another 100 seconds).
			Eventually(func() string {
				sha256sum, err := util.CalculateSha256ByTaskID([]*util.PodExec{seedClientPod}, testFile.GetTaskID(util.WithTaskIDTag("prefetch-proxy-bytes-100-200")))
				if err != nil {
					return err.Error()
				}
				// Eventually, the sha256sum should be equal.
				return sha256sum
			}, 100*time.Second, 10*time.Second).Should(Equal(testFile.GetSha256()))
		})
	})
})
