package main

import (
	"fmt"
	"github.com/rhobro/visio/internal/fv"
	"github.com/rhobro/visio/internal/platform"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func init() {
	//u, _ := url.Parse("socks5://localhost:9050")
	http.DefaultTransport = &http.Transport{
		MaxIdleConns: 1,
		//Proxy:        http.ProxyURL(u),
	}
	rand.Seed(time.Now().UnixNano())

	platform.Init()
}

var pth = "/Users/robro/Desktop/composite.mp4"

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			file, err := fv.UploadFile(pth)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(file)
		}()
	}

	wg.Wait()
}
