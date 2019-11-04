package dataServer

import (
	"OSS/app/dataServer/config"
	"OSS/app/dataServer/heartbeat"
	"OSS/app/dataServer/locate"
	"OSS/app/dataServer/objects"
	"OSS/app/dataServer/temp"
	"log"
	"net/http"
)

func Run(cfgFile string) {
	err := config.InitCfg(cfgFile)
	if err != nil {
		panic(err)
	}

	go heartbeat.Heartbeat()
	go locate.Locate()

	log.Println("Listening on", config.ServerCfg.Server.Address)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(config.ServerCfg.Server.Address, nil))
}