package objects

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/locate"
	"OSS/comm/es"
	"OSS/comm/httpstream"
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

	server := locate.Locate(url.PathEscape(metaData.Hash))
	if server == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stream, err := httpstream.NewGetStream(server, url.PathEscape(metaData.Hash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.Copy(w, stream)
}
