package main

import (
	"encoding/json"
	"fmt"
	"github.com/etherlabsio/go-m3u8/m3u8"
	"github.com/rhobro/goutils/pkg/httputil"
	"github.com/rhobro/visio/internal/platform"
	"github.com/rhobro/visio/pkg/visio"
	"log"
	"math/rand"
	"net/http"
	"path"
	"strings"
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

func main() {
	roots := []string{
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeiftigxs45jqsipv6evhkx7avfckllw6t47hx7ebjjp3hycm644k4y/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeigxbgufn5k47xsogftqkmpikxdzj5tthif3ql74tbz6mfsdgbsq7u/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeibwj2lakyolvthy7x7hlcjvk4bs62xi7igr37vjgbtxlfvaw4mc4e/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeidv72l7n3cfrfyvlfpwoidamwvzsarlctzcbsqfmm5h42z32ev5tm/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeibsloou2vj7eoj25dr3lzg5nnowurjq6afetuylyuz2vc4v6nrlgu/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeie2djlm5s4wmsr63aszrmpdb23im5v5pokf6safbvdaa6esutsqoy/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeign57zea6rrg5ba2qk3n4tm4boywp7x6toqycjq3eklunumueqycy/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeicxplvwbjeb2fxndla73h72ais4554mp2g2kp22rebfjwur7acfje/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeid7dcxplhfsuitqvczo4pciiqumjwdn3ah2waboz63cekercinjri/root.m3u8",
		"https://demuxipfsrevproxy.onrender.com/ipfs/bafybeib2apxj2sptnnf5o3fjst7hwolqkdtypnrfoej26uyl4bb3yvxtzi/root.m3u8",
	}

	newP := make(visio.Video)

	for _, pURL := range roots {
		// request master root playlist
		rq, _ := http.NewRequest(http.MethodGet, pURL, nil)
		rq.Header.Set("User-Agent", httputil.RandUA())
		rsp, err := http.DefaultClient.Do(rq)
		if err != nil {
			log.Fatal(err)
		}
		superP, err := m3u8.Read(rsp.Body)
		if err != nil {
			log.Fatal(err)
		}
		rsp.Body.Close()

		// map and request source playlists
		base := pURL[:strings.LastIndex(pURL, "/")+1]
		for _, sourceEntryP := range superP.Playlists() {
			sourceURL := base + sourceEntryP.URI
			key := path.Base(sourceEntryP.URI)
			key = key[:strings.Index(key, ".m3u8")]
			base := sourceURL[:strings.LastIndex(sourceURL, "/")+1]

			// request playlist
			rq, _ := http.NewRequest(http.MethodGet, sourceURL, nil)
			rq.Header.Set("User-Agent", httputil.RandUA())
			rsp, err := http.DefaultClient.Do(rq)
			if err != nil {
				log.Fatal(err)
			}
			sourceP, err := m3u8.Read(rsp.Body)
			if err != nil {
				log.Fatal(err)
			}
			rsp.Body.Close()

			// loop through .ts
			for _, seg := range sourceP.Segments() {
				// add to restructured playlist
				newP[key] = append(newP[key], base+seg.Segment)
			}
		}
	}

	bd, _ := json.MarshalIndent(&newP, "", "\t")
	fmt.Println(string(bd))
	//rq, _ := http.NewRequest(http.MethodPost, "http://localhost:1580/add", bytes.NewReader(bd))
	//rq.Header.Set("id", "test")
	//rsp, err := http.DefaultClient.Do(rq)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//bd, _ = io.ReadAll(rsp.Body)
	//fmt.Println(string(bd))
}

//var path = flag.Arg(0)
var pth = "/Users/robro/Desktop/composite.mp4"

func maind() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			file, err := visio.UploadFile(pth)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(file)
		}()
	}

	wg.Wait()
}
