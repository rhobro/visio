package visio

import (
	"io"
	"math/rand"
	"os"
	"strings"
	"visio/internal/fv"
)

func Upload(r *io.Reader, id string) (string, error) {
	uploader, err := fv.New()
	if err != nil {
		return "", err
	}

	asset, err := uploader.Upload(&fv.Params{
		Body: r,
		ID: id,
	})
	if err != nil {
		return "", err
	}

	return asset.M3U8URL, nil
}

func UploadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var rd io.Reader = f
	return Upload(&rd, id())
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func id() string {
	id := strings.Builder{}

	for i := 0; i < 10 + rand.Intn(10); i++ {
		id.WriteByte(chars[rand.Intn(len(chars))])
	}

	return id.String() + "-mp4"
}
