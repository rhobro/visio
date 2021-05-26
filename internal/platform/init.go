package platform

import (
	"github.com/rhobro/goutils/pkg/fileio"
	"github.com/rhobro/goutils/pkg/services/cfgcat"
	"github.com/rhobro/goutils/pkg/services/storaje"
	"github.com/rhobro/visio/internal/server/cache"
	"log"
)

func Init() {
	cfgcat.Init("CR7ZCKLIe0OJIFp0hHbqsA/WKLihHgrhEiW-7xYfrz0Eg", false)

	err := storaje.Init(
		cfgcat.C.GetStringValue("storjSatellite", "", nil),
		cfgcat.C.GetStringValue("storjKey", "", nil))
	if err != nil {
		log.Fatalf("unable to connect to Storj: %s", err)
	}

	fileio.Init("", "visio_server_*")

	cache.Start()
}
