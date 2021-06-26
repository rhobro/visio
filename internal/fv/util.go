package fv

import (
	"github.com/rhobro/goutils/pkg/services/storaje"
	"net/url"
	"strconv"
)

func IdExists(id string) bool {
	if id == "" {
		return true
	}

	ids := storaje.ListObjects("videos")

	for ids.Next() {
		if ids.Item().Key == id+".json" {
			return true
		}
	}
	return false
}

// Asset successful upload details
type Asset struct {
	ID           string
	Status       int
	Errors       bool
	M3U8URL      string `json:"stream_url"`
	ThumbnailURL string
	Ready        bool
}

// JSON matched struct for upload URL
type demux struct {
	ID  string `json:"asset_id"`
	URL string `json:"url"`
}

// params to upload chunks to fv
type uploadParams struct {
	chunkN           int
	currentChunkSize int64
	chunkSize        int64
	nChunks          int
	totalSize        int64

	fType string
	id    string
	fName string
}

func (up uploadParams) values() url.Values {
	// add url encoded values
	v := url.Values{}

	v.Set("resumableChunkNumber", strconv.Itoa(up.chunkN))
	v.Set("resumableChunkSize", strconv.FormatInt(up.chunkSize, 10))
	v.Set("resumableCurrentChunkSize", strconv.FormatInt(up.currentChunkSize, 10))
	v.Set("resumableTotalChunks", strconv.Itoa(up.nChunks))
	v.Set("resumableTotalSize", strconv.FormatInt(up.totalSize, 10))

	v.Set("resumableType", up.fType)
	v.Set("resumableIdentifier", up.id)
	v.Set("resumableFilename", up.fName)
	v.Set("resumableRelativePath", up.fName)

	return v
}
