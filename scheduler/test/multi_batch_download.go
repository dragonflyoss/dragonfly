package test

import (
	"d7y.io/dragonfly/v2/scheduler/test/common"
	"d7y.io/dragonfly/v2/scheduler/test/mock_client"
	. "github.com/onsi/ginkgo"
	"reflect"
	"time"
)

var _ = FDescribe("Multi Client Multi Batch Download Test", func() {
	tl := common.NewE2ELogger()

	var (
		clientNum  = 20
		stopChList []chan struct{}
	)

	Describe("Create Multi Client", func() {
		It("create first batch client should be successfully", func() {
			for i := 0; i < clientNum; i++ {
				client := mock_client.NewMockClient("127.0.0.1:8002", "http://dragonfly.com?type=multi_batch", "mb", tl)
				go client.Start()
				stopCh := client.GetStopChan()
				stopChList = append(stopChList, stopCh)
			}
		})

		It("create second batch client should be successfully", func() {
			time.Sleep(time.Second * 5)
			for i := 0; i < clientNum; i++ {
				client := mock_client.NewMockClient("127.0.0.1:8002", "http://dragonfly.com?type=multi_batch", "mb", tl)
				go client.Start()
				stopCh := client.GetStopChan()
				stopChList = append(stopChList, stopCh)
			}
		})

		It("create third batch client should be successfully", func() {
			time.Sleep(time.Second * 5)
			for i := 0; i < clientNum; i++ {
				client := mock_client.NewMockClient("127.0.0.1:8002", "http://dragonfly.com?type=multi_batch", "mb", tl)
				go client.Start()
				stopCh := client.GetStopChan()
				stopChList = append(stopChList, stopCh)
			}
		})
	})

	Describe("Wait Clients Finish", func() {
		It("all clients should be stopped successfully", func() {
			time.Sleep(time.Second*5)
			timer := time.After(time.Minute * 10)
			caseList := []reflect.SelectCase{
				{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(timer), Send: reflect.Value{}},
			}
			for _, stopCh := range stopChList {
				caseList = append(caseList, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(stopCh), Send: reflect.Value{}})
			}
			closedNumber := 0
			for {
				selIndex, _, _ := reflect.Select(caseList)
				caseList = append(caseList[:selIndex], caseList[selIndex+1:]...)
				if selIndex == 0 {
					tl.Fatalf("download file failed")
				} else {
					closedNumber++
					if closedNumber >= clientNum {
						break
					}
				}
			}
			tl.Log("multiple batch all client download file finished")
		})
	})
})
