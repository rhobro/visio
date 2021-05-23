package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"visio/pkg/visio"
)

func init() {
	http.DefaultTransport = &http.Transport{
		MaxIdleConns: 1,
	}
	rand.Seed(time.Now().UnixNano())
}

//var path = flag.Arg(0)
var path = "/Users/robro/Desktop/composite.mp4"

func main() {
	m3u8, err := visio.UploadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m3u8)
}
