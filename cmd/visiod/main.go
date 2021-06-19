package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/platform"
	"github.com/rhobro/visio/internal/server/cache"
	"github.com/rhobro/visio/internal/server/handler"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	platform.Init()
	cache.Start()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", handler.Upload)
	r.HandleFunc("/x/{id}/root.m3u8", handler.Master)
	r.HandleFunc("/x/{id}/{src}.m3u8", handler.Playlist)
	r.HandleFunc("/x/{id}/{src}/{n}.ts", handler.Segment)

	// wait to stop
	go quitter()

	err := http.ListenAndServe(":1580", r)
	if err != nil {
		log.Fatalf("listening and serving: %s", err)
	}
}

func quitter() {
	rd := bufio.NewScanner(os.Stdin)

	for rd.Scan() {
		if rd.Text() == "q" {
			// check if purposeful
			fmt.Print("Are you sure you want to quit? (Y/n): ")
			rd.Scan()

			if strings.ToLower(rd.Text()) == "y" {
				platform.Close()
			}
		}
	}
}
