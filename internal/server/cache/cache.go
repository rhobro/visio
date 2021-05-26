package cache

import (
	"github.com/rhobro/goutils/pkg/fileio"
	"log"
	"os"
	"time"
)

func Start() {
	go manager()
}

func manager() {
	for {
		// video data
		for k, v := range videos {
			if time.Now().Sub(v.t) >= lifetime {
				// expired
				vMu.Lock()
				delete(videos, k)
				vMu.Unlock()
			}
		}

		// segment cache
		// check total cache size
		size, err := fileio.DirSize(fileio.TmpDir)
		if err != nil {
			log.Fatalf("cannot compute size of %s: %s", fileio.TmpDir, err)
		}
		if size > maxCacheStorage {
			// calculate amount required to shed
			delta := size - maxCacheStorage

			// prune until delta < 0
			for delta > 0 {
				// find oldest file
				oldestKey := ""

				for k, v := range segments {
					if oldestKey == "" {
						oldestKey = k
						continue
					}

					// if older
					if v.t.Before(segments[oldestKey].t) {
						oldestKey = k
					}
				}

				// size file
				stat, err := os.Stat(segments[oldestKey].file)
				if err != nil {
					log.Fatalf("can't measure stats of %s: %s", segments[oldestKey].file, err)
				}
				delta -= stat.Size()

				// delete cache entry and file
				sMu.Lock()
				err = os.Remove(segments[oldestKey].file)
				if err != nil {
					log.Fatalf("can't remove %s: %s", segments[oldestKey].file, err)
				}
				delete(segments, oldestKey)
				sMu.Unlock()
			}
		}
	}
}
