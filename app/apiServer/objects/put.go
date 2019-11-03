package objects

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/heartbeat"
	"OSS/comm/es"
	"OSS/comm/httpstream"
	"OSS/comm/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("Wrong request : need hash val")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status, err := storeObject(r.Body, url.PathEscape(hash))
	if err != nil {
		log.Printf("storeObject error : %v, code : %d", err, status)
		w.WriteHeader(status)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size := utils.GetSizeFromHeader(r.Header)
	err = es.AddVersion(config.ServerCfg.Server.ES, name, hash, size)
	if err != nil {
		log.Println("AddVersion failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func storeObject(reader io.Reader, name string) (int, error) {
	dataServer := heartbeat.GetRandDataServer()
	if dataServer == "" {
		log.Println("dataServer is empty")
		return http.StatusInternalServerError, fmt.Errorf("dataServer empty")
	}

	stream, err := httpstream.NewPutStream(dataServer, name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	io.Copy(stream, reader)
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
