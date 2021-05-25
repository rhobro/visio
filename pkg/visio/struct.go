package visio

import (
	"path"
	"strings"
)

type Video map[string][]string

func (s *Video) GetSource(src string) []string {
	return (*s)[src]
}

func (s *Video) GetSegmentURL(src string, idx int) string {
	if idx >= len((*s)[src]) {
		return ""
	}
	return (*s)[src][idx]
}

func (s *Video) IsValid() bool {
	for res, urls := range *s {
		// check keys are names of .m3u8 files
		if !(strings.Contains(res, ".m3u8") && len(res) > 5) {
			return false
		}

		for _, url := range urls {
			if path.Ext(url) != ".ts" {
				return false
			}
		}
	}

	return true
}
