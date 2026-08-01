package http

import (
	"strconv"
	"strings"
)

func FuzzParseRange(data []byte) int {
	s := string(data)
	parts := strings.SplitN(s, "\n", 2)
	if len(parts) < 2 {
		return 0
	}
	rangeStr := parts[0]
	sizeStr := strings.TrimSpace(parts[1])
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0
	}
	_, _ = ParseRange(rangeStr, size)
	_, _ = ParseURLMetaRange(rangeStr, size)
	return 1
}
