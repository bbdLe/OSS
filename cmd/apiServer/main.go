package main

import (
	"OSS/app/apiServer"
	"flag"
)

var (
	cfgPath string
)

func main() {
	flag.StringVar(&cfgPath, "config", "config.toml", "config file")
	flag.Parse()

	apiServer.Run(cfgPath)
}
