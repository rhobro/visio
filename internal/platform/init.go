package platform

import (
	configcat "github.com/configcat/go-sdk/v7"
	"github.com/rhobro/goutils/pkg/fileio"
	"github.com/rhobro/goutils/pkg/services/cfgcat"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"log"
	"net/http"
	"os"
)

var Running = true

func Init() {
	cfgcat.InitCustom(configcat.Config{
		SDKKey:    "CR7ZCKLIe0OJIFp0hHbqsA/WKLihHgrhEiW-7xYfrz0Eg",
		Transport: &http.Transport{},
	}, false)

	err := storaje.Init(
		cfgcat.C.GetStringValue("storjSatellite", "", nil),
		cfgcat.C.GetStringValue("storjKey", "", nil))
	if err != nil {
		log.Fatalf("unable to connect to Storj: %s", err)
	}

	fileio.Init("", "visio_server_*")

	// test
	//u, _ := url.Parse("http://localhost:9090")
	//http.DefaultTransport = &http.Transport{
	//	Proxy: http.ProxyURL(u),
	//}
}

func Close() {
	Running = false
	cfgcat.C.Close()
	storaje.C.Close()
	fileio.Close()
	os.Exit(0)
}
