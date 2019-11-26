package objects

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/heartbeat"
	"OSS/app/apiServer/locate"
	"OSS/comm/es"
	"OSS/comm/rs"
	"OSS/comm/utils"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version := 0
	versionID := r.URL.Query()["version"]

	var err error
	if len(versionID) != 0 {
		version, err = strconv.Atoi(versionID[0])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	metaData, err := es.GetMetadata(config.ServerCfg.Server.ES, name, version)
	if err != nil {
		log.Println("GetMetadata failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if metaData.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	locateInfos := locate.Locate(metaData.Hash)
	if len(locateInfos) < rs.DataShards {
		log.Printf("locateInfos less than rs.DataShards : %d(%s)\n", len(locateInfos), url.PathEscape(metaData.Hash))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dataServers := make([]string, 0)
	if len(locateInfos) < rs.AllShares {
		dataServers = heartbeat.ChooseRandomDataServers(rs.AllShares - len(locateInfos), locateInfos)
	}

	stream, err := rs.NewRSGetStream(locateInfos, dataServers, url.PathEscape(metaData.Hash), metaData.Size)
	defer stream.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n := utils.GetOffsetFromHeader(r.Header)
	if n != 0 {
		stream.Seek(n, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("%d-%d/%d", n, metaData.Size - 1, metaData.Size))
		w.WriteHeader(http.StatusPartialContent)
	}

	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}

	if acceptGzip {
		w.Header().Set("Accept-Encoding", "gzip")
		w2 := gzip.NewWriter(w)
		defer w2.Close()
		_, err := io.Copy(w2, stream)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	} else {
		_, err = io.Copy(w, stream)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
