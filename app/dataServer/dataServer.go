package dataServer

import (
	"OSS/app/dataServer/config"
	"OSS/app/dataServer/heartbeat"
	"OSS/app/dataServer/locate"
	"time"
)

func Run(cfgFile string) {
	err := config.InitCfg(cfgFile)
	if err != nil {
		panic(err)
	}

	go heartbeat.Heartbeat()
	go locate.Locate()
	for {
		time.Sleep(time.Second)
	}
}