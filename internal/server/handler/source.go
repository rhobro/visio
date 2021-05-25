package handler

import (
	"github.com/gorilla/mux"
	"net/http"
)

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
}
