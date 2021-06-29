package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/platform"
	"github.com/rhobro/visio/internal/server/handler"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	platform.Init()
}

func main() {
	// TODO use fasthttp
	r := mux.NewRouter()
	r.HandleFunc("/upload", handler.Upload).Methods(http.MethodPost)
	r.HandleFunc("/x/{id}/root.m3u8", handler.Master).Methods(http.MethodGet)
	r.HandleFunc("/x/{id}/{src}.m3u8", handler.Playlist).Methods(http.MethodGet)
	r.HandleFunc("/x/{id}/{src}/{n}.ts", handler.Segment).Methods(http.MethodGet)

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
