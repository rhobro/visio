package mp4

import (
	"os/exec"
	"strconv"
)

func Split(path string, chunkKBS int) error {
	// TODO use go bindings for mp4box
	cmd := exec.Command("/Applications/GPAC.app/Contents/MacOS/mp4box", "-splits", strconv.Itoa(chunkKBS), path)
	return cmd.Run()
}
