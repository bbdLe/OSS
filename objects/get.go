package objects

import (
	"OSS/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(path.Join(config.ServerCfg.Server.StoragePath,
		"objects", strings.Split(r.URL.EscapedPath(), "/")[2]))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

