package digest

import (
	"strings"
)

func FuzzParse(data []byte) int {
	s := string(data)
	s = strings.TrimSpace(s)
	_, err := Parse(s)
	if err != nil {
		return 0
	}
	return 1
}
