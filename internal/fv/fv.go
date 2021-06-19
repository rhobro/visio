package fv

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/rhobro/goutils/pkg/httputil"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"rogchap.com/v8go"
	"strconv"
	"strings"
)

const (
	fvURL    = "https://file.video"
	assetURL = "https://demux.onrender.com/asset"
)

const chunkSize = 2 << (20 - 1)

type FV struct {
	auth  string
	asset *demux
}

// New FV upload object
func New() (*FV, error) {
	var auth string

	// get auth key
	rq, _ := http.NewRequest(http.MethodGet, fvURL, nil)
	rq.Header.Set("User-Agent", httputil.RandUA())
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	page, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}

	// find index.js file
	page.Find(`link[rel="preload"][as="script"]`).Each(func(_ int, sl *goquery.Selection) {
		href := sl.AttrOr("href", "")
		if path.Base(href) == "index.js" {
			rq, _ := http.NewRequest(http.MethodGet, fvURL+href, nil)
			rq.Header.Set("User-Agent", httputil.RandUA())
			rsp, err := http.DefaultClient.Do(rq)
			if err != nil {
				log.Print(err)
				return
			}
			defer rsp.Body.Close()
			bd, err := io.ReadAll(rsp.Body)
			if err != nil {
				log.Print(err)
				return
			}

			// find section with JS expression
			body := string(bd)
			body = body[strings.Index(body, `"Basic "`)+16:]
			body = body[:strings.Index(body, `,"utf-8"`)]

			// run JS expression
			ctx, err := v8go.NewContext(nil)
			if err != nil {
				log.Print(err)
				return
			}
			val, err := ctx.RunScript(body, "decentriflix.js")
			if err != nil {
				log.Print(err)
				return
			}

			// b64 encode + concat
			auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(val.String()))
			ctx.Close()
		}
	})

	return &FV{
		auth: auth,
	}, nil
}

// Upload io.Reader in Params
// Returns an Asset object with upload details
func (fv *FV) Upload(up *Params) (*Asset, error) {
	// get url from Asset
	rq, _ := http.NewRequest(http.MethodPost, assetURL, nil)
	rq.Header.Set("User-Agent", httputil.RandUA())
	rq.Header.Set("Authorization", fv.auth)
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	bd, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	// unmarshal
	err = json.Unmarshal(bd, &fv.asset)
	if err != nil {
		return nil, err
	}

	// upload sequence
	err = fv.upload(up)
	if err != nil {
		return nil, err
	}

	// check for errors
	errored, err := fv.isErrored()
	if err != nil {
		return nil, err
	}
	if errored {
		return nil, errors.New("upload has errored")
	}

	// wait until ready
	return fv.waitAsset()
}

// Performs raw chunking and uploading with mpart
func (fv *FV) upload(up *Params) error {
	// upload parameters
	params := uploadParams{
		chunkN:    1,
		chunkSize: chunkSize,
		nChunks:   1,
		fType:     "video/mp4",
		id:        up.ID,
	}

	// multipart
	buf := &bytes.Buffer{}
	mPart := multipart.NewWriter(buf)

	file, _ := mPart.CreateFormFile("file", up.ID+".mp4")
	size, err := io.Copy(file, up.Body)
	if err != nil {
		return err
	}
	params.currentChunkSize = size
	params.totalSize = size

	rChunkN, _ := mPart.CreateFormField("resumableChunkNumber")
	_, err = rChunkN.Write([]byte("1"))
	if err != nil {
		return err
	}

	rChunkSize, _ := mPart.CreateFormField("resumableChunkSize")
	_, err = rChunkSize.Write([]byte(strconv.FormatInt(params.chunkSize, 10)))
	if err != nil {
		return err
	}

	rNChunks, _ := mPart.CreateFormField("resumableTotalChunks")
	_, err = rNChunks.Write([]byte(strconv.Itoa(params.nChunks)))
	if err != nil {
		return err
	}

	rCurrentChunkSize, _ := mPart.CreateFormField("resumableCurrentChunkSize")
	_, err = rCurrentChunkSize.Write([]byte(strconv.FormatInt(params.currentChunkSize, 10)))
	if err != nil {
		return err
	}

	rTotalSize, _ := mPart.CreateFormField("resumableTotalSize")
	_, err = rTotalSize.Write([]byte(strconv.FormatInt(params.totalSize, 10)))
	if err != nil {
		return err
	}

	rType, _ := mPart.CreateFormField("resumableType")
	_, err = rType.Write([]byte(params.fType))
	if err != nil {
		return err
	}

	rID, _ := mPart.CreateFormField("resumableIdentifier")
	_, err = rID.Write([]byte(params.id))
	if err != nil {
		return err
	}

	rFName, _ := mPart.CreateFormField("resumableFilename")
	_, err = rFName.Write([]byte(params.fName))
	if err != nil {
		return err
	}

	rRelPath, _ := mPart.CreateFormField("resumableRelativePath")
	_, err = rRelPath.Write([]byte(params.fName))
	if err != nil {
		return err
	}

	// close form
	mPart.Close()
	// url prep
	u, err := url.Parse(fv.asset.URL)
	if err != nil {
		return err
	}
	u.RawQuery = params.values().Encode()

	rq, _ := http.NewRequest(http.MethodPost, u.String(), buf)
	rq.Header.Add("Content-Type", mPart.FormDataContentType())
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	return nil
}

// checks if upload errored
func (fv *FV) isErrored() (bool, error) {
	rq, _ := http.NewRequest(http.MethodGet, "https://file.video/api/upload/"+fv.asset.ID, nil)
	rq.Header.Set("User-Agent", httputil.RandUA())
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return false, err
	}
	defer rsp.Body.Close()

	// unmarshal and parse
	var status struct {
		Upload struct {
			Status bool
			Errors bool
		}
	}
	bd, err := io.ReadAll(rsp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bd, &status)
	if err != nil {
		return false, err
	}

	return !status.Upload.Status && status.Upload.Errors, nil
}

// Wait until asset is prepared in the background. Returns the Asset information when ready
func (fv *FV) waitAsset() (*Asset, error) {
	var assets struct {
		Asset Asset
	}

	for !assets.Asset.Ready {
		rq, _ := http.NewRequest(http.MethodGet, "https://file.video/api/asset/"+fv.asset.ID, nil)
		rq.Header.Set("User-Agent", httputil.RandUA())
		rsp, err := http.DefaultClient.Do(rq)
		if err != nil {
			return nil, err
		}

		bd, err := io.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		rsp.Body.Close()
		err = json.Unmarshal(bd, &assets)
		if err != nil {
			return nil, err
		}
	}

	return &assets.Asset, nil
}
