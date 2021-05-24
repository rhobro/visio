package handler

import (
	"bytes"
	"encoding/json"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/pkg/visio"
	"io"
	"log"
	"net/http"
	"path"
)

func Add(rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := rq.Header.Get("id")
	if id == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// if data present and readable
	if rq.Body == nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	body, err := io.ReadAll(rq.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if data is JSON formatted correctly
	var data visio.Source
	err = json.Unmarshal(body, &data)
	if err != nil {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// if structured properly with the correct type of data
	if !data.IsValid() {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// upload to Storj after formatting as minified JSON
	refmt, err := json.Marshal(&data)
	if err != nil {
		log.Printf("marshaling JSON: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = storaje.Upload(bytes.NewReader(refmt), "videos", id)
	if err != nil {
		log.Printf("uploading to Storj: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if no errors, return URL to video
	url := path.Join(rq.URL.Host, id, "root.m3u8")
	_, _ = rw.Write([]byte(url))
}
