package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"d7y.io/dragonfly/v2/test/e2e/manager"

	"d7y.io/dragonfly/v2/internal/idgen"
	"d7y.io/dragonfly/v2/pkg/rpc/base"

	machineryv1tasks "github.com/RichardKnop/machinery/v1/tasks"

	"d7y.io/dragonfly/v2/manager/types"

	"d7y.io/dragonfly/v2/test/e2e/e2eutil"
	. "github.com/onsi/ginkgo" //nolint
	. "github.com/onsi/gomega" //nolint
)

var _ = Describe("Preheat with manager", func() {
	Context("preheat", func() {
		It("preheat should be ok", func() {
			var cdnPods [3]*e2eutil.PodExec
			for i := 0; i < 3; i++ {
				out, err := e2eutil.KubeCtlCommand("-n", dragonflyNamespace, "get", "pod", "-l", "component=cdn",
					"-o", fmt.Sprintf("jsonpath='{range .items[%d]}{.metadata.name}{end}'", i)).CombinedOutput()
				podName := strings.Trim(string(out), "'")
				Expect(err).NotTo(HaveOccurred())
				fmt.Println(podName)
				Expect(strings.HasPrefix(podName, "dragonfly-cdn-")).Should(BeTrue())
				cdnPods[i] = e2eutil.NewPodExec(dragonflyNamespace, podName, "cdn")
			}

			for _, v := range e2eutil.GetFileList() {
				url := e2eutil.GetFileURL(v)
				fmt.Println("download url " + url)

				// get original file digest
				out, err := e2eutil.DockerCommand("sha256sum", v).CombinedOutput()
				fmt.Println(string(out))
				Expect(err).NotTo(HaveOccurred())
				sha256sum1 := strings.Split(string(out), " ")[0]

				// preheat file
				out, err = cdnPods[0].Command("curl", "-s", "-X", "POST", "-H", "Content-Type:application/json",
					"-d", fmt.Sprintf(`{"type":"file","url":"%s"}`, url), manager.PreheatURL).CombinedOutput()
				fmt.Println(string(out))
				Expect(err).NotTo(HaveOccurred())

				// generate task_id, also the filename
				// TODO(zzy987) this tag may not always be right.
				taskID := idgen.TaskID(url, &base.UrlMeta{Tag: "d7y/manager"})
				fmt.Println(taskID)

				// wait for success
				preheatJob := &types.Preheat{}
				err = json.Unmarshal(out, preheatJob)
				Expect(err).NotTo(HaveOccurred())
				waitForDone(preheatJob, cdnPods[0])

				// get downloaded file digest
				var sha256sum2 string
				for _, cdn := range cdnPods {
					out, err = cdn.Command("ls", manager.CDNCachePath).CombinedOutput()
					Expect(err).NotTo(HaveOccurred())
					dir := taskID[0:3]
					if !strings.Contains(string(out), dir) {
						continue
					}

					out, err = cdn.Command("ls", fmt.Sprintf("%s/%s", manager.CDNCachePath, dir)).CombinedOutput()
					Expect(err).NotTo(HaveOccurred())
					file := taskID
					if !strings.Contains(string(out), file) {
						continue
					}

					out, err = cdn.Command("sha256sum", fmt.Sprintf("%s/%s/%s", manager.CDNCachePath, dir, file)).CombinedOutput()
					fmt.Println(string(out))
					Expect(err).NotTo(HaveOccurred())
					sha256sum2 = strings.Split(string(out), " ")[0]
					break
				}
				if sha256sum2 == "" {
					fmt.Println("preheat file not found")
				}
				Expect(sha256sum1).To(Equal(sha256sum2))
			}
		})

		It("concurrency 100 preheat should be ok", func() {
			url := e2eutil.GetFileURL(hostnameFilePath)
			fmt.Println("download url " + url)
			dataFilePath := "post_data"
			fd, err := os.Create(dataFilePath)
			Expect(err).NotTo(HaveOccurred())
			_, err = fd.WriteString(fmt.Sprintf(`{"type":"file","url":"%s"}`, url))
			fd.Close()
			Expect(err).NotTo(HaveOccurred())

			out, err := e2eutil.ABCommand("-c", "100", "-n", "200", "-T", "application/json", "-p", dataFilePath, fmt.Sprintf("127.0.0.1:%s/%s", manager.Port, manager.PreheatPath), url).CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
			os.Remove(dataFilePath)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

func waitForDone(preheat *types.Preheat, pod *e2eutil.PodExec) bool {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
			out, err := pod.Command("curl", "-s", fmt.Sprintf("%s/%s", manager.PreheatURL, preheat.ID)).CombinedOutput()
			fmt.Println(string(out))
			Expect(err).NotTo(HaveOccurred())
			err = json.Unmarshal(out, preheat)
			Expect(err).NotTo(HaveOccurred())
			switch preheat.Status {
			case machineryv1tasks.StateSuccess:
				return true
			case machineryv1tasks.StateFailure:
				return false
			default:
			}
		}
	}
}
