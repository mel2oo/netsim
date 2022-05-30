package config

import (
	"netsim/pkg/logger"

	"github.com/spf13/viper"
)

type Config struct {
	Logger    logger.Config `mapstructure:"logger"`
	TLS       TLS           `mapstructure:"tls"`
	Reslover  Reslover      `mapstructure:"reslover"`
	Forwarder Forwarder     `mapstructure:"forwarder"`
	Listener  []Listener    `mapstructure:"listener"`
}

type TLS struct {
	Ca   string `mapstructure:"ca"`
	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
}

type Reslover struct {
	Dns       string   `mapstructure:"dns"`
	DnsServer []string `mapstructure:"dnsserver"`
	Timeout   int      `mapstructure:"timeout"`
	MaxTTL    int      `mapstructure:"maxTTL"`
	MinTTL    int      `mapstructure:"minTTL"`
	CacheSize int      `mapstructure:"cache"`
}

type Forwarder struct {
	Mode     string `mapstructure:"mode"`
	DTimeout int    `mapstructure:"dial-timeout"`
	RTimeout int    `mapstructure:"relay-timeout"`
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
