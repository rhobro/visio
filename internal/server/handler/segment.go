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
	source, ok := params["source"]
	if !ok || source == "" {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
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

	seg, err := cache.Segment("videos", id, source, n)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// respond with segment
	rw.Header().Set("Content-Type", "video/MP2T")
	_, err = io.Copy(rw, seg)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
