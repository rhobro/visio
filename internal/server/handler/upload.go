package handler

import (
	"bytes"
	"encoding/json"
	"github.com/rhobro/goutils/pkg/fileio"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/internal/fv"
	"github.com/rhobro/visio/pkg/mp4"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const chunkSizeKBs = 30000

// Upload video and host on file.video
func Upload(rw http.ResponseWriter, r *http.Request) {
	// is ID present and if ID already been used
	id := r.Header.Get("ID")
	if fv.IdExists(id) {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// if data present
	if r.Body == nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	dir, err := fileio.TmpPath(id)
	videoPath := filepath.Join(dir, "video.mp4")

	// save file in temp directory
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	f, err := os.Create(videoPath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(f, r.Body)
	f.Close()

	// split
	err = mp4.Split(videoPath, chunkSizeKBs)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	// delete original
	err = os.Remove(videoPath)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// order of m3u8 files
	files, err := os.ReadDir(dir)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	var roots []string

	// loop through each file and upload
	for _, entry := range files {
		m3u8URL, err := fv.Upload(r.Body, id)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		roots = append(roots, m3u8URL)
	}

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
