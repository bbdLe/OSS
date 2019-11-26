package objects

import (
	"OSS/app/dataServer/config"
	"compress/gzip"
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
	defer f.Close()
	r2, err := gzip.NewReader(f)
	if err != nil {
		log.Println("gzip.NewReader failed :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r2.Close()
	io.Copy(w, r2)
}
