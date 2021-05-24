package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/pkg/visio"
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

	// request
	video, err := storaje.Download("videos", id+".json")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	body, err := io.ReadAll(video)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// unmarshal into visio.Source
	var src visio.Source
	err = json.Unmarshal(body, &src)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// request clip
	rq, _ := http.NewRequest(http.MethodGet, src[source][n], nil)
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rsp.Body.Close()

	// respond with clip
	rw.Header().Set("Content-Type", "video/MP2T")
	_, err = io.Copy(rw, rsp.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
