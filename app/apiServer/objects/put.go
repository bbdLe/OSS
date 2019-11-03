package objects

import (
	"OSS/app/apiServer/heartbeat"
	"OSS/comm/httpstream"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	dataServer := heartbeat.GetRandDataServer()
	if dataServer == "" {
		log.Println("dataServer is empty")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream, err := httpstream.NewPutStream(dataServer, strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil {
		log.Println("NewPutStream", dataServer)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.Copy(stream, r.Body)
	err = stream.Close()
	if err != nil {
		log.Println("stream err :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
