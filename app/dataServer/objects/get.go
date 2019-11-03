package objects

import (
	"OSS/app/dataServer/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func isExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func get(w http.ResponseWriter, r *http.Request) {
	file := path.Join(config.ServerCfg.Server.StoragePath, "objects",
		strings.Split(r.URL.EscapedPath(), "/")[2])
	if !isExist(file) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := os.Open(file)
	if err != nil {
		log.Printf("open %s failed %v", file, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, f)
}
