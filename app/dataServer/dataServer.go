package dataServer

import (
	"OSS/app/dataServer/config"
)

func Run(cfgFile string) {
	err := config.InitCfg(cfgFile)
	if err != nil {
		panic(err)
	}
}
