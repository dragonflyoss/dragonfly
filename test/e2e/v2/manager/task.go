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
	"d7y.io/dragonfly/v2/test/e2e/v2/util"
)

var _ = Describe("GetTask and DeleteTask with Manager", func() {
	Context("1MiB file", Label("getTask", "deleteTask", "file"), func() {
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

		It("getTask and deleteTask should be ok", func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type: "file",
					URL:  testFile.GetDownloadURL(),
				},
				SchedulerClusterIDs: []uint{1},
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
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			req, err = structure.StructToMap(types.CreateGetTaskJobRequest{
				Type: internaljob.GetTaskJob,
				Args: types.GetTaskArgs{
					TaskID: testFile.GetTaskID(),
				},
			})
			Expect(err).NotTo(HaveOccurred())
			out, err = managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://127.0.0.1:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job = &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done = waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())
			Expect(validateTaskResponse(job)).To(BeTrue())

			req, err = structure.StructToMap(types.CreateDeleteTaskJobRequest{
				Type: internaljob.DeleteTaskJob,
				Args: types.DeleteTaskArgs{
					TaskID: testFile.GetTaskID(),
				},
			})
			Expect(err).NotTo(HaveOccurred())
			out, err = managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://127.0.0.1:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job = &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done = waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())
			Expect(validateTaskResponse(job)).To(BeTrue())

			exist := util.CheckFilesExist(seedClientPods, testFile.GetTaskID())
			Expect(exist).Should(BeFalse())
		})
	})

	Context("10MiB file", Label("getTask", "deleteTask", "file"), func() {
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

		It("getTask and deleteTask should be ok", func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			req, err := structure.StructToMap(types.CreatePreheatJobRequest{
				Type: internaljob.PreheatJob,
				Args: types.PreheatArgs{
					Type: "file",
					URL:  testFile.GetDownloadURL(),
				},
				SchedulerClusterIDs: []uint{1},
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
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			sha256sum, err := util.CalculateSha256ByTaskID(seedClientPods, testFile.GetTaskID())
			Expect(err).NotTo(HaveOccurred())
			Expect(testFile.GetSha256()).To(Equal(sha256sum))

			req, err = structure.StructToMap(types.CreateGetTaskJobRequest{
				Type: internaljob.GetTaskJob,
				Args: types.GetTaskArgs{
					TaskID: testFile.GetTaskID(),
				},
			})
			Expect(err).NotTo(HaveOccurred())
			out, err = managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://127.0.0.1:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job = &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done = waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())
			Expect(validateTaskResponse(job)).To(BeTrue())

			req, err = structure.StructToMap(types.CreateDeleteTaskJobRequest{
				Type: internaljob.DeleteTaskJob,
				Args: types.DeleteTaskArgs{
					TaskID: testFile.GetTaskID(),
				},
			})
			Expect(err).NotTo(HaveOccurred())
			out, err = managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://127.0.0.1:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job = &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done = waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())
			Expect(validateTaskResponse(job)).To(BeTrue())

			exist := util.CheckFilesExist(seedClientPods, testFile.GetTaskID())
			Expect(exist).Should(BeFalse())
		})
	})

	Context("100MiB file", Label("getTask", "deleteTask", "file"), func() {
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

		It("getTask and deleteTask should be failed", func() {
			managerPod, err := util.ManagerExec(0)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			seedClientPods := make([]*util.PodExec, 3)
			for i := range 3 {
				seedClientPods[i], err = util.SeedClientExec(i)
				fmt.Println(err)
				Expect(err).NotTo(HaveOccurred())
			}

			exist := util.CheckFilesExist(seedClientPods, testFile.GetTaskID())
			Expect(exist).Should(BeFalse())

			req, err := structure.StructToMap(types.CreateGetTaskJobRequest{
				Type: internaljob.GetTaskJob,
				Args: types.GetTaskArgs{
					TaskID: testFile.GetTaskID(),
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
			Expect(validateTaskResponse(job)).To(BeFalse())

			req, err = structure.StructToMap(types.CreateDeleteTaskJobRequest{
				Type: internaljob.DeleteTaskJob,
				Args: types.DeleteTaskArgs{
					TaskID: testFile.GetTaskID(),
				},
			})
			Expect(err).NotTo(HaveOccurred())
			out, err = managerPod.CurlCommand("POST", map[string]string{"Content-Type": "application/json"}, req,
				"http://127.0.0.1:8080/api/v1/jobs").CombinedOutput()
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())
			fmt.Println(string(out))

			job = &models.Job{}
			err = json.Unmarshal(out, job)
			fmt.Println(err)
			Expect(err).NotTo(HaveOccurred())

			done = waitForDone(job, managerPod)
			Expect(done).Should(BeTrue())

			Expect(validateTaskResponse(job)).To(BeFalse())
		})
	})
})

func validateTaskResponse(job *models.Job) bool {
	Expect(job.Result).NotTo(BeNil())

	groupJobStateData, err := json.Marshal(job.Result)
	Expect(err).NotTo(HaveOccurred())

	groupJobState := internaljob.GroupJobState{}
	err = json.Unmarshal(groupJobStateData, &groupJobState)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(groupJobState.JobStates)).Should(BeNumerically("==", 3))

	for _, state := range groupJobState.JobStates {
		for _, result := range state.Results {
			resultData, err := json.Marshal(result)
			Expect(err).NotTo(HaveOccurred())

			switch job.Type {
			case internaljob.GetTaskJob:
				getTaskResponse := internaljob.GetTaskResponse{}
				err = json.Unmarshal(resultData, &getTaskResponse)
				Expect(err).NotTo(HaveOccurred())
				if len(getTaskResponse.Peers) > 0 {
					return true
				}
			case internaljob.DeleteTaskJob:
				deleteTaskResponse := internaljob.DeleteTaskResponse{}
				err = json.Unmarshal(resultData, &deleteTaskResponse)
				Expect(err).NotTo(HaveOccurred())
				if len(deleteTaskResponse.SuccessTasks) > 0 || len(deleteTaskResponse.FailureTasks) > 0 {
					return true
				}
			}
		}
	}

	return false
}
