package temp

import (
	"OSS/app/dataServer/config"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	output, err := exec.Command("uuidgen").Output()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	uuid := strings.TrimSuffix(string(output), "\n")
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println("ParseInt failed :", r.Header.Get("size"), ",", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	t := tempInfo{uuid, name, size}
	if err = t.WriteToFile(); err != nil {
		log.Println("WriteToFile failed : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	f, err := os.Create(path.Join(config.ServerCfg.Server.StoragePath, "temp", t.Uuid) + ".dat")
	if err != nil {
		log.Println("create file fail :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	w.Write([]byte(uuid))
}

func (t *tempInfo) WriteToFile() error {
	f, err := os.Create(path.Join(config.ServerCfg.Server.StoragePath, "temp", t.Uuid))
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	_, err = f.Write(b)

	return err
}

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}
