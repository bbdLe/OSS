package locate

import (
	"OSS/app/dataServer/config"
	"OSS/comm/rabbitmq"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	mutex sync.Mutex
	objects = make(map[string]int)
)

func Add(hash string) {
	mutex.Lock()
	defer mutex.Unlock()
	objects[hash] = 1
}

func Del(hash string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(objects, hash)
}

func locate(hash string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := objects[hash]
	return ok
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
		if locate(file) {
			mq.Send(msg.ReplyTo, config.ServerCfg.Server.Address)
		}
	}
}

func CollectObject() {
	files, _ := filepath.Glob(config.ServerCfg.Server.StoragePath + "/objects/*")
	for i := range files {
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
