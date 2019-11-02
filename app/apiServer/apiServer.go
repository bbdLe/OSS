package apiServer

import (
	"OSS/app/apiServer/config"
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

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(config.ServerCfg.Server.Address, nil))
}
