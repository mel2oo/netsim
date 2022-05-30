package config

import (
	"netsim/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   logger.Config `mapstructure:"logger"`
	Listener []Listener    `mapstructure:"listener"`
}

type Listener struct {
	Transport string `mapstructure:"transport"`
	Protocol  string `mapstructure:"protocol"`
	Address   string `mapstructure:"address"`
	Size      int    `mapstructure:"size"`
	RTimeout  int    `mapstructure:"rtimeout"`
	WTimeout  int    `mapstructure:"wtimeout"`
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
