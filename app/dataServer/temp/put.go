package temp

import (
	"OSS/app/dataServer/config"
	"OSS/app/dataServer/locate"
	"compress/gzip"
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

	commitFile(uuidDataPath, info)
}

func commitFile(tempFile string, info *tempInfo) {
	f, _ := os.Open(tempFile)
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	hash := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	f.Seek(0, io.SeekStart)
	log.Println(path.Join(config.ServerCfg.Server.StoragePath, "objects", info.Name))
	w, _ := os.Create(path.Join(config.ServerCfg.Server.StoragePath, "objects", info.Name) + "." + hash)
	w2 := gzip.NewWriter(w)
	io.Copy(w2, f)
	w2.Close()
	os.Remove(tempFile)
	locate.Add(info.hash(), info.id())
}
