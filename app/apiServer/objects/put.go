package objects

import (
	"OSS/app/apiServer/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	f, err := os.Create(path.Join(config.ServerCfg.Server.StoragePath,
		"objects", strings.Split(r.URL.EscapedPath(), "/")[2]))
	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
}
