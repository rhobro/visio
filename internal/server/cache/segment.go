// currently disabling caching

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

const maxCacheStorage = 1 << 30

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
	video, err := Video("videos", id+".json")
	if err != nil {
		return nil, err
	}
	// request segment
	rq, _ := http.NewRequest(http.MethodGet, video.SegmentURL(source+".m3u8", i), nil)
	rq.Header.Set("User-Agent", httputil.RandUA())
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}

	// add to cache process
	/*cacheStream <- cacheItem{
		key:    key,
		bucket: bucket,
		id:     id,
		source: source,
		i:      i,
	}*/

	return rsp.Body, err
}

var cacheStream = make(chan cacheItem, 10)

const cacheStreamSize = 10

type cacheItem struct {
	key    string
	bucket string
	id     string
	source string
	i      int
}

func init() {
	//go cacher()
}

func cacher() {
	for item := range cacheStream {
		// get video data
		video, err := Video("videos", item.id+".json")
		if err != nil {
			continue
		}
		// request segment
		rq, _ := http.NewRequest(http.MethodGet, video.SegmentURL(item.source+".m3u8", item.i), nil)
		rq.Header.Set("User-Agent", httputil.RandUA())
		rsp, err := http.DefaultClient.Do(rq)
		if err != nil {
			continue
		}

		// cache
		tmpF, err := fileio.TmpPath(item.key)
		if err != nil {
			continue
		}
		f, err := os.Create(tmpF)
		if err != nil {
			continue
		}
		// copy in
		_, err = io.Copy(f, rsp.Body)
		if err != nil {
			continue
		}
		f.Close()

		// cache path in segments
		sMu.Lock()
		segments[item.key] = segment{
			t:    time.Now(),
			file: tmpF,
		}
		sMu.Unlock()
	}
}
