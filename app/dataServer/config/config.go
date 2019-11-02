package config

import "OSS/comm/configreader"

type server struct {
	Address string
	StoragePath string `mapstructure:"storage_path"`
	RabbitMq string
}

type config struct {
	Server server
}

var ServerCfg config

func InitCfg(fileName string) error {
	return configreader.InitConfig(fileName, &ServerCfg)
}
