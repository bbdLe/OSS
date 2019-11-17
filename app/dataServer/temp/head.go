package temp

import (
	"OSS/app/dataServer/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, err := os.Open(path.Join(config.ServerCfg.Server.StoragePath, "temp", uuid) + ".dat")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", info.Size()))
}
