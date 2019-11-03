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

func put(w http.ResponseWriter, r *http.Request) {
	file := path.Join(config.ServerCfg.Server.StoragePath, "objects",
		strings.Split(r.URL.EscapedPath(), "/")[2])
	f, err := os.Create(file)
	if err != nil {
		log.Println("Create", file, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(f, r.Body)
}
