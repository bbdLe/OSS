package objects

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/heartbeat"
	"OSS/app/apiServer/locate"
	"OSS/comm/es"
	"OSS/comm/rs"
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

	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
