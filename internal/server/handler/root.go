package handler

import (
	"github.com/etherlabsio/go-m3u8/m3u8"
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/server/cache"
	"net/http"
	"strconv"
	"strings"
)

const heightToWidth = 16. / 9.

var m3u8Version = 3

/*
#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:1
#EXT-X-MEDIA-SEQUENCE:0
#EXTINF:1.149,
1080p/0.ts
#EXT-X-ENDLIST
*/

func RootM3U8(rw http.ResponseWriter, r *http.Request) {
	// extract parameters
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	id += ".json"

	// request, cached if possible
	video, err := cache.Video("videos", id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create root.m3u8
	root := m3u8.NewPlaylist()
	root.Version = &m3u8Version

	// add sources for each resolution
	programID := "0"
	for res := range *video {
		height, _ := strconv.Atoi(res[:strings.Index(res, ".")-1])
		bandwidth := 4000000
		switch height {
		case 1080:
			bandwidth = 6000000
		case 720:
			bandwidth = 2000000
		case 360:
			bandwidth = 500000
		}

		root.AppendItem(&m3u8.PlaylistItem{
			ProgramID: &programID,
			Resolution: &m3u8.Resolution{
				Width:  int(float64(height) * heightToWidth),
				Height: height,
			},
			Bandwidth: bandwidth,

			URI: res,
		})
	}

	// return m3u8 file
	_, err = rw.Write([]byte(root.String()))
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
