package main

import (
	"bufio"
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/platform"
	"github.com/rhobro/visio/internal/server/cache"
	"github.com/rhobro/visio/internal/server/handler"
	"log"
	"net/http"
	"os"
)

func init() {
	platform.Init()
	cache.Start()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/add", handler.Add)
	r.HandleFunc("/{id}/root.m3u8", handler.RootM3U8)
	r.HandleFunc("/{id}/{source}.m3u8", handler.SourceM3U8)
	r.HandleFunc("/{id}/{source}/{n}.ts", handler.TS)

	// wait to stop
	go func() {
		rd := bufio.NewScanner(os.Stdin)

		for rd.Scan() {
			if rd.Text() == "q" {
				platform.Close()
			}
		}
	}()

	err := http.ListenAndServe(":1580", r)
	if err != nil {
		log.Fatalf("listening and serving: %s", err)
	}
}
