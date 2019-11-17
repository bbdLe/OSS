package temp

import (
	"OSS/app/dataServer/config"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, err := os.Open(path.Join(config.ServerCfg.Server.StoragePath, "temp", uuid) + ".dat")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
