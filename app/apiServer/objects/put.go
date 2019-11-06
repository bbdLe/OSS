package objects

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/heartbeat"
	"OSS/comm/es"
	"OSS/comm/httpstream"
	"OSS/comm/utils"
	"crypto/sha256"
	"encoding/base64"
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

	size := utils.GetSizeFromHeader(r.Header)

	status, err := storeObject(r.Body, url.PathEscape(hash), size)
	if err != nil {
		log.Printf("storeObject error : %v, code : %d", err, status)
		w.WriteHeader(status)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	err = es.AddVersion(config.ServerCfg.Server.ES, name, hash, size)
	if err != nil {
		log.Println("AddVersion failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func storeObject(reader io.Reader, name string, size int64) (int, error) {
	dataServer := heartbeat.GetRandDataServer()
	if dataServer == "" {
		log.Println("dataServer is empty")
		return http.StatusInternalServerError, fmt.Errorf("dataServer empty")
	}

	stream, err := httpstream.NewTempPutStream(dataServer, name, size)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	r := io.TeeReader(reader, stream)
	s := sha256.New()
	io.Copy(s, r)
	hash := base64.StdEncoding.EncodeToString(s.Sum(nil))
	if hash != name {
		stream.Commit(false)
		log.Println("Hash wrong")
		return http.StatusBadRequest, fmt.Errorf("Hash Worng")
	} else {
		stream.Commit(true)
		return http.StatusOK, nil
	}
}
