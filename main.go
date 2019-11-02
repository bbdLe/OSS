package main

import (
	"OSS/config"
	"flag"
	"fmt"
	"log"
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

	fmt.Printf("%v", config.ServerCfg.Server)
}
