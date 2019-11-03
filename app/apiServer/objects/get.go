package objects

import (
	"OSS/app/apiServer/locate"
	"OSS/comm/httpstream"
	"io"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	server := locate.Locate(object)
	if server == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	stream, err := httpstream.NewGetStream(server,  object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.Copy(w, stream)
}
