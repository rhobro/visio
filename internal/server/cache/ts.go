package cache

import (
	"github.com/rhobro/goutils/pkg/fileio"
	"github.com/rhobro/goutils/pkg/httputil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const maxCacheStorage = 1 << 20

type segment struct {
	t    time.Time
	file string
}

var sMu sync.RWMutex
var segments = make(map[string]segment)

func Segment(bucket, id, source string, i int) (io.Reader, error) {
	key := filepath.Join(bucket, id, source, strconv.Itoa(i)+".ts")

	// cached?
	sMu.RLock()
	for k, v := range segments {
		if k == key {
			// open file
			f, err := os.Open(v.file)
			if err != nil {
				return nil, err
			}
			v.t = time.Now() // keep cached
			return f, nil
		}
	}
	sMu.RUnlock()

	// get video data
	video, err := Video("videos", id)
	if err != nil {
		return nil, err
	}
	// request segment
	rq, _ := http.NewRequest(http.MethodGet, video.GetSegmentURL(source, i), nil)
	rq.Header.Set("User-Agent", httputil.RandUA())
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	// cache
	tmpF, err := fileio.TmpPath(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(tmpF)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// copy in
	_, err = io.Copy(f, rsp.Body)
	if err != nil {
		return nil, err
	}

	// cache path in segments
	sMu.Lock()
	segments[key] = segment{
		t:    time.Now(),
		file: tmpF,
	}
	sMu.Unlock()

	return Segment(bucket, id, source, i) // recurse to return after cached
}
