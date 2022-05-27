package config

type Config struct {
	Listener Listener `mapstructure:"listener"`
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
