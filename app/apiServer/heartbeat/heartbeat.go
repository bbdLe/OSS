package heartbeat

import (
	"OSS/app/apiServer/config"
	"OSS/comm/rabbitmq"
	"math/rand"
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

func getAllDataServer() []string {
	mutex.Lock()
	defer mutex.Unlock()
	servers := make([]string, 0, len(dataServers))
	for s, _ := range dataServers {
		servers = append(servers, s)
	}
	return servers
}

func GetRandDataServer() string {
	servers := getAllDataServer()
	n := len(servers)
	if n == 0 {
		return ""
	} else {
		return servers[rand.Intn(n)]
	}
}

func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}

	servers := getAllDataServer()
	for i := range servers {
		server := servers[i]
		if _, exclude := reverseExcludeMap[server]; !exclude {
			candidates = append(candidates, server)
		}
	}

	length := len(candidates)
	if len(candidates) < n {
		return
	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}

	return
}
