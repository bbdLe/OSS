package objects

import (
	"OSS/app/apiServer/config"
	"OSS/comm/es"
	"log"
	"net/http"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	metaData, err := es.SearchLatestVersion(config.ServerCfg.Server.ES, name)
	if err != nil {
		log.Println("SearchLatestVersion failed :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 不存在或者已经被删除
	if metaData.Version == 0 || metaData.Hash == "" {
		return
	}

	err = es.PutMetaData(config.ServerCfg.Server.ES, name, metaData.Version + 1, 0,"")
	if err != nil {
		log.Println("PutMetaData fail :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
