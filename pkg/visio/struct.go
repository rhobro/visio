package visio

import (
	"path"
	"strings"
)

type Video map[string][]string

func (v *Video) Sources() []string {
	sources := make([]string, len(*v), len(*v))

	var i int
	for k := range *v {
		sources[i] = k
		i++
	}
	return sources
}

func (v *Video) Source(src string) []string {
	return (*v)[src]
}

func (v *Video) SegmentURL(src string, idx int) string {
	if idx >= len((*v)[src]) {
		return ""
	}
	return (*v)[src][idx]
}

func (v *Video) IsValid() bool {
	for res, urls := range *v {
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
