package config

import (
	"github.com/spf13/viper"
)

type server struct {
	Address string
	StoragePath string `mapstructure:"storage_path"`
}

type config struct {
	Server server
}

var ServerCfg config

func InitConfig(filepath string, cfg interface{}) error {
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}


