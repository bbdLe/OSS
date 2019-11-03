package objects

import (
	"OSS/app/apiServer/config"
	"OSS/comm/es"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	from := 0
	size := 1000
	for {
		metas, err := es.SearchAllVersion(config.ServerCfg.Server.ES, name, from, size)
		fmt.Println(len(metas))
		if err != nil {
			log.Println("SearchAllVersion error :", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, meta := range metas {
			b, _ := json.Marshal(meta)
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			break
		}
		from += size
	}
}
