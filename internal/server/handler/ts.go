package handler

import (
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/server/cache"
	"io"
	"net/http"
	"strconv"
)

func TS(rw http.ResponseWriter, r *http.Request) {
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
	nStr, ok := params["n"]
	if !ok || nStr == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	n, err := strconv.Atoi(nStr)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// request, cached if possible
	video, err := cache.LoadSource("videos", id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// request segment
	rq, _ := http.NewRequest(http.MethodGet, video.GetSegmentURL(source, n), nil)
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rsp.Body.Close()

	// respond with segment
	rw.Header().Set("Content-Type", "video/MP2T")
	_, err = io.Copy(rw, rsp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
