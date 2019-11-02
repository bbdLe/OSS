package dataServer

import (
	"OSS/app/dataServer/config"
	"OSS/app/dataServer/heartbeat"
	"time"
)

func Run(cfgFile string) {
	err := config.InitCfg(cfgFile)
	if err != nil {
		panic(err)
	}

	go heartbeat.Heartbeat()
	for {
		time.Sleep(time.Second)
	}
}