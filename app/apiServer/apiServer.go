package apiServer

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/heartbeat"
	"OSS/app/apiServer/locate"
	"OSS/app/apiServer/objects"
	"log"
	"net/http"
)

func Run(cfgFile string) {
	err := config.InitCfg(cfgFile)
	if err != nil {
		panic(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Listening on :", config.ServerCfg.Server.Address)

	go heartbeat.Heartbeat()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/version/", objects.VersionHandler)
	log.Fatal(http.ListenAndServe(config.ServerCfg.Server.Address, nil))
}
