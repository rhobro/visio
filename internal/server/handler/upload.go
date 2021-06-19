package handler

import (
	"bytes"
	"encoding/json"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/internal/fv"
	"net/http"
	"path"
)

// Upload video and host on file.video
func Upload(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.Header.Get("ID")
	if id == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// if data present
	if r.Body == nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// order of m3u8 files
	var roots []string

	m3u8URL, err := fv.Upload(r.Body, id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	roots = append(roots, m3u8URL)

	// visify into visio compatible
	video, err := fv.Visify(roots)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// format as minified JSON
	refmt, err := json.Marshal(&video)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// upload to Storj
	err = storaje.Upload(bytes.NewReader(refmt), "videos", id+".json")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if no errors, return URL to video
	url := path.Join(r.URL.Host, "x", id, "root.m3u8")
	_, _ = rw.Write([]byte(url))
}
