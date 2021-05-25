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

const lifetime = 5 * time.Hour

func init() {
	go manager()
}

type item struct {
	t   time.Time
	src *visio.Source
}

var mtx sync.RWMutex
var items = make(map[string]*item)

func manager() {
	for {
		for k, v := range items {
			if time.Now().Sub(v.t) >= lifetime {
				// expired
				mtx.Lock()
				delete(items, k)
				mtx.Unlock()
			}
		}
		time.Sleep(time.Minute)
	}
}

func LoadSource(bucket, id string) (*visio.Source, error) {
	key := path.Join(bucket, id)

	// cached?
	mtx.RLock()
	for k, v := range items {
		if k == key {
			return v.src, nil
		}
	}
	mtx.RUnlock()

	// not found - must download
	video, err := storaje.Download("videos", id)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(video)
	if err != nil {
		return nil, err
	}

	// unmarshal into visio.Source
	src := new(visio.Source)
	err = json.Unmarshal(body, src)
	if err != nil {
		return nil, err
	}

	// cache
	mtx.Lock()
	items[key] = &item{
		t:   time.Now(),
		src: src,
	}
	mtx.Unlock()

	return src, nil
}
