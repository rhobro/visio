package fv

import (
	"github.com/etherlabsio/go-m3u8/m3u8"
	"github.com/rhobro/goutils/pkg/httputil"
	"github.com/rhobro/visio/pkg/visio"
	"net/http"
	"path"
	"strings"
)

// Visify takes a list of m3u8 root URLs hosted on file.video and converts them
// into a visio.Video
func Visify(roots []string) (*visio.Video, error) {
	v := make(visio.Video)

	for _, rootURL := range roots {
		// request master root playlist
		rq, _ := http.NewRequest(http.MethodGet, rootURL, nil)
		rq.Header.Set("User-Agent", httputil.RandUA())
		rsp, err := http.DefaultClient.Do(rq)
		if err != nil {
			return nil, err
		}

		root, err := m3u8.Read(rsp.Body)
		if err != nil {
			return nil, err
		}
		rsp.Body.Close()

		// map and request source playlists
		base := rootURL[:strings.LastIndex(rootURL, "/")+1]

		for _, playlist := range root.Playlists() {
			playlistURL := base + playlist.URI
			key := path.Base(playlist.URI)
			key = key[:strings.Index(key, ".m3u8")]
			base := playlistURL[:strings.LastIndex(playlistURL, "/")+1]

			// request playlist
			rq, _ := http.NewRequest(http.MethodGet, playlistURL, nil)
			rq.Header.Set("User-Agent", httputil.RandUA())
			rsp, err := http.DefaultClient.Do(rq)
			if err != nil {
				return nil, err
			}
			playlist, err := m3u8.Read(rsp.Body)
			if err != nil {
				return nil, err
			}
			rsp.Body.Close()

			// loop through .ts
			for _, seg := range playlist.Segments() {
				// add to restructured playlist
				v[key] = append(v[key], base+seg.Segment)
			}
		}
	}

	return &v, nil
}
