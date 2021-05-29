package handler

import (
	"fmt"
	"github.com/etherlabsio/go-m3u8/m3u8"
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/server/cache"
	"net/http"
)

var playlistType = "VOD"
var allowCache = true

// Playlist creation for particular source
func Playlist(rw http.ResponseWriter, r *http.Request) {
	// extract parameters
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	src, ok := params["src"]
	if !ok || src == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	// get video data
	video, err := cache.Video("videos", id+".json")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create root.m3u8
	root := m3u8.NewPlaylist()
	root.Version = &m3u8Version
	root.Cache = &allowCache
	root.Type = &playlistType
	root.Sequence = 0

	// find lowest-res version
	/*sources := video.Sources()
	lowRes := "source.m3u8"
	lowResValue := mp4.ResSource
	for _, src := range sources {
		if mp4.Resolution(src) < lowResValue {
			lowRes = src
		}
	}*/

	// for n segments form src index m3u8
	for i := 0; i < len(video.Source(src+".m3u8")); i++ {
		// download low-res version
		/*rq, _ := http.NewRequest(http.MethodGet, video.SegmentURL(lowRes, i), nil)
		rq.Header.Set("User-Agent", httputil.RandUA())
		rsp, err := http.DefaultClient.Do(rq)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rsp.Body.Close()*/

		root.AppendItem(&m3u8.SegmentItem{
			Duration: 10, // TODO find
			Segment:  fmt.Sprintf("%s/%d.ts", src, i),
		})
	}

	_, err = rw.Write([]byte(root.String()))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
}
