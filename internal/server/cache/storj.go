package cache

import (
	"encoding/json"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/pkg/visio"
	"io"
	"path"
	"sync"
	"time"
)

const lifetime = 1 * time.Hour

type video struct {
	t   time.Time
	src *visio.Video
}

var vMu sync.RWMutex
var videos = make(map[string]*video)

func Video(bucket, id string) (*visio.Video, error) {
	key := path.Join(bucket, id)

	// cached?
	vMu.RLock()
	for k, v := range videos {
		if k == key {
			v.t = time.Now() // keep cached
			return v.src, nil
		}
	}
	vMu.RUnlock()

	// not found - must download
	vid, err := storaje.Download("videos", id)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(vid)
	if err != nil {
		return nil, err
	}

	// unmarshal into visio.Video
	src := new(visio.Video)
	err = json.Unmarshal(body, src)
	if err != nil {
		return nil, err
	}

	// cache
	vMu.Lock()
	videos[key] = &video{
		t:   time.Now(),
		src: src,
	}
	vMu.Unlock()

	return src, nil
}
