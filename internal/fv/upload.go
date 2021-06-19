package fv

import (
	"io"
	"math/rand"
	"os"
	"strings"
)

func Upload(r io.Reader, id string) (string, error) {
	uploader, err := New()
	if err != nil {
		return "", err
	}

	// upload data
	asset, err := uploader.Upload(&Params{
		Body: r,
		ID:   id,
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
	return Upload(f, randID())
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randID() string {
	id := strings.Builder{}

	for i := 0; i < 10+rand.Intn(10); i++ {
		id.WriteByte(chars[rand.Intn(len(chars))])
	}

	return id.String() + "-mp4"
}
