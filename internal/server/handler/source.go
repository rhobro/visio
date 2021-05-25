package handler

import (
	"github.com/etherlabsio/go-m3u8/m3u8"
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/server/cache"
	"net/http"
)

var playlistType = "VOD"
var allowCache = true

func SourceM3U8(rw http.ResponseWriter, r *http.Request) {
	// extract parameters
	params := mux.Vars(r)
	id, ok := params["id"]
	if !ok || id == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	id += ".json"
	source, ok := params["source"]
	if !ok || source == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	source += ".m3u8"

	// get video data
	video, err := cache.Video("videos", id)
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

	// for n segments form source index m3u8
	for i := 0; i < len(video.GetSource(source)); i++ {

	}
}
