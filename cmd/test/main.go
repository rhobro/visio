package main

import (
	"fmt"
	"github.com/etherlabsio/go-m3u8/m3u8"
)

func main() {
	pl, _ := m3u8.ReadString(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=4000000,RESOLUTION=0x0
81ce24e1-fbb0-4e14-b0e5-82335e7b463e/source.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=6000000,RESOLUTION=1920x1080
81ce24e1-fbb0-4e14-b0e5-82335e7b463e/1080p.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=2000000,RESOLUTION=1280x720
81ce24e1-fbb0-4e14-b0e5-82335e7b463e/720p.m3u8
#EXT-X-STREAM-INF:PROGRAM-ID=0,BANDWIDTH=500000,RESOLUTION=640x360
81ce24e1-fbb0-4e14-b0e5-82335e7b463e/360p.m3u8`)

	fmt.Println(pl.Items)
	fmt.Println()
	fmt.Println(*pl.Version)
	fmt.Println()
	fmt.Println(pl.Playlists())

	_ = m3u8.PlaylistItem{
		ProgramID: 0,
		Resolution: &m3u8.Resolution{
			Width:  1920,
			Height: 1080,
		},
		Bandwidth: 6000000, // TODO research

		URI: "source.m3u8",
	}
}

var bandwidths = map[string]int{}
