package objects

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/heartbeat"
	"OSS/app/apiServer/locate"
	"OSS/comm/es"
	"OSS/comm/rs"
	"OSS/comm/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println("ParseInt size failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing hash from header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if locate.Exist(name) {
		err := es.AddVersion(config.ServerCfg.Server.ES, name, hash, size)
		if err != nil {
			log.Println("es add version failed :", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	ds := heartbeat.ChooseRandomDataServers(rs.AllShares, nil)
	if len(ds) != rs.AllShares {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	stream, err := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if err != nil {
		log.Println("NewRSResumablePutStream failed :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/temp/" + url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}