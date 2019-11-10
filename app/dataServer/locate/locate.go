package locate

import (
	"OSS/app/dataServer/config"
	"OSS/comm/rabbitmq"
	"OSS/comm/types"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	mutex sync.Mutex
	objects = make(map[string]int)
)

func Add(hash string, id int) {
	mutex.Lock()
	defer mutex.Unlock()
	objects[hash] = id
}

func Del(hash string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(objects, hash)
}

func locate(hash string) int {
	mutex.Lock()
	defer mutex.Unlock()
	id, ok := objects[hash]
	if !ok {
		return -1
	} else {
		return id
	}
}

func Locate() {
	mq := rabbitmq.New(config.ServerCfg.Server.RabbitMq)
	defer mq.Close()
	mq.Bind("dataServer")
	c := mq.Cosume()
	for msg := range c {
		file, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		id := locate(file)
		if id != -1 {
			mq.Send(msg.ReplyTo, types.LocateMessage{Addr : config.ServerCfg.Server.Address,
				Id : id})
		}
	}
}

func CollectObject() {
	files, _ := filepath.Glob(config.ServerCfg.Server.StoragePath + "/objects/*")
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files)
		}

		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		objects[hash] = id
	}
}
