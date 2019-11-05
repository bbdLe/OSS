package temp

import (
	"OSS/app/dataServer/config"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	uuidPath := path.Join(config.ServerCfg.Server.StoragePath, "temp", uuid)
	info, err := readTempInfoFromFile(uuidPath)
	if err != nil {
		log.Println("readTempInfoFromFile err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uuidDataPath := uuidPath + ".dat"
	f, err := os.OpenFile(uuidDataPath, os.O_WRONLY | os.O_APPEND, 0)
	if err != nil {
		log.Println("OpenFile file : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Println("io.Copy failed :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stat, err := os.Stat(uuidDataPath)
	if err != nil {
		log.Println("os.Stat failed ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if stat.Size() > info.Size {
		os.Remove(uuidPath)
		os.Remove(uuidDataPath)
		log.Printf("except size : %d, but receive %d\n", info.Size, stat.Size())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func readTempInfoFromFile(uuidPath string) (*tempInfo, error) {
	f, err := os.Open(uuidPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var info tempInfo
	err = json.Unmarshal(b, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
