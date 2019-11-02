package main

import (
	"OSS/config"
	"OSS/objects"
	"flag"
	"log"
	"net/http"
)

var (
	cfgPath string
)

func main() {
	flag.StringVar(&cfgPath, "config", "etc/config.toml", "config file")
	flag.Parse()

	if err := config.InitConfig(cfgPath, &config.ServerCfg); err != nil {
		log.Fatal(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Listening on :", config.ServerCfg.Server.Address)

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(config.ServerCfg.Server.Address, nil))
}
