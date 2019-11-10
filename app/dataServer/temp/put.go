package temp

import (
	"OSS/app/dataServer/config"
	"OSS/app/dataServer/locate"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
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

	log.Println("recv commit")
	commitFile(uuidDataPath, info)
}

func commitFile(tempFile string, info *tempInfo) {
	f, _ := os.Open(tempFile)
	h := sha256.New()
	io.Copy(h, f)
	hash := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	err := os.Rename(tempFile, path.Join(config.ServerCfg.Server.StoragePath, "objects", info.Name) + "." + hash)
	if err != nil {
		log.Println(err)
	}
	locate.Add(info.hash(), info.id())
}
