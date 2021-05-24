package main

import (
	"github.com/gorilla/mux"
	"github.com/rhobro/visio/internal/platform"
	handler2 "github.com/rhobro/visio/internal/server/handler"
	"log"
	"net/http"
)

func init() {
	platform.Init()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/add", handler2.Add)
	r.HandleFunc("/{id}/root.m3u8", handler2.RootM3U8)
	r.HandleFunc("/{id}/{source}.m3u8", handler2.SourceM3U8)
	r.HandleFunc("/{id}/{source}/{n}.ts", handler2.TS)

	err := http.ListenAndServe(":1580", r)
	if err != nil {
		log.Fatalf("listening and serving: %s", err)
	}
}
