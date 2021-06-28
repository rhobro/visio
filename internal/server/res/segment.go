// currently disabling caching

package res

import (
	"github.com/rhobro/goutils/pkg/httputil"
	"io"
	"net/http"
)

func Segment(id, source string, i int) (io.ReadCloser, error) {
	// get video data
	video, err := Video("videos", id)
	if err != nil {
		return nil, err
	}

	// request segment
	rq, _ := http.NewRequest(http.MethodGet, video.SegmentURL(source, i), nil)
	rq.Header.Set("User-Agent", httputil.RandUA())
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}

	return rsp.Body, err
}
