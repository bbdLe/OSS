package heartbeat

import (
	"OSS/app/apiServer/config"
	"OSS/comm/rabbitmq"
	"strconv"
	"sync"
	"time"
)

var (
	dataServers map[string]time.Time
	mutex sync.Mutex
)

func init() {
	dataServers = make(map[string]time.Time)
}

func Heartbeat() {
	mq := rabbitmq.New(config.ServerCfg.Server.RabbitMq)
	defer mq.Close()
	mq.Bind("apiServer")
	go removeExpiredDataServer()
	c := mq.Cosume()
	for msg := range c {
		dataAddr, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		dataServers[dataAddr] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(time.Second * 5)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(time.Second * 10).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}
