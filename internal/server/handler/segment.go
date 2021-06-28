package handler

import (
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/server/res"
	"io"
	"net/http"
	"strconv"
)

// Segment requests and returns the segment required
func Segment(rw http.ResponseWriter, r *http.Request) {
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

	seg, err := res.Segment(id, src, n)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer seg.Close()

	// respond with segment
	_, err = io.Copy(rw, seg)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "video/MP2T")
}
