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
	"360p":   Res360p,
	"540p":   Res540p,
	"720p":   Res720p,
	"1080p":  Res1080p,
	"2160p":  Res2160p,
	"4320p":  Res4320p,
	"source": ResSource,
}

func Resolution(res string) Res {
	return resolutions[res]
}
