package res

import (
	"encoding/json"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/pkg/visio"
	"io"
)

func Video(bucket, id string) (*visio.Video, error) {
	vid, err := storaje.Download("videos", id+".json")
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(vid)
	if err != nil {
		return nil, err
	}

	// unmarshal into visio.Video
	src := new(visio.Video)
	err = json.Unmarshal(body, src)
	if err != nil {
		return nil, err
	}

	return src, nil
}
