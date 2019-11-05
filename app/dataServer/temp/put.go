package temp

import (
	"OSS/app/dataServer/config"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	uuidPath := path.Join(config.ServerCfg.Server.StoragePath, "temp", uuid)
	info, err := readTempInfoFromFile(uuidPath)
	if err != nil {
		log.Println("readTempInfoFromFile failed :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uuidDataPath := uuidPath + ".dat"
	stat, err := os.Stat(uuidDataPath)
	if err != nil {
		log.Println("open filed failed :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	os.Remove(uuidPath)
	if stat.Size() != info.Size {
		log.Printf("except file size : %d, but size is %d\n", info.Size, stat.Size())
		os.Remove(uuidDataPath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	commitFile(uuidDataPath, info)
}

func commitFile(tempFile string, info *tempInfo) {
	os.Rename(tempFile, path.Join(config.ServerCfg.Server.StoragePath, "objects", info.Name))
}
