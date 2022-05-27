package config

import (
	"netsim/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   logger.Config `mapstructure:"logger"`
	Listener Listener      `mapstructure:"listener"`
}

type Listener struct {
	Tcp TCP `mapstructure:"tcp"`
	Udp UDP `mapstructure:"udp"`
}

type TCP struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type UDP struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

func Load(path string) (v *Config, err error) {
	viper.SetConfigFile(path)

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&v); err != nil {
		return nil, err
	}

	return
}
