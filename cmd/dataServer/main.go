package main

import (
	"OSS/app/dataServer"
	"flag"
)

var (
	cfgPath string
)

func main() {
	flag.StringVar(&cfgPath, "config", "config.toml", "config file")
	flag.Parse()

	dataServer.Run(cfgPath)
}
