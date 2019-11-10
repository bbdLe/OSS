package objects

import (
	"OSS/app/dataServer/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func isExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	files, _ := filepath.Glob(path.Join(config.ServerCfg.Server.StoragePath, "objects",
		name) + ".*")
	if len(files) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	file := files[0]

	f, err := os.Open(file)
	if err != nil {
		log.Printf("open %s failed %v", file, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, f)
}
