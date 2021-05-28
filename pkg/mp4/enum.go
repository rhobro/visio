package mp4

type Res int

const (
	Res360p Res = iota
	Res540p
	Res720p
	Res1080p
	Res2160p
	Res4320p
	ResSource
)

var resolutions = map[string]Res{
	"360p.m3u8":   Res360p,
	"540p.m3u8":   Res540p,
	"720p.m3u8":   Res720p,
	"1080p.m3u8":  Res1080p,
	"2160p.m3u8":  Res2160p,
	"4320p.m3u8":  Res4320p,
	"source.m3u8": ResSource,
}

func Resolution(res string) Res {
	return resolutions[res]
}
