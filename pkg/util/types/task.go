package types

import (
	"crypto/sha1"
	"fmt"
	logger "d7y.io/dragonfly/v2/pkg/dflog"
	"d7y.io/dragonfly/v2/pkg/rpc/base"
	"net/url"
	"strings"
)

func GenerateTaskId(rawUrl string, filter string, meta *base.UrlMeta) (taskId string) {
	taskUrl, err := url.Parse(rawUrl)
	if err != nil {
		logger.Warnf("GenerateTaskId rawUrl[%s] is invalid url", rawUrl)
		taskId = fmt.Sprintf("%x", sha1.Sum([]byte(rawUrl)))
		return
	}
	if filter != "" {
		queries := taskUrl.Query()
		fields := strings.Split(filter, "&")
		if len(fields) > 0 {
			for _, key := range fields {
				queries.Del(key)
			}
		}
		taskUrl.RawQuery = queries.Encode()
	}
	taskRawString := taskUrl.String()
	if meta != nil {
		taskRawString += "_" + meta.Range + "_" + meta.Md5
	}
	taskId = fmt.Sprintf("%x", sha1.Sum([]byte(taskRawString)))
	return
}
